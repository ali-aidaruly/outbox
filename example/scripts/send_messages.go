package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// RequestPayload defines the structure of the JSON request.
type RequestPayload struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func main() {
	// Define command-line arguments
	port := flag.Int("port", 8080, "HTTP server port")
	startNumber := flag.Int("start", 1, "Starting number for message ID counter")
	requests := flag.Int("count", 10, "Number of messages to send")
	flag.Parse()

	// Base URL for requests
	url := fmt.Sprintf("http://localhost:%d/create-message", *port)

	for i := *startNumber; i < *startNumber+*requests; i++ {
		// Construct the request payload
		payload := RequestPayload{
			Type: "order.created",
			Data: map[string]string{
				"message-id": fmt.Sprintf("%d", i),
			},
		}

		// Encode to JSON
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Send HTTP request
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		// Read response
		fmt.Printf("Sent message %d, Response: %s\n", i, resp.Status)
		resp.Body.Close()
	}

	log.Println("✅ All messages sent successfully")
}
