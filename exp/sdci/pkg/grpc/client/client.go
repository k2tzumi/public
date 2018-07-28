// Package client implements the client-side GRPC interface of the coordinator
// and workers.
package client // import "cirello.io/exp/sdci/pkg/grpc/client"

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/grpc/api"
	"cirello.io/exp/sdci/pkg/infra/git"
	"cirello.io/exp/sdci/pkg/infra/slack"
	"google.golang.org/grpc"
)

// Client executes build requests from server.
type Client struct {
	runner api.RunnerClient
}

// New instantiates a new client.
func New(cc *grpc.ClientConn) *Client {
	return &Client{
		runner: api.NewRunnerClient(cc),
	}
}

// Run listens for build request made by the server.
func (c *Client) Run(ctx context.Context, buildsDir, repoFullName string) error {
	cl, err := c.runner.Run(ctx)
	if err != nil {
		return errors.E(err, "cannot dial to server")
	}

	err = cl.Send(&api.JobRequest{
		Command: &api.JobRequest_Build{
			Build: &api.BuildRequest{
				RepoFullName: repoFullName,
			},
		},
	})
	if err != nil {
		return errors.E(err, "cannot handshake with the server")
	}

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				err = cl.Send(&api.JobRequest{
					Command: &api.JobRequest_KeepAlive{
						KeepAlive: &api.KeepAlive{},
					}})
				if err != nil {
					err = errors.E(err, "cannot handshake with the server")
					log.Println(err)
					return
				}
			}
		}
	}()

	for {
		job, err := cl.Recv()
		if err != nil {
			log.Println("I died:", err)
			return errors.E(err, "cannot receive remote command")
		}
		c.build(ctx, cl, buildsDir, job)
	}
}

func (c *Client) markInProgress(cl api.Runner_RunClient, job *api.JobResponse) error {
	err := cl.Send(&api.JobRequest{
		Command: &api.JobRequest_MarkInProgress{
			MarkInProgress: job.Build,
		},
	})
	return err
}

func (c *Client) markComplete(cl api.Runner_RunClient, job *api.JobResponse) error {
	err := cl.Send(&api.JobRequest{
		Command: &api.JobRequest_MarkComplete{
			MarkComplete: job.Build,
		},
	})
	return err
}

func (c *Client) build(ctx context.Context, cl api.Runner_RunClient, buildsDir string, job *api.JobResponse) {
	repoFullName := job.Build.RepoFullName
	commitHash := job.Build.CommitHash
	commitMessage := job.Build.CommitMessage
	if err := c.markInProgress(cl, job); err != nil {
		log.Println(repoFullName, "cannot mark job as in-progress:", err)
		return
	}
	defer func() {
		if err := c.markComplete(cl, job); err != nil {
			log.Println(repoFullName, "cannot mark job as completed:", err)
		}
	}()
	log.Println(repoFullName, "checking out code...")
	dir, name := filepath.Split(repoFullName)
	baseDir := filepath.Join(buildsDir, repoFullName)
	repoDir := filepath.Join(baseDir, "src", "github.com", dir, name)
	if err := git.Checkout(ctx, job.Recipe.Clone, repoDir, commitHash); err != nil {
		log.Println(repoFullName, "cannot checkout code:", err)
		return
	}
	log.Println(repoFullName, "building...")
	output, err := run(ctx, job.Recipe, repoDir, baseDir)
	job.Build.Success = err == nil
	job.Build.Log = output
	log.Println(repoFullName, "building result:", err)
	msg := fmt.Sprintln("build", job.Build.ID, "for", repoFullName,
		"commit:`", commitMessage, "`",
		"("+commitHash+")", "done")
	if err != nil {
		msg = fmt.Sprint("-  errored with:", err)
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
