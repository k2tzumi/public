package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func startHub() (addr string, srv *hub) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	hub := &hub{
		outbound: make(chan *Packet),
	}
	grpcServer := grpc.NewServer()
	RegisterProxyServer(grpcServer, hub)
	fmt.Println("backend listening address:", "127.0.0.1:0")
	go grpcServer.Serve(l)
	return l.Addr().String(), srv
}

type hub struct {
	inbound  sync.Map // map of connID to chan *Packet
	outbound chan *Packet
}

func (s *hub) Proxy(srv Proxy_ProxyServer) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			p, err := srv.Recv()
			if err != nil {
				log.Println("srv.recv:", err)
				break
			}
			fmt.Println("recv", p, "piping")
			pipe, _ := s.inbound.Load(p.GetPeers().GetConnID())
			inboundPipe := pipe.(chan *Packet)
			inboundPipe <- p
			fmt.Println("recv", p, "piped")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for p := range s.outbound {
			if err := srv.Send(p); err != nil {
				log.Println("srv.send:", err)
			}
		}
	}()
	wg.Wait()
	return nil
}
