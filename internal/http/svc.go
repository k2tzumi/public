package http

import (
	"log"
	"net"
	"net/http"

	"cirello.io/bloomfilterd/internal/storage"
	"github.com/coreos/etcd/raft/raftpb"
)

type Service struct {
	d *daemon
	l net.Listener
}

func (dsvc *Service) Serve() {
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Println(err)
		return
	}
	dsvc.l = l
	http.Handle("/", dsvc.d)
	log.Println(http.Serve(dsvc.l, nil))
}

func (dsvc *Service) Stop() {
	dsvc.l.Close()
}

func New(propose chan string, confChange chan raftpb.ConfChange, t storage.Type) *Service {
	return &Service{
		d: &daemon{
			propose:    propose,
			confChange: confChange,
			storage:    storage.Must(t),
		},
	}
}
