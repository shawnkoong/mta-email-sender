package main

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

func handleMessage(message kafka.Message, emailSender EmailSender) {
	var data map[string]map[string]map[string][]string // "alerts": {"email": {"route": alerts}}
	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Printf("Failed to parse message: %s\n", err.Error())
		return
	}
	tracker := getEmailTracker()
	for key, emailMap := range data {
		fmt.Printf("Key: %s, Value: %v\n", key, emailMap)
		for email, routeAlerts := range emailMap {
			go handleEmail(email, routeAlerts, tracker, emailSender)
		}
	}
	//for key, emailMap := range data {
	//	tracker := getEmailTracker()
	//	emailsSent := tracker.emailMap
	//	// wrap in goroutine
	//	for email, routeAlerts := range emailMap {
	//		//sent, ok := emailsSent[email]
	//		//if !ok {
	//		//	emailsSent[email] = NewRouteTracker()
	//		//}
	//		//for line, alert := range routeAlerts {
	//		//}
	//	}
	//	fmt.Printf("Key: %s, Value: %v\n", key, emailMap)
	//}
}

// function to handle sending to one user's email
func handleEmail(email string, routeAlerts map[string][]string, tracker *EmailTracker, emailSender EmailSender) {
	routeMap, ok := tracker.get(email)
	if !ok {
		routeTracker := NewRouteTracker()
		tracker.update(email, routeTracker)
		routeMap, _ = tracker.get(email)
	}
	for route, alerts := range routeAlerts {
		if routeMap.checkLastTimeSent(route) {
			var sb strings.Builder
			for _, alert := range alerts {
				sb.WriteString("<p>")
				sb.WriteString(alert)
				sb.WriteString("<p>")
			}
			body := sb.String()
			err := emailSender.SendEmail("MTA Chingu Update", body, []string{email}, nil)
			if err != nil {
				fmt.Printf("Error sending email: %v\n", err.Error())
			}
		}
	}
}
