package main

import (
	"os"

	"github.com/3bbbeau/tfvars-atlantis-config/cmd"
)

func main() {
	cmd, err := cmd.New()
	if err != nil {
		os.Exit(1)
	}

	err = cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
