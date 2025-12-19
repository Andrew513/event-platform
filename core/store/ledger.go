package store

import "sync"

type LedgerStore struct {
	mu sync.Mutex
	balances map[string]float64
}

func NewLedgerStore() *LedgerStore {
	return &LedgerStore {
		balances: make(map[string]float64),
	}
}

func (s *LedgerStore) GetBalance(account string) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.balances[account]
}

func (s *LedgerStore) ApplyDelta(account string, delta float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.balances[account] += delta
}