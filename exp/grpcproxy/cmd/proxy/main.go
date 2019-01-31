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
