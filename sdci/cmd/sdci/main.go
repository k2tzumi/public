// Command sdci implements a simple and dirty CI service.
package main // import "cirello.io/exp/sdci/cmd/sdci"
import (
	"context"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/models"
	"cirello.io/exp/sdci/pkg/ui/dashboard"
	"cirello.io/exp/sdci/pkg/ui/webhooks"
	"cirello.io/exp/sdci/pkg/worker"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sdci: ")
	buildsDir := buildsDir()
	db := openDB()
	configuration := loadConfiguration()
	ctx, coord := startCoordinator(db, configuration)
	startWorkers(ctx, buildsDir, coord, configuration)
	startWebhooksServer(ctx, coord)
	startDashboard(ctx, db)
	if err := coord.Wait(); err != nil {
		log.Fatalln("coordinator error:", err)
	}
}

func buildsDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalln("cannot load current user information:", err)
	}
	return filepath.Join(currentUser.HomeDir, ".sdci", "builds-%v")
}

func openDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "sdci.db")
	if err != nil {
		log.Fatalln("cannot open database:", err)
	}
	return db
}

func loadConfiguration() models.Configuration {
	fd, err := os.Open("sdci-config.yaml")
	if err != nil {
		log.Fatal("cannot open configuration file:", err)
	}
	configuration, err := models.LoadConfiguration(fd)
	if err != nil {
		log.Fatal(err)
	}
	return configuration
}

func startCoordinator(db *sqlx.DB, configuration models.Configuration) (context.Context, *coordinator.Coordinator) {
	ctx, coord := coordinator.New(db, configuration)
	if err := coord.Error(); err != nil {
		log.Fatalln("coordinator error on start:", err)
	}
	return ctx, coord
}

func startWorkers(ctx context.Context, buildsDir string,
	coord *coordinator.Coordinator, configuration models.Configuration) {
	if err := worker.Start(ctx, buildsDir, coord, configuration); err != nil {
		log.Fatalln("coordinator error on start:", err)
	}
}

func startWebhooksServer(ctx context.Context, coord *coordinator.Coordinator) {
	webhookListener, err := net.Listen("tcp", ":6500")
	if err != nil {
		log.Fatalln("cannot start web server:", err)
	}
	webhookServer := webhooks.New(coord)
	go func() {
		log.Println(webhookServer.ServeContext(ctx, webhookListener))
	}()
}

func startDashboard(ctx context.Context, db *sqlx.DB) {
	dashboardListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("cannot start dashboard server:", err)
	}
	dashboardServer := dashboard.New(models.NewBuildDAO(db))
	if err := dashboardServer.ServeContext(ctx, dashboardListener); err != nil {
		log.Println("cannot server dashboard:", err)
	}
}
