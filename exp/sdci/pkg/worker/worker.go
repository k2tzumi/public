package worker

import (
	"context"
	"log"
	"path/filepath"

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
		if err := git.Checkout(job.Recipe.Clone, repoDir, job.Commit); err != nil {
			log.Println("cannot checkout code:", err)
			continue
		}
		log.Println("building...")
		err := run(context.Background(), job.Recipe, repoDir)
		log.Println("building result:", err)
	}
}
