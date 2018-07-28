// Command sdci implements a simple and dirty CI service.
package main // import "cirello.io/exp/sdci/cmd/sdci"
import (
	"log"
	"os"

	"cirello.io/exp/sdci/pkg/ui/cli"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetPrefix("sdci: ")
	log.SetFlags(0)
	fn := "sdci.db"
	if envFn := os.Getenv("SDCI_DB"); envFn != "" {
		fn = envFn
	}
	db := openDB(fn)
	cli.Run(db)
}

func openDB(fn string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", fn)
	if err != nil {
		log.Fatalln("cannot open database:", err)
	}
	return db
}
