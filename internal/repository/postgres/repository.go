package postgres

import (
	"context"
	"fmt"
	"github.com/ali-aidaruly/outbox"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateMessage(ctx context.Context, msg outbox.Message) error {
	query := `INSERT INTO outbox (type, payload) VALUES ($1, $2) RETURNING id, created_at`

	_, err := r.db.Exec(ctx, query, msg.Type, msg.Payload)
	if err != nil {
		return fmt.Errorf("")
	}

	return nil
}
