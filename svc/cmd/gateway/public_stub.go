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

// +build !prod

package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/alecthomas/template"
	"golang.org/x/crypto/acme/autocert"
)

func publicSites() {
	log.Println("bootstrapping sites")
	pkgRedirect := http.HandlerFunc(pkgRedirect)

	m := &autocert.Manager{
		Cache:      autocert.DirCache("./httpd-sites.secrets"),
		Prompt:     autocert.AcceptTOS,
		Email:      "user@example.com",
		HostPolicy: autocert.HostWhitelist(frontPkgDomain),
	}
	log.Println("starting sites:80")
	go func() {
		log.Println("sites:80",
			http.ListenAndServe(publicBindIP+":http",
				m.HTTPHandler(nil)))
	}()
	s := &http.Server{
		Addr: publicBindIP + ":https",
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Host {
			case frontPkgDomain:
				pkgRedirect.ServeHTTP(w, r)
			default:
				http.NotFound(w, r)
			}
		}),
	}
	log.Println("starting sites:443")
	log.Println("sites:443", s.ListenAndServeTLS("", ""))
}

func pkgRedirect(w http.ResponseWriter, r *http.Request) {
	pkgName := strings.TrimPrefix(r.URL.Path, "/")
	if strings.HasPrefix(pkgName, "s/") {
		target := strings.TrimPrefix(pkgName, "s/")
		resp, err := http.Get(links)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}
		defer resp.Body.Close()

		var l map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		if link, ok := l[target]; ok {
			http.Redirect(w, r, link, http.StatusSeeOther)
			return
		}

		http.NotFound(w, r)
		return
	}
	tmpl, err := template.New("pkgHTML").Parse(pkgHTML)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	rootPkgName := pkgName
	if i := strings.Index(rootPkgName, "/"); i >= 0 {
		rootPkgName = rootPkgName[:i]
	}

	if err := tmpl.Execute(w, struct {
		RootPackage        string
		Package            string
		BaseGithubAccount  string
		FrontPackageDomain string
	}{rootPkgName, pkgName, baseGithubAcct, frontPkgDomain}); err != nil {
		log.Println(err)
		http.NotFound(w, r)
	}
	return
}
