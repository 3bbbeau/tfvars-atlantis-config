package repocfg

import (
	"fmt"

	"github.com/runatlantis/atlantis/server/core/config/raw"
	"gopkg.in/yaml.v2"
)

var ErrNoExistingConfig = fmt.Errorf("no existing config found")

// Options represents the top-level configuration for a new Atlantis RepoCfg
type Options struct {
	Automerge               bool
	Autoplan                bool
	DefaultTerraformVersion string
	Parallel                bool
	MultiEnv                bool
	UseWorkspaces           bool
}

// Component represents a Terraform component and its associated Terraform variable files
type Component struct {
	Path     string
	VarFiles []string
}

// ExtRawRepoCfg is an embedded type for a raw.RepoCfg
type ExtRawRepoCfg struct {
	raw.RepoCfg `yaml:",inline"`
}

// NewRepoCfg returns a new Atlantis RepoCfg from a slice of components
// and options.
//
// Reference: https://www.runatlantis.io/docs/repo-level-atlantis-yaml.html
func NewRepoCfg(components []Component, opts Options) (*ExtRawRepoCfg, error) {
	repoCfg := &ExtRawRepoCfg{
		RepoCfg: raw.RepoCfg{
			Version:       ptr(3),
			Automerge:     &opts.Automerge,
			ParallelPlan:  &opts.Parallel,
			ParallelApply: &opts.Parallel,
		},
	}

	var projects []raw.Project
	for _, c := range components {
		generated, err := ProjectsFrom(c, opts)
		if err != nil {
			return nil, fmt.Errorf("failed while creating projects with component %+v: %w", c, err)
		}

		for _, p := range generated {
			projects = append(projects, p.Project)
		}
	}

	repoCfg.Projects = append(repoCfg.Projects, projects...)

	workflows := map[string]raw.Workflow{}
	for _, c := range components {
		generated, err := WorkflowsFrom(c, opts)
		if err != nil {
			return nil, fmt.Errorf("failed while creating workflows with component %+v: %w", c, err)
		}
		for _, wf := range generated {
			workflows[wf.Name] = wf.Workflow
		}
	}

	repoCfg.Workflows = workflows

	return repoCfg, nil
}

func (rc *ExtRawRepoCfg) MarshalYAML() (interface{}, error) {
	m := yaml.MapSlice{
		{Key: "version", Value: rc.Version},
		{Key: "automerge", Value: rc.Automerge},
		{Key: "parallel_plan", Value: rc.ParallelPlan},
		{Key: "parallel_apply", Value: rc.ParallelApply},
		{Key: "projects", Value: rc.Projects},
	}

	workflows := yaml.MapSlice{}
	for name, wf := range rc.Workflows {
		workflows = append(workflows, yaml.MapItem{
			Key: name,
			Value: yaml.MapSlice{
				{Key: "plan", Value: wf.Plan},
				{Key: "apply", Value: wf.Apply},
			},
		})
	}

	m = append(m, yaml.MapItem{Key: "workflows", Value: workflows})
	return m, nil
}
