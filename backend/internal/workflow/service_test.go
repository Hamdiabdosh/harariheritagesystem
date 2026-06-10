package workflow

import (
	"testing"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

func TestReviewReturnRequiresComment(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)
	_, err := svc.ReviewReturn(t.Context(), models.RecordTypeImmovable, uuid.New(), uuid.New(), "   ")
	if err != ErrCommentRequired {
		t.Fatalf("expected ErrCommentRequired, got %v", err)
	}
}

func TestTrimComment(t *testing.T) {
	text := "  hello  "
	if got := trimComment(&text); got != "hello" {
		t.Fatalf("expected hello, got %q", got)
	}
}
