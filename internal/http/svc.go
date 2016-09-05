package http

import (
	"log"
	"net"
	"net/http"

	"cirello.io/bloomfilterd/internal/filter"
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

func New() *Service {
	return &Service{
		d: &daemon{
			filters: make(map[string]*filter.Bloomfilter),
		},
	}
}
