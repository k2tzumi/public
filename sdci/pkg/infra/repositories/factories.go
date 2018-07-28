// Package repositories implements the models repositories of the application,
// it derives the underlying implementation according to the database driver
// injected in.
package repositories // import "cirello.io/exp/sdci/pkg/infra/repositories"

import (
	sqlite3 "cirello.io/exp/sdci/pkg/infra/repositories/internal/sqlite3"
	"cirello.io/exp/sdci/pkg/models"
	"github.com/jmoiron/sqlx"
	sqlite3Driver "github.com/mattn/go-sqlite3"
)

// Builds creates a Builds repository.
func Builds(db *sqlx.DB) models.BuildRepository {
	switch db.Driver().(type) {
	case *sqlite3Driver.SQLiteDriver:
		return sqlite3.NewBuildDAO(db)
	default:
		panic("invalid DB driver")
	}
}
