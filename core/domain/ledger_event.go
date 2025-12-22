package domain

import "time"

// this is the kafka message schema for ledger events
type LedgerEvent struct {
	EventID string `json:"event_id"`
	AccountID string `json:"account_id"`
	Type string `json:"type"`
	Amount float64 `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}