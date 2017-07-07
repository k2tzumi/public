package storage

import (
	"math/rand"
	"time"
)

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{make(map[string]string)}
}

func (db *Storage) Get(key string) string {
	// It adds some latency, so it becomes more obvious when the result
	// is being returned by the storage.
	time.Sleep(time.Duration(rand.Intn(1500)+1500) * time.Millisecond)
	return db.data[key]
}

func (db *Storage) Set(key string, value string) {
	db.data[key] = value
}
