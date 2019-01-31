// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hub

import (
	"net/http"
	"sync"

	"cirello.io/exp/grpcproxy/pkg/internal/proto"
	"google.golang.org/grpc"
)

// Hub concentrates traffic from gateway and proxy.
type Hub struct {
	svc *hub
}

// New creates a new Hub.
func New() *Hub {
	h := &Hub{svc: &hub{
		proxyOutbound: make(chan *proto.Packet),
	}}
	return h
}

// ServeHTTP serves HTTP/2.0 requests of hub.
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv := grpc.NewServer()
	proto.RegisterHubServer(srv, h.svc)
	srv.ServeHTTP(w, r)
}

type hub struct {
	proxyInbound  sync.Map // map of connID to chan *proto.Packet
	proxyOutbound chan *proto.Packet

	mu     sync.Mutex
	connID int64
}

func (h *hub) nextConnID() int64 {
	h.mu.Lock()
	h.connID++
	h.mu.Unlock()
	return h.connID
}

func (h *hub) loadProxyInboundChannel(connID int64) chan *proto.Packet {
	ch, ok := h.proxyInbound.Load(connID)
	if !ok {
		panic("bug found. proxy inbound channels should be declared ahead of time")
	}
	return ch.(chan *proto.Packet)
}

func (h *hub) Gateway(srv proto.Hub_GatewayServer) error {
	connID := h.nextConnID()
	packetInboundChannel := make(chan *proto.Packet)
	h.proxyInbound.Store(connID, packetInboundChannel)
	errCh := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			p, err := srv.Recv()
			if err != nil {
				errCh <- err
				return
			}

			p.ConnID = connID
			h.proxyOutbound <- p
		}
	}()
	go func() {
		for p := range packetInboundChannel {
			if err := srv.Send(p); err != nil {
				errCh <- err
				return
			}
		}
	}()
	wg.Wait()
	close(errCh)
	return <-errCh
}

func (h *hub) Proxy(srv proto.Hub_ProxyServer) error {
	errCh := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			p, err := srv.Recv()
			if err != nil {
				errCh <- err
				return
			}
			h.loadProxyInboundChannel(p.GetConnID()) <- p
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for p := range h.proxyOutbound {
			err := srv.Send(p)
			if err != nil {
				errCh <- err
				return
			}
		}
	}()
	wg.Wait()
	close(errCh)
	return <-errCh
}
