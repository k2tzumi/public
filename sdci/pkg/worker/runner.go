package worker // import "cirello.io/exp/sdci/pkg/worker"

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
)

const execScript = `#!/bin/bash

set -e

%s
`

func run(ctx context.Context, recipe *coordinator.Recipe, repoDir string) error {
	tmpfile, err := ioutil.TempFile(repoDir, "agent")
	if err != nil {
		return errors.E(errors.FailedPrecondition, err,
			"agent cannot create temporary file")
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	fmt.Fprintf(tmpfile, execScript, recipe.Commands)
	tmpfile.Close()
	cmd := exec.CommandContext(ctx, "/bin/sh", tmpfile.Name())
	cmd.Dir = repoDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, recipe.Environment...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.E(cmd.Run(), "failed when running builder")
}
