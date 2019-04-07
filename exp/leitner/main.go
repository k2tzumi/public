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

// Command leitner is a spaced repetition service based on Leitner's algorithm.
package main // import "cirello.io/exp/leitner"

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	_ "gopkg.in/urfave/cli.v1"
)

var repetitionScale = [...]int{
	1,
	2,
	4,
	8,
	16,
	32,
	64,
}

func daysSince(t time.Time) int {
	return int(time.Since(t).Hours()) / 24
}

type meme struct {
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	Level      int       `json:"level"`
	PromotedAt time.Time `json:"promoted_at"`
}

func (m meme) repeat() bool {
	d := daysSince(m.PromotedAt)
	return d%repetitionScale[m.Level] == 0
}

func main() {
	log.SetFlags(0)
	app := cli.NewApp()
	app.Action = exercise
	app.Commands = []cli.Command{promote()}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func promote() cli.Command {
	return cli.Command{
		Name: "promote",
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				return xerrors.New("missing args")
			}
			topics, err := load()
			if err != nil {
				return xerrors.Errorf("cannot load topics: %w", err)
			}
			id, direction := atoi(c.Args().First()), c.Args().Get(1)
			if id > len(topics) {
				return xerrors.New("invalid meme ID")
			}
			if direction != "up" && direction != "down" {
				return xerrors.New("invalid promotion direction")
			}
			meme := topics[id]
			if direction == "up" {
				meme.Level++
			} else if direction == "down" {
				meme.Level--
			}
			if meme.Level < 0 || meme.Level > len(repetitionScale)-1 {
				return xerrors.Errorf("cannot level %s", direction)
			}
			meme.PromotedAt = time.Now()
			topics[id] = meme
			return store(topics)
		},
	}
}

func exercise(c *cli.Context) error {
	topics, err := load()
	if err != nil {
		return xerrors.Errorf("cannot load topics: %w", err)
	}
	for i, meme := range topics {
		if !meme.repeat() {
			continue
		}
		fmt.Println("ID:", i)
		fmt.Println("Question:", meme.Question)
		fmt.Println("Answer:", meme.Answer)
	}
	return nil
}

func load() ([]meme, error) {
	var topics []meme
	fd, err := os.Open("topics.json")
	if err != nil {
		return topics, xerrors.Errorf("cannot open topics.json: %w", err)
	}
	err = json.NewDecoder(fd).Decode(&topics)
	return topics, err
}

func atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func store(topics []meme) error {
	fd, err := os.Create("topics.json")
	if err != nil {
		return xerrors.Errorf("cannot create topics.json: %w", err)
	}
	err = json.NewEncoder(fd).Encode(topics)
	return err
}
