package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"storage"
)

var (
	listenAddr = flag.String("listen", "localhost:8080", "backend listening address")
)

func init() {
	flag.Parse()
}

func main() {
	backend := NewBackend(storage.NewStorage())
	fmt.Println("backend listening address:", *listenAddr)
	bind(backend, *listenAddr)
}

func bind(s *Backend, port string) {
	rpc.Register(s)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatalln("error:", e)
	}
	http.Serve(l, nil)
}
