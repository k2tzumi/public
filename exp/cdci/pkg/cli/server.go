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
	"net"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
	"cirello.io/exp/cdci/pkg/server"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
)

func (c *commands) serverMode() cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"dispatcher", "queue"},
		Usage:       "start server mode",
		Description: "start server mode",
		Action: func(ctx *cli.Context) error {
			tasks := make(chan *api.Recipe, 1)
			tasks <- &api.Recipe{
				Id:          1,
				Environment: []string{"WHO=world"},
				Commands:    "echo Hello, $WHO;",
			}
			// close(tasks)
			l, err := net.Listen("tcp", ":9999")
			if err != nil {
				return errors.E(err, "failed to listen")
			}
			s := grpc.NewServer()
			api.RegisterRunnerServer(s, server.New(tasks))
			return errors.E("failed to serve", s.Serve(l))
		},
	}
}
