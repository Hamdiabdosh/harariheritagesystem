package photos

import (
	"path/filepath"
	"testing"
)

func TestResolveAbsolutePath(t *testing.T) {
	root := "/var/media"
	got := ResolveAbsolutePath(root, "abc.jpg")
	want := filepath.Join(root, "abc.jpg")
	if got != want {
		t.Fatalf("flat path: got %q want %q", got, want)
	}

	legacy := "immovable/record-id/abc.jpg"
	got = ResolveAbsolutePath(root, legacy)
	want = filepath.Join(root, "immovable", "record-id", "abc.jpg")
	if got != want {
		t.Fatalf("legacy path: got %q want %q", got, want)
	}
}
