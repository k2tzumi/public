// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package supervisor provides supervisor trees for Go applications.

DEPRECATED: consider using cirello.io/oversight

This package implements supervisor trees, similar to what Erlang runtime offers.
It is built on top of context package, with all of its advantages, namely the
possibility trickle down context-related values and cancelation signals.

A supervisor tree can be composed either of services or other supervisors - each
supervisor can have its own set of configurations. Any instance of
supervisor.Service can be added to a tree.

	Supervisor
	     ├─▶ Supervisor (if one service dies, only one is restarted)
	     │       ├─▶ Service
	     │       └─▶ Service
	     ├─▶ Group (if one service dies, all others are restarted too)
	     │       └─▶ Service
	     │           Service
	     │           Service
	     └─▶ Service

Example:
	package main

	import (
		"fmt"
		"os"
		"os/signal"
		"time"

		"cirello.io/supervisor"
		"context"
	)

	type Simpleservice int

	func (s *Simpleservice) String() string {
		return fmt.Sprintf("simple service %d", int(*s))
	}

	func (s *Simpleservice) Serve(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("do something...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	func main(){
		var supervisor supervisor.Supervisor

		svc := Simpleservice(1)
		supervisor.Add(&svc)

		// Simply, if not special context is needed:
		// supervisor.Serve()
		// Or, using context.Context to propagate behavior:
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		ctx, cancel := context.WithCancel(context.Background())
		go func(){
			<-c
			fmt.Println("halting supervisor...")
			cancel()
		}()

		supervisor.Serve(ctx)
	}

TheJerf's blog post about Suture is a very good and helpful read to understand
how this package has been implemented.

This is package is inspired by github.com/thejerf/suture

http://www.jerf.org/iri/post/2930
*/
package supervisor // import "cirello.io/supervisor"
