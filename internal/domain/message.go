package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID      uuid.UUID
	Type    string
	Payload []byte

	CreatedAt   time.Time
	LeaseUntil  *time.Time
	PublishedAt *time.Time
}
