package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"

	"client"

	"github.com/golang/groupcache"
)

var (
	listenAddr   = flag.String("listen", "http://localhost:8001", "groupcache listen address")
	frontendAddr = flag.String("frontend", "localhost:9001", "frontend listen address")
	backend      = flag.String("backend", "localhost:8080", "backend listen address")
	gcpeers      = flag.String("peers", "http://localhost:8001,http://localhost:8002,http://localhost:8003", "groupcache peers")
)

func init() {
	flag.Parse()
}

func main() {
	peers := groupcache.NewHTTPPool(*listenAddr)
	client := client.NewClient(*backend)

	var stringcache = groupcache.NewGroup(
		"BackendCache",
		64<<20,
		groupcache.GetterFunc(func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result := client.Get(key)
			dest.SetBytes([]byte(result))
			return nil
		}),
	)

	peers.Set(strings.Split(*gcpeers, ",")...)

	frontendServer := NewServer(stringcache)
	go bind(frontendServer, *frontendAddr)

	fmt.Println("cachegroup slave listening address:", *listenAddr)
	fmt.Println("frontend  listening address:", *frontendAddr)
	fmt.Println("peers pool:", strings.Split(*gcpeers, ","))
	log.Fatalln(http.ListenAndServe(strings.Replace(*listenAddr, "http://", "", 1), http.HandlerFunc(peers.ServeHTTP)))
}

func bind(s *Frontend, port string) {
	rpc.Register(s)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatalln(e)
	}
	http.Serve(l, nil)
}
