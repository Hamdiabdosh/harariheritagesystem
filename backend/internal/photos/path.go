package photos

import (
	"path/filepath"
	"strings"
)

// ResolveAbsolutePath maps a stored file_path to an on-disk path under mediaRoot.
// Supports flat names (new) and legacy nested paths (recordType/recordId/file.ext).
func ResolveAbsolutePath(mediaRoot, filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	return filepath.Join(mediaRoot, filepath.FromSlash(strings.TrimPrefix(filePath, "/")))
}
