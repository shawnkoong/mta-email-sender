package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func handleMessage(message kafka.Message) {
	var data map[string]interface{}
	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Printf("Failed to parse message: %s\n", err.Error())
		return
	}
	for key, value := range data {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
}
