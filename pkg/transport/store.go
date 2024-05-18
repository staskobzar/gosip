package transport

import (
	"gosip/pkg/logger"
	"sync"
)

// Store conturency safe storage type
type Store[T any] struct {
	mu   sync.RWMutex
	pool map[string]T
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		pool: make(map[string]T),
	}
}

func (s *Store[T]) Put(key string, val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pool[key] = val
	logger.Log("[%T] store add %q. store size %d", val, key, len(s.pool))
}

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

func (s *Store[T]) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pool, key)
	logger.Log("store del %q. store size %d", key, len(s.pool))
}
