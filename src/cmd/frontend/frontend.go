package main

import (
	"api"

	"github.com/golang/groupcache"
)

type Frontend struct {
	cacheGroup *groupcache.Group
}

func NewServer(cacheGroup *groupcache.Group) *Frontend {
	return &Frontend{cacheGroup}
}

func (s *Frontend) Get(args *api.Load, reply *api.ValueResult) error {
	var data []byte
	err := s.cacheGroup.Get(
		nil,
		args.Key,
		groupcache.AllocatingByteSliceSink(&data),
	)
	reply.Value = string(data)
	return err
}
