package repocfg

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/runatlantis/atlantis/server/core/config/raw"
)

// Tests creating a new workflow from a Terraform component
func Test_WorkflowsFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		component Component
		want      []ExtRawWorkflow
	}{
		{
			name: "new-workflow",
			component: Component{
				Path:     "test",
				VarFiles: []string{"test/vars/dev.tfvars", "test/vars/stg.tfvars"},
			},
			want: []ExtRawWorkflow{
				{
					Name: "test-dev",
					Args: &ExtraArgs{
						"extra_args": []string{"-var-file=vars/dev.tfvars"},
					},
					Workspace: "dev",
					Workflow: raw.Workflow{
						Plan: &raw.Stage{
							Steps: []raw.Step{
								{
									Key: ptr("init"),
								},
								{
									Map: map[string]map[string][]string{
										"plan": {"extra_args": []string{"-var-file=vars/dev.tfvars"}},
									},
								},
							},
						},
						Apply: &raw.Stage{
							Steps: []raw.Step{
								{
									Key: ptr("apply"),
								},
							},
						},
					},
				},
				{
					Name: "test-stg",
					Args: &ExtraArgs{
						"extra_args": []string{"-var-file=vars/stg.tfvars"},
					},
					Workspace: "stg",
					Workflow: raw.Workflow{
						Plan: &raw.Stage{
							Steps: []raw.Step{
								{
									Key: ptr("init"),
								},
								{
									Map: map[string]map[string][]string{
										"plan": {"extra_args": []string{"-var-file=vars/stg.tfvars"}},
									},
								},
							},
						},
						Apply: &raw.Stage{
							Steps: []raw.Step{
								{
									Key: ptr("apply"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		got, err := WorkflowsFrom(tc.component, Options{})
		if err != nil {
			t.Errorf("WorkflowsFrom(): %s", err)
		}
		if !cmp.Equal(got, tc.want) {
			t.Errorf(`WorkflowFrom()
				diff %s`, cmp.Diff(got, tc.want))
		}
	}
}
