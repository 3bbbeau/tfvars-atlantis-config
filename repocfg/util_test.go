package repocfg

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tests the ptr helper function. Given a generic value it should return a
// pointer to that value.
func Test_Ptr(t *testing.T) {
	t.Parallel()

	intV := 3
	strV := "test"
	boolV := true

	for _, T := range []any{intV, strV, boolV} {
		got := ptr(T)
		if *got != T {
			t.Errorf(`ptr()
			got %v
			want %v
			diff %s`, got, T, cmp.Diff(got, T))
		}
	}
}

// Tests the friendlyName helper function. Given a path and an environment it
// should provide a contextual name to be used for Atlantis projects.
func Test_FriendlyName(t *testing.T) {
	t.Parallel()

	path := "my/path/to/some/terraform/component"
	environment := "my/path/to/some/terraform/component/dev.tfvars"

	want := "my-path-to-some-terraform-component-dev"
	got := friendlyName(path, environment)
	if got != want {
		t.Errorf(`friendlyName()
		got %v
		want %v
		diff %s`, got, want, cmp.Diff(got, want))
	}
}

func Test_PathWithoutExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want string
	}{
		{
			path: "dev.tfvars",
			want: "dev",
		},
		{
			path: "dev",
			want: "",
		},
	}

	for _, tc := range tests {
		got := pathWithoutExtension(tc.path)
		if got != tc.want {
			t.Errorf(`pathWithoutExtension()
			got %v
			want %v
			diff %s`, got, tc.want, cmp.Diff(got, tc.want))
		}
	}
}
