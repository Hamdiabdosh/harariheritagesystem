package audit

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type mockRepo struct {
	history []models.StatusHistoryEntry
}

func (m *mockRepo) ListByRecord(_ context.Context, _ models.RecordType, _ uuid.UUID) ([]models.StatusHistoryEntry, error) {
	return m.history, nil
}

func TestServiceListByRecord(t *testing.T) {
	recordID := uuid.New()
	expected := []models.StatusHistoryEntry{
		{
			ID:            uuid.New(),
			RecordType:    models.RecordTypeImmovable,
			RecordID:      recordID,
			ChangedBy:     uuid.New(),
			ChangedByName: "Supervisor One",
			ToStatus:      string(models.StatusUnderReview),
		},
	}

	repo := &mockRepo{history: expected}
	svc := NewService(repo)

	got, err := svc.ListByRecord(t.Context(), models.RecordTypeImmovable, recordID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].ChangedByName != expected[0].ChangedByName {
		t.Fatalf("expected changed_by_name %q, got %q", expected[0].ChangedByName, got[0].ChangedByName)
	}
}
