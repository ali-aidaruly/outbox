package outbox

import (
	"database/sql"
	"github.com/ali-aidaruly/outbox/internal/repository"
	"log/slog"
	"time"
)

type Client struct {
	cfg    Config
	repo   repository.Repository
	logger *slog.Logger
}

func NewClient(db *sql.DB, opts ...Option) (*Client, error) {
	cfg := Config{
		HardDelete:    false,
		LeaseDuration: 30 * time.Second,
		DBType:        repository.DBPostgres,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	repo, err := repository.NewRepository(db, cfg.DBType)
	if err != nil {
		return nil, err
	}

	client := &Client{
		repo: repo,
		cfg:  cfg,
	}

	return client, nil
}
