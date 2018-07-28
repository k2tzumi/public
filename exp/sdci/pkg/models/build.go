package models

import (
	"time"

	"cirello.io/exp/sdci/pkg/grpc/api"
)

// Build defines the necessary data to run a build successfully.
type Build struct {
	*api.Build
	*api.Recipe
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

// BuildRepository manipulate a collection of Build requests.
type BuildRepository interface {
	Bootstrap() error
	Register(build *Build) (*Build, error)
	MarkInProgress(build *Build) error
	MarkComplete(build *Build) error
	GetLastBuild(repoFullName string) (*Build, error)
	ListByRepoFullName(repoFullName string) ([]*Build, error)
	SweepExpired(timeout time.Duration) (int64, error)
}
