package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"cirello.io/errors"
	"cirello.io/exp/grpcproxy/pkg/internal/proto"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
)

// Proxy takes packets from the hub and deliver to some remote host.
type Proxy struct {
	hub proto.HubClient

	conns sync.Map // map of connID to net.Conn
	dials singleflight.Group
}

// New creates a Proxy.
func New(cc *grpc.ClientConn) *Proxy {
	proxy := &Proxy{
		hub: proto.NewHubClient(cc),
	}
	return proxy
}

// Serve dials to target hosts and forward the packets hub sends.
func (p *Proxy) Serve() error {
	cl, err := p.hub.Proxy(context.Background())
	if err != nil {
		return errors.E(err, "cannot talk to hub")
	}
	for {
		packet, err := cl.Recv()
		if err != nil {
			return errors.E(err, "cannot receive packets from hub")
		}
		conn, err := p.loadConn(packet.GetConnID(), packet.GetTargetAddress(), cl.Send)
		if err != nil {
			log.Println("cannot dial to target", err)
			cl.Send(&proto.Packet{
				ConnID: packet.GetConnID(),
				State:  proto.Packet_Closed,
			})
			p.conns.Delete(packet.GetConnID())
			continue
		}
		if packet.State == proto.Packet_Handshake {
			continue
		} else if packet.State == proto.Packet_Closed {
			p.conns.Delete(packet.GetConnID())
			conn.Close()
			continue
		}
		conn.Write(packet.GetBody())
	}
}

func (p *Proxy) loadConn(connID int64, targetAddr string, send func(*proto.Packet) error) (net.Conn, error) {
	c, ok := p.conns.Load(connID)
	if ok {
		return c.(net.Conn), nil
	}
	c, err, _ := p.dials.Do(fmt.Sprint(connID),
		func() (interface{}, error) {
			conn, err := net.Dial("tcp", targetAddr)
			if err != nil {
				return nil, err
			}
			go func() {
				pt := &proxyTranslator{
					connID: connID,
					send:   send,
				}
				_, err := io.Copy(pt, conn)
				if _, ok := p.conns.Load(connID); !ok {
					return
				}
				p.conns.Delete(connID)
				var errBody []byte
				if err != nil {
					errBody = []byte(err.Error())
				}
				send(&proto.Packet{
					ConnID: connID,
					State:  proto.Packet_Closed,
					Body:   errBody,
				})
				conn.Close()
			}()
			p.conns.Store(connID, conn)
			return conn, nil
		})
	if err != nil {
		return nil, err
	}
	return c.(net.Conn), nil
}

type proxyTranslator struct {
	connID int64
	send   func(*proto.Packet) error
}

func (pt *proxyTranslator) Write(p []byte) (int, error) {
	err := pt.send(&proto.Packet{
		ConnID: pt.connID,
		Body:   p,
	})
	l := 0
	if err == nil {
		l = len(p)
	}
	return l, err
}
