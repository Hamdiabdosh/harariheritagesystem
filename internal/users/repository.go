package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qirs-mezgeb/api/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context, filters ListFilters) ([]models.User, int, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	query := `
		SELECT id, full_name, email, password_hash, role, language, is_active, deleted_at, created_at, updated_at,
		       COUNT(*) OVER() AS total
		FROM users
		WHERE deleted_at IS NULL
	`
	args := []any{}
	argIdx := 1

	if filters.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, filters.Role)
		argIdx++
	}
	if filters.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filters.IsActive)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filters.Limit, (filters.Page-1)*filters.Limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	total := 0

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.Language,
			&user.IsActive,
			&user.DeletedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&total,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate users: %w", err)
	}

	return users, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, full_name, email, password_hash, role, language, is_active, deleted_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Language,
		&user.IsActive,
		&user.DeletedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT id, full_name, email, password_hash, role, language, is_active, deleted_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, strings.ToLower(email)).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Language,
		&user.IsActive,
		&user.DeletedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	const query = `
		INSERT INTO users (full_name, email, password_hash, role, language, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		user.FullName,
		strings.ToLower(user.Email),
		user.PasswordHash,
		user.Role,
		user.Language,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*models.User, error) {
	setClauses := make([]string, 0, 4)
	args := []any{id}
	argIdx := 2

	if input.FullName != nil {
		setClauses = append(setClauses, fmt.Sprintf("full_name = $%d", argIdx))
		args = append(args, *input.FullName)
		argIdx++
	}
	if input.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, *input.Role)
		argIdx++
	}
	if input.Language != nil {
		setClauses = append(setClauses, fmt.Sprintf("language = $%d", argIdx))
		args = append(args, *input.Language)
		argIdx++
	}
	if input.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *input.IsActive)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET %s, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, full_name, email, password_hash, role, language, is_active, deleted_at, created_at, updated_at
	`, strings.Join(setClauses, ", "))

	var user models.User
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Language,
		&user.IsActive,
		&user.DeletedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &user, nil
}

func (r *Repository) Deactivate(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE users
		SET is_active = FALSE, deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) UpdateLanguage(ctx context.Context, id uuid.UUID, language models.Language) (*models.User, error) {
	return r.Update(ctx, id, UpdateInput{Language: &language})
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
