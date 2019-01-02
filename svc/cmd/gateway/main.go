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

// Command gateway implements an edge reverse-proxy used for serving Go packages and static content.
package main

import (
	"os"
	"sync"
)

const gatewayTokenCookie = "gateway-jwt"

var (
	publicBindIP   = os.Getenv("GATEWAY_PUBLIC_BIND_IP")
	servicesBindIP = os.Getenv("GATEWAY_SERVICES_BIND_IP")
	googleClientID = os.Getenv("GATEWAY_GOOGLE_CLIENT_ID")
	links          = os.Getenv("GATEWAY_LINKS_GIST")
	baseGithubAcct = os.Getenv("GATEWAY_BASE_GITHUB_ACCOUNT")
	frontPkgDomain = os.Getenv("GATEWAY_FRONT_PACKAGE_DOMAIN")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		services()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		publicSites()
	}()
	wg.Wait()
}
