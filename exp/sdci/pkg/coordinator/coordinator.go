package coordinator // import "cirello.io/exp/sdci/pkg/coordinator"
import (
	"time"

	"cirello.io/errors"
	"github.com/jmoiron/sqlx"
)

// Recipe defines the execution steps and environment.
type Recipe struct {
	Clone        string `db:"clone" yaml:"clone"`
	SlackWebhook string `db:"slack_webhook" yaml:"slack_webhook"`
	Environment  string `db:"environment" yaml:"environment"`
	Commands     string `db:"commands" yaml:"commands"`
}

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

// Coordinator takes and dispatches build requests.
type Coordinator struct {
	db  *sqlx.DB
	in  chan *Build
	out chan *Build
	err error
}

// New creates a new coordinator
func New(db *sqlx.DB) *Coordinator {
	c := &Coordinator{
		db: db,
		// TODO: replace with proper queues
		in:  make(chan *Build, 10),
		out: make(chan *Build),
	}
	c.bootstrap()
	go c.forward()
	return c
}

func (c *Coordinator) bootstrap() {
	_, err := c.db.Exec(`
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
	if err != nil {
		c.err = errors.E(err, "cannot bootstrap database")
	}
}

func (c *Coordinator) forward() {
	if c.err != nil {
		return
	}
	for j := range c.in {
		res, err := c.db.NamedExec(`
			INSERT INTO builds
			(repo_full_name, commit_hash, commit_message, environment, commands)
			VALUES
			(:repo_full_name, :commit_hash, :commit_message, :environment, :commands)
		`, j)
		if err != nil {
			c.err = errors.E(err, "cannot add job to database")
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			c.err = errors.E(err, "cannot load ID from the added job")
			return
		}
		j.ID = id
		c.out <- j
	}
}

// Error returns the last found error.
func (c *Coordinator) Error() error {
	return c.err
}

// Enqueue puts a build into the building pipeline.
func (c *Coordinator) Enqueue(b *Build) {
	c.in <- b
}

// Next returns the next job in the pipe. If nil, the client must stop reading.
func (c *Coordinator) Next() *Build {
	return <-c.out
}

// MarkInProgress determines a build has started and update its build
// information in the database.
func (c *Coordinator) MarkInProgress(b *Build) error {
	b.StartedAt = time.Now()
	_, err := c.db.NamedExec(`
	UPDATE builds
		SET in_progress = true, started_at = :started_at
		WHERE id = :id
	`, b)
	return errors.E(err, "cannot set job to in-progress")
}

// MarkComplete determines a build has completed and update its build
// information in the database
func (c *Coordinator) MarkComplete(b *Build) error {
	b.CompletedAt = time.Now()
	b.InProgress = false
	_, err := c.db.NamedExec(`
		UPDATE builds
		SET
			in_progress = false,
			success = :success,
			log = :log,
			completed_at = :completed_at
		WHERE id = :id
	`, b)
	return errors.E(err, "cannot update job in the database")
}
