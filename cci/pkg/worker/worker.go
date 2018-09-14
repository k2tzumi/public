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

// Package worker implements the build worker.
package worker // import "cirello.io/cci/pkg/worker"

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cirello.io/cci/pkg/coordinator"
	"cirello.io/cci/pkg/infra/git"
	"cirello.io/cci/pkg/infra/slack"
	"cirello.io/cci/pkg/models"
	"cirello.io/errors"
)

// Start the builders.
func Start(ctx context.Context, buildsDir string,
	configuration models.Configuration, coord *coordinator.Coordinator) error {
	for repoFullName, recipe := range configuration {
		total := int(recipe.Concurrency)
		for i := 0; i < total; i++ {
			buildsDir := fmt.Sprintf(buildsDir, i)
			if err := os.MkdirAll(buildsDir,
				os.ModePerm&0700); err != nil {
				return errors.E(err, "cannot create .cci build directory")
			}
			go worker(ctx, buildsDir, repoFullName, coord, i)
		}
	}
	return nil
}

func worker(ctx context.Context, buildsDir, repoFullName string, c *coordinator.Coordinator, i int) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-c.Next(repoFullName):
			if job == nil {
				log.Println("no more jobs in the pipe, halting worker")
				return
			}
			build(ctx, buildsDir, c, job)
		}
	}
}

func build(ctx context.Context, buildsDir string, c *coordinator.Coordinator, job *models.Build) {
	repoFullName := job.RepoFullName
	dir, name := filepath.Split(repoFullName)
	baseDir := filepath.Join(buildsDir, repoFullName)
	repoDir := filepath.Join(baseDir, "src", "github.com", dir, name)
	if err := c.MarkInProgress(job); err != nil {
		log.Println(repoFullName, "cannot mark job as in-progress:", err)
		return
	}
	defer func() {
		if err := c.MarkComplete(job); err != nil {
			log.Println(repoFullName, "cannot mark job as completed:", err)
		}
	}()
	slackStart(job)
	log.Println(repoFullName, "checking out code...")
	if err := git.Checkout(ctx, job.Recipe.Clone, repoDir, job.CommitHash); err != nil {
		log.Println(repoFullName, "cannot checkout code:", err)
		return
	}
	log.Println(repoFullName, "building...")
	output, err := run(ctx, job.Recipe, repoDir, baseDir)
	job.Success = err == nil
	job.Log = output
	log.Println(repoFullName, "building result:", err)
	slackEnd(job, output, err)
}

func slackStart(job *models.Build) {
	repoFullName := job.RepoFullName
	commitHash := job.CommitHash
	msg := fmt.Sprintln("build", job.ID, "for", repoFullName,
		"("+commitHash+")", "started")
	if err := slack.Send(job.Recipe.SlackWebhook, msg); err != nil {
		log.Println(repoFullName, "cannot send slack message:", err)
	}
}

func slackEnd(job *models.Build, output string, err error) {
	repoFullName := job.RepoFullName
	commitHash := job.CommitHash
	msg := fmt.Sprintln("build", job.ID, "for", repoFullName,
		"("+commitHash+")", "done")
	if err != nil {
		msg = fmt.Sprint(msg, "-  errored with:", err)
	}
	slackMessages := []string{msg}
	slackMessages = append(slackMessages, splitMsg(output, "```")...)
	for _, msg := range slackMessages {
		if err := slack.Send(job.Recipe.SlackWebhook, msg); err != nil {
			log.Println(repoFullName, "cannot send slack message:", err)
		}
	}
}

func splitMsg(msg, split string) []string {
	var msgs []string
	const maxsize = 2048
	current := 0
	r := strings.NewReader(msg)
	scanner := bufio.NewScanner(r)
	var buf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		current += len(line)
		fmt.Fprintln(&buf, line)
		if current > maxsize {
			msgs = append(msgs, split+"\n"+buf.String()+"\n"+split)
			buf.Reset()
			current = 0
		}
	}
	if str := buf.String(); str != "" {
		msgs = append(msgs, split+"\n"+buf.String()+"\n"+split)
	}
	return msgs
}
