package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	kafkaPort := os.Getenv("KAFKA_PORT")
	if kafkaPort == "" {
		log.Fatal("Missing port for kafka")
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:" + kafkaPort},
		Topic:   "emails",
	})
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}
		fmt.Printf("Received message: %v\n", message.Value)

		err = reader.CommitMessages(context.Background(), message)
		if err != nil {
			fmt.Printf("Error commiting message: %v\n", err)
		}
	}
}
