// Package worker implements the build worker.
package worker // import "cirello.io/exp/sdci/pkg/worker"

import (
	"context"
	"fmt"
	"log"
	"os"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/grpc/api"
	"cirello.io/exp/sdci/pkg/grpc/client"
	"google.golang.org/grpc"
)

// Start the builders.
func Start(ctx context.Context, grpcServerAddr, buildsDir string) error {
	cc, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.E(err, "cannot dial to GRPC server")
	}
	cl := api.NewRunnerClient(cc)
	resp, err := cl.Configuration(context.Background(), &api.ConfigurationRequest{})
	if err != nil {
		return errors.E(err, "cannot load configuration")
	}
	for repoFullName, recipe := range resp.Configuration {
		total := int(recipe.Concurrency)
		for i := 0; i < total; i++ {
			buildsDir := fmt.Sprintf(buildsDir, i)
			if err := os.MkdirAll(buildsDir,
				os.ModePerm&0700); err != nil {
				return errors.E(err, "cannot create .sdci build directory")
			}
			go worker(ctx, cc, buildsDir, repoFullName, i)
		}
	}
	return nil
}

func worker(ctx context.Context, cc *grpc.ClientConn, buildsDir, repoFullName string, i int) {
	c := client.New(cc)
	log.Println("starting worker for", repoFullName, i)
	defer log.Println("done with ", repoFullName, i)
	err := c.Run(ctx, buildsDir, repoFullName)
	if err != nil {
		log.Println("cannot run worker", repoFullName, i, ":", err)
	}
}
