package photos

import (
	"testing"
)

func TestDetectImageExtension(t *testing.T) {
	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	ext, err := detectImageExtension(jpeg)
	if err != nil || ext != ".jpg" {
		t.Fatalf("expected .jpg, got %s err=%v", ext, err)
	}

	png := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	ext, err = detectImageExtension(png)
	if err != nil || ext != ".png" {
		t.Fatalf("expected .png, got %s err=%v", ext, err)
	}

	_, err = detectImageExtension([]byte("not-an-image"))
	if err != ErrUnsupportedType {
		t.Fatalf("expected ErrUnsupportedType, got %v", err)
	}
}
