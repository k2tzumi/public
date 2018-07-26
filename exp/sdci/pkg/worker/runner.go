package worker // import "cirello.io/exp/sdci/pkg/worker"

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/models"
)

const execScript = `#!/bin/bash

set -e

%s
`

func run(ctx context.Context, recipe *models.Recipe, repoDir string) (string, error) {
	tmpfile, err := ioutil.TempFile(repoDir, "agent")
	if err != nil {
		return "", errors.E(errors.FailedPrecondition, err,
			"agent cannot create temporary file")
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	fmt.Fprintf(tmpfile, execScript, recipe.Commands)
	tmpfile.Close()
	cmd := exec.CommandContext(ctx, "/bin/sh", tmpfile.Name())
	cmd.Dir = repoDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, strings.Split(recipe.Environment, "\n")...)
	var buf crbuffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err = cmd.Run()
	return buf.String(), errors.E(err, "failed when running builder")
}
