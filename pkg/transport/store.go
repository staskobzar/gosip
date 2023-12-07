package transport

import (
	"gosip/pkg/logger"
	"sync"
)

// Store conturency safe storage type
type Store[T any] struct {
	mu   sync.RWMutex
	pool map[sID]T
}

type sID string

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		pool: make(map[sID]T),
	}
}

func (s *Store[T]) Put(key sID, val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	logger.Log("store add %q", key)
	s.pool[key] = val
}

func (s *Store[T]) Get(key sID) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.pool[key]
	logger.Log("store get %q found %v", key, ok)
	return val, ok
}

func (s *Store[T]) Del(key sID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	logger.Log("store del %q", key)
	delete(s.pool, key)
}
