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

package agent // import "cirello.io/exp/cdci/pkg/agent"

import (
	"context"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
	"cirello.io/exp/cdci/pkg/runner"
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
			case *api.RunRequest_Recipe:
				result, err := runner.Run(context.TODO(), v.Recipe)
				result.AgentID = a.agentID
				spew.Dump(result, err)
				pipe.Send(resultMessage(result))
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

func resultMessage(result *api.Result) *api.RunResponse {
	return &api.RunResponse{
		Response: &api.RunResponse_Result{
			Result: result,
		},
	}
}
