package storage

import (
	"errors"

	"cirello.io/bloomfilterd/internal/filter"
)

var (
	errFilterAlreadyExists = errors.New("filter already exists")
	errInvalidStorage      = errors.New("invalid storage type")
)

type Type int

const (
	Memory Type = iota
	Raft
)

type Engine interface {
	Add(name string, size uint64, hashcount int) error
	List() []string
	Filter(name string) *filter.Bloomfilter
}

func New(t Type) (Engine, error) {
	switch t {
	case Memory:
		return &Mem{
			filters: make(map[string]*filter.Bloomfilter),
		}, nil
	case Raft:
		panic("not implemented")
	default:
		return nil, errInvalidStorage
	}
}

func Must(t Type) Engine {
	e, err := New(t)
	if err != nil {
		panic(err)
	}
	return e
}
