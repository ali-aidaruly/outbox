package outbox

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (o *Outbox) processMessages(ctx context.Context) error {
	logger.Debug("Fetching next batch of messages",
		slog.Int("batch_size", o.cfg.pollBatchSize),
		slog.Duration("lease_duration", o.cfg.leaseDuration),
	)

	messages, err := o.repo.FetchNextBatch(ctx, o.cfg.pollBatchSize, o.cfg.leaseDuration) // Fetch batch
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	if len(messages) == 0 {
		logger.Debug("No messages to process")
		return nil
	}

	logger.Debug("Fetched messages", slog.Int("count", len(messages)))

	var successfulIDs []uuid.UUID

	for _, msg := range messages {
		if err := o.publishMessage(msg.Type, msg.Payload); err != nil {
			logger.Error("Failed to publish message",
				slog.String("message_id", msg.ID.String()),
				slog.String("error", err.Error()),
			)
			continue
		}
		logger.Debug("Successfully published message", slog.String("message_id", msg.ID.String()))
		successfulIDs = append(successfulIDs, msg.ID)
	}

	if len(successfulIDs) == 0 {
		logger.Debug("No messages were successfully published")
		return nil
	}

	logger.Debug("Marking messages as processed", slog.Int("count", len(successfulIDs)))

	err = o.repo.MarkMessagesProcessed(ctx, successfulIDs)
	if err != nil {
		return fmt.Errorf("failed to mark messages: %w", err)
	}

	logger.Debug("Successfully marked messages as processed", slog.Int("count", len(successfulIDs)))

	return nil
}

func (o *Outbox) publishMessage(msgType string, payload []byte) error {
	logger.Debug("Publishing message",
		slog.String("type", msgType),
		slog.String("payload", string(payload)),
	)

	err := o.natsConn.Publish(msgType, payload)
	if err != nil {
		return fmt.Errorf("failed to publish message to NATS: %w", err)
	}

	logger.Debug("Message published successfully", slog.String("type", msgType))
	return nil
}
