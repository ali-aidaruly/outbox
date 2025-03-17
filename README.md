# Outbox Pattern Library for Go

## Overview

This library provides an implementation of the **Outbox Pattern** for Golang applications.  
It ensures **reliable event publishing** by storing messages in a database before sending them to a **message broker (NATS)**.

## Basic Usage

### **1. Initialize Outbox**
```go
db, _ := sql.Open("pgx", "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable")
natsConn, _ := nats.Connect("nats://localhost:4322")

outbox, _ := outbox.NewOutbox(db, natsConn, 
    outbox.WithPollInterval(500*time.Millisecond),
    outbox.WithPollBatchSize(10),
    outbox.WithLeaseDuration(5*time.Second),
)
```
### **2. Message creation**
```go
ctx := context.Background()
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    log.Fatal("Failed to start transaction:", err)
}

_, err = tx.ExecContext(ctx, "INSERT INTO orders (id, status) VALUES ($1, $2)", orderID, "pending")
if err != nil {
    tx.Rollback()
    log.Fatal("Failed to insert order:", err)
}

msgID, err := outbox.CreateMessage(ctx, tx, "order.created", []byte(`{"order_id": 123}`))
if err != nil {
    tx.Rollback()
    log.Fatal("Failed to insert outbox message:", err)
}

if err := tx.Commit(); err != nil {
    log.Fatal("Failed to commit transaction:", err)
}
```
### **3. Start Worker**
```go
outbox.StartWorker(ctx)
```

### **Key Features**
- **At-Least-Once Delivery** → Guarantees message publishing even if failures occur.
- **FIFO Processing** → Ensures events are processed in the order they were created.
- **Multi-Worker Safe** → Uses **`FOR UPDATE SKIP LOCKED`** to prevent multiple workers from processing the same message.
- **Configurable Lease Locking** → Prevents duplicate processing by locking messages temporarily.
- **Database Agnostic** → Supports **PostgreSQL (default)** and is extensible for other databases.
## Flexible Configuration

The library allows users to customize various parameters to fit their needs.

```go
outbox, _ := outbox.NewOutbox(db, natsConn, 
    outbox.WithPollInterval(500*time.Millisecond), // Adjust how often the worker polls for new messages
    outbox.WithPollBatchSize(10),                  // Set the number of messages processed per batch
    outbox.WithLeaseDuration(5*time.Second),       // Prevent duplicate processing by locking messages
    outbox.WithHardDelete(true),                   // Enable hard delete (remove messages after processing)
    outbox.WithDebugLogs(),                        // Enable detailed logging for debugging
)
```

### **How It Works**
1. **Producer** inserts messages into the `outbox` table instead of publishing directly.
2. **Worker** periodically fetches unprocessed messages and publishes them to **NATS**.
3. **After successful publishing**, messages are marked as processed (`published_at IS NOT NULL`).
4. **Multiple workers can safely process messages** without conflicts.

---

## Database Migrations

Before using the library, you must apply the required database migrations.
The migration files are located in:

migrations/postgres/

Run these migrations on your PostgreSQL database before starting the application.  
Indexes are optional but recommended for performance.

## Testing the Library

To test the full outbox pattern workflow, follow these steps:

Navigate to the `example/` directory and start the required services:
```sh
cd outbox/example

make up

make migrate

make run HTTP_PORT=9090 POLL_INTERVAL=10s POLL_BATCH_SIZE=3 LEASE_DURATION=10s
make run HTTP_PORT=9091
make run HTTP_PORT=9092

make run-consumer

go run scripts/send_messages.go --port=9090 --start=1 --count=100
go run scripts/send_messages.go --port=9091 --start=101 --count=100
go run scripts/send_messages.go --port=9092 --start=201 --count=100
go run scripts/send_messages.go --port=9090 --start=301 --count=100
```

