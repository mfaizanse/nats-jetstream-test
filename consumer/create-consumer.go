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
	jsCtx, err := nc.JetStream()
	if err != nil {
		log.Fatalln(err)
	}

	// define stream and prefix
	// streamName := "sap"
	consumerName := "dummy_123"
	namespacedSubjectName := "dumm123456"
	jsSubject := ">"

	// define event handler for consumer
	eventHandlerCallback := func(m *nats.Msg) {

		msg := string(m.Data)
		fmt.Printf("Subject %s,  Message: %s\n", m.Subject, msg)

		m.Ack()
	}

	_, err = jsCtx.Subscribe(
		jsSubject,
		eventHandlerCallback,
		nats.Durable(consumerName),
		nats.Description(namespacedSubjectName),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.IdleHeartbeat(1 * time.Minute),
		nats.EnableFlowControl(),
		nats.DeliverAll(),
		nats.MaxAckPending(10),
		nats.MaxDeliver(100),
		nats.AckWait(30 * time.Second),
	)
	if err != nil {
		log.Fatalln(err)
	}

	//jsCtx.AddConsumer(streamName, &nats.ConsumerConfig{
	//	Durable: consumerName,
	//	Description: namespacedSubjectName,
	//	AckPolicy: nats.AckExplicitPolicy,
	//	FlowControl: true,
	//
	//})

	time.Sleep(20 * time.Minute)
	log.Println("Closing...")
}
