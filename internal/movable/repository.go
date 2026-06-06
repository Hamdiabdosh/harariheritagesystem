package movable

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qirs-mezgeb/api/internal/audit"
	"github.com/qirs-mezgeb/api/internal/models"
)

const recordColumns = `
	id, record_id, registrar_id, status,
	name_amharic, name_local, category,
	location_name, woreda, kebele, house_number, current_use, previous_id,
	owner_type, owner_name, storage_location, storage_location_other,
	made_by, period_made, age_method, acquisition_methods,
	height_cm, width_cm, length_cm, diameter_cm, thickness_cm, weight_kg,
	num_pages, num_chapters, num_illustrations,
	color_type, has_decoration, materials, material_other, description,
	notable_because, notable_other, significance,
	condition, has_threat, threat_description, maintenance_done, maintenance_by,
	maintenance_date, maintenance_count, preventive_level, accessibility, notes,
	related_docs,
	informant_name, informant_sex, informant_age, informant_occupation,
	caretaker_name, caretaker_role, registrar_date,
	approved_at, approved_by, created_at, updated_at
`

type Repository struct {
	pool  *pgxpool.Pool
	audit *audit.Repository
	now   func() time.Time
}

func NewRepository(pool *pgxpool.Pool, auditRepo *audit.Repository) *Repository {
	return &Repository{pool: pool, audit: auditRepo, now: time.Now}
}

