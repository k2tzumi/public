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
