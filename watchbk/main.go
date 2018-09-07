package main

import (
	"log"
	"os"
	"time"

	buildkite "github.com/buildkite/go-buildkite/buildkite"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	log.SetPrefix("bk: ")
	log.SetFlags(0)
	app := cli.NewApp()
	app.Name = "bk"
	app.Usage = "operate Buildkite from the CLI"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			EnvVar: "BK_TOKEN",
		},
		cli.BoolFlag{
			Name:   "debug",
			Hidden: true,
		},
	}
	app.Action = func(c *cli.Context) error {
		config, err := buildkite.NewTokenConfig(c.GlobalString("token"), c.GlobalBool("debug"))
		if err != nil {
			log.Fatalf("client config failed: %s", err)
		}
		client := buildkite.NewClient(config.Client())
		for {
			builds, _, err := client.Builds.List(&buildkite.BuildsListOptions{
				CreatedFrom: time.Now().Add(-24 * time.Hour),
				State: []string{
					"running", "scheduled",
				},
			})
			if err != nil {
				log.Fatalf("list builds failed: %s", err)
			}
			if len(builds) == 0 {
				break
			}
			time.Sleep(time.Second)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
