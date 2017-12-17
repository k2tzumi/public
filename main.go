// Copyright 2017 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Command runner is a very ugly and simple structured command executer that
monitor file changes to trigger process restarts.

Create a file name Procfile in the root of the project you want to run, and add
the following content:

	workdir: $GOPATH/src/github.com/example/go-app
	observe: *.go *.js
	ignore: /vendor
	build-server: make server
	web: restart=always waitfor=localhost:8888 ./server serve

Special process types:

- workdir: the working directory. Environment variables are expanded. It follows
the same rules for exec.Command.Dir.

- observe: a space separated list of file patterns to scan for. It uses
filepath.Match internally.

- ignore: a space separated list of ignored directories relative to workdir,
typically vendor directories.

- build*: process type name prefixed by "build" are always executed first and in
order of declaration. On failure, they halt the initialization.

- waitfor (in process type): target hostname and port that the runner will probe
before starting the process type.

- restart (in process type): "always" will restart the process type every time;
"fail" will restart the process type on failure.
*/
package main // import "cirello.io/runner"

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"cirello.io/runner/procfile"
	"cirello.io/runner/runner"
)

// DefaultProcfile is the file that runner will open by default if no custom
// is given.
const DefaultProcfile = "Procfile"

var (
	convertToJSON = flag.Bool("convert", false, "takes a declared Procfile and prints as JSON to standard output")
	basePort      = flag.Int("port", 5000, "IP port used to set $`PORT` for each process type")
	formation     = flag.String("formation", "", "formation allows to start more than one instance of a process type, format: `procTypeA=# procTypeB=# ... procTypeN=#`")
	envFn         = flag.String("env", ".env", "environment `file` to be loaded for all processes.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "runner - simple Procfile runner\n\n")
		fmt.Fprintf(os.Stderr, "usage: %s [-convert] [Procfile]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("runner: ")

	fn := DefaultProcfile
	if argFn := flag.Arg(0); argFn != "" {
		fn = argFn
	}

	fd, err := os.Open(fn)
	if err != nil {
		log.Fatalln(err)
	}

	if *basePort < 1 || *basePort > 65535 {
		log.Fatalln("invalid IP port")
	}

	var s runner.Runner

	switch filepath.Ext(fn) {
	case ".json":
		if err := json.NewDecoder(fd).Decode(&s); err != nil {
			log.Fatalln("cannot parse spec file (json):", err)
		}
	default:
		s, err = procfile.Parse(fd)
		if err != nil {
			log.Fatalln("cannot parse spec file (procfile):", err)
		}
	}

	if *convertToJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		if err := enc.Encode(&s); err != nil {
			log.Fatalln("cannot encode procfile into JSON:", err)
		}
		return
	}

	s.WorkDir = os.ExpandEnv(s.WorkDir)
	if s.WorkDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln("cannot load current workdir", err)
		}
		s.WorkDir = wd
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-c
		log.Println("shutting down")
		cancel()
	}()

	s.BasePort = *basePort

	if fd, err := os.Open(*envFn); err == nil {
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := strings.Split(strings.TrimSpace(scanner.Text()), "=")
			if len(line) != 2 {
				continue
			}

			s.BaseEnvironment = append(s.BaseEnvironment, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading environment file (%v): %v", *envFn, err)
		}
	}

	if err := s.Start(ctx); err != nil {
		log.Fatalln("cannot serve:", err)
	}
}
