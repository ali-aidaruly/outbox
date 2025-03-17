package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ali-aidaruly/outbox"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"
)

var (
	httpPort      = flag.Int("http-port", 8080, "HTTP server port")
	pollInterval  = flag.Duration("poll-interval", 500*time.Millisecond, "Polling interval for worker")
	pollBatchSize = flag.Int("poll-batch-size", 10, "Number of messages fetched per batch")
	leaseDuration = flag.Duration("with-lease-duration", 5*time.Second, "Lease duration for locking messages")
	outboxService *outbox.Outbox
)

type application struct {
	outboxService *outbox.Outbox
	db            *sql.DB
}

func main() {
	flag.Parse()

	db, err := sql.Open("pgx", "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	natsConn, err := nats.Connect("nats://localhost:4322")
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer natsConn.Close()

	outboxService, err = outbox.NewOutbox(db, natsConn,
		outbox.WithPollInterval(*pollInterval),
		outbox.WithPollBatchSize(*pollBatchSize),
		outbox.WithLeaseDuration(*leaseDuration),
		outbox.WithDebugLogs(),
	)
	if err != nil {
		log.Fatal("Failed to initialize outbox:", err)
	}

	app := application{
		outboxService: outboxService,
		db:            db,
	}

	ctx := context.Background()
	go outboxService.StartWorker(ctx)

	http.HandleFunc("/create-message", app.createMessageHandler)
	addr := fmt.Sprintf(":%d", *httpPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	msg := Message{}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tx, err := app.db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(msg.Data)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}

	msgID, err := app.outboxService.CreateMessage(ctx, tx, msg.Type, payload)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message created with ID: " + msgID.String()))
}
