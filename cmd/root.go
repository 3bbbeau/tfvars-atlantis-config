package cmd

import (
	"fmt"

	"github.com/3bbbeau/tfvars-atlantis-config/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// New returns the new root command for this utility.
//
// It initializes the root command and adds all subcommands.
func New() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "tfvars-atlantis-config",
		Short: "Generates Atlantis Config for Terraform projects",
		Long: `tfvars-atlantis-config is a utility that generates Atlantis configurations for Terraform projects
		that use tfvars files per environment.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogger(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				logger.FromContext(cmd.Context()).Error("help cmd", zap.Error(err))
			}
		},
	}
	cmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	cmd.AddCommand(NewVersionCmd())
	gCmd, err := NewGenerateCmd()
	if err != nil {
		return nil, fmt.Errorf("creating generate command: %w", err)
	}
	cmd.AddCommand(gCmd)
	cmd.AddCommand(NewMultiEnvCmd())

	return cmd, nil
}

// setLogger sets the logger for the command's lifespan, based on the debug flag.
func setLogger(cmd *cobra.Command) {
	ctx := cmd.Context()

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil || !debug {
		z := zap.Must(zap.NewProductionConfig().Build())
		ctx = logger.WithContext(ctx, z)
	} else {
		z := zap.Must(zap.NewDevelopmentConfig().Build())
		ctx = logger.WithContext(ctx, z)
	}

	cmd.SetContext(ctx)
}
