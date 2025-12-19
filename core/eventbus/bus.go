package eventbus

import (
	"context"
	"fmt"

	"github.com/Andrew513/event-platform/core/processor"
)

type EventBus struct {
	ch        chan processor.Event
	processor processor.Processor
}

func NewEventBus(size int, p processor.Processor) *EventBus {
	return &EventBus{
		ch:        make(chan processor.Event, size),
		processor: p,
	}
}

func (b *EventBus) Submit(e processor.Event) error {
	b.ch <- e
	return nil
}

func (b *EventBus) SubmitAndWait(e processor.Event) error {
	
}

func (b *EventBus) Start() {
	go func() {
		for e := range b.ch {
			// ctx.background() = function that returns a non-nil, empty Context, starting point for
			// all context and serve as root of a context tree
			ctx := context.Background()
			if err := b.processor.Process(ctx, e); err != nil {
				fmt.Printf("[eventbus] error processing event %s: %v\n", e.EventID, err)
			}
		}
	}()
}
