package processor

import "time"

type Event struct {
	EventID   string
	Key       string
	Type      string
	Payload   []byte
	Timestamp time.Time
}
