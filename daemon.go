package main

import (
	"errors"
	"sync"

	"cirello.io/bloomfilterd/internal/filter"
)

var (
	errFilterAlreadyExists = errors.New("filter already exists")
)

type daemon struct {
	mu      sync.Mutex
	filters map[string]*filter.Bloomfilter
}

func (d *daemon) add(name string, size uint32, hashcount int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.filters[name]; ok {
		return errFilterAlreadyExists
	}

	d.filters[name] = filter.New(size, hashcount)
	return nil
}

func (d *daemon) list() []string {
	d.mu.Lock()
	defer d.mu.Unlock()

	names := make([]string, 0, len(d.filters))
	for name := range d.filters {
		names = append(names, name)
	}
	return names
}

func (d *daemon) filter(name string) *filter.Bloomfilter {
	d.mu.Lock()
	defer d.mu.Unlock()
	if f, ok := d.filters[name]; ok {
		return f
	}
	return nil
}
