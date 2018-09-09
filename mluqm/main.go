// Copyright 2018 github.com/ucirello and https://cirello.io. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// Command mluqm runs a AI-powered Ur-Quan Masters compatible client.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	log.SetPrefix("mluqm: ")
	log.SetFlags(0)
	l, err := net.Listen("tcp4", ":21837")
	check(err)
	log.Println("listening...")
	for {
		log.Println("waiting for incoming connection")
		inboundConn, err := l.Accept()
		check(err)
		outboundConn, err := net.Dial("tcp4", "localhost:21838")
		check(err)
		log.Println("conn accepted, conjoining")

		inwardsDumpR, inwardsDumpW := io.Pipe()
		go func() {
			for {
				out, err := parsePackets("in", inwardsDumpR)
				fmt.Println(out)
				if err != nil {
					break
				}
			}
		}()
		teeInwards := io.TeeReader(outboundConn, inwardsDumpW)
		go io.Copy(inboundConn, teeInwards)

		outwardsDumpR, outwardsDumpW := io.Pipe()
		go func() {
			for {
				out, err := parsePackets("out", outwardsDumpR)
				fmt.Println(out)
				if err != nil {
					break
				}
			}
		}()
		teeOutwards := io.TeeReader(inboundConn, outwardsDumpW)
		io.Copy(outboundConn, teeOutwards)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
