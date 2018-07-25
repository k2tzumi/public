// Command sdci implements a simple and dirty CI service.
package main // import "cirello.io/exp/sdci/cmd/sdci"
import (
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/web"
	"cirello.io/exp/sdci/pkg/worker"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
	yaml "gopkg.in/yaml.v2"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sdci: ")
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalln("cannot load current user information:", err)
	}
	buildsDir := filepath.Join(currentUser.HomeDir, ".sdci", "builds", "src", "github.com")
	if err := os.MkdirAll(buildsDir, os.ModePerm&0700); err != nil {
		log.Fatalln("cannot create .sdci:", err)
	}
	db, err := sqlx.Open("sqlite3", "sdci.db")
	if err != nil {
		log.Fatalln("cannot open database:", err)
	}
	coord := coordinator.New(db)
	go worker.Build(buildsDir, coord)
	recipes, err := loadRecipes()
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("cannot start web server:", err)
	}
	// TODO: handle coordinator disruptions.
	s := web.New(recipes, coord)
	log.Fatalln(s.Serve(l))
}

func loadRecipes() (map[string]*coordinator.Recipe, error) {
	fd, err := os.Open("sdci-config.yaml")
	if err != nil {
		return nil, errors.E(err, "cannot open configuration file")
	}
	var recipes map[string]*coordinator.Recipe
	err = yaml.NewDecoder(fd).Decode(&recipes)
	return recipes, errors.E(err, "cannot parse configuration")
}
