package agent // import "cirello.io/exp/cdci/pkg/agent"

import (
	"context"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
)

// Agent runs jobs locally.
type Agent struct {
	agentID int64
	client  api.RunnerClient
}

// New prepares a new agent.
func New(agentID int64, conn *grpc.ClientConn) *Agent {
	return &Agent{
		agentID: agentID,
		client:  api.NewRunnerClient(conn),
	}
}

// Run reacts to requests from the server, in absence of work, it return pings.
func (a *Agent) Run() error {
	pipe, err := a.client.Run(context.Background())
	if err != nil {
		return errors.E(err, "cannot talk to server")
	}

	spew.Dump("I connected to server")

	ctx := pipe.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			r, err := pipe.Recv()
			if err != nil {
				spew.Dump("oops", err)
				return errors.E(err,
					"cannot receive request from server")
			}

			switch v := r.GetAction().(type) {
			case *api.RunRequest_Ping:
				pipe.Send(pongMessage(a.agentID))
				spew.Dump("I just pong")
			default:
				spew.Dump("something else???", v)
			}
		}
	}
}

func pongMessage(agentID int64) *api.RunResponse {
	return &api.RunResponse{
		Response: &api.RunResponse_Pong{
			Pong: &api.Pong{AgentID: agentID},
		},
	}
}
