package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func startGateway(hub *hub) {
	l, err := net.Listen("tcp", "127.0.0.1:8765")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	connID := int64(0)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("broker:", err)
		}
		connID++
		log.Println("connID", connID)
		b := &gateway{
			peers: &Peers{
				ConnID:        connID,
				RightHandAddr: "localhost:9999",
			},
			hub: hub,
		}
		go b.consume(conn)
		go io.Copy(b, conn)
	}
}

type gateway struct {
	peers *Peers
	hub   *hub
}

func (b *gateway) Write(p []byte) (int, error) {
	b.hub.outbound <- &Packet{
		Peers: b.peers,
		Body:  p,
	}
	return len(p), nil
}

func (b *gateway) consume(conn net.Conn) {
	inboundPipe := make(chan *Packet)
	b.hub.inbound.Store(b.peers.GetConnID(), inboundPipe)
	for packet := range inboundPipe {
		fmt.Println("gateway.consume", packet)
		conn.Write(packet.GetBody())
	}
}
