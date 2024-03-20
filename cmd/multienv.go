package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/3bbbeau/tfvars-atlantis-config/logger"
	"github.com/spf13/cobra"
)

const (
	// ATLANTIS_WORKSPACE is the environment variable that contains the name of
	// the workspace that Atlantis is currently running for.
	ATLANTIS_WORKSPACE = "WORKSPACE"
)

// NewMultiEnvCmd creates a new `multienv` command
func NewMultiEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multienv",
		Short: "Returns a string representing the multienv for Atlantis",
		Long: `Returns a string representing the multienv for Atlantis, e.g.:
		EnvVar1Name=value1,EnvVar2Name=value2,EnvVar3Name=value3

		Reference: https://www.runatlantis.io/docs/custom-workflows.html#multiple-environment-variables-multienv-command

		This is useful when you want to configure providers via environment variables
		on a per-workspace/environment basis.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace := os.Getenv(ATLANTIS_WORKSPACE)
			if len(workspace) == 0 {
				return fmt.Errorf("environment variable %s is not set. no workspace for multienv", ATLANTIS_WORKSPACE)
			}

			multienv, err := multienv(os.Getenv(ATLANTIS_WORKSPACE))
			if err != nil {
				if errors.Is(err, ErrNoEnvVars) {
					logger.FromContext(cmd.Context()).Debug("no matching prefixed environment variables found")
					// Return nil to indicate success, but no multienv string
					return nil
				}
				return err
			}

			// Atlanis expects the multienv string to be written to stdout
			fmt.Print(multienv)
			return nil
		},
	}
	return cmd
}

var ErrNoEnvVars error = fmt.Errorf("no matching prefixed environment variables found")

// Generates the Atlantis multienv string for multi-environment
// Terraform projects, e.g:
//
//	EnvVar1Name=value1,EnvVar2Name=value2,EnvVar3Name=value3
//
// Given a prefix for the workspace name, strips the
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
func multienv(prefix string) (string, error) {
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
			split := strings.SplitN(v, "=", 2)

			// Strips the prefix from the environment variable name, e.g. "DEV_" from
			// "DEV_AWS_ACCESS_KEY_ID" and let it equal to the original environment variable
			strippedEnviron = append(strippedEnviron, fmt.Sprintf("%s=%s", strings.TrimPrefix(split[0], prefix), split[1]))
		}
	}
	if len(strippedEnviron) == 0 {
		return "", ErrNoEnvVars
	}
	return strings.Join(strippedEnviron, ","), nil
}
