package cli

import (
	"net"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
	"cirello.io/exp/cdci/pkg/server"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) serverMode() cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"dispatcher", "queue"},
		Usage:       "start server mode",
		Description: "start server mode",
		Action: func(ctx *cli.Context) error {
			l, err := net.Listen("tcp", ":9999")
			if err != nil {
				return errors.E(err, "failed to listen")
			}
			s := grpc.NewServer()
			api.RegisterRunnerServer(s, &server.Server{})
			return errors.E("failed to serve", s.Serve(l))
		},
	}
}
