package cmd

import (
	"bytes"
	"testing"
)

// Tests the NewVersionCmd function returns the version of the utility.
func Test_NewVersionCmd(t *testing.T) {
	t.Parallel()

	got := new(bytes.Buffer)
	cmd, _ := New()
	cmd.SetArgs([]string{"version"})
	cmd.SetOutput(got)
	cmd.Execute() // nolint:errcheck
	want := "devel\n"

	if got.String() != want {
		t.Errorf(`NewVersionCmd()
		got %s
		want %v`, got.String(), want)
	}
}
