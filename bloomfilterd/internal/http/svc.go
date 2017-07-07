package http

import (
	"log"
	"net"
	"net/http"

	"cirello.io/bloomfilterd/internal/storage"

	"github.com/coreos/etcd/raft/raftpb"
)

type Service struct {
	listen string

	d *daemon
	l net.Listener
}

func (s *Service) Serve() {
	l, err := net.Listen("tcp", s.listen)
	if err != nil {
		log.Println(err)
		return
	}
	s.l = l
	http.Handle("/", s.d)
	log.Println(http.Serve(s.l, nil))
}

func (s *Service) Stop() {
	s.l.Close()
}

func New(propose chan string, confChange chan raftpb.ConfChange, opts ...Option) *Service {
	svc := &Service{
		d: &daemon{
			propose:    propose,
			confChange: confChange,
		},
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

type Option func(*Service)

func MemoryStorage(s *Service) {
	s.d.storage = storage.Must(storage.Memory)
}

func Listen(listen string) Option {
	return func(s *Service) {
		s.listen = listen
	}
}
