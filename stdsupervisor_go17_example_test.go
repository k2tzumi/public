// +build go1.7

package supervisor_test

import (
	"context"
	"fmt"
	"sync"

	"cirello.io/supervisor"
)

type Simpleservice struct {
	id int
	sync.WaitGroup
}

func (s *Simpleservice) Serve(ctx context.Context) {
	fmt.Println(s.String())
	s.Done()
	<-ctx.Done()
}

func (s *Simpleservice) String() string {
	return fmt.Sprintf("simple service %d", s.id)
}

func ExampleServeContext() {
	svc := &Simpleservice{id: 1}
	svc.Add(1)
	supervisor.Add(svc)

	ctx, cancel := context.WithCancel(context.Background())
	go supervisor.ServeContext(ctx)

	svc.Wait()
	cancel()

	// output:
	// simple service 1
}

func ExampleServeGroupContext() {
	svc1 := &Simpleservice{id: 1}
	svc1.Add(1)
	supervisor.Add(svc1)
	svc2 := &Simpleservice{id: 2}
	svc2.Add(1)
	supervisor.Add(svc2)

	ctx, cancel := context.WithCancel(context.Background())
	go supervisor.ServeGroupContext(ctx)

	svc1.Wait()
	svc2.Wait()
	cancel()

	// unordered output:
	// simple service 1
	// simple service 2
}

func ExampleServe() {
	svc := &Simpleservice{id: 1}
	svc.Add(1)
	supervisor.Add(svc)

	var cancel context.CancelFunc
	ctx, cancel := context.WithCancel(context.Background())
	supervisor.SetDefaultContext(ctx)
	go supervisor.Serve()

	svc.Wait()
	cancel()

	// output:
	// simple service 1
}

func ExampleServeGroup() {
	svc1 := &Simpleservice{id: 1}
	svc1.Add(1)
	supervisor.Add(svc1)
	svc2 := &Simpleservice{id: 2}
	svc2.Add(1)
	supervisor.Add(svc2)

	ctx, cancel := context.WithCancel(context.Background())
	supervisor.SetDefaultContext(ctx)
	go supervisor.ServeGroup()

	svc1.Wait()
	svc2.Wait()
	cancel()

	// unordered output:
	// simple service 1
	// simple service 2
}
