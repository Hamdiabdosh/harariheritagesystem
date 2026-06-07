package movable

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

func TestFormatRecordID(t *testing.T) {
	got := FormatRecordID(2026, 42)
	want := "ET-HR-AN-V-2026-0042"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestValidateForSubmit(t *testing.T) {
	ownerType := models.MovableOwnerType("private")
	storage := models.StorageLocation("museum")
	category := "III"
	woreda := "Amir-Nur"
	kebele := "01"

	record := &models.MovableRecord{
		NameAmharic:       "የሀረሪ ስሜት",
		Category:          &category,
		Woreda:            &woreda,
		Kebele:            &kebele,
		OwnerType:         &ownerType,
		StorageLocation:   &storage,
		Materials:         []string{"gold"},
	}

	if err := validateForSubmit(record); err != nil {
		t.Fatalf("expected valid record, got %v", err)
	}

	record.Woreda = nil
	if err := validateForSubmit(record); err == nil {
		t.Fatal("expected validation error for missing woreda")
	}
	record.Woreda = &woreda

	record.Materials = nil
	if err := validateForSubmit(record); err == nil {
		t.Fatal("expected validation error for missing materials")
	}
}

func TestUpdateForbiddenForSupervisor(t *testing.T) {
	repo := &mockRepo{
		record: &models.MovableRecord{
			ID:          uuid.New(),
			RegistrarID: uuid.New(),
			Status:      models.StatusDraft,
		},
	}
	svc := NewService(repo, nil, nil, nil)

	_, err := svc.Update(t.Context(), repo.record.ID, uuid.New(), models.RoleSupervisor, models.MovableRecordInput{})
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

type mockRepo struct {
	record *models.MovableRecord
}

func (m *mockRepo) Create(ctx context.Context, record *models.MovableRecord) error {
	record.ID = uuid.New()
	record.RecordID = FormatRecordID(2026, 1)
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.MovableRecord, error) {
	if m.record == nil || m.record.ID != id {
		return nil, nil
	}
	copy := *m.record
	return &copy, nil
}

func (m *mockRepo) List(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]models.MovableRecord, int, error) {
	return nil, 0, nil
}

func (m *mockRepo) Update(ctx context.Context, record *models.MovableRecord) error {
	return nil
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error {
	return nil
}

func (m *mockRepo) FinalApprove(ctx context.Context, id uuid.UUID, fromStatus models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error) {
	return time.Now(), nil
}
