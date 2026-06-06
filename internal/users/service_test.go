package users

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type mockRepository struct {
	users map[uuid.UUID]*models.User
}

func newMockRepository() *mockRepository {
	return &mockRepository{users: make(map[uuid.UUID]*models.User)}
}

func (m *mockRepository) seed(user *models.User) {
	m.users[user.ID] = user
}

func (m *mockRepository) List(ctx context.Context, filters ListFilters) ([]models.User, int, error) {
	users := make([]models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, len(users), nil
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	copy := *user
	return &copy, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			copy := *user
			return &copy, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, user *models.User) error {
	for _, existing := range m.users {
		if existing.Email == user.Email {
			return ErrEmailAlreadyExists
		}
	}
	user.ID = uuid.New()
	m.users[user.ID] = user
	return nil
}

func (m *mockRepository) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*models.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Language != nil {
		user.Language = *input.Language
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	copy := *user
	return &copy, nil
}

func (m *mockRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *mockRepository) UpdateLanguage(ctx context.Context, id uuid.UUID, language models.Language) (*models.User, error) {
	return m.Update(ctx, id, UpdateInput{Language: &language})
}

func TestCreateUserValidation(t *testing.T) {
	svc := NewService(newMockRepository())

	_, err := svc.Create(context.Background(), CreateInput{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "short",
		Role:     models.RoleRegistrar,
	})
	if err != ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestCreateDuplicateEmail(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateInput{
		FullName: "First User",
		Email:    "dup@example.com",
		Password: "password1",
		Role:     models.RoleRegistrar,
	})
	if err != nil {
		t.Fatalf("create first user: %v", err)
	}

	_, err = svc.Create(context.Background(), CreateInput{
		FullName: "Second User",
		Email:    "dup@example.com",
		Password: "password2",
		Role:     models.RoleRegistrar,
	})
	if err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestUpdateCannotDeactivateSelf(t *testing.T) {
	managerID := uuid.New()
	repo := newMockRepository()
	repo.seed(&models.User{
		ID:       managerID,
		FullName: "Manager",
		Email:    "manager@test.gov.et",
		Role:     models.RoleManager,
		IsActive: true,
	})
	svc := NewService(repo)

	isActive := false
	_, err := svc.Update(context.Background(), managerID, managerID, UpdateInput{IsActive: &isActive})
	if err != ErrCannotDeactivateSelf {
		t.Fatalf("expected ErrCannotDeactivateSelf, got %v", err)
	}
}

func TestUpdateCannotChangeOwnRole(t *testing.T) {
	managerID := uuid.New()
	repo := newMockRepository()
	repo.seed(&models.User{
		ID:       managerID,
		FullName: "Manager",
		Email:    "manager@test.gov.et",
		Role:     models.RoleManager,
		IsActive: true,
	})
	svc := NewService(repo)

	newRole := models.RoleRegistrar
	_, err := svc.Update(context.Background(), managerID, managerID, UpdateInput{Role: &newRole})
	if err != ErrCannotChangeOwnRole {
		t.Fatalf("expected ErrCannotChangeOwnRole, got %v", err)
	}
}

func TestDeactivateCannotDeactivateSelf(t *testing.T) {
	managerID := uuid.New()
	repo := newMockRepository()
	repo.seed(&models.User{
		ID:       managerID,
		FullName: "Manager",
		Email:    "manager@test.gov.et",
		Role:     models.RoleManager,
		IsActive: true,
	})
	svc := NewService(repo)

	err := svc.Deactivate(context.Background(), managerID, managerID)
	if err != ErrCannotDeactivateSelf {
		t.Fatalf("expected ErrCannotDeactivateSelf, got %v", err)
	}
}
