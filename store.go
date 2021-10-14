package main

import "sync"

// Store is an interface for simple in-mem store application.
type Store interface {
	// Set inserts/updates given key with given val
	Set(key, val string)

	// Get returns val and existence for given key
	Get(key string) (string, bool)

	// Delete deletes given key from store
	Delete(key string)

	// Flush deletes all data from store
	Flush()

	// Count returns data count
	Count() int
}

type store struct {
	lock sync.RWMutex
	data map[string]string
}

// NewStore creates a store.
func NewStore() Store {
	return &store{
		lock: sync.RWMutex{},
		data: map[string]string{},
	}
}

func (s *store) Set(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
}

func (s *store) Get(key string) (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	val, found := s.data[key]

	return val, found
}

func (s *store) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
}

func (s *store) Flush() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = map[string]string{}
}

func (s *store) Count() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.data)
}
