package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ali-aidaruly/outbox/internal/domain"
	"github.com/ali-aidaruly/outbox/internal/repository/postgres"
	"github.com/google/uuid"
)

//go:generate mockgen --source=repository.go --package=mocks --destination=mocks/mock_repository.go

type repository interface {
	CreateMessage(ctx context.Context, tx *sql.Tx, msgType string, payload []byte) (uuid.UUID, error)
	FetchNextBatch(ctx context.Context, batchSize int, lease time.Duration) ([]domain.Message, error)
	MarkMessagesProcessed(ctx context.Context, ids []uuid.UUID) error
}

func newRepository(db *sql.DB, dbType string) (repository, error) {
	switch dbType {
	case DBPostgres:
		return postgres.NewRepository(db), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
