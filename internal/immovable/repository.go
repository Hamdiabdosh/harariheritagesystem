package immovable

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
	name_amharic, name_local, category, current_use, current_use_other, previous_id,
	woreda, kebele, house_number, street_number, gate,
	owner_type, owner_name, map_reference, gps_east, gps_north, elevation_m,
	built_by, construction_period, age_method, height_m, length_m, width_m,
	num_doors, num_windows, num_rooms, material, description, harari_house_grades, neighborhood_type,
	overall_condition, damage_roof, damage_cornice, damage_wall, damage_floor, damage_door,
	damage_cupboard, damage_upper_floor, damage_dera, damage_pillar,
	value_historical, value_craftsmanship, value_artistic, value_scientific, value_cultural,
	has_threat, maintenance_done, maintenance_reason, maintenance_by, maintenance_date, maintenance_count,
	preventive_level, accessibility, notes,
	related_docs, has_oral_history,
	caretaker_name, caretaker_role, informant_name, informant_sex, informant_age, registrar_date,
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

func (r *Repository) Create(ctx context.Context, record *models.ImmovableRecord) error {
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
		INSERT INTO immovable_records (
			record_id, registrar_id, status,
			name_amharic, name_local, category, current_use, current_use_other, previous_id,
			woreda, kebele, house_number, street_number, gate,
			owner_type, owner_name, map_reference, gps_east, gps_north, elevation_m,
			built_by, construction_period, age_method, height_m, length_m, width_m,
			num_doors, num_windows, num_rooms, material, description, harari_house_grades, neighborhood_type,
			overall_condition, damage_roof, damage_cornice, damage_wall, damage_floor, damage_door,
			damage_cupboard, damage_upper_floor, damage_dera, damage_pillar,
			value_historical, value_craftsmanship, value_artistic, value_scientific, value_cultural,
			has_threat, maintenance_done, maintenance_reason, maintenance_by, maintenance_date, maintenance_count,
			preventive_level, accessibility, notes,
			related_docs, has_oral_history,
			caretaker_name, caretaker_role, informant_name, informant_sex, informant_age, registrar_date
		) VALUES (
			$1, $2, $3,
			$4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26,
			$27, $28, $29, $30, $31, $32, $33,
			$34, $35, $36, $37, $38, $39,
			$40, $41, $42, $43,
			$44, $45, $46, $47, $48,
			$49, $50, $51, $52, $53, $54,
			$55, $56, $57,
			$58, $59,
			$60, $61, $62, $63, $64, $65
		)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, recordArgs(record)...).Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert immovable record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.ImmovableRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM immovable_records WHERE id = $1
	`, recordColumns)
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
		return nil, fmt.Errorf("get immovable record: %w", err)
	}

	return record, nil
}

func (r *Repository) List(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]models.ImmovableRecord, int, error) {
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
		FROM immovable_records
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
		query += fmt.Sprintf(" AND (name_amharic ILIKE $%d OR record_id ILIKE $%d OR woreda ILIKE $%d OR kebele ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
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
		return nil, 0, fmt.Errorf("list immovable records: %w", err)
	}
	defer rows.Close()

	records := make([]models.ImmovableRecord, 0)
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
		return nil, 0, fmt.Errorf("iterate immovable records: %w", err)
	}

	return records, total, nil
}

func (r *Repository) Update(ctx context.Context, record *models.ImmovableRecord) error {
	const query = `
		UPDATE immovable_records SET
			name_amharic = $2, name_local = $3, category = $4, current_use = $5, current_use_other = $6, previous_id = $7,
			woreda = $8, kebele = $9, house_number = $10, street_number = $11, gate = $12,
			owner_type = $13, owner_name = $14, map_reference = $15, gps_east = $16, gps_north = $17, elevation_m = $18,
			built_by = $19, construction_period = $20, age_method = $21, height_m = $22, length_m = $23, width_m = $24,
			num_doors = $25, num_windows = $26, num_rooms = $27, material = $28, description = $29, harari_house_grades = $30, neighborhood_type = $31,
			overall_condition = $32, damage_roof = $33, damage_cornice = $34, damage_wall = $35, damage_floor = $36, damage_door = $37,
			damage_cupboard = $38, damage_upper_floor = $39, damage_dera = $40, damage_pillar = $41,
			value_historical = $42, value_craftsmanship = $43, value_artistic = $44, value_scientific = $45, value_cultural = $46,
			has_threat = $47, maintenance_done = $48, maintenance_reason = $49, maintenance_by = $50, maintenance_date = $51, maintenance_count = $52,
			preventive_level = $53, accessibility = $54, notes = $55,
			related_docs = $56, has_oral_history = $57,
			caretaker_name = $58, caretaker_role = $59, informant_name = $60, informant_sex = $61, informant_age = $62, registrar_date = $63,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	args := updateArgs(record)
	err := r.pool.QueryRow(ctx, query, args...).Scan(&record.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("update immovable record: %w", err)
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
		UPDATE immovable_records
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
	if err := r.audit.InsertTx(ctx, tx, models.RecordTypeImmovable, id, changedBy, &fromStatusStr, string(toStatus), note); err != nil {
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
		UPDATE immovable_records
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
	if err := r.audit.InsertTx(ctx, tx, models.RecordTypeImmovable, id, approvedBy, &fromStatusStr, string(models.StatusApproved), note); err != nil {
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
		FROM immovable_records
		WHERE record_id LIKE 'ET-HR-AN-I-' || $1::text || '-%'
	`

	var sequence int
	prefix := fmt.Sprintf("%d", year)
	if err := tx.QueryRow(ctx, query, prefix).Scan(&sequence); err != nil {
		return "", fmt.Errorf("generate record id: %w", err)
	}

	return FormatRecordID(year, sequence), nil
}

func recordArgs(record *models.ImmovableRecord) []any {
	return []any{
		record.RecordID, record.RegistrarID, record.Status,
		record.NameAmharic, record.NameLocal, record.Category, record.CurrentUse, record.CurrentUseOther, record.PreviousID,
		record.Woreda, record.Kebele, record.HouseNumber, record.StreetNumber, record.Gate,
		enumPtr(record.OwnerType), record.OwnerName, record.MapReference, record.GPSEast, record.GPSNorth, record.ElevationM,
		record.BuiltBy, record.ConstructionPeriod, enumPtr(record.AgeMethod), record.HeightM, record.LengthM, record.WidthM,
		record.NumDoors, record.NumWindows, record.NumRooms, record.Material, record.Description, record.HarariHouseGrades, record.NeighborhoodType,
		enumPtr(record.OverallCondition), enumPtr(record.DamageRoof), enumPtr(record.DamageCornice), enumPtr(record.DamageWall), enumPtr(record.DamageFloor), enumPtr(record.DamageDoor),
		enumPtr(record.DamageCupboard), enumPtr(record.DamageUpperFloor), enumPtr(record.DamageDera), enumPtr(record.DamagePillar),
		record.ValueHistorical, record.ValueCraftsmanship, record.ValueArtistic, record.ValueScientific, record.ValueCultural,
		record.HasThreat, record.MaintenanceDone, record.MaintenanceReason, record.MaintenanceBy, record.MaintenanceDate, record.MaintenanceCount,
		enumPtr(record.PreventiveLevel), enumPtr(record.Accessibility), record.Notes,
		record.RelatedDocs, record.HasOralHistory,
		record.CaretakerName, record.CaretakerRole, record.InformantName, enumPtr(record.InformantSex), record.InformantAge, record.RegistrarDate,
	}
}

func updateArgs(record *models.ImmovableRecord) []any {
	return []any{
		record.ID,
		record.NameAmharic, record.NameLocal, record.Category, record.CurrentUse, record.CurrentUseOther, record.PreviousID,
		record.Woreda, record.Kebele, record.HouseNumber, record.StreetNumber, record.Gate,
		enumPtr(record.OwnerType), record.OwnerName, record.MapReference, record.GPSEast, record.GPSNorth, record.ElevationM,
		record.BuiltBy, record.ConstructionPeriod, enumPtr(record.AgeMethod), record.HeightM, record.LengthM, record.WidthM,
		record.NumDoors, record.NumWindows, record.NumRooms, record.Material, record.Description, record.HarariHouseGrades, record.NeighborhoodType,
		enumPtr(record.OverallCondition), enumPtr(record.DamageRoof), enumPtr(record.DamageCornice), enumPtr(record.DamageWall), enumPtr(record.DamageFloor), enumPtr(record.DamageDoor),
		enumPtr(record.DamageCupboard), enumPtr(record.DamageUpperFloor), enumPtr(record.DamageDera), enumPtr(record.DamagePillar),
		record.ValueHistorical, record.ValueCraftsmanship, record.ValueArtistic, record.ValueScientific, record.ValueCultural,
		record.HasThreat, record.MaintenanceDone, record.MaintenanceReason, record.MaintenanceBy, record.MaintenanceDate, record.MaintenanceCount,
		enumPtr(record.PreventiveLevel), enumPtr(record.Accessibility), record.Notes,
		record.RelatedDocs, record.HasOralHistory,
		record.CaretakerName, record.CaretakerRole, record.InformantName, enumPtr(record.InformantSex), record.InformantAge, record.RegistrarDate,
	}
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

func scanRecord(row scannable) (*models.ImmovableRecord, error) {
	var record models.ImmovableRecord
	var ownerType, ageMethod, overallCondition sqlNullString
	var damageRoof, damageCornice, damageWall, damageFloor, damageDoor sqlNullString
	var damageCupboard, damageUpperFloor, damageDera, damagePillar sqlNullString
	var preventiveLevel, accessibility, informantSex sqlNullString

	err := row.Scan(
		&record.ID, &record.RecordID, &record.RegistrarID, &record.Status,
		&record.NameAmharic, &record.NameLocal, &record.Category, &record.CurrentUse, &record.CurrentUseOther, &record.PreviousID,
		&record.Woreda, &record.Kebele, &record.HouseNumber, &record.StreetNumber, &record.Gate,
		&ownerType, &record.OwnerName, &record.MapReference, &record.GPSEast, &record.GPSNorth, &record.ElevationM,
		&record.BuiltBy, &record.ConstructionPeriod, &ageMethod, &record.HeightM, &record.LengthM, &record.WidthM,
		&record.NumDoors, &record.NumWindows, &record.NumRooms, &record.Material, &record.Description, &record.HarariHouseGrades, &record.NeighborhoodType,
		&overallCondition, &damageRoof, &damageCornice, &damageWall, &damageFloor, &damageDoor,
		&damageCupboard, &damageUpperFloor, &damageDera, &damagePillar,
		&record.ValueHistorical, &record.ValueCraftsmanship, &record.ValueArtistic, &record.ValueScientific, &record.ValueCultural,
		&record.HasThreat, &record.MaintenanceDone, &record.MaintenanceReason, &record.MaintenanceBy, &record.MaintenanceDate, &record.MaintenanceCount,
		&preventiveLevel, &accessibility, &record.Notes,
		&record.RelatedDocs, &record.HasOralHistory,
		&record.CaretakerName, &record.CaretakerRole, &record.InformantName, &informantSex, &record.InformantAge, &record.RegistrarDate,
		&record.ApprovedAt, &record.ApprovedBy, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	assignEnum(&record.OwnerType, ownerType)
	assignEnum(&record.AgeMethod, ageMethod)
	assignEnum(&record.OverallCondition, overallCondition)
	assignEnum(&record.DamageRoof, damageRoof)
	assignEnum(&record.DamageCornice, damageCornice)
	assignEnum(&record.DamageWall, damageWall)
	assignEnum(&record.DamageFloor, damageFloor)
	assignEnum(&record.DamageDoor, damageDoor)
	assignEnum(&record.DamageCupboard, damageCupboard)
	assignEnum(&record.DamageUpperFloor, damageUpperFloor)
	assignEnum(&record.DamageDera, damageDera)
	assignEnum(&record.DamagePillar, damagePillar)
	assignEnum(&record.PreventiveLevel, preventiveLevel)
	assignEnum(&record.Accessibility, accessibility)
	assignEnum(&record.InformantSex, informantSex)

	return &record, nil
}

func scanRecordWithTotal(rows pgx.Rows) (*models.ImmovableRecord, int, error) {
	var record models.ImmovableRecord
	var ownerType, ageMethod, overallCondition sqlNullString
	var damageRoof, damageCornice, damageWall, damageFloor, damageDoor sqlNullString
	var damageCupboard, damageUpperFloor, damageDera, damagePillar sqlNullString
	var preventiveLevel, accessibility, informantSex sqlNullString
	var total int

	err := rows.Scan(
		&record.ID, &record.RecordID, &record.RegistrarID, &record.Status,
		&record.NameAmharic, &record.NameLocal, &record.Category, &record.CurrentUse, &record.CurrentUseOther, &record.PreviousID,
		&record.Woreda, &record.Kebele, &record.HouseNumber, &record.StreetNumber, &record.Gate,
		&ownerType, &record.OwnerName, &record.MapReference, &record.GPSEast, &record.GPSNorth, &record.ElevationM,
		&record.BuiltBy, &record.ConstructionPeriod, &ageMethod, &record.HeightM, &record.LengthM, &record.WidthM,
		&record.NumDoors, &record.NumWindows, &record.NumRooms, &record.Material, &record.Description, &record.HarariHouseGrades, &record.NeighborhoodType,
		&overallCondition, &damageRoof, &damageCornice, &damageWall, &damageFloor, &damageDoor,
		&damageCupboard, &damageUpperFloor, &damageDera, &damagePillar,
		&record.ValueHistorical, &record.ValueCraftsmanship, &record.ValueArtistic, &record.ValueScientific, &record.ValueCultural,
		&record.HasThreat, &record.MaintenanceDone, &record.MaintenanceReason, &record.MaintenanceBy, &record.MaintenanceDate, &record.MaintenanceCount,
		&preventiveLevel, &accessibility, &record.Notes,
		&record.RelatedDocs, &record.HasOralHistory,
		&record.CaretakerName, &record.CaretakerRole, &record.InformantName, &informantSex, &record.InformantAge, &record.RegistrarDate,
		&record.ApprovedAt, &record.ApprovedBy, &record.CreatedAt, &record.UpdatedAt,
		&total,
	)
	if err != nil {
		return nil, 0, err
	}

	assignEnum(&record.OwnerType, ownerType)
	assignEnum(&record.AgeMethod, ageMethod)
	assignEnum(&record.OverallCondition, overallCondition)
	assignEnum(&record.DamageRoof, damageRoof)
	assignEnum(&record.DamageCornice, damageCornice)
	assignEnum(&record.DamageWall, damageWall)
	assignEnum(&record.DamageFloor, damageFloor)
	assignEnum(&record.DamageDoor, damageDoor)
	assignEnum(&record.DamageCupboard, damageCupboard)
	assignEnum(&record.DamageUpperFloor, damageUpperFloor)
	assignEnum(&record.DamageDera, damageDera)
	assignEnum(&record.DamagePillar, damagePillar)
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
