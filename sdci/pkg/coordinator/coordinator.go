package coordinator // import "cirello.io/exp/sdci/pkg/coordinator"
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strings"
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
	// TODO: remove this log line when the relationship between coordinator
	// and its peers is solved.
	log.Println("coordinator err trapped:", err)
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
func (c *Coordinator) Enqueue(b *models.Build, sig string, body []byte) error {
	recipe, ok := c.configuration[b.RepoFullName]
	if !ok {
		return errors.Errorf("cannot find recipe for", b.RepoFullName)
	}
	if recipe.GithubSecret != "" &&
		!isValidSecret(sig, []byte(recipe.GithubSecret), body) {
		return errors.E("invalid signature")
	}
	b.Recipe = &recipe
	c.in <- b
	return nil
}

func isValidSecret(sig string, secret, body []byte) bool {
	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))
	if len(sig) != signatureLength || !strings.HasPrefix(sig, signaturePrefix) {
		return false
	}
	actual := make([]byte, 20)
	hex.Decode(actual, []byte(strings.TrimPrefix(sig, signaturePrefix)))
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	signedBody := []byte(computed.Sum(nil))
	return hmac.Equal(signedBody, actual)
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

// GetLastBuildStatus loads last known build status for a repository.
func (c *Coordinator) GetLastBuildStatus(repoFullName string) models.Status {
	build, err := c.buildDAO.GetLastBuild(repoFullName)
	err = errors.E(err, "cannot load last build for repository")
	c.setError(err)
	if err != nil {
		return models.Unknown
	}
	return build.Status()
}
