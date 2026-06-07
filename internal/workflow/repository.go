package workflow

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qirs-mezgeb/api/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateComment(ctx context.Context, comment *models.RecordComment) error {
	const query = `
		INSERT INTO record_comments (record_type, record_id, author_id, comment_text)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		comment.RecordType,
		comment.RecordID,
		comment.AuthorID,
		comment.CommentText,
	).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}

	return nil
}

func (r *Repository) ListComments(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordComment, error) {
	const query = `
		SELECT
			c.id,
			c.record_type,
			c.record_id,
			c.author_id,
			u.full_name,
			c.comment_text,
			c.created_at
		FROM record_comments c
		INNER JOIN users u ON u.id = c.author_id
		WHERE c.record_type = $1 AND c.record_id = $2
		ORDER BY c.created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, recordType, recordID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	comments := make([]models.RecordComment, 0)
	for rows.Next() {
		var comment models.RecordComment
		if err := rows.Scan(
			&comment.ID,
			&comment.RecordType,
			&comment.RecordID,
			&comment.AuthorID,
			&comment.AuthorName,
			&comment.CommentText,
			&comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate comments: %w", err)
	}

	return comments, nil
}

func (r *Repository) GetComment(ctx context.Context, id uuid.UUID) (*models.RecordComment, error) {
	const query = `
		SELECT id, record_type, record_id, author_id, comment_text, created_at
		FROM record_comments
		WHERE id = $1
	`

	var comment models.RecordComment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.RecordType,
		&comment.RecordID,
		&comment.AuthorID,
		&comment.CommentText,
		&comment.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get comment: %w", err)
	}

	return &comment, nil
}
