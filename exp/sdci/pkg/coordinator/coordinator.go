package coordinator // import "cirello.io/exp/sdci/pkg/coordinator"
import (
	"cirello.io/errors"
	"github.com/jmoiron/sqlx"
)

// Recipe defines the execution steps and environment.
type Recipe struct {
	Clone       string `db:"clone"`
	Slack       string `db:"slack"`
	Channel     string `db:"channel"`
	Environment string `db:"environment"`
	Commands    string `db:"commands"`
}

// Build defines the necessary data to run a build successfully.
type Build struct {
	ID            int64   `db:"id"`
	RepoFullName  string  `db:"repo_full_name"`
	CommitHash    string  `db:"commit_hash"`
	CommitMessage string  `db:"commit_message"`
	Recipe        *Recipe `db:"recipe"`
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
	go c.forward()
	return c
}

func (c *Coordinator) forward() {
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
