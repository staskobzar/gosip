package transport

import (
	"sync"

	"gosip/pkg/logger"
)

// Store concurency safe storage type.
type Store[T any] struct {
	mu   sync.RWMutex
	pool map[string]T
}

// NewStore creates new store for connection or listener.
func NewStore[T any]() *Store[T] {
	return &Store[T]{
		pool: make(map[string]T),
	}
}

// Put connection or listener into store.
func (s *Store[T]) Put(key string, val T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pool[key] = val
	logger.Log("[%T] store add %q. store size %d", val, key, len(s.pool))
}

// Get connection or listener from store by key.
func (s *Store[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var val T

	if len(key) > 0 {
		val, ok := s.pool[key]
		logger.Log("store get %q found %v", key, ok)

		return val, ok
	}

	// lookup first available ip
	for _, val := range s.pool {
		return val, true
	}

	return val, false
}

// Del removes connection or listener from store.
func (s *Store[T]) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pool, key)
	logger.Log("store del %q. store size %d", key, len(s.pool))
}

// Len returns number of elements in the store
func (s *Store[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.pool)
}
