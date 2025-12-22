package main

import (
	"context"
	// "encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {

	// create new kafka reader 
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "test-topic",
		GroupID: "test-consumer-group",
	})

	defer reader.Close()

	ctx := context.Background()
	log.Println("Kafka consumer started...")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"Received message: topic=%s partition=%d offset=%d value=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value),
		)
	}
}