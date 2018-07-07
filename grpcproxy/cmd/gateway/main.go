package main

import (
	"crypto/tls"
	"log"
	"net"

	"cirello.io/exp/grpcproxy/pkg/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("gateway: ")
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	cc, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.Listen("tcp", ":55459")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("serving at ", l.Addr())
	gw := gateway.New(cc, "127.0.0.1:9999")
	log.Println(gw.Serve(l))
}
