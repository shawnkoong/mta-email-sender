package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

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
		Topic:   "emails-topic",
		GroupID: "1",
	})

	//for using cloudkarafka
	//mechanism, err := scram.Mechanism(scram.SHA256, "username", "password")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//dialer := &kafka.Dialer{
	//	SASLMechanism: mechanism,
	//}
	//r := kafka.NewReader(kafka.ReaderConfig{
	//	Brokers: []string{"localhost:" + kafkaPort},
	//	Topic:   "emails",
	//	Dialer:  dialer,
	//})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		case <-signals:
			cancel()
		default:
			message, err := reader.FetchMessage(ctx)
			if err != nil {
				if err.Error() == context.Canceled.Error() {
					break
				}
				log.Printf("error reading message: %s\n", err.Error())
				continue
			}
			handleMessage(message)
			err = reader.CommitMessages(context.Background(), message)
			if err != nil {
				fmt.Printf("Error commiting message: %v\n", err)
			}
		}
	}

	if err := reader.Close(); err != nil {
		log.Fatal("failed to close:", err)
	}
}
