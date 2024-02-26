package repocfg

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/runatlantis/atlantis/server/core/config/raw"
)

// Tests the ProjectFrom method for a project
func Test_ProjectsFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		component Component
		options   Options
		want      []ExtRawProject
	}{
		{
			component: Component{
				Path:     "test",
				VarFiles: []string{"env.tfvars"},
			},
			options: Options{
				Autoplan:                true,
				DefaultTerraformVersion: "8.8.8",
				Parallel:                true,
				UseWorkspaces:           true,
			},
			want: []ExtRawProject{
				{
					Project: raw.Project{
						Name:             ptr("test-env"),
						Dir:              ptr("test"),
						Workflow:         ptr("test-env"),
						Workspace:        ptr("env"),
						TerraformVersion: ptr("8.8.8"),
						Autoplan: &raw.Autoplan{
							Enabled: ptr(true),
							WhenModified: []string{
								"*.tf",
								"env.tfvars",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		got, err := ProjectsFrom(tc.component, tc.options)
		if err != nil {
			t.Errorf("ProjectFrom() error: %s", err)
		}

		if !cmp.Equal(got, tc.want) {
			t.Errorf(`ProjectFrom()
				diff %s`, cmp.Diff(got, tc.want))
		}
	}
}

// Tests the AutoPlan method for a project
func Test_AutoPlan(t *testing.T) {
	t.Parallel()

	want := &raw.Autoplan{
		Enabled: ptr(true),
		WhenModified: []string{
			"*.tf",
			"env.tfvars",
		},
	}

	got := new(ExtRawProject)

	got.AutoPlan("env.tfvars")

	if !cmp.Equal(got.Autoplan, want) {
		t.Errorf(`AutoPlan()
		diff %s`, cmp.Diff(got, want))
	}
}

// Tests the DefaultTerraformVersion method for a project
func Test_DefaultTerraformVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		version string
		want    raw.Project
	}{
		{
			version: "8.8.8",
			want: raw.Project{
				TerraformVersion: ptr("8.8.8"),
			},
		},
		{
			version: "invalid",
			want: raw.Project{
				TerraformVersion: new(string),
			},
		},
	}

	for _, tc := range tests {
		new := new(ExtRawProject)
		new.DefaultTerraformVersion(tc.version)
		got := new.Project.TerraformVersion

		if !cmp.Equal(got, tc.want.TerraformVersion) {
			t.Errorf(`DefaultTerraformVersion()
			diff %s`, cmp.Diff(got, tc.want))
		}
	}
}
