package worker // import "cirello.io/exp/sdci/pkg/worker"

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/grpc/api"
)

const execScript = `#!/bin/bash

set -e

%s
`

func run(ctx context.Context, recipe *api.Recipe, repoDir, baseDir string) (string, error) {
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
	cmd.Env = append(cmd.Env, fmt.Sprintf("SDCI_BUILD_BASE_DIRECTORY=%s", baseDir))
	recipeEnvVars := strings.Split(recipe.Environment, "\n")
	for _, v := range recipeEnvVars {
		cmd.Env = append(cmd.Env, os.Expand(v, expandVar(cmd.Env)))
	}
	var buf crbuffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err = cmd.Run()
	return buf.String(), errors.E(err, "failed when running builder")
}

func expandVar(currentEnv []string) func(string) string {
	return func(s string) string {
		for _, e := range currentEnv {
			if strings.HasPrefix(e, s+"=") {
				return strings.TrimPrefix(e, s+"=")
			}
		}
		return ""
	}
}
