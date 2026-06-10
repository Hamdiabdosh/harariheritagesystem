package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/qirs-mezgeb/api/internal/models"
)

type mockRepository struct {
	user         *models.User
	userByID     *models.User
	refreshToken *models.RefreshToken
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.user == nil || m.user.Email != email {
		return nil, nil
	}
	return m.user, nil
}

func (m *mockRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if m.userByID == nil || m.userByID.ID != id {
		return nil, nil
	}
	return m.userByID, nil
}

func (m *mockRepository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	m.refreshToken = token
	return nil
}

func (m *mockRepository) GetRefreshToken(ctx context.Context, id uuid.UUID) (*models.RefreshToken, error) {
	if m.refreshToken == nil || m.refreshToken.ID != id {
		return nil, nil
	}
	return m.refreshToken, nil
}

func (m *mockRepository) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	if m.refreshToken == nil || m.refreshToken.ID != id {
		return ErrInvalidRefreshToken
	}
	m.refreshToken = nil
	return nil
}

func testUser(password string) *models.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &models.User{
		ID:           uuid.New(),
		FullName:     "Test User",
		Email:        "registrar@test.gov.et",
		PasswordHash: string(hash),
		Role:         models.RoleRegistrar,
		Language:     models.LanguageAm,
		IsActive:     true,
	}
}

func newTestService(repo RepositoryInterface, now time.Time) *Service {
	svc := NewService(repo, "test-access-secret-key-32chars!!", "test-refresh-secret-key-32chars!")
	svc.now = func() time.Time { return now }
	return svc
}

func TestLogin(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		repo    *mockRepository
		email   string
		pass    string
		wantErr error
	}{
		{
			name: "valid credentials returns token pair",
			repo: &mockRepository{user: testUser("password123")},
			email: "registrar@test.gov.et",
			pass:  "password123",
		},
		{
			name:    "wrong password returns unauthorized",
			repo:    &mockRepository{user: testUser("password123")},
			email:   "registrar@test.gov.et",
			pass:    "wrongpassword",
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "deactivated account returns forbidden",
			repo: func() *mockRepository {
				user := testUser("password123")
				user.IsActive = false
				return &mockRepository{user: user}
			}(),
			email:   "registrar@test.gov.et",
			pass:    "password123",
			wantErr: ErrAccountDeactivated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo, now)
			pair, err := svc.Login(context.Background(), tt.email, tt.pass)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if pair.AccessToken == "" || pair.RefreshToken == "" {
				t.Fatal("expected access and refresh tokens")
			}
			if pair.User.Email != tt.email {
				t.Fatalf("expected user email %s, got %s", tt.email, pair.User.Email)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	now := time.Now()
	user := testUser("password123")

	t.Run("valid refresh token returns new access token", func(t *testing.T) {
		repo := &mockRepository{user: user, userByID: user}
		svc := newTestService(repo, now)

		pair, err := svc.Login(context.Background(), user.Email, "password123")
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		accessToken, err := svc.Refresh(context.Background(), pair.RefreshToken)
		if err != nil {
			t.Fatalf("refresh failed: %v", err)
		}
		if accessToken == "" {
			t.Fatal("expected non-empty access token")
		}

		claims, err := svc.ParseAccessToken(accessToken)
		if err != nil {
			t.Fatalf("parse access token failed: %v", err)
		}
		if claims.UserID != user.ID {
			t.Fatalf("expected user id %s, got %s", user.ID, claims.UserID)
		}
	})

	t.Run("expired refresh token returns unauthorized", func(t *testing.T) {
		repo := &mockRepository{user: user, userByID: user}
		svc := newTestService(repo, now)

		pair, err := svc.Login(context.Background(), user.Email, "password123")
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		svc.now = func() time.Time { return now.Add(8 * 24 * time.Hour) }

		_, err = svc.Refresh(context.Background(), pair.RefreshToken)
		if err != ErrInvalidRefreshToken {
			t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
		}
	})
}
