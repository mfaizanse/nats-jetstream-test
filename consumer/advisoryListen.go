package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS test
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	// Create JetStream Context
	// js, err := nc.JetStream()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// define stream and prefix
	streamName := "sap"
	fmt.Printf("Stream: %s\n", streamName)

	// Subscribe to maxDeliver advisory
	_, err = nc.Subscribe(fmt.Sprintf("$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.%s.>", streamName), func(msg *nats.Msg) {
		fmt.Println("Received advisory msg for maxDeliver !!!")
		fmt.Printf(string(msg.Data))
	})
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(20 * time.Minute)
	log.Println("Closing...")
}
