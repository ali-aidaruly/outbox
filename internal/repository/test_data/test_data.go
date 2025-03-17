package testdata

import (
	"time"

	"github.com/ali-aidaruly/outbox/internal/domain"
	"github.com/ali-aidaruly/outbox/internal/helper"
	"github.com/google/uuid"
)

var (
	defaultType              = "default-type"
	defaultData              = []byte(`{"default-data": "data"}`)
	defaultCreatedAt         = helper.MustParseTime(time.RFC3339Nano, "2023-01-02T18:00:00Z")
	defaultExpiredLeaseUntil = helper.MustParseTime(time.RFC3339Nano, "2023-01-02T18:00:00Z")

	TestMessages = []domain.Message{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Type:        "type-1",
			Payload:     []byte(`{"some-data-1": "data-1"}`),
			CreatedAt:   helper.MustParseTime(time.RFC3339Nano, "2023-01-02T15:00:00Z"),
			LeaseUntil:  nil,
			PublishedAt: nil,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Type:        "type-1",
			Payload:     []byte(`{"some-data-2": "data-2"}`),
			CreatedAt:   helper.MustParseTime(time.RFC3339Nano, "2023-01-02T16:00:00Z"),
			LeaseUntil:  nil,
			PublishedAt: nil,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Type:        "type-3",
			Payload:     []byte(`{"some-data-3": "data-3"}`),
			CreatedAt:   helper.MustParseTime(time.RFC3339Nano, "2023-01-02T17:00:00Z"),
			LeaseUntil:  nil,
			PublishedAt: nil,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Type:        defaultType,
			Payload:     defaultData,
			CreatedAt:   defaultCreatedAt,
			LeaseUntil:  &defaultExpiredLeaseUntil,
			PublishedAt: nil,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			Type:        defaultType,
			Payload:     defaultData,
			CreatedAt:   helper.MustParseTime(time.RFC3339Nano, "2023-01-02T19:00:00Z"),
			LeaseUntil:  &defaultExpiredLeaseUntil,
			PublishedAt: nil,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			Type:        defaultType,
			Payload:     defaultData,
			CreatedAt:   helper.MustParseTime(time.RFC3339Nano, "2023-01-02T20:00:00Z"),
			LeaseUntil:  &defaultExpiredLeaseUntil,
			PublishedAt: nil,
		},
	}
)
