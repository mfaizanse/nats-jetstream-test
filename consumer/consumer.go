package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalln(err)
	}

	// streamName := "STREAM1"
	subjectPrefix := "EVENTING"

	eventType := subjectPrefix + ".>"
	// Simple Async Ephemeral Consumer
	js.Subscribe(eventType, func(m *nats.Msg) {
		fmt.Printf("Subject %s,  Message: %s\n", m.Subject, string(m.Data))
	})

	time.Sleep(10 * time.Minute)
	log.Println("Closing...")
}
