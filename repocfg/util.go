package repocfg

import (
	"path/filepath"
	"strings"
)

// Ptr returns a pointer to type T
func ptr[T any](v T) *T { return &v }

// friendlyName creates a contextual name used for Atlantis projects
func friendlyName(path, environment string) string {
	name := []string{strings.ReplaceAll(path, "/", "-"), pathWithoutExtension(filepath.Base(environment))}
	return strings.TrimSuffix(strings.Join(name, "-"), "-")
}

// pathWithoutExtension removes the file extension from a base path.
// if the path has no extension ("."), it returns an empty string.
func pathWithoutExtension(path string) string {
	return strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), filepath.Base(path))
}
