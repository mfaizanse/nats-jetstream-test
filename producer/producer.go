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

	streamName := "STREAM1"
	subjectPrefix := "EVENTING"

	streamInfo, err := js.StreamInfo(streamName)
	if err == nil {
		log.Println("Found stream: " + streamInfo.Config.Name)
	} else {
		// Create a Stream
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{subjectPrefix + ".>"},
		})
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Created stream: " + streamName)
	}

	for i := 0; i < 500; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		eventType1 := subjectPrefix + ".one"
		// eventType2 := subjectPrefix + ".one.two"
		// eventType3 := subjectPrefix + ".one.two.three"

		// ack, err := js.Publish(eventType1, []byte(msg), nats.MsgId(strconv.Itoa(i+1))) // For message deduplication
		ack, err := js.Publish(eventType1, []byte(msg))
		if err != nil {
			log.Println("Failed to publish on subject: " + eventType1)
			log.Println(err)
		}

		// ack, err = js.Publish(eventType2, []byte(msg))
		// if err != nil {
		// 	log.Println("Failed to publish on subject: " + eventType2)
		// 	log.Println(err)
		// }

		// ack, err = js.Publish(eventType3, []byte(msg))
		// if err != nil {
		// 	log.Println("Failed to publish on subject: " + eventType3)
		// 	log.Println(err)
		// }

		log.Println("Published to Stream: " + ack.Stream)
		time.Sleep(1 * time.Second)
	}
}
