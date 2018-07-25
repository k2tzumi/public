package worker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/git"
)

// Build consumes coordinator jobs.
func Build(buildsDir string, c *coordinator.Coordinator) {
	for {
		job := c.Next()
		if job == nil {
			log.Println("no more jobs in the pipe, halting worker")
			return
		}
		log.Println("checking out code...")
		dir, name := filepath.Split(job.RepoFullName)
		repoDir := filepath.Join(buildsDir, dir, name)
		if err := git.Checkout(job.Recipe.Clone, repoDir, job.CommitHash); err != nil {
			log.Println("cannot checkout code:", err)
			continue
		}
		log.Println("building...")
		output, err := run(context.Background(), job.Recipe, repoDir)
		log.Println("building result:", err)
		msg := fmt.Sprintln("build", job.ID, "for", job.RepoFullName,
			"commit:`", job.CommitMessage, "`",
			"("+job.CommitHash+")", "done")
		if err != nil {
			msg = fmt.Sprint("-  errored with:", err)
		}
		slackMessages := []string{msg}
		slackMessages = append(slackMessages, splitMsg(output, "```")...)
		for _, msg := range slackMessages {
			if err := slackSend(job.SlackWebhook, msg); err != nil {
				log.Println("cannot send slack message:", err)
			}
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
	return msgs
}

func slackSend(webhookURL string, msg string) error {
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(struct {
		Text string `json:"text"`
	}{Text: msg})
	if err != nil {
		return errors.E(err, "cannot encode slack message")
	}
	response, err := http.Post(webhookURL, "application/json", &payload)
	if err != nil {
		return errors.E(err, "cannot send slack message")
	}
	if _, err := io.Copy(ioutil.Discard, response.Body); err != nil {
		return errors.E(err, "cannot drain response body")
	}
	return nil
}
