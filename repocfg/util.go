package repocfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Ptr returns a pointer to type T
func ptr[T any](v T) *T { return &v }

// friendlyName creates a contextual name used for Atlantis projects and workflows
func friendlyName(path, environment string) string {
	name := []string{strings.ReplaceAll(path, "/", "-"), pathWithoutExtension(filepath.Base(environment))}
	return strings.TrimSuffix(strings.Join(name, "-"), "-")
}

// pathWithoutExtension removes the file extension from a base path.
// if the path has no extension ("."), it returns an empty string.
func pathWithoutExtension(path string) string {
	return strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), filepath.Base(path))
}

var ErrNoEnvVars error = fmt.Errorf("no matching prefixed environment variables found")

// Generates the Atlantis multienv string within stages for multi-environment
// Terraform projects, e.g:
//
//	EnvVar1Name=value1,EnvVar2Name=value2,EnvVar3Name=value3
//
// Given a prefix for the environment name of the workflow, strips the
// environment name from existing environment vars while keeping their value.
//
// This is useful when you want to configure providers via environment variables
// on a per-environment basis.
//
// Example:
//
//	DEV_AWS_ACCESS_KEY_ID="foo"
//	DEV_AWS_SECRET_ACCESS_KEY="bar"
//	->
//	AWS_ACCESS_KEY_ID=$DEV_AWS_ACCESS_KEY_ID
//	AWS_SECRET_ACCESS_KEY=$DEV_AWS_SECRET_ACCESS_KEY
func prefixedEnviron(prefix string) (*string, error) {
	prefix = strings.ToUpper(prefix)
	if !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	strippedEnviron := []string{}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, prefix) {
			// Limits the split to only two parts, separating the key from the first
			// occurrence of '=', otherwise if the value contains '=' character(s) the
			// string would be split into more than two parts.
			name := strings.SplitN(v, "=", 2)[0]

			// Strips the prefix from the environment variable name, e.g. "DEV_" from
			// "DEV_AWS_ACCESS_KEY_ID" and let it equal to the original environment variable
			strippedEnviron = append(strippedEnviron, fmt.Sprintf("%s=$%s", strings.TrimPrefix(name, prefix), name))
		}
	}
	if len(strippedEnviron) == 0 {
		return nil, ErrNoEnvVars
	}
	return ptr(fmt.Sprintf("echo %s", strings.Join(strippedEnviron, ","))), nil
}
