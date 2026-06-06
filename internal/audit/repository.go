package audit

import (
	"context"
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

func (r *Repository) InsertTx(
	ctx context.Context,
	tx pgx.Tx,
	recordType models.RecordType,
	recordID uuid.UUID,
	changedBy uuid.UUID,
	fromStatus *string,
	toStatus string,
	note *string,
) error {
	const query = `
		INSERT INTO status_history (record_type, record_id, changed_by, from_status, to_status, note)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if _, err := tx.Exec(ctx, query, recordType, recordID, changedBy, fromStatus, toStatus, note); err != nil {
		return fmt.Errorf("insert status history: %w", err)
	}

	return nil
}

func (r *Repository) ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.StatusHistoryEntry, error) {
	const query = `
		SELECT
			h.id,
			h.record_type,
			h.record_id,
			h.changed_by,
			u.full_name,
			h.from_status,
			h.to_status,
			h.note,
			h.created_at
		FROM status_history h
		INNER JOIN users u ON u.id = h.changed_by
		WHERE h.record_type = $1 AND h.record_id = $2
		ORDER BY h.created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, recordType, recordID)
	if err != nil {
		return nil, fmt.Errorf("list status history: %w", err)
	}
	defer rows.Close()

	history := make([]models.StatusHistoryEntry, 0)
	for rows.Next() {
		var entry models.StatusHistoryEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.RecordType,
			&entry.RecordID,
			&entry.ChangedBy,
			&entry.ChangedByName,
			&entry.FromStatus,
			&entry.ToStatus,
			&entry.Note,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan status history: %w", err)
		}
		history = append(history, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate status history: %w", err)
	}

	return history, nil
}
