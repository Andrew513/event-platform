package store

import (
	"sync"
)

type IdempotencyStore struct {
	mu sync.Mutex
	processed map[string]struct{}
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		processed: make(map[string]struct{}),
	}
}

func (s *IdempotencyStore) MarkIfNew(eventID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.processed[eventID]; exists {
		return false
	}

	s.processed[eventID] = struct{}{}
	return true
}