func (r *Repository) Create(ctx context.Context, record *models.MovableRecord) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	year := r.now().Year()
	recordID, err := r.nextRecordID(ctx, tx, year)
	if err != nil {
		return err
	}
	record.RecordID = recordID
	record.Status = models.StatusDraft

	const query = `
		INSERT INTO movable_records (
			record_id, registrar_id, status,
			name_amharic, name_local, category,
			location_name, woreda, kebele, house_number, current_use, previous_id,
			owner_type, owner_name, storage_location, storage_location_other,
			made_by, period_made, age_method, acquisition_methods,
			height_cm, width_cm, length_cm, diameter_cm, thickness_cm, weight_kg,
			num_pages, num_chapters, num_illustrations,
			color_type, has_decoration, materials, material_other, description,
			notable_because, notable_other, significance,
			condition, has_threat, threat_description, maintenance_done, maintenance_by,
			maintenance_date, maintenance_count, preventive_level, accessibility, notes,
			related_docs,
			informant_name, informant_sex, informant_age, informant_occupation,
			caretaker_name, caretaker_role, registrar_date
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26,
			$27, $28, $29,
			$30, $31, $32, $33, $34,
			$35, $36, $37,
			$38, $39, $40, $41, $42,
			$43, $44, $45, $46, $47,
			$48,
			$49, $50, $51, $52,
			$53, $54, $55
		)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, recordArgs(record)...).Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert movable record: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.MovableRecord, error) {
	query := fmt.Sprintf(`SELECT %s FROM movable_records WHERE id = $1`, recordColumns)
	args := []any{id}

	if role == models.RoleRegistrar {
		query += " AND registrar_id = $2"
		args = append(args, userID)
	}

	row := r.pool.QueryRow(ctx, query, args...)
	record, err := scanRecord(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get movable record: %w", err)
	}

	return record, nil
}

func (r *Repository) List(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]models.MovableRecord, int, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	query := fmt.Sprintf(`
		SELECT %s, COUNT(*) OVER() AS total
		FROM movable_records
		WHERE 1=1
	`, recordColumns)
	args := []any{}
	argIdx := 1

	if role == models.RoleRegistrar {
		query += fmt.Sprintf(" AND registrar_id = $%d", argIdx)
		args = append(args, userID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.Woreda != "" {
		query += fmt.Sprintf(" AND woreda = $%d", argIdx)
		args = append(args, filters.Woreda)
		argIdx++
	}
	if filters.Search != "" {
		query += fmt.Sprintf(` AND (
			name_amharic ILIKE $%d OR record_id ILIKE $%d OR
			COALESCE(woreda, '') ILIKE $%d OR COALESCE(kebele, '') ILIKE $%d OR
			COALESCE(location_name, '') ILIKE $%d
		)`, argIdx, argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+filters.Search+"%")
		argIdx++
	}
	if filters.DateFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *filters.DateTo)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filters.Limit, (filters.Page-1)*filters.Limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list movable records: %w", err)
	}
	defer rows.Close()

	records := make([]models.MovableRecord, 0)
	total := 0

	for rows.Next() {
		record, count, err := scanRecordWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = count
		records = append(records, *record)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate movable records: %w", err)
	}

	return records, total, nil
}

func (r *Repository) Update(ctx context.Context, record *models.MovableRecord) error {
	const query = `
		UPDATE movable_records SET
			name_amharic = $2, name_local = $3, category = $4,
			location_name = $5, woreda = $6, kebele = $7, house_number = $8, current_use = $9, previous_id = $10,
			owner_type = $11, owner_name = $12, storage_location = $13, storage_location_other = $14,
			made_by = $15, period_made = $16, age_method = $17, acquisition_methods = $18,
			height_cm = $19, width_cm = $20, length_cm = $21, diameter_cm = $22, thickness_cm = $23, weight_kg = $24,
			num_pages = $25, num_chapters = $26, num_illustrations = $27,
			color_type = $28, has_decoration = $29, materials = $30, material_other = $31, description = $32,
			notable_because = $33, notable_other = $34, significance = $35,
			condition = $36, has_threat = $37, threat_description = $38, maintenance_done = $39, maintenance_by = $40,
			maintenance_date = $41, maintenance_count = $42, preventive_level = $43, accessibility = $44, notes = $45,
			related_docs = $46,
			informant_name = $47, informant_sex = $48, informant_age = $49, informant_occupation = $50,
			caretaker_name = $51, caretaker_role = $52, registrar_date = $53,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query, updateArgs(record)...).Scan(&record.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrRecordNotFound
	}
	if err != nil {
		return fmt.Errorf("update movable record: %w", err)
	}

	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const updateQuery = `
		UPDATE movable_records
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND status = $3
		RETURNING id
	`

	var updatedID uuid.UUID
	err = tx.QueryRow(ctx, updateQuery, id, toStatus, fromStatus).Scan(&updatedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidStatusTransition
	}
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	fromStatusStr := string(fromStatus)
	if err := r.audit.InsertTx(ctx, tx, models.RecordTypeMovable, id, changedBy, &fromStatusStr, string(toStatus), note); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) FinalApprove(ctx context.Context, id uuid.UUID, fromStatus models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const updateQuery = `
		UPDATE movable_records
		SET status = 'approved', approved_at = NOW(), approved_by = $2, updated_at = NOW()
		WHERE id = $1 AND status = $3
		RETURNING approved_at
	`

	var approvedAt time.Time
	err = tx.QueryRow(ctx, updateQuery, id, approvedBy, fromStatus).Scan(&approvedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrInvalidStatusTransition
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("final approve: %w", err)
	}

	fromStatusStr := string(fromStatus)
	if err := r.audit.InsertTx(ctx, tx, models.RecordTypeMovable, id, approvedBy, &fromStatusStr, string(models.StatusApproved), note); err != nil {
		return time.Time{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return time.Time{}, fmt.Errorf("commit transaction: %w", err)
	}

	return approvedAt, nil
}

func (r *Repository) nextRecordID(ctx context.Context, tx pgx.Tx, year int) (string, error) {
	const query = `
		SELECT COUNT(*) + 1
		FROM movable_records
		WHERE record_id LIKE 'ET-HR-AN-V-' || $1::text || '-%'
	`

	var sequence int
	if err := tx.QueryRow(ctx, query, fmt.Sprintf("%d", year)).Scan(&sequence); err != nil {
		return "", fmt.Errorf("generate record id: %w", err)
	}

	return FormatRecordID(year, sequence), nil
}

func recordArgs(record *models.MovableRecord) []any {
	return []any{
		record.RecordID, record.RegistrarID, record.Status,
		record.NameAmharic, record.NameLocal, record.Category,
		record.LocationName, record.Woreda, record.Kebele, record.HouseNumber, record.CurrentUse, record.PreviousID,
		enumPtr(record.OwnerType), record.OwnerName, enumPtr(record.StorageLocation), record.StorageLocationOther,
		record.MadeBy, record.PeriodMade, enumPtr(record.AgeMethod), record.AcquisitionMethods,
		record.HeightCM, record.WidthCM, record.LengthCM, record.DiameterCM, record.ThicknessCM, record.WeightKG,
		record.NumPages, record.NumChapters, record.NumIllustrations,
		record.ColorType, record.HasDecoration, record.Materials, record.MaterialOther, record.Description,
		record.NotableBecause, record.NotableOther, record.Significance,
		enumPtr(record.Condition), record.HasThreat, record.ThreatDescription, record.MaintenanceDone, record.MaintenanceBy,
		record.MaintenanceDate, record.MaintenanceCount, enumPtr(record.PreventiveLevel), enumPtr(record.Accessibility), record.Notes,
		record.RelatedDocs,
		record.InformantName, enumPtr(record.InformantSex), record.InformantAge, record.InformantOccupation,
		record.CaretakerName, record.CaretakerRole, record.RegistrarDate,
	}
}

func updateArgs(record *models.MovableRecord) []any {
	args := recordArgs(record)
	return append([]any{record.ID}, args[3:]...)
}

func enumPtr[T ~string](value *T) *string {
	if value == nil {
		return nil
	}
	s := string(*value)
	return &s
}

type scannable interface {
	Scan(dest ...any) error
}

func scanRecord(row scannable) (*models.MovableRecord, error) {
	var record models.MovableRecord
	var ownerType, storageLocation, ageMethod, condition sqlNullString
	var preventiveLevel, accessibility, informantSex sqlNullString

	err := row.Scan(
		&record.ID, &record.RecordID, &record.RegistrarID, &record.Status,
		&record.NameAmharic, &record.NameLocal, &record.Category,
		&record.LocationName, &record.Woreda, &record.Kebele, &record.HouseNumber, &record.CurrentUse, &record.PreviousID,
		&ownerType, &record.OwnerName, &storageLocation, &record.StorageLocationOther,
		&record.MadeBy, &record.PeriodMade, &ageMethod, &record.AcquisitionMethods,
		&record.HeightCM, &record.WidthCM, &record.LengthCM, &record.DiameterCM, &record.ThicknessCM, &record.WeightKG,
		&record.NumPages, &record.NumChapters, &record.NumIllustrations,
		&record.ColorType, &record.HasDecoration, &record.Materials, &record.MaterialOther, &record.Description,
		&record.NotableBecause, &record.NotableOther, &record.Significance,
		&condition, &record.HasThreat, &record.ThreatDescription, &record.MaintenanceDone, &record.MaintenanceBy,
		&record.MaintenanceDate, &record.MaintenanceCount, &preventiveLevel, &accessibility, &record.Notes,
		&record.RelatedDocs,
		&record.InformantName, &informantSex, &record.InformantAge, &record.InformantOccupation,
		&record.CaretakerName, &record.CaretakerRole, &record.RegistrarDate,
		&record.ApprovedAt, &record.ApprovedBy, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	assignEnum(&record.OwnerType, ownerType)
	assignEnum(&record.StorageLocation, storageLocation)
	assignEnum(&record.AgeMethod, ageMethod)
	assignEnum(&record.Condition, condition)
	assignEnum(&record.PreventiveLevel, preventiveLevel)
	assignEnum(&record.Accessibility, accessibility)
	assignEnum(&record.InformantSex, informantSex)

	return &record, nil
}

func scanRecordWithTotal(rows pgx.Rows) (*models.MovableRecord, int, error) {
	var record models.MovableRecord
	var ownerType, storageLocation, ageMethod, condition sqlNullString
	var preventiveLevel, accessibility, informantSex sqlNullString
	var total int

	err := rows.Scan(
		&record.ID, &record.RecordID, &record.RegistrarID, &record.Status,
		&record.NameAmharic, &record.NameLocal, &record.Category,
		&record.LocationName, &record.Woreda, &record.Kebele, &record.HouseNumber, &record.CurrentUse, &record.PreviousID,
		&ownerType, &record.OwnerName, &storageLocation, &record.StorageLocationOther,
		&record.MadeBy, &record.PeriodMade, &ageMethod, &record.AcquisitionMethods,
		&record.HeightCM, &record.WidthCM, &record.LengthCM, &record.DiameterCM, &record.ThicknessCM, &record.WeightKG,
		&record.NumPages, &record.NumChapters, &record.NumIllustrations,
		&record.ColorType, &record.HasDecoration, &record.Materials, &record.MaterialOther, &record.Description,
		&record.NotableBecause, &record.NotableOther, &record.Significance,
		&condition, &record.HasThreat, &record.ThreatDescription, &record.MaintenanceDone, &record.MaintenanceBy,
		&record.MaintenanceDate, &record.MaintenanceCount, &preventiveLevel, &accessibility, &record.Notes,
		&record.RelatedDocs,
		&record.InformantName, &informantSex, &record.InformantAge, &record.InformantOccupation,
		&record.CaretakerName, &record.CaretakerRole, &record.RegistrarDate,
		&record.ApprovedAt, &record.ApprovedBy, &record.CreatedAt, &record.UpdatedAt,
		&total,
	)
	if err != nil {
		return nil, 0, err
	}

	assignEnum(&record.OwnerType, ownerType)
	assignEnum(&record.StorageLocation, storageLocation)
	assignEnum(&record.AgeMethod, ageMethod)
	assignEnum(&record.Condition, condition)
	assignEnum(&record.PreventiveLevel, preventiveLevel)
	assignEnum(&record.Accessibility, accessibility)
	assignEnum(&record.InformantSex, informantSex)

	return &record, total, nil
}

type sqlNullString struct {
	value *string
}

func (n *sqlNullString) Scan(src any) error {
	if src == nil {
		n.value = nil
		return nil
	}
	switch v := src.(type) {
	case string:
		s := v
		n.value = &s
	case []byte:
		s := string(v)
		n.value = &s
	default:
		return fmt.Errorf("unsupported scan type %T", src)
	}
	return nil
}

func assignEnum[T ~string](target **T, source sqlNullString) {
	if source.value == nil {
		*target = nil
		return
	}
	v := T(*source.value)
	*target = &v
}
