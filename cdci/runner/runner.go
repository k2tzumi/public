package runner

import (
	"context"
	"os"
	"os/exec"

	"cirello.io/cdci/api"
)

// Run executes a recipe.
func Run(ctx context.Context, recipe *api.Recipe) *api.RunResponse {
	resp := &api.RunResponse{}
	for _, r := range recipe.Steps {
		for _, c := range r.Commands {
			cmd := exec.CommandContext(ctx, "/bin/sh", "-c", c)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, recipe.Environment...)
			cmd.Env = append(cmd.Env, r.Environment...)
			out, err := cmd.CombinedOutput()
			resp.Output += string(out)
			if err != nil {
				resp.Output += "error: " + err.Error()
			}
		}
	}

	return resp
}
