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
					Projects: []raw.Project{
						{
							Name: ptr("test-dev"),
							Dir:  ptr("test"),
						},
						{
							Name: ptr("test-stg"),
							Dir:  ptr("test"),
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
