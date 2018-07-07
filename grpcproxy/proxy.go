package main

import (
	"context"
	fmt "fmt"
	io "io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func proxy(hubAddr string) {
	cc, err := grpc.Dial(hubAddr, grpc.WithInsecure()) // for testing purpose, no security.
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	c := NewProxyClient(cc)
	proxyClient, err := c.Proxy(context.Background())
	if err != nil {
		log.Fatalf("cannot run proxy call: %v", err)
	}

	var conns sync.Map
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			recv, err := proxyClient.Recv()
			if err != nil {
				log.Println("proxy.recv:", err)
			}

			lc, ok := conns.Load(recv.GetPeers().GetConnID())
			if !ok {
				fmt.Println("dialing", recv.GetPeers().GetRightHandAddr(), recv.GetPeers().GetConnID())
				c, err := net.Dial("tcp", recv.GetPeers().GetRightHandAddr())
				if err != nil {
					log.Println("proxy.dial:", err)
				}
				conns.Store(recv.GetPeers().GetConnID(), c)
				lc, _ = conns.Load(recv.GetPeers().GetConnID())
				go func() {
					fwd := &packetBackSender{
						send:  proxyClient.Send,
						peers: recv.GetPeers(),
					}
					io.Copy(fwd, c)
				}()
			}
			fmt.Println("writing", recv.GetPeers().GetConnID(), len(recv.GetBody()))
			conn := lc.(net.Conn)
			conn.Write(recv.GetBody())
		}
	}()
	wg.Wait()
	resp, err := proxyClient.Recv()
	fmt.Println("client got from server:", resp, err)
}

type packetBackSender struct {
	send  func(*Packet) error
	peers *Peers
}

func (b *packetBackSender) Write(p []byte) (int, error) {
	err := b.send(&Packet{
		Peers: b.peers,
		Body:  p,
	})
	return len(p), err
}
