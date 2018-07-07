package main

import (
	"crypto/tls"
	"log"

	"cirello.io/exp/grpcproxy/pkg/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("proxy: ")
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	cc, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatal(err)
	}
	proxy := proxy.New(cc)
	log.Println(proxy.Serve())
}
