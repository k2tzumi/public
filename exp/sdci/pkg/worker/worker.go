package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"

	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/git"
	"github.com/nlopes/slack"
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
		api := slack.New(job.Recipe.Slack)
		channels, err := api.GetChannels(false)
		if err != nil {
			log.Println("cannot load slack channels:", err)
			continue
		}
		for _, c := range channels {
			if c.Name != job.Recipe.Channel {
				continue
			}
			msg := fmt.Sprintln("build for", job.RepoFullName,
				"commit:", job.CommitMessage,
				"("+job.CommitHash+")", "done")
			if err != nil {
				msg = fmt.Sprintln(msg, "errored with", err)
			}
			_, _, err := api.PostMessage(c.ID, msg,
				slack.PostMessageParameters{
					Username: "sdci",
				})
			if err != nil {
				log.Println("cannot send slack message:", err)
				continue
			}

			chunks := splitmsg(output)
			for _, chunk := range chunks {
				_, _, err := api.PostMessage(c.ID, "```"+chunk+"```",
					slack.PostMessageParameters{
						Username: "sdci",
						Markdown: true,
					})
				if err != nil {
					log.Println("cannot send slack message:", err)
					continue
				}
			}
			break
		}
	}
}

func splitmsg(s string) []string {
	sub := ""
	subs := []string{}
	const maxMessageSize = 3900
	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%maxMessageSize == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}
