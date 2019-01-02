// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
)

const execScript = `#!/bin/bash -e

%s
`

// Run executes a recipe.
func Run(ctx context.Context, recipe *api.Recipe) (*api.Result, error) {
	result := &api.Result{}
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
	result.Output += string(out)
	if err != nil {
		result.Output += "error: " + err.Error()
	}
	result.Success = true
	return result, nil
}
