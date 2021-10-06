package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
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

	// define stream and prefix
	streamName := "STREAM1"
	subjectPrefix := "EVENTING"
	fmt.Printf("Stream: %s, Subject Prefix: %s\n", streamName, subjectPrefix)

	// define subject name for consumer
	eventType := subjectPrefix + ".one"

	// Queue subscribe using JetStream context
	maxInFlightMessages := 5
	for i := 0; i < maxInFlightMessages; i++ {
		consumerName := fmt.Sprintf("test%d", i+1)
		queueName := "group1"

		_, err := js.QueueSubscribe(
			eventType,
			queueName,
			func(m *nats.Msg) {

				msg := string(m.Data)
				fmt.Printf("Consumer: %s, Subject %s,  Message: %s\n", consumerName, m.Subject, msg)

				m.Ack()
			},
			nats.ManualAck(),
			nats.AckExplicit(),
			nats.MaxAckPending(10),
			nats.RateLimit(10 * 1024),
		)

		if err != nil {
			log.Fatalln(err)
		}
	}


	time.Sleep(20 * time.Minute)
	log.Println("Closing...")
}
