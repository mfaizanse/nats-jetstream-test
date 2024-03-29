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

	// Subscribe to maxDeliver advisory
	_, err = nc.Subscribe(fmt.Sprintf("$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.%s.>", streamName), func(msg *nats.Msg) {
		fmt.Println("Received advisory msg for maxDeliver !!!")
		fmt.Printf(string(msg.Data))
	})
	if err != nil {
		log.Fatalln(err)
	}

	// define subject name for consumer
	eventType := subjectPrefix + ".>"

	// define event handler for consumer
	eventHandlerCallback := func(m *nats.Msg) {
		// fmt.Println(m.Header)

		msg := string(m.Data)
		fmt.Printf("Subject %s,  Message: %s\n", m.Subject, msg)

		time.Sleep(10 * time.Second)

		m.Ack()
	}

	//syncCallback := func(m *nats.Msg) {
	//	eventHandlerCallback(m)
	//}

	asyncCallback := func(m *nats.Msg) {
		go eventHandlerCallback(m)
	}

	// Create/bind Async Durable Consumer
	_, err = js.Subscribe(
		eventType,
		asyncCallback,
		nats.Durable("consumer1"),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.IdleHeartbeat(30 * time.Second),
		nats.EnableFlowControl(),
		nats.MaxAckPending(10),
		nats.MaxDeliver(3),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// sub.SetPendingLimits(5, 10)


	time.Sleep(20 * time.Minute)
	log.Println("Closing...")
}
