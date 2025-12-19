package processor

import "context"

// interface for different usage
// cryptoledger processor
// testprocessor
// mockprocessor
// retryprocessor(wrapper)
type Processor interface {
	Process(ctx context.Context, event Event) error
}