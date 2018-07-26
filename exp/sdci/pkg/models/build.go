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
	StartedAt     time.Time `db:"started_at"`
	Success       bool      `db:"success"`
	Log           string    `db:"log"`
	CompletedAt   time.Time `db:"completed_at"`
	*Recipe
}

// Status define the current build status of a Build
type Status int

// Possible build status
const (
	Unknown Status = iota
	Success
	Failed
	InProgress
)

// Status returns the current status of the build.
func (b *Build) Status() Status {
	switch {
	case b.StartedAt.IsZero():
		return Unknown
	case b.CompletedAt.IsZero():
		return InProgress
	case b.Success:
		return Success
	default:
		return Failed
	}
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
	ops := []string{
		`CREATE TABLE IF NOT EXISTS builds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_full_name text,
			commit_hash text,
			commit_message text,
			environment text,
			commands text,
			started_at datetime default (datetime('now')) not null,
			success bool default false not null,
			log bigtext default '' not null,
			completed_at datetime default '' not null
		);`,
		`CREATE INDEX builds_repo_full_nam ON builds (repo_full_name)`,
		`CREATE INDEX builds_started_at ON builds (started_at)`,
	}
	for _, op := range ops {
		_, err := b.db.Exec(op)
		if err != nil {
			errors.E(err, "cannot bootstrap database")
		}
	}
	return nil
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
		SET started_at = :started_at
		WHERE id = :id
	`, build)
	return errors.E(err, "cannot mark build to in-progress")
}

// MarkComplete determines a build has completed and update its build
// information in the database.
func (b *BuildDAO) MarkComplete(build *Build) error {
	build.CompletedAt = time.Now()
	_, err := b.db.NamedExec(`
		UPDATE builds
		SET
			success = :success,
			log = :log,
			completed_at = :completed_at
		WHERE id = :id
	`, build)
	return errors.E(err, "cannot mark build to complete")
}

// GetLastBuild loads last known build for a repository.
func (b *BuildDAO) GetLastBuild(repoFullName string) (*Build, error) {
	var build Build
	err := b.db.Get(&build, `
		SELECT
			id,
			repo_full_name,
			commit_hash,
			commit_message,
			environment,
			commands,
			started_at,
			success,
			log,
			completed_at
		FROM builds
		WHERE repo_full_name = :repoFullName
		ORDER BY started_at DESC
		LIMIT 1
	`, repoFullName)
	return &build, errors.E(err, "cannot load last known build")
}
