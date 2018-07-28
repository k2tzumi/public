// Package cli implements the command line interface primitives used by
// cmd/sdci.
package cli // import "cirello.io/exp/sdci/pkg/ui/cli"

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	cli "gopkg.in/urfave/cli.v1"
)

// Run executes the application in CLI mode
func Run(db *sqlx.DB) {
	app := cli.NewApp()
	app.Name = "bookmarkd"
	app.Usage = "bookmark manager"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		standaloneCmd(db),
	}
	sort.Slice(app.Commands, func(i, j int) bool {
		return strings.Compare(app.Commands[i].Name, app.Commands[j].Name) < 0
	})
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func standaloneCmd(db *sqlx.DB) cli.Command {
	return cli.Command{
		Name:        "standalone",
		Description: "start service in standalone mode",
		Action: func(ctx *cli.Context) error {
			return standalone(db)
		},
	}
}
