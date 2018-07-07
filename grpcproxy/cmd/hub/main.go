package main

import (
	"log"
	"net/http"

	"cirello.io/exp/grpcproxy/pkg/hub"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("hub: ")
	hub := hub.New()
	server := &http.Server{
		Addr:    ":8080",
		Handler: hub,
	}
	log.Println(server.ListenAndServeTLS("fake-server.crt", "fake-server.key"))
}
