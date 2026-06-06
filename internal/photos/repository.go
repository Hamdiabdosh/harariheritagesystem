package photos

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

func (r *Repository) Create(ctx context.Context, photo *models.RecordPhoto) error {
	const query = `
		INSERT INTO record_photos (id, record_type, record_id, file_path, file_name, file_size_bytes, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query,
		photo.ID,
		photo.RecordType,
		photo.RecordID,
		photo.FilePath,
		photo.FileName,
		photo.FileSizeBytes,
		photo.UploadedBy,
	).Scan(&photo.CreatedAt)
	if err != nil {
		return fmt.Errorf("create photo record: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.RecordPhoto, error) {
	const query = `
		SELECT id, record_type, record_id, file_path, file_name, file_size_bytes, uploaded_by, created_at
		FROM record_photos
		WHERE id = $1
	`

	var photo models.RecordPhoto
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&photo.ID,
		&photo.RecordType,
		&photo.RecordID,
		&photo.FilePath,
		&photo.FileName,
		&photo.FileSizeBytes,
		&photo.UploadedBy,
		&photo.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get photo: %w", err)
	}

	return &photo, nil
}

func (r *Repository) ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordPhoto, error) {
	const query = `
		SELECT id, record_type, record_id, file_path, file_name, file_size_bytes, uploaded_by, created_at
		FROM record_photos
		WHERE record_type = $1 AND record_id = $2
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, recordType, recordID)
	if err != nil {
		return nil, fmt.Errorf("list photos: %w", err)
	}
	defer rows.Close()

	photos := make([]models.RecordPhoto, 0)
	for rows.Next() {
		var photo models.RecordPhoto
		if err := rows.Scan(
			&photo.ID,
			&photo.RecordType,
			&photo.RecordID,
			&photo.FilePath,
			&photo.FileName,
			&photo.FileSizeBytes,
			&photo.UploadedBy,
			&photo.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate photos: %w", err)
	}

	return photos, nil
}

func (r *Repository) CountByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM record_photos
		WHERE record_type = $1 AND record_id = $2
	`

	var count int
	if err := r.pool.QueryRow(ctx, query, recordType, recordID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count photos: %w", err)
	}

	return count, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM record_photos WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete photo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPhotoNotFound
	}

	return nil
}
