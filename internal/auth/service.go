package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/qirs-mezgeb/api/internal/models"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type RepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, id uuid.UUID) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
}

type AccessClaims struct {
	UserID uuid.UUID   `json:"user_id"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	User         models.UserPublic `json:"user"`
}

type Service struct {
	repo               RepositoryInterface
	jwtSecret          []byte
	jwtRefreshSecret   []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	now                func() time.Time
}

func NewService(repo RepositoryInterface, jwtSecret, jwtRefreshSecret string) *Service {
	return &Service{
		repo:               repo,
		jwtSecret:          []byte(jwtSecret),
		jwtRefreshSecret:   []byte(jwtRefreshSecret),
		accessTokenExpiry:  accessTokenDuration,
		refreshTokenExpiry: refreshTokenDuration,
		now:                time.Now,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if !user.IsUsable() {
		return nil, ErrAccountDeactivated
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokenPair(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	tokenID, err := parseRefreshTokenID(refreshToken)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	stored, err := s.repo.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return "", err
	}
	if stored == nil || s.now().After(stored.ExpiresAt) {
		return "", ErrInvalidRefreshToken
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored.TokenHash), hashRefreshToken(refreshToken)); err != nil {
		return "", ErrInvalidRefreshToken
	}

	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil {
		return "", err
	}
	if user == nil || !user.IsUsable() {
		return "", ErrAccountDeactivated
	}

	return s.createAccessToken(user)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	tokenID, err := parseRefreshTokenID(refreshToken)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	stored, err := s.repo.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return err
	}
	if stored == nil {
		return ErrInvalidRefreshToken
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored.TokenHash), hashRefreshToken(refreshToken)); err != nil {
		return ErrInvalidRefreshToken
	}

	return s.repo.DeleteRefreshToken(ctx, tokenID)
}

func (s *Service) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, ErrInvalidAccessToken
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidAccessToken
	}

	return claims, nil
}

func (s *Service) issueTokenPair(ctx context.Context, user *models.User) (*TokenPair, error) {
	accessToken, err := s.createAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToPublic(),
	}, nil
}

func (s *Service) createAccessToken(user *models.User) (string, error) {
	now := s.now()
	claims := AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

func (s *Service) createRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	tokenID := uuid.New()
	secret, err := randomTokenSecret()
	if err != nil {
		return "", err
	}

	plainToken := tokenID.String() + "." + secret
	hash, err := bcrypt.GenerateFromPassword(hashRefreshToken(plainToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash refresh token: %w", err)
	}

	record := &models.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		TokenHash: string(hash),
		ExpiresAt: s.now().Add(s.refreshTokenExpiry),
	}

	if err := s.repo.CreateRefreshToken(ctx, record); err != nil {
		return "", err
	}

	return plainToken, nil
}

func parseRefreshTokenID(refreshToken string) (uuid.UUID, error) {
	parts := strings.SplitN(refreshToken, ".", 2)
	if len(parts) != 2 {
		return uuid.Nil, fmt.Errorf("invalid refresh token format")
	}

	return uuid.Parse(parts[0])
}

func randomTokenSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate token secret: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashRefreshToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
