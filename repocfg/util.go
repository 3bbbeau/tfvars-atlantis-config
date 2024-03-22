package repocfg

import (
	"path/filepath"
	"strings"
)

// Ptr returns a pointer to type T
func ptr[T any](v T) *T { return &v }

// friendlyName creates a contextual name used for Atlantis projects
func friendlyName(path, varFile string) string {
	environment := pathWithoutExtension(filepath.Base(varFile))

	// avoid constructing a joined path if the context is the current directory
	if filepath.Base(path) == "." {
		return environment
	}

	name := []string{strings.ReplaceAll(path, "/", "-"), environment}
	return strings.TrimSuffix(strings.Join(name, "-"), "-")
}

// pathWithoutExtension removes the file extension from a base path.
// if the path has no extension ("."), it returns an empty string.
func pathWithoutExtension(path string) string {
	return strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), filepath.Base(path))
}
