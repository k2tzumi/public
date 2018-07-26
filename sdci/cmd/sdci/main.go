// Command sdci implements a simple and dirty CI service.
package main // import "cirello.io/exp/sdci/cmd/sdci"
import (
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/models"
	"cirello.io/exp/sdci/pkg/web"
	"cirello.io/exp/sdci/pkg/worker"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
	// TODO: work out webserver and worker stop when coordinator fails.
	buildsDir := filepath.Join(currentUser.HomeDir, ".sdci", "builds-%v", "src", "github.com")
	fd, err := os.Open("sdci-config.yaml")
	if err != nil {
		log.Fatal("cannot open configuration file:", err)
	}
	configuration, err := models.LoadConfiguration(fd)
	if err != nil {
		log.Fatal(err)
	}
	coord := coordinator.New(db, configuration)
	if err := coord.Error(); err != nil {
		log.Println("coordinator error on start:", err)
	}
	defer func() {
		if err := coord.Error(); err != nil {
			log.Println("coordinator error:", err)
		}
	}()
	worker.Start(buildsDir, coord, configuration)
	l, err := net.Listen("tcp", ":6500")
	if err != nil {
		log.Fatalln("cannot start web server:", err)
	}
	s := web.New(coord)
	log.Fatalln(s.Serve(l))
}
