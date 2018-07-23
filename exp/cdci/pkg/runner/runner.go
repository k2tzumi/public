package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"cirello.io/errors"
	"cirello.io/exp/cdci/api"
)

const execScript = `#!/bin/bash -e

%s
`

// Run executes a recipe.
func Run(ctx context.Context, recipe *api.Recipe) (*api.RunResponse, error) {
	resp := &api.RunResponse{}
	tmpfile, err := ioutil.TempFile("", "agent")
	if err != nil {
		return nil, errors.E(errors.FailedPrecondition, err,
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
	resp.Output += string(out)
	if err != nil {
		resp.Output += "error: " + err.Error()
	}
	resp.Success = true
	return resp, nil
}
