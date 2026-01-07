package main

import (
	"context"
	"log"
	"encoding/json"
	"os"
	"time"

	"github.com/Andrew513/event-platform/core/domain"
	"github.com/segmentio/kafka-go"
	"github.com/Andrew513/event-platform/adapters/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	brokers := []string{getenv("KAFKA_BROKER", "localhost:9092")}
	topic := getenv("KAFKA_TOPIC", "ledger.commands")
	groupID := getenv("KAFKA_GROUP", "ledger-consumer")

	dbURL := getenv("DATABASE_URL", "postgres://localhost:5432/event_platform?sslmode=disable")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	store := postgres.NewLedgerStore(pool)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic: topic,
		GroupID: groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
		CommitInterval: 0, // close auto commit, full manual commit
	})

	defer reader.Close()

	log.Printf("consumer started: brokers=%v topic=%s group=%s db=%ss", brokers, topic, groupID, dbURL)

	for{
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		var e domain.LedgerEvent
		if err := json.Unmarshal(msg.Value, &e); err != nil {
			log.Printf("invalid message  (commit and skip): err=%v msg=%s", err, string(msg.Value))
			_ = reader.CommitMessages(ctx, msg)
			continue
		}

		if e.Timestamp.IsZero() {
			e.Timestamp = time.Now()
		}

		already, err := store.ApplyEventTx(ctx, e)
		if err != nil {
			log.Printf("apply failed (no commit): event_id=%s account=%s type=%s amount=%.2f err=%v",
				e.EventID, e.AccountID, e.Type, e.Amount, err)
			continue
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("commit offset failed: err=%v", err)
			continue
		}

		log.Printf("applied ok: event_id=%s account=%s type=%s amount=%.2f already_processed=%v",
			e.EventID, e.AccountID, e.Type, e.Amount, already)

	}
}