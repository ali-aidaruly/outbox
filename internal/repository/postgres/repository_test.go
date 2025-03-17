package postgres_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/ali-aidaruly/outbox/internal/repository/postgres"

	testdata "github.com/ali-aidaruly/outbox/internal/repository/test_data"
	"github.com/google/uuid"

	"github.com/ali-aidaruly/outbox/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateMessage(t *testing.T) {

	type input struct {
		msgType string
		payload []byte
	}

	var (
		repo = postgres.NewRepository(testDB)
		ctx  = context.Background()
	)

	tt := []struct {
		name           string
		input          input
		shouldRollBack bool
		exists         bool
		err            string
	}{
		{
			name: "happy-path",
			input: input{
				msgType: "type-1543435",
				payload: []byte(`{"data": "some-data"}`),
			},
			shouldRollBack: false,
			exists:         true,
			err:            "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			require.NoError(t, err)

			id, err := repo.CreateMessage(ctx, tx, tc.input.msgType, tc.input.payload)
			if tc.err != "" {
				require.Error(t, err)
				require.Equal(t, tc.err, err.Error())
				return
			}
			require.NoError(t, err)

			if tc.shouldRollBack {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()

				err = repo.MarkMessagesProcessed(ctx, []uuid.UUID{id}) // cleanup
				require.NoError(t, err)
			}

			var exists bool
			query := `SELECT EXISTS(SELECT 1 FROM outbox WHERE id = $1)`
			err = testDB.QueryRowContext(ctx, query, id).Scan(&exists)

			require.NoError(t, err)
			require.Equal(t, tc.exists, exists)
		})
	}
}

func TestRepository_FetchNextBatch(t *testing.T) {

	type input struct {
		batchSize int
		lease     time.Duration
	}

	var (
		repo = postgres.NewRepository(testDB)
		ctx  = context.Background()
	)

	tt := []struct {
		name            string
		input           input
		want            []domain.Message
		prepareOutputFn func(res []domain.Message, expected []domain.Message) []domain.Message
		err             string
	}{
		{
			name: "with lease_until with batch_size 3",
			input: input{
				batchSize: 3,
				lease:     time.Hour * 100,
			},
			want: []domain.Message{
				testdata.TestMessages[0],
				testdata.TestMessages[1],
				testdata.TestMessages[2],
			},
			prepareOutputFn: func(res []domain.Message, expected []domain.Message) []domain.Message {
				for i := range res {
					res[i].LeaseUntil = expected[i].LeaseUntil
				}
				return res
			},
			err: "",
		},
		{
			name: "lease_until next and with batch size 10",
			input: input{
				batchSize: 10,
				lease:     time.Hour * 100,
			},
			want: []domain.Message{
				testdata.TestMessages[3],
				testdata.TestMessages[4],
				testdata.TestMessages[5],
			},
			prepareOutputFn: func(res []domain.Message, expected []domain.Message) []domain.Message {
				for i := range res {
					res[i].LeaseUntil = expected[i].LeaseUntil
				}
				return res
			},
			err: "",
		},
		{
			name: "next is empty",
			input: input{
				batchSize: 10,
				lease:     time.Hour * 100,
			},
			want: []domain.Message{},
			err:  "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := repo.FetchNextBatch(ctx, tc.input.batchSize, tc.input.lease)
			if tc.err != "" {
				require.Error(t, err)
				require.Equal(t, tc.err, err.Error())
				return
			}

			for _, w := range res {
				require.Greater(t, *w.LeaseUntil, time.Now())
			}

			if tc.prepareOutputFn != nil {
				tc.prepareOutputFn(res, tc.want)
			}
			sort.Slice(res, func(i, j int) bool {
				return res[i].CreatedAt.Before(res[j].CreatedAt)
			})

			require.NoError(t, err)
			require.Equal(t, tc.want, res)
		})
	}
}

func TestRepository_MarkMessagesProcessed(t *testing.T) {

	type input struct {
		ids []uuid.UUID
	}

	var (
		repo = postgres.NewRepository(testDB)
		ctx  = context.Background()
	)

	tt := []struct {
		name  string
		input input
		err   string
	}{
		{
			name: "update one",
			input: input{
				ids: []uuid.UUID{
					testdata.TestMessages[0].ID,
				},
			},
			err: "",
		},
		{
			name: "update many",
			input: input{
				ids: []uuid.UUID{
					testdata.TestMessages[1].ID,
					testdata.TestMessages[2].ID,
					testdata.TestMessages[3].ID,
				},
			},
			err: "",
		},
		{
			name: "update none",
			input: input{
				ids: []uuid.UUID{},
			},
			err: "",
		},
		{
			name: "update non-existent",
			input: input{
				ids: []uuid.UUID{
					uuid.New(),
					uuid.New(),
					uuid.New(),
				},
			},
			err: "",
		},
		{
			name: "update existent with non-existent",
			input: input{
				ids: []uuid.UUID{
					testdata.TestMessages[4].ID,
					uuid.New(),
					uuid.New(),
					uuid.New(),
				},
			},
			err: "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := repo.MarkMessagesProcessed(ctx, tc.input.ids)
			if tc.err != "" {
				require.Error(t, err)
				require.Equal(t, tc.err, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
