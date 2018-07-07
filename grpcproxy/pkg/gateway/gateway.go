package gateway

import (
	"context"
	"io"
	"log"
	"net"
	"sync"

	"cirello.io/errors"
	"cirello.io/exp/grpcproxy/pkg/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Gateway connects to the hub and forward local packets.
type Gateway struct {
	hub           proto.HubClient
	targetAddress string
}

// New creates a Gateway.
func New(cc *grpc.ClientConn, targetAddress string) *Gateway {
	gw := &Gateway{
		hub:           proto.NewHubClient(cc),
		targetAddress: targetAddress,
	}
	return gw
}

// Serve forwards its accepted connections to the hub.
func (gw *Gateway) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return errors.E(err, "cannot serve gateway")
		}
		go func() {
			if err := gw.forward(conn, gw.targetAddress); err != nil {
				log.Println("cannot forward connection", err)
			}
		}()
	}
}

func (gw *Gateway) forward(conn net.Conn, targetAddress string) error {
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := gw.hub.Gateway(ctx)
	if err != nil {
		return errors.E(err, "cannot connect to hub")
	}
	gt := &gatewayTranslator{
		cl:            cl,
		targetAddress: targetAddress,
	}
	cl.Send(&proto.Packet{
		TargetAddress: targetAddress,
		State:         proto.Packet_Handshake,
	})
	go func() {
		defer cancel()
		io.Copy(gt, conn)
		if !gt.isClosed() {
			cl.Send(&proto.Packet{
				TargetAddress: targetAddress,
				State:         proto.Packet_Closed,
			})
		}
	}()
	gt.inboundConsumer(conn)
	return gt.err
}

type gatewayTranslator struct {
	cl            proto.Hub_GatewayClient
	targetAddress string
	err           error

	mu     sync.Mutex
	closed bool
}

func (gt *gatewayTranslator) isClosed() bool {
	gt.mu.Lock()
	closed := gt.closed
	gt.mu.Unlock()
	return closed
}

func (gt *gatewayTranslator) inboundConsumer(conn net.Conn) {
	for {
		p, err := gt.cl.Recv()
		if status.Code(err) == codes.Canceled {
			return
		} else if err != nil {
			gt.err = err
			return
		}
		conn.Write(p.GetBody())
		if p.State == proto.Packet_Closed {
			gt.mu.Lock()
			gt.closed = true
			gt.mu.Unlock()
			return
		}
	}
}

func (gt *gatewayTranslator) Write(p []byte) (int, error) {
	packet := &proto.Packet{
		TargetAddress: gt.targetAddress,
		Body:          p,
	}
	err := gt.cl.Send(packet)
	l := 0
	if err == nil {
		l = len(p)
	}
	return l, err
}
