package processor

import (
	"context"
	"fmt"
)

// LoggingProcessor is a simple processor that logs the event details.
type LoggingProcessor struct{}

func NewLoggingProcessor() *LoggingProcessor {
	return &LoggingProcessor{}
}

func (p *LoggingProcessor) Process(ctx context.Context, event Event) error {
	fmt.Printf("[processor] handling event: %+v\n", event)
	return nil
}
