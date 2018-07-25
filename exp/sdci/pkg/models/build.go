package models

import (
	"time"

	"cirello.io/errors"
	"github.com/jmoiron/sqlx"
)

// Build defines the necessary data to run a build successfully.
type Build struct {
	ID            int64     `db:"id"`
	RepoFullName  string    `db:"repo_full_name"`
	CommitHash    string    `db:"commit_hash"`
	CommitMessage string    `db:"commit_message"`
	InProgress    bool      `db:"in_progress"`
	StartedAt     time.Time `db:"started_at"`
	Success       bool      `db:"success"`
	Log           string    `db:"log"`
	CompletedAt   time.Time `db:"completed_at"`
	*Recipe
}

// BuildDAO provides access to the collection of Builds.
type BuildDAO struct {
	db *sqlx.DB
}

// NewBuildDAO creates a new Build data access object.
func NewBuildDAO(db *sqlx.DB) *BuildDAO {
	return &BuildDAO{
		db: db,
	}
}

// Bootstrap creates the necessary table to operate builds.
func (b *BuildDAO) Bootstrap() error {
	_, err := b.db.Exec(`
		CREATE TABLE IF NOT EXISTS builds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_full_name text,
			commit_hash text,
			commit_message text,
			environment text,
			commands text,
			started_at bigint,
			in_progress bool,
			success bool,
			log bigtext,
			completed_at bigint
		);
	`)
	return errors.E(err, "cannot bootstrap database")
}

// Register stores an new build in the database.
func (b *BuildDAO) Register(build *Build) (*Build, error) {
	res, err := b.db.NamedExec(`
		INSERT INTO builds
		(repo_full_name, commit_hash, commit_message, environment, commands)
		VALUES
		(:repo_full_name, :commit_hash, :commit_message, :environment, :commands)
	`, build)
	if err != nil {
		return build, errors.E(err, "cannot add job to database")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return build, errors.E(err, "cannot load ID from the added job")
	}
	build.ID = id
	return build, nil
}

// MarkInProgress determines a build has started and update its build
// information in the database.
func (b *BuildDAO) MarkInProgress(build *Build) error {
	build.StartedAt = time.Now()
	_, err := b.db.NamedExec(`
		UPDATE builds
		SET in_progress = true, started_at = :started_at
		WHERE id = :id
	`, build)
	return errors.E(err, "cannot mark build to in-progress")
}

// MarkComplete determines a build has completed and update its build
// information in the database
func (b *BuildDAO) MarkComplete(build *Build) error {
	build.CompletedAt = time.Now()
	build.InProgress = false
	_, err := b.db.NamedExec(`
		UPDATE builds
		SET
			in_progress = false,
			success = :success,
			log = :log,
			completed_at = :completed_at
		WHERE id = :id
	`, build)
	return errors.E(err, "cannot mark build to complete")
}
