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
	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sdci: ")
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalln("cannot load current user information:", err)
	}
	db, err := sqlx.Open("sqlite3", "sdci.db")
	if err != nil {
		log.Fatalln("cannot open database:", err)
	}
	// TODO: organize the relationship between coordinator, web and workers.
	coord := coordinator.New(db)
	defer func() {
		if err := coord.Error(); err != nil {
			log.Println("coordinator error:", err)
		}
	}()
	buildsDir := filepath.Join(currentUser.HomeDir, ".sdci", "builds-%v", "src", "github.com")
	worker.Start(buildsDir, coord, 1)
	recipes, err := loadRecipes()
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("cannot start web server:", err)
	}
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
