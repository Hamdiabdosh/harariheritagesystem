package immovable

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

func TestFormatRecordID(t *testing.T) {
	got := FormatRecordID(2026, 1)
	want := "ET-HR-AN-I-2026-0001"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestCanTransitionSubmit(t *testing.T) {
	if !canTransition(models.StatusDraft, models.StatusPendingReview) {
		t.Fatal("draft should transition to pending_review")
	}
	if !canTransition(models.StatusReturned, models.StatusPendingReview) {
		t.Fatal("returned should transition to pending_review")
	}
	if canTransition(models.StatusApproved, models.StatusPendingReview) {
		t.Fatal("approved should not transition")
	}
}

func TestValidateForSubmit(t *testing.T) {
	ownerType := models.ImmovableOwnerType("private")
	period := "19th century"
	ageMethod := models.AgeMethod("estimated")

	record := &models.ImmovableRecord{
		NameAmharic:        "ሀረሪ ቤት",
		Category:           []string{"I"},
		Woreda:             "ሀረሪ",
		Kebele:             "01",
		CurrentUse:         []string{"residence"},
		OwnerType:          &ownerType,
		ConstructionPeriod: &period,
		AgeMethod:          &ageMethod,
	}

	if err := validateForSubmit(record); err != nil {
		t.Fatalf("expected valid record, got %v", err)
	}

	record.NameAmharic = ""
	err := validateForSubmit(record)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUpdateForbiddenForSupervisor(t *testing.T) {
	repo := &mockRepo{
		record: &models.ImmovableRecord{
			ID:          uuid.New(),
			RegistrarID: uuid.New(),
			Status:      models.StatusDraft,
		},
	}
	svc := NewService(repo, nil, nil)

	_, err := svc.Update(t.Context(), repo.record.ID, uuid.New(), models.RoleSupervisor, models.ImmovableRecordInput{})
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

type mockRepo struct {
	record *models.ImmovableRecord
}

func (m *mockRepo) Create(ctx context.Context, record *models.ImmovableRecord) error {
	record.ID = uuid.New()
	record.RecordID = FormatRecordID(2026, 1)
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.ImmovableRecord, error) {
	if m.record == nil || m.record.ID != id {
		return nil, nil
	}
	copy := *m.record
	return &copy, nil
}

func (m *mockRepo) List(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]models.ImmovableRecord, int, error) {
	return nil, 0, nil
}

func (m *mockRepo) Update(ctx context.Context, record *models.ImmovableRecord) error {
	return nil
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error {
	return nil
}

func (m *mockRepo) FinalApprove(ctx context.Context, id uuid.UUID, fromStatus models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error) {
	return time.Now(), nil
}
