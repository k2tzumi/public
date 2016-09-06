package storage

import (
	"sync"

	"cirello.io/bloomfilterd/internal/filter"
)

type Mem struct {
	mu      sync.Mutex
	filters map[string]*filter.Bloomfilter
}

func (m *Mem) Add(name string, size uint64, hashcount int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.filters[name]; ok {
		return errFilterAlreadyExists
	}

	m.filters[name] = filter.New(size, hashcount)
	return nil
}

func (m *Mem) List() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	names := make([]string, 0, len(m.filters))
	for name := range m.filters {
		names = append(names, name)
	}

	return names
}

func (m *Mem) Filter(name string) *filter.Bloomfilter {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f, ok := m.filters[name]; ok {
		return f
	}

	return nil
}
