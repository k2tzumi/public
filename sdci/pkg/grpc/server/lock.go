package server

import (
	"sync"
	"time"

	"cirello.io/errors"
)

type lock struct {
	mu         sync.Mutex
	locked     bool
	seq        int
	lastUpdate time.Time
}

func (l *lock) tryLock() (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.locked {
		return 0, errors.E("already locked")
	}
	l.locked = true
	l.seq++
	l.lastUpdate = time.Now()
	return l.seq, nil
}

func (l *lock) isOwner(seq int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.seq == seq
}

func (l *lock) refresh(seq int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.seq == seq {
		l.lastUpdate = time.Now()
	}
}

func (l *lock) release(seq int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.locked {
		return errors.E("already unlocked")
	}
	if l.seq != seq {
		return errors.E("not current lock owner")
	}
	l.locked = false
	return nil
}

func (l *lock) expire(ttl time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.lastUpdate) > ttl {
		l.locked = false
		l.seq++
		return true
	}
	return false
}
