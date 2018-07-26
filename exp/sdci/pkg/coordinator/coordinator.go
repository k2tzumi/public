package coordinator // import "cirello.io/exp/sdci/pkg/coordinator"
import (
	"sync"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/models"
	"github.com/jmoiron/sqlx"
)

// Coordinator takes and dispatches build requests.
type Coordinator struct {
	configuration models.Configuration
	buildDAO      *models.BuildDAO
	in            chan *models.Build
	out           mappedChans

	errMu sync.Mutex
	err   error
}

// New creates a new coordinator
func New(db *sqlx.DB, configuration models.Configuration) *Coordinator {
	c := &Coordinator{
		configuration: configuration,
		buildDAO:      models.NewBuildDAO(db),
		in:            make(chan *models.Build, 10),
	}
	c.bootstrap()
	go c.forward()
	return c
}

func (c *Coordinator) setError(err error) {
	if err == nil {
		return
	}
	c.errMu.Lock()
	c.err = err
	c.errMu.Unlock()
}

func (c *Coordinator) bootstrap() {
	if err := c.buildDAO.Bootstrap(); err != nil {
		c.setError(errors.E(err, "cannot bootstrap BuildDAO"))
	}
}

func (c *Coordinator) forward() {
	if c.err != nil {
		return
	}
	for build := range c.in {
		build, err := c.buildDAO.Register(build)
		if err != nil {
			c.setError(errors.E(err, "cannot register build"))
			return
		}
		c.out.ch(build.RepoFullName) <- build
	}
}

// Error returns the last found error.
func (c *Coordinator) Error() error {
	c.errMu.Lock()
	defer c.errMu.Unlock()
	return c.err
}

// Enqueue puts a build into the building pipeline.
func (c *Coordinator) Enqueue(b *models.Build) {
	recipe, ok := c.configuration[b.RepoFullName]
	if !ok {
		c.setError(errors.Errorf("cannot find recipe for", b.RepoFullName))
		return
	}
	b.Recipe = recipe
	c.in <- b
}

// Next returns the next job in the pipe. If nil, the client must stop reading.
func (c *Coordinator) Next(repoFullName string) *models.Build {
	return <-c.out.ch(repoFullName)
}

// MarkInProgress determines a build has started and update its build
// information in the database.
func (c *Coordinator) MarkInProgress(build *models.Build) error {
	err := errors.E(c.buildDAO.MarkInProgress(build))
	c.setError(err)
	return err
}

// MarkComplete determines a build has completed and update its build
// information in the database.
func (c *Coordinator) MarkComplete(build *models.Build) error {
	err := errors.E(c.buildDAO.MarkComplete(build))
	c.setError(err)
	return err
}
