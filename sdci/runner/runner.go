package runner // import "cirello.io/exp/sdci/runner"

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"cirello.io/errors"
)

const execScript = `#!/bin/bash

set -e

%s
`

// Recipe defines the execution steps and environment.
type Recipe struct {
	Environment []string
	Commands    string
}

// Run executes a recipe.
func Run(ctx context.Context, recipe *Recipe) (string, error) {
	tmpfile, err := ioutil.TempFile("", "agent")
	if err != nil {
		return "", errors.E(errors.FailedPrecondition, err,
			"agent cannot create temporary file")
	}
	defer os.Remove(tmpfile.Name()) // clean up
	defer tmpfile.Close()
	fmt.Fprintf(tmpfile, execScript, recipe.Commands)
	tmpfile.Close()
	cmd := exec.CommandContext(ctx, "/bin/sh", tmpfile.Name())
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, recipe.Environment...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
