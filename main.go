package main // import "cirello.io/bloomfilterd"

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"cirello.io/bloomfilterd/internal/filter"
	"cirello.io/suture"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var supervisor suture.Supervisor
	dsvc := &daemonSvc{
		d: &daemon{
			filters: make(map[string]*filter.Bloomfilter),
		},
	}
	supervisor.Add(dsvc)
	supervisor.ServeBackground()

	<-c
	fmt.Println("stopping...")
	supervisor.Stop()
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
