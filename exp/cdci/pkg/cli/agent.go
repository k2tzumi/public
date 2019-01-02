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

package cli

import (
	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/agent"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
)

func (c *commands) agentMode() cli.Command {
	return cli.Command{
		Name:        "agent",
		Aliases:     []string{"worker", "builder"},
		Usage:       "start agent mode",
		Description: "start agent mode",
		Action: func(ctx *cli.Context) error {
			// Set up a connection to the server.
			conn, err := grpc.Dial("127.0.0.1:9999",
				grpc.WithInsecure())
			if err != nil {
				return errors.E(err, "did not connect")
			}
			defer conn.Close()
			agent := agent.New(1, conn)
			return errors.E(agent.Run(), "error running agent")
		},
	}
}
