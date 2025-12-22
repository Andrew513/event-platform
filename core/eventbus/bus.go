package eventbus

import (
	"context"
	"fmt"

	"github.com/Andrew513/event-platform/core/processor"
)

type job struct {
	event processor.Event
	ack chan error // nil means fire-and-forget
}

type EventBus struct {
	ch        chan job
	processor processor.Processor
}

func NewEventBus(size int, p processor.Processor) *EventBus {
	return &EventBus{
		ch:        make(chan job, size),
		processor: p,
	}
}

func (b *EventBus) Submit(e processor.Event) error {
	// only care about if job is in the channel, not the result (async/eventual)
	b.ch <- job{event: e, ack: nil}
	return nil
}

// SubmitAndWait submits an event and waits for processing to complete
func (b *EventBus) SubmitAndWait(e processor.Event) error {
	ack := make(chan error, 1)
	b.ch <- job{event: e, ack: ack}
	// this will wait until the processor sends the result to the ack channel
	// and then return it
	return <- ack
}

func (b *EventBus) Start() {
	go func() {
		for j := range b.ch {
			// ctx.background() = function that returns a non-nil, empty Context, starting point for
			// all context and serve as root of a context tree
			ctx := context.Background()
			err := b.processor.Process(ctx, j.event)
			if err != nil {
				fmt.Printf("[eventbus] error processing event %s: %v\n", j.event.EventID, err)
			}
			if j.ack != nil {
				j.ack <- err
				close(j.ack)
			}
		}
	}()
}
