package main // import "cirello.io/bloomfilterd"

import (
	"log"
	"net"
	"net/http"

	"cirello.io/bloomfilterd/internal/filter"
	"cirello.io/suture"
)

func main() {
	var supervisor suture.Supervisor
	dsvc := &daemonSvc{
		d: &daemon{
			filters: make(map[string]*filter.Bloomfilter),
		},
	}
	supervisor.Add(dsvc)
	supervisor.Serve()
}

type daemonSvc struct {
	d *daemon
	l net.Listener
}

func (dsvc *daemonSvc) Serve() {
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Println(err)
		return
	}
	dsvc.l = l
	http.Handle("/", dsvc.d)
	log.Println(http.Serve(dsvc.l, nil))
}
func (dsvc *daemonSvc) Stop() {
	dsvc.l.Close()
}
