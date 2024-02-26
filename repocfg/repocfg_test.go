package repocfg

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/runatlantis/atlantis/server/core/config/raw"
)

func Test_NewFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		components []Component
		want       *ExtRawRepoCfg
	}{
		{
			name: "WithNestedVars",
			components: []Component{
				{
					Path:     "test",
					VarFiles: []string{"test/vars/dev.tfvars", "test/vars/nested/stg.tfvars"},
				},
			},
			want: &ExtRawRepoCfg{
				RepoCfg: raw.RepoCfg{
					Version:       ptr(3),
					Automerge:     ptr(false),
					ParallelPlan:  ptr(false),
					ParallelApply: ptr(false),
					Workflows: map[string]raw.Workflow{
						"test-dev": {
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
						"test-stg": {
							Plan: &raw.Stage{
								Steps: []raw.Step{
									{
										Key: ptr("init"),
									},
									{
										Map: map[string]map[string][]string{
											"plan": {"extra_args": []string{"-var-file=vars/nested/stg.tfvars"}},
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
					Projects: []raw.Project{
						{
							Name:     ptr("test-dev"),
							Dir:      ptr("test"),
							Workflow: ptr("test-dev"),
						},
						{
							Name:     ptr("test-stg"),
							Dir:      ptr("test"),
							Workflow: ptr("test-stg"),
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		got, _ := NewRepoCfg(tc.components, Options{})
		if !cmp.Equal(got, tc.want) {
			t.Errorf(`NewFrom()
				diff %s`, cmp.Diff(got, tc.want))
		}
	}
}
