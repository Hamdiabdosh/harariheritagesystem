package dashboard

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type mockRepo struct {
	stats *Stats
	items []RecordSummary
	total int
}

func (m *mockRepo) GetStats(_ context.Context, _ uuid.UUID, _ models.Role) (*Stats, error) {
	return m.stats, nil
}

func (m *mockRepo) ListRecords(_ context.Context, _ ListFilters, _ uuid.UUID, _ models.Role) ([]RecordSummary, int, error) {
	return m.items, m.total, nil
}

func TestListRecordsPagination(t *testing.T) {
	repo := &mockRepo{items: []RecordSummary{}, total: 45}
	svc := NewService(repo)

	result, err := svc.ListRecords(t.Context(), ListFilters{Page: 2, Limit: 20}, uuid.New(), models.RoleManager)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPages != 3 {
		t.Fatalf("expected 3 total pages, got %d", result.TotalPages)
	}
	if result.Page != 2 {
		t.Fatalf("expected page 2, got %d", result.Page)
	}
}
