package coordinator // import "cirello.io/exp/sdci/pkg/coordinator"

// Recipe defines the execution steps and environment.
type Recipe struct {
	Clone       string
	Slack       string
	Channel     string
	Environment []string
	Commands    string
}

// Build defines the necessary data to run a build successfully.
type Build struct {
	RepoFullName  string
	CommitHash    string
	CommitMessage string
	Recipe        *Recipe
}

// Coordinator takes and dispatches build requests.
type Coordinator struct {
	jobs chan *Build
}

// New creates a new coordinator
func New() *Coordinator {
	return &Coordinator{
		// TODO: replace with proper queues
		jobs: make(chan *Build, 10),
	}
}

// Enqueue puts a build into the building pipeline.
func (c *Coordinator) Enqueue(b *Build) {
	c.jobs <- b
}

// Next returns the next job in the pipe. If nil, the client must stop reading.
func (c *Coordinator) Next() *Build {
	return <-c.jobs
}
