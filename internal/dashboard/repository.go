package dashboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qirs-mezgeb/api/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetStats(ctx context.Context, userID uuid.UUID, role models.Role) (*Stats, error) {
	immovableWhere, immovableArgs := buildTableFilters("immovable_records", userID, role, ListFilters{})
	movableWhere, movableArgs := buildTableFilters("movable_records", userID, role, ListFilters{})

	var totalImmovable int
	immovableCountQuery := fmt.Sprintf("SELECT COUNT(*) FROM immovable_records WHERE %s", immovableWhere)
	if err := r.pool.QueryRow(ctx, immovableCountQuery, immovableArgs...).Scan(&totalImmovable); err != nil {
		return nil, fmt.Errorf("count immovable records: %w", err)
	}

	var totalMovable int
	movableCountQuery := fmt.Sprintf("SELECT COUNT(*) FROM movable_records WHERE %s", movableWhere)
	if err := r.pool.QueryRow(ctx, movableCountQuery, movableArgs...).Scan(&totalMovable); err != nil {
		return nil, fmt.Errorf("count movable records: %w", err)
	}

	byStatus, err := r.countByStatus(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		TotalImmovable: totalImmovable,
		TotalMovable:   totalMovable,
		ByStatus:       byStatus,
	}

	if pending, ok, err := r.pendingMyAction(ctx, userID, role); err != nil {
		return nil, err
	} else if ok {
		stats.PendingMyAction = &pending
	}

	return stats, nil
}

func (r *Repository) countByStatus(ctx context.Context, userID uuid.UUID, role models.Role) (StatusCounts, error) {
	immovableWhere, immovableArgs, nextIdx := buildTableFiltersIndexed("immovable_records", userID, role, ListFilters{}, 1)
	movableWhere, movableArgs, _ := buildTableFiltersIndexed("movable_records", userID, role, ListFilters{}, nextIdx)

	query := fmt.Sprintf(`
		SELECT status, COUNT(*)::int
		FROM (
			SELECT status FROM immovable_records WHERE %s
			UNION ALL
			SELECT status FROM movable_records WHERE %s
		) combined
		GROUP BY status
	`, immovableWhere, movableWhere)

	args := append(immovableArgs, movableArgs...)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return StatusCounts{}, fmt.Errorf("count by status: %w", err)
	}
	defer rows.Close()

	counts := StatusCounts{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return StatusCounts{}, fmt.Errorf("scan status count: %w", err)
		}
		switch models.RecordStatus(status) {
		case models.StatusDraft:
			counts.Draft = count
		case models.StatusPendingReview:
			counts.PendingReview = count
		case models.StatusUnderReview:
			counts.UnderReview = count
		case models.StatusReturned:
			counts.Returned = count
		case models.StatusApproved:
			counts.Approved = count
		}
	}

	if err := rows.Err(); err != nil {
		return StatusCounts{}, fmt.Errorf("iterate status counts: %w", err)
	}

	return counts, nil
}

func (r *Repository) pendingMyAction(ctx context.Context, userID uuid.UUID, role models.Role) (int, bool, error) {
	var status models.RecordStatus
	switch role {
	case models.RoleSupervisor:
		status = models.StatusPendingReview
	case models.RoleManager:
		status = models.StatusUnderReview
	default:
		return 0, false, nil
	}

	filters := ListFilters{Status: status}
	immovableWhere, immovableArgs, nextIdx := buildTableFiltersIndexed("immovable_records", userID, role, filters, 1)
	movableWhere, movableArgs, _ := buildTableFiltersIndexed("movable_records", userID, role, filters, nextIdx)

	query := fmt.Sprintf(`
		SELECT COUNT(*)::int
		FROM (
			SELECT id FROM immovable_records WHERE %s
			UNION ALL
			SELECT id FROM movable_records WHERE %s
		) combined
	`, immovableWhere, movableWhere)

	args := append(immovableArgs, movableArgs...)
	var count int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, false, fmt.Errorf("count pending action: %w", err)
	}

	return count, true, nil
}

