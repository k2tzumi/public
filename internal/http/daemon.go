package http

import (
	"cirello.io/bloomfilterd/internal/filter"
	"cirello.io/bloomfilterd/internal/storage"
	"github.com/coreos/etcd/raft/raftpb"
)

type daemon struct {
	propose    chan string
	confChange chan raftpb.ConfChange

	storage storage.Engine
}

func (d *daemon) add(name string, size uint32, hashcount int) error {
	return d.storage.Add(name, size, hashcount)
}

func (d *daemon) list() []string {
	return d.storage.List()
}

func (d *daemon) filter(name string) *filter.Bloomfilter {
	return d.storage.Filter(name)
}
