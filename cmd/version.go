package cmd

import (
	"github.com/spf13/cobra"
)

// Build-time version string, overriden by ldflags.
var v string = "devel"

// NewVersionCmd returns the command for retrieving the version of this utility.
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version of terragrunt-atlantis-config",
		Long:  "Version of terragrunt-atlantis-config",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(v)
		},
	}

	return cmd
}
