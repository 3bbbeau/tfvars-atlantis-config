package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/3bbbeau/tfvars-atlantis-config/logger"
	"github.com/3bbbeau/tfvars-atlantis-config/repocfg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// Flags represents the flags for the `generate` command
type Flags struct {
	AutoMerge               bool
	AutoPlan                bool
	DefaultTerraformVersion string
	MultiEnv                bool
	Output                  string
	Parallel                bool
	Root                    string
	UseWorkspaces           bool
}

// NewFlags returns a default Flags struct
func NewFlags() (*Flags, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("default flags for git root: %w", err)
	}

	return &Flags{
		AutoMerge:               false,
		AutoPlan:                false,
		DefaultTerraformVersion: "",
		Root:                    pwd,
		MultiEnv:                false,
		Output:                  "",
		Parallel:                false,
		UseWorkspaces:           false,
	}, nil
}

// AddFlags registers flags for the `generate` command
func (flags *Flags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&flags.AutoPlan, "autoplan", flags.AutoPlan, "Enable auto plan. Default is disabled")
	cmd.Flags().BoolVar(&flags.AutoMerge, "automerge", flags.AutoMerge, "Enable auto merge. Default is disabled")
	cmd.Flags().BoolVar(&flags.Parallel, "parallel", flags.Parallel, "Enables plans and applys to happen in parallel. Default is disabled")
	cmd.Flags().StringVar(&flags.Output, "output", flags.Output, "Path of the file where configuration will be generated. Default is stdout")
	cmd.Flags().StringVar(&flags.Root, "root", flags.Root, "Path to the root directory of the git repo you want to build config for. Default is current dir")
	cmd.Flags().BoolVar(&flags.MultiEnv, "multienv", flags.MultiEnv, "Enable injection of environment specific environment variables to each workflow. Default is disabled")
	cmd.Flags().StringVar(&flags.DefaultTerraformVersion, "terraform-version", flags.DefaultTerraformVersion, "Default terraform version to run for Atlantis. Default is determined by the Terraform version constraints.")
	cmd.Flags().BoolVar(&flags.UseWorkspaces, "use-workspaces", flags.UseWorkspaces, "Use workspaces for projects. Default is disabled")

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		logger.FromContext(cmd.Context()).Sugar().Debugf("Set flag: %s = %v", f.Name, f.Value)
	})
}

// toOptions converts the flags provided for usage to Options within the repocfg package
func (flags *Flags) toOptions() repocfg.Options {
	return repocfg.Options{
		Automerge:               flags.AutoMerge,
		Autoplan:                flags.AutoPlan,
		DefaultTerraformVersion: flags.DefaultTerraformVersion,
		MultiEnv:                flags.MultiEnv,
		Parallel:                flags.Parallel,
		UseWorkspaces:           flags.UseWorkspaces,
	}
}

// NewGenerateCmd creates a new `generate` command, while applying all flags
// with their defaults overlayed by the flags passed in by the caller.
func NewGenerateCmd() (*cobra.Command, error) {
	flags, err := NewFlags()
	if err != nil {
		return nil, fmt.Errorf("new flags: %w", err)
	}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Makes Atlantis config for tfvar backed Terraform projects",
		Long:  "Logs Yaml representing Atlantis config to stderr",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := generate(cmd, flags)
			if err != nil {
				return err
			}
			return nil
		},
	}

	flags.AddFlags(cmd)

	return cmd, nil
}

// generate is used to generate an Atlantis config, called from the `generate`
// command, creating an Atlantis configuration written to flags.OutputPath
func generate(cmd *cobra.Command, flags *Flags) error {
	logger := logger.FromContext(cmd.Context())

	tree, err := discover(cmd.Context(), flags)
	if err != nil {
		return err
	}

	opts := flags.toOptions()
	cfg, err := repocfg.NewRepoCfg(tree, opts)
	if err != nil {
		return err
	}

	cfgBytes, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("repocfg: %w", err)
	}

	switch flags.Output {
	case "":
		fmt.Println(string(cfgBytes))
	default:
		logger.Sugar().Debugf("writing config to %s", flags.Output)
		err := os.WriteFile(flags.Output, cfgBytes, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	TF_EXT          = ".tf"
	TFVARS_EXT      = ".tfvars"
	TFVARS_JSON_EXT = ".tfvars.json"
)

// discover walks the directory tree and creates a slice of Terraform components
// and their dependencies based on the .tfvars files in the directory and
// subdirectories.
//
// However, if nested directories contain a .tf file, the .tfvars files in the
// subdirectories are ignored, as they are assumed to be part of a different
// parent node (Terraform component).
//
// For example, if the function is invoked at the path defined by the `root` flag:
//
//	.
//	├── components
//	│   ├── component1
//	│   │   ├── main.tf
//	│   │   └── dev.tfvars
//	│   └── component2
//	│       ├── extraVars
//	│       │   └── stg.tfvars
//	│       ├── main.tf
//	│       └── dev.tfvars
//
// Then, the resulting slice of components would look like:
//
//	Path:	"components/component1"
//	VarFiles:	"components/component1/dev.tfvars"
//
//	Path:	"components/component2"
//	VarFiles:	[]string{"components/component2/dev.tfvars", "components/component2/extraVars/stg.tfvars"}
func discover(ctx context.Context, flags *Flags) ([]repocfg.Component, error) {
	logger := logger.FromContext(ctx)

	var discovered []repocfg.Component

	// Walk the directory tree from the root
	err := filepath.Walk(flags.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Each parent should be a path containing a Terraform component
		if !info.IsDir() && strings.HasSuffix(info.Name(), TF_EXT) {
			logger.Debug(fmt.Sprintf("found .tf file: %s", path))

			// Add the directory of the Terraform component as the parent
			parent, err := filepath.Rel(flags.Root, filepath.Dir(path))
			if err != nil {
				return err
			}

			// Each edge should be a path containing Terraform variables files
			// relative to the Terraform component
			err = filepath.Walk(filepath.Dir(path), func(subpath string, subinfo os.FileInfo, suberr error) error {
				if suberr != nil {
					return suberr
				}

				// Filter out Terraform variable files
				if !subinfo.IsDir() && (strings.HasSuffix(subinfo.Name(), TFVARS_EXT) || strings.HasSuffix(subinfo.Name(), TFVARS_JSON_EXT)) {
					// Ignore nested Terraform variable files that might belong to nested components
					if filepath.Dir(subpath) != filepath.Dir(path) && tfExistsInDir(filepath.Dir(subpath)) {
						logger.Sugar().Debugf(`ignoring nested %s: %s
						which belongs to the component at: %s
						not: %s`, filepath.Ext(subpath), subpath, filepath.Dir(subpath), filepath.Dir(path),
						)

						return filepath.SkipDir
					}

					logger.Sugar().Debugf("component parent is %s", parent)
					found := repocfg.Component{
						Path: parent,
					}

					logger.Sugar().Debugf("found %s file: %s", filepath.Ext(subpath), subpath)
					child, err := filepath.Rel(flags.Root, subpath)
					if err != nil {
						return err
					}

					// VarFile is a member of this component
					logger.Sugar().Debugf("component %s has var file %s", parent, child)
					found.VarFiles = append(found.VarFiles, child)

					discovered = append(discovered, found)
				}
				return err
			})
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return discovered, nil
}

// tfExistsInDir returns true if a directory contains a Terraform file
func tfExistsInDir(dir string) bool {
	files, err := filepath.Glob(filepath.Join(dir, "*"+TF_EXT))
	if err != nil {
		return false
	}

	return len(files) > 0
}
