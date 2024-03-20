package repocfg

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/runatlantis/atlantis/server/core/config/raw"
)

// ExtRawProject extends the raw.Project type to add additional methods
type ExtRawProject struct {
	raw.Project
}

// ErrProjectFrom represents an error when creating an Atlantis project from a
// Terraform component along with some useful contextual information.
type ErrNewProject struct {
	Err       error
	Component Component
	Project   *ExtRawProject
}

// Stringer implementation for ErrProjectFrom
func (e ErrNewProject) Error() string {
	return fmt.Sprintf(`failed to create project:
		component: %+v,
		project: %+v,
		error: %s`, e.Component, e.Project, e.Err)
}

// ProjectsFrom creates an Atlantis project from a Terraform component and
// its associated Terraform variable files.
func ProjectsFrom(c Component, opts Options) ([]ExtRawProject, error) {
	var projects []ExtRawProject

	for _, v := range c.VarFiles {
		p := ExtRawProject{
			Project: raw.Project{
				Name: ptr(friendlyName(c.Path, v)),

				// The directory of this project relative to the repo root.
				Dir: ptr(c.Path),
			},
		}

		if opts.UseWorkspaces {
			// Terraform workspaces are represented by the Terraform variable file
			// names.
			//
			// Example:
			//
			// dev.tfvars -> dev
			// stg.tfvars.json -> stg
			p.Workspace = ptr(pathWithoutExtension(v))
		}

		// Generate a default Terraform version for the project if enabled
		if opts.DefaultTerraformVersion != "" {
			p.DefaultTerraformVersion(opts.DefaultTerraformVersion)
		}

		// Generate autoplan configuration for the project if enabled
		if opts.Autoplan {
			p.AutoPlan(v)
		}

		// We can validate the project using the Atlantis validate method
		err := p.Validate()
		if err != nil {
			return nil, fmt.Errorf("failed to validate project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// AutoPlan sets the autoplan configuration for the project
//
// This will trigger a plan when the Terraform files or the variable files
// are modified.
func (p *ExtRawProject) AutoPlan(v string) {
	autoplan := &raw.Autoplan{
		Enabled: ptr(true),

		// Paths are relative to the project directory
		WhenModified: []string{
			"*.tf",
			filepath.Base(v),
		},
	}
	p.Autoplan = autoplan
}

// DefaultTerraformVersion sets the default Terraform version for the project
// if the version is valid.
func (p *ExtRawProject) DefaultTerraformVersion(v string) {
	ver, err := version.NewSemver(v)
	if err != nil {
		p.TerraformVersion = new(string)
	} else {
		p.TerraformVersion = ptr(ver.String())
	}
}
