package coordinator // import "cirello.io/exp/sdci/pkg/coordinator"
import (
	"sync"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/models"
	"github.com/jmoiron/sqlx"
)

// Coordinator takes and dispatches build requests.
type Coordinator struct {
	in  chan *models.Build
	out chan *models.Build

	buildDAO *models.BuildDAO

	mu  sync.Mutex
	err error
}

// New creates a new coordinator
func New(db *sqlx.DB) *Coordinator {
	c := &Coordinator{
		buildDAO: models.NewBuildDAO(db),
		// TODO: replace with proper queues
		in:  make(chan *models.Build, 10),
		out: make(chan *models.Build),
	}
	c.bootstrap()
	go c.forward()
	return c
}

func (c *Coordinator) setError(err error) {
	c.mu.Lock()
	c.err = err
	c.mu.Unlock()
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
		c.out <- build
	}
}

// Error returns the last found error.
func (c *Coordinator) Error() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

// Enqueue puts a build into the building pipeline.
func (c *Coordinator) Enqueue(b *models.Build) {
	c.in <- b
}

// Next returns the next job in the pipe. If nil, the client must stop reading.
func (c *Coordinator) Next() *models.Build {
	return <-c.out
}
