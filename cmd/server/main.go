package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrew513/event-platform/adapters/cryptoledger"
	"github.com/Andrew513/event-platform/core/eventbus"
	"github.com/Andrew513/event-platform/core/processor"
	"github.com/Andrew513/event-platform/core/store"
)

func main() {
	store := store.NewLedgerStore()
	p := cryptoledger.NewCryptoLedgerProcessor(store)

	bus := eventbus.NewEventBus(100, p)
	bus.Start()

	depositPayload, _ := json.Marshal(cryptoledger.LedgerEvent{
		Account: "user-1",
		Amount:  100,
	})

	bus.Submit(processor.Event{
		EventID:   "e-1",
		Key:       "user-1",
		Type:      string(cryptoledger.Deposit),
		Payload:   depositPayload,
		Timestamp: time.Now(),
	})

	withdrawPayload, _ := json.Marshal(cryptoledger.LedgerEvent{
		Account: "user-1",
		Amount:  30,
	})

	bus.Submit(processor.Event{
		EventID:   "e-2",
		Key:       "user-1",
		Type:      string(cryptoledger.Withdrawal),
		Payload:   withdrawPayload,
		Timestamp: time.Now(),
	})

	// Wait for events to be processed
	time.Sleep(100 * time.Millisecond)

	fmt.Println("balance: ", store.GetBalance("user-1"))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

}
