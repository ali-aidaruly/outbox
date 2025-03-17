package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/nats-io/nats.go"
)

const (
	// DBPostgres defines the database type for PostgreSQL.
	DBPostgres = "postgres"
)

// Outbox handles message storage and processing using the outbox pattern.
type Outbox struct {
	cfg  Config
	repo repository

	natsConn *nats.Conn
}

// NewOutbox initializes a new Outbox instance with the provided database and NATS connection.
// Additional options can be passed to override default configuration settings.
func NewOutbox(db *sql.DB, natsConn *nats.Conn, opts ...Option) (*Outbox, error) {
	cfg := Config{
		dbType:     DBPostgres,
		hardDelete: false,

		pollInterval:  time.Millisecond * 500, // Default polling interval for worker.
		pollBatchSize: 10,                     // Number of messages to process per batch.
		leaseDuration: 5 * time.Second,        // Default lease duration before retrying a message.
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	repo, err := newRepository(db, cfg.dbType)
	if err != nil {
		return nil, err
	}

	client := &Outbox{
		repo: repo,
		cfg:  cfg,

		natsConn: natsConn,
	}

	return client, nil
}

// StartWorker launches the outbox worker, which periodically polls for unprocessed messages
// and attempts to publish them to NATS.
func (o *Outbox) StartWorker(ctx context.Context) {
	logger.Info("Worker started", slog.Duration("poll_interval", o.cfg.pollInterval))

	go func() {
		ticker := time.NewTicker(o.cfg.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Worker stopping due to context cancellation")
				return
			case <-ticker.C:
				if err := o.processMessages(ctx); err != nil {
					logger.Error("error while processing messages", slog.String("error", err.Error()))
				}
			}
		}
	}()
}

// CreateMessage inserts a new message into the outbox for future processing.
// It validates that the message type and payload are non-empty before saving.
func (o *Outbox) CreateMessage(ctx context.Context, tx *sql.Tx, msgType string, payload []byte) (uuid.UUID, error) {

	if msgType == "" {
		logger.Warn("Failed to create message: msg type is required")
		return uuid.Nil, fmt.Errorf("msg type is required")
	}

	if len(payload) == 0 {
		logger.Warn("Failed to create message: payload is required")
		return uuid.Nil, fmt.Errorf("payload is required")
	}

	id, err := o.repo.CreateMessage(ctx, tx, msgType, payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating message: %w", err)
	}

	logger.Debug("Message successfully created",
		slog.String("id", id.String()),
		slog.String("type", msgType),
		slog.Int("payload_size", len(payload)),
	)

	return id, nil
}
