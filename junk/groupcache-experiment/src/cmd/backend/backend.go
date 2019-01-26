package main

import (
	"api"
	"storage"
)

func NewBackend(db *storage.Storage) *Backend {
	return &Backend{db}
}

type Backend struct {
	db *storage.Storage
}

func (b *Backend) Get(args *api.Load, reply *api.ValueResult) error {
	data := b.db.Get(args.Key)
	reply.Value = string(data)
	return nil
}

func (b *Backend) Set(args *api.Store, reply *api.NullResult) error {
	b.db.Set(args.Key, args.Value)
	*reply = 0
	return nil
}
