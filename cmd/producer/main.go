package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	})

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte("hello from go"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("message sent")
}