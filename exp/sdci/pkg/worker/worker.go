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
	"os"
	"path/filepath"
	"strings"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/git"
	"cirello.io/exp/sdci/pkg/models"
)

// Start the builders.
func Start(ctx context.Context, buildsDir string, c *coordinator.Coordinator, configuration models.Configuration) error {
	for repoFullName, recipe := range configuration {
		total := recipe.Concurrency
		if total == 0 {
			total = 1
		}
		for i := 0; i < total; i++ {
			log.Println("starting worker for", repoFullName, i)
			buildsDir := fmt.Sprintf(buildsDir, i)
			if err := os.MkdirAll(buildsDir,
				os.ModePerm&0700); err != nil {
				return errors.E(err, "cannot create .sdci build directory")
			}
			go worker(ctx, buildsDir, repoFullName, c, i)
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
	if err := c.MarkInProgress(job); err != nil {
		log.Println(job.RepoFullName, "cannot mark job as in-progress:", err)
		return
	}
	defer func() {
		if err := c.MarkComplete(job); err != nil {
			log.Println(job.RepoFullName, "cannot mark job as completed:", err)
		}
	}()
	log.Println(job.RepoFullName, "checking out code...")
	dir, name := filepath.Split(job.RepoFullName)
	baseDir := filepath.Join(buildsDir, job.RepoFullName)
	repoDir := filepath.Join(baseDir, "src", "github.com", dir, name)
	if err := git.Checkout(ctx, job.Recipe.Clone, repoDir, job.CommitHash); err != nil {
		log.Println(job.RepoFullName, "cannot checkout code:", err)
		return
	}
	log.Println(job.RepoFullName, "building...")
	output, err := run(context.Background(), job.Recipe, repoDir, baseDir)
	job.Success = err == nil
	job.Log = output
	log.Println(job.RepoFullName, "building result:", err)
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
			log.Println(job.RepoFullName, "cannot send slack message:", err)
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
