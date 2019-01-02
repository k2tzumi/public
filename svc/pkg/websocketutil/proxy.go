// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
//
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

// Package websocketutil implements calls to assist handling with Websocket
// connections.
package websocketutil

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// IsWebsocketRequest detects if the HTTP request has the websocket upgrade
// header.
func IsWebsocketRequest(req *http.Request) bool {
	connHeader := ""
	connHeaders := req.Header["Connection"]
	if len(connHeaders) > 0 {
		connHeader = connHeaders[0]
	}

	upgradeWebsocket := false
	if strings.ToLower(connHeader) == "upgrade" {
		upgradeHeaders := req.Header["Upgrade"]
		if len(upgradeHeaders) > 0 {
			upgradeWebsocket = (strings.ToLower(upgradeHeaders[0]) == "websocket")
		}
	}

	return upgradeWebsocket
}

// Proxy returns a http.Handler capable of forwarding Websocket connections.
func Proxy(target string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := net.Dial("tcp", target)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("error dialing websocket backend %s: %v", target, err)
			return
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			log.Printf("not a hijacker: %T", w)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		nc, _, err := hj.Hijack()
		if err != nil {
			log.Printf("hijack error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer nc.Close()
		defer d.Close()

		err = r.Write(d)
		if err != nil {
			log.Printf("error copying request to target: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(d, nc)
		go cp(nc, d)
		<-errc
	})
}
