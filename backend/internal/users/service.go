package users

import (
	"context"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/qirs-mezgeb/api/internal/models"
)

type RepositoryInterface interface {
	List(ctx context.Context, filters ListFilters) ([]models.User, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*models.User, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
	UpdateLanguage(ctx context.Context, id uuid.UUID, language models.Language) (*models.User, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filters ListFilters) (*PaginatedUsers, error) {
	users, total, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	items := make([]UserItem, len(users))
	for i := range users {
		items[i] = toUserItem(&users[i])
	}

	limit := filters.Limit
	if limit < 1 {
		limit = 20
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &PaginatedUsers{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*UserItem, error) {
	if err := validatePassword(input.Password); err != nil {
		return nil, err
	}
	if !input.Role.IsValid() {
		return nil, ErrInvalidRole
	}
	if !isValidLanguage(input.Language) {
		return nil, ErrInvalidLanguage
	}

	fullName := strings.TrimSpace(input.FullName)
	if fullName == "" || len(fullName) > 100 {
		return nil, ErrInvalidFullName
	}

	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	language := input.Language
	if language == "" {
		language = models.LanguageAm
	}

	user := &models.User{
		FullName:     fullName,
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: string(hash),
		Role:         input.Role,
		Language:     language,
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	item := toUserItem(user)
	return &item, nil
}

func (s *Service) Update(ctx context.Context, actorID, targetID uuid.UUID, input UpdateInput) (*UserItem, error) {
	target, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, ErrUserNotFound
	}

	if actorID == targetID {
		if input.Role != nil && *input.Role != target.Role {
			return nil, ErrCannotChangeOwnRole
		}
		if input.IsActive != nil && !*input.IsActive {
			return nil, ErrCannotDeactivateSelf
		}
	}

	if input.Role != nil && !input.Role.IsValid() {
		return nil, ErrInvalidRole
	}
	if input.Language != nil && !isValidLanguage(*input.Language) {
		return nil, ErrInvalidLanguage
	}
	if input.FullName != nil {
		trimmed := strings.TrimSpace(*input.FullName)
		if trimmed == "" || len(trimmed) > 100 {
			return nil, ErrInvalidFullName
		}
		input.FullName = &trimmed
	}

	updated, err := s.repo.Update(ctx, targetID, input)
	if err != nil {
		return nil, err
	}

	item := toUserItem(updated)
	return &item, nil
}

func (s *Service) Deactivate(ctx context.Context, actorID, targetID uuid.UUID) error {
	if actorID == targetID {
		return ErrCannotDeactivateSelf
	}

	target, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target == nil {
		return ErrUserNotFound
	}

	return s.repo.Deactivate(ctx, targetID)
}

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (*UserItem, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	item := toUserItem(user)
	return &item, nil
}

func (s *Service) UpdateMyLanguage(ctx context.Context, userID uuid.UUID, language models.Language) (*models.Language, error) {
	if !isValidLanguage(language) {
		return nil, ErrInvalidLanguage
	}

	user, err := s.repo.UpdateLanguage(ctx, userID, language)
	if err != nil {
		return nil, err
	}

	return &user.Language, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}

	for _, r := range password {
		if unicode.IsDigit(r) {
			return nil
		}
	}

	return ErrInvalidPassword
}

func isValidLanguage(language models.Language) bool {
	return language == models.LanguageAm || language == models.LanguageEn || language == ""
}
