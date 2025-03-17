package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ali-aidaruly/outbox/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateMessage(ctx context.Context, tx *sql.Tx, msgType string, payload []byte) (uuid.UUID, error) {
	var id uuid.UUID

	query := `INSERT INTO outbox (type, payload) VALUES ($1, $2) RETURNING id`
	err := tx.QueryRowContext(ctx, query, msgType, payload).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *Repo) FetchNextBatch(ctx context.Context, batchSize int, lease time.Duration) ([]domain.Message, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		commitErr := tx.Commit()
		if commitErr != nil {
			err = fmt.Errorf("tx commit failed: %w", commitErr)
		}
	}()

	query := `
		UPDATE outbox 
		SET lease_until = NOW() + ($2 * INTERVAL '1 millisecond')
		WHERE id IN (
			SELECT id FROM outbox 
			WHERE (lease_until IS NULL OR lease_until < NOW()) AND published_at IS NULL
			ORDER BY created_at ASC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, type, payload, created_at, lease_until;
	`

	rows, err := tx.QueryContext(ctx, query, batchSize, int(lease.Milliseconds()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]domain.Message, 0, batchSize)
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.Type, &msg.Payload, &msg.CreatedAt, &msg.LeaseUntil); err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *Repo) MarkMessagesProcessed(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE outbox SET published_at = NOW() WHERE id = ANY($1)`

	_, err := r.db.ExecContext(ctx, query, pq.Array(ids))
	if err != nil {
		return fmt.Errorf("db exec failed: %w", err)
	}

	return nil
}
