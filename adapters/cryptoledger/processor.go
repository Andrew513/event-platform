package cryptoledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Andrew513/event-platform/core/processor"
	"github.com/Andrew513/event-platform/core/store"
)

type CryptoLedgerProcessor struct {
	store *store.LedgerStore
}

func NewCryptoLedgerProcessor(store *store.LedgerStore) *CryptoLedgerProcessor {
	return &CryptoLedgerProcessor{
		store: store,
	}
}

func (p *CryptoLedgerProcessor) Process(ctx context.Context, event processor.Event) error {
	var le LedgerEvent
	if err := json.Unmarshal(event.Payload, &le); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	fmt.Printf("[processor] Processing event %s, type: %s, account: %s, amount: %.2f\n", event.EventID, event.Type, le.Account, le.Amount)
	switch EventType(event.Type) {
	case Deposit:
		p.store.ApplyDelta(le.Account, le.Amount)
		fmt.Printf("[processor] After deposit, balance: %.2f\n", p.store.GetBalance(le.Account))
	case Withdrawal:
		balance := p.store.GetBalance(le.Account)
		if balance < le.Amount {
			return errors.New("insufficient funds")
		}
		p.store.ApplyDelta(le.Account, -le.Amount)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)

	}
	return nil

}
