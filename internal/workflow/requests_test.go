package workflow

import "testing"

func TestOptionalCommentRequestPrefersCommentText(t *testing.T) {
	text := "primary"
	legacy := "fallback"
	req := optionalCommentRequest{CommentText: &text, Comment: &legacy}
	got := req.value()
	if got == nil || *got != "primary" {
		t.Fatalf("expected primary, got %v", got)
	}
}

func TestOptionalCommentRequestLegacyComment(t *testing.T) {
	legacy := "legacy note"
	req := optionalCommentRequest{Comment: &legacy}
	got := req.value()
	if got == nil || *got != "legacy note" {
		t.Fatalf("expected legacy note, got %v", got)
	}
}

func TestRequiredCommentRequestCommentText(t *testing.T) {
	req := requiredCommentRequest{CommentText: "return reason"}
	got, ok := req.value()
	if !ok || got != "return reason" {
		t.Fatalf("expected return reason, got %q ok=%v", got, ok)
	}
}

func TestRequiredCommentRequestLegacyComment(t *testing.T) {
	req := requiredCommentRequest{Comment: "legacy return"}
	got, ok := req.value()
	if !ok || got != "legacy return" {
		t.Fatalf("expected legacy return, got %q ok=%v", got, ok)
	}
}

func TestRequiredCommentRequestEmpty(t *testing.T) {
	req := requiredCommentRequest{CommentText: "   "}
	if got, ok := req.value(); ok {
		t.Fatalf("expected empty, got %q", got)
	}
}
