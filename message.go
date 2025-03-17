package outbox

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID      uuid.UUID
	Type    string
	Payload []byte

	CreatedAt   time.Time
	LeaseUntil  *time.Time
	PublishedAt *time.Time
}