func (r *Repository) ListRecords(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]RecordSummary, int, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	includeImmovable := filters.Type == "" || filters.Type == models.RecordTypeImmovable
	includeMovable := filters.Type == "" || filters.Type == models.RecordTypeMovable

	parts := make([]string, 0, 2)
	args := make([]any, 0)
	argIdx := 1

	if includeImmovable {
		where, partArgs, nextIdx := buildTableFiltersIndexed("immovable_records", userID, role, filters, argIdx)
		parts = append(parts, fmt.Sprintf(`
			SELECT
				id,
				'immovable'::record_type AS record_type,
				record_id,
				name_amharic,
				status,
				woreda,
				kebele,
				registrar_id,
				created_at,
				updated_at
			FROM immovable_records
			WHERE %s
		`, where))
		args = append(args, partArgs...)
		argIdx = nextIdx
	}

	if includeMovable {
		where, partArgs, nextIdx := buildTableFiltersIndexed("movable_records", userID, role, filters, argIdx)
		parts = append(parts, fmt.Sprintf(`
			SELECT
				id,
				'movable'::record_type AS record_type,
				record_id,
				name_amharic,
				status,
				woreda,
				kebele,
				registrar_id,
				created_at,
				updated_at
			FROM movable_records
			WHERE %s
		`, where))
		args = append(args, partArgs...)
		argIdx = nextIdx
	}

	if len(parts) == 0 {
		return []RecordSummary{}, 0, nil
	}

	combined := strings.Join(parts, " UNION ALL ")
	query := fmt.Sprintf(`
		SELECT id, record_type, record_id, name_amharic, status, woreda, kebele, registrar_id, created_at, updated_at,
			COUNT(*) OVER() AS total
		FROM (%s) records
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, combined, argIdx, argIdx+1)

	args = append(args, filters.Limit, (filters.Page-1)*filters.Limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list records: %w", err)
	}
	defer rows.Close()

	items := make([]RecordSummary, 0)
	total := 0

	for rows.Next() {
		var item RecordSummary
		if err := rows.Scan(
			&item.ID,
			&item.RecordType,
			&item.RecordID,
			&item.NameAmharic,
			&item.Status,
			&item.Woreda,
			&item.Kebele,
			&item.RegistrarID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&total,
		); err != nil {
			return nil, 0, fmt.Errorf("scan record summary: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate record summaries: %w", err)
	}

	return items, total, nil
}

func buildTableFilters(table string, userID uuid.UUID, role models.Role, filters ListFilters) (string, []any) {
	where, args, _ := buildTableFiltersIndexed(table, userID, role, filters, 1)
	return where, args
}

func buildTableFiltersIndexed(table string, userID uuid.UUID, role models.Role, filters ListFilters, startIdx int) (string, []any, int) {
	clauses := []string{"1=1"}
	args := make([]any, 0)
	argIdx := startIdx

	if role == models.RoleRegistrar {
		clauses = append(clauses, fmt.Sprintf("registrar_id = $%d", argIdx))
		args = append(args, userID)
		argIdx++
	}

	if filters.Status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filters.Status)
		argIdx++
	}

	if filters.Woreda != "" {
		if table == "movable_records" {
			clauses = append(clauses, fmt.Sprintf("COALESCE(woreda, '') = $%d", argIdx))
		} else {
			clauses = append(clauses, fmt.Sprintf("woreda = $%d", argIdx))
		}
		args = append(args, filters.Woreda)
		argIdx++
	}

	if filters.Kebele != "" {
		if table == "movable_records" {
			clauses = append(clauses, fmt.Sprintf("COALESCE(kebele, '') = $%d", argIdx))
		} else {
			clauses = append(clauses, fmt.Sprintf("kebele = $%d", argIdx))
		}
		args = append(args, filters.Kebele)
		argIdx++
	}

	if filters.Search != "" {
		pattern := "%" + filters.Search + "%"
		if table == "movable_records" {
			clauses = append(clauses, fmt.Sprintf(`(
				name_amharic ILIKE $%d OR record_id ILIKE $%d OR
				COALESCE(woreda, '') ILIKE $%d OR COALESCE(kebele, '') ILIKE $%d OR
				COALESCE(location_name, '') ILIKE $%d
			)`, argIdx, argIdx, argIdx, argIdx, argIdx))
		} else {
			clauses = append(clauses, fmt.Sprintf(`(
				name_amharic ILIKE $%d OR record_id ILIKE $%d OR
				woreda ILIKE $%d OR kebele ILIKE $%d
			)`, argIdx, argIdx, argIdx, argIdx))
		}
		args = append(args, pattern)
		argIdx++
	}

	if filters.DateFrom != nil {
		clauses = append(clauses, fmt.Sprintf("created_at >= $%d", argIdx))
		args = append(args, *filters.DateFrom)
		argIdx++
	}

	if filters.DateTo != nil {
		clauses = append(clauses, fmt.Sprintf("created_at <= $%d", argIdx))
		args = append(args, *filters.DateTo)
		argIdx++
	}

	return strings.Join(clauses, " AND "), args, argIdx
}
