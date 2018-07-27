// Package server implements the server-side GRPC interface of the coordinator
// and workers.
package server // import "cirello.io/exp/sdci/pkg/grpc/server"
import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/grpc/api"
	"cirello.io/exp/sdci/pkg/models"
)

// Server interprets worker calls to coordinator
type Server struct {
	configuration models.Configuration
	coord         *coordinator.Coordinator

	locks       sync.Map // map of repoNames+workerID to lock state
	lockRefresh *sync.Cond
}

// New instantiates a new server
func New(configuration models.Configuration, coord *coordinator.Coordinator) *Server {
	s := &Server{
		configuration: configuration,
		coord:         coord,
		lockRefresh:   sync.NewCond(&sync.Mutex{}),
	}
	for repoName, recipe := range configuration {
		for i := 0; i < int(recipe.Concurrency); i++ {
			lockName := fmt.Sprintf("%v-%v", repoName, i)
			s.locks.Store(lockName, &lock{})
		}
	}
	go s.expireLocks()
	return s
}

// Run coordinates both the delivery of build job to workers and their actual
// liveness.
func (s *Server) Run(srv api.Runner_RunServer) error {
	// TODO: convert errors to GRPC's.
	// Run has two phases: a) handshake in which the worker declares which
	// repo it is observing; b) a continuous keep-alive to prove that the
	// worker is still alive.
	// The handshake takes the declared repository and tries to grab an
	// execution slot for the work. Even there are more workers than slot,
	// only n-slotted workers will actually have work to do.
	req, err := srv.Recv()
	if err != nil {
		return errors.E(err, "error receiving call from client")
	}
	buildReq := req.GetBuild()
	if buildReq == nil {
		return errors.E("client did not send the handshake message")
	}

	repoName := buildReq.GetRepoFullName()
	_ = repoName

	lockIndex, lockSeq, err := s.waitForLock(repoName)
	if err != nil {
		return errors.E(err, "cannot grab lock")
	}

	ctx, cancel := context.WithCancel(srv.Context())
	defer cancel()

	go func() {
		for {
			req, err := srv.Recv()
			if err != nil {
				err := errors.E(err, "error receiving call from client")
				log.Println(err)
				return
			}
			switch req.GetCommand().(type) {
			case *api.JobRequest_KeepAlive:
				if err := s.refreshLock(lockIndex, lockSeq); err != nil {
					err := errors.E(err, "lost lock")
					log.Println(err)
					cancel()
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case build := <-s.coord.Next(repoName):
			if err := s.dispatchBuild(srv, repoName, lockIndex, lockSeq, build); err != nil {
				cancel()
				s.coord.Recover(repoName, build)
				return err
			}
		}
	}
}

func (s *Server) dispatchBuild(srv api.Runner_RunServer, repoName, lockIndex string, lockSeq int, build *models.Build) error {
	if err := s.refreshLock(lockIndex, lockSeq); err != nil {
		return errors.E(err, "lost lock before build dispatch")
	}
	jobResp := &api.JobResponse{
		Build:  build.Build,
		Recipe: build.Recipe,
	}
	err := srv.Send(jobResp)
	return errors.E(err, "cannot dispatch build to client")
}

type lock struct {
	mu         sync.Mutex
	locked     bool
	seq        int
	lastUpdate time.Time
}

func (l *lock) tryLock() (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.locked {
		return 0, errors.E("already locked")
	}
	l.locked = true
	l.seq++
	l.lastUpdate = time.Now()
	return l.seq, nil
}

func (l *lock) refresh(seq int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.seq == seq {
		l.lastUpdate = time.Now()
	}
}

func (l *lock) release(seq int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.locked {
		return errors.E("already unlocked")
	}
	if l.seq != seq {
		return errors.E("not current lock owner")
	}
	l.locked = false
	return nil
}

func (l *lock) expire(ttl time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.lastUpdate) > ttl {
		l.locked = false
		l.seq++
	}
}

func (s *Server) waitForLock(repoName string) (lockName string, seq int, err error) {
	// TODO: implement using lock server
	cfg, ok := s.configuration[repoName]
	if !ok {
		return "", -1, errors.Errorf("invalid repository: %s", repoName)
	}
	for {
		for i := 0; i < int(cfg.Concurrency); i++ {
			lockName := fmt.Sprintf("%v-%v", repoName, i)
			v, ok := s.locks.Load(lockName)
			if !ok {
				return "", -1, errors.Errorf("cannot find lock for %s", lockName)
			}
			l := v.(*lock)
			seq, err := l.tryLock()
			if err != nil {
				continue
			}
			return lockName, seq, nil
		}
		s.lockRefresh.Wait()
	}

}

func (s *Server) expireLocks() {
	const ttl = 5 * time.Minute
	t := time.Tick(time.Second)
	for range t {
		s.locks.Range(func(k, v interface{}) bool {
			l := v.(*lock)
			l.expire(ttl)
			return true
		})
	}
}

func (s *Server) refreshLock(lockName string, seq int) error {
	v, ok := s.locks.Load(lockName)
	if !ok {
		return errors.Errorf("cannot find lock for %s", lockName)
	}
	l := v.(*lock)
	if l.seq != seq {
		return errors.E("not current lock owner")
	}
	l.refresh(seq)
	return nil
}
