package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4322")
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	counter := 1
	mux := sync.Mutex{}
	_, err = nc.Subscribe("order.created", func(msg *nats.Msg) {
		mux.Lock()
		fmt.Printf("%d Received message: %s\n", counter, string(msg.Data))
		counter++
		mux.Unlock()
	})
	if err != nil {
		log.Fatal("Failed to subscribe to NATS:", err)
	}

	log.Println("Listening on NATS subject: order.created")
	select {} // Keep running
}
