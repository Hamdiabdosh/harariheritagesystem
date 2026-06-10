package dashboard

import (
	"testing"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

func TestBuildTableFiltersIndexed_chainsPlaceholders(t *testing.T) {
	userID := uuid.MustParse("a63c44c9-bbd2-44d6-8af6-9c9d62098be1")
	filters := ListFilters{Status: models.StatusPendingReview}

	immWhere, immArgs, nextIdx := buildTableFiltersIndexed("immovable_records", userID, models.RoleRegistrar, filters, 1)
	movWhere, movArgs, _ := buildTableFiltersIndexed("movable_records", userID, models.RoleRegistrar, filters, nextIdx)

	if immWhere != "1=1 AND registrar_id = $1 AND status = $2" {
		t.Fatalf("unexpected immovable where: %s", immWhere)
	}
	if movWhere != "1=1 AND registrar_id = $3 AND status = $4" {
		t.Fatalf("unexpected movable where: %s", movWhere)
	}
	if len(immArgs) != 2 || len(movArgs) != 2 {
		t.Fatalf("expected 2 args each, got imm=%d mov=%d", len(immArgs), len(movArgs))
	}
}
