package outbox

import (
	"context"
	"errors"
	"testing"

	"github.com/ali-aidaruly/outbox/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)

	o := Outbox{
		repo: mockRepo,
	}

	ctx := context.Background()
	validMsgType := "test.event"
	validPayload := []byte(`{"key": "value"}`)

	t.Run("err: msg type is empty", func(t *testing.T) {
		_, err := o.CreateMessage(ctx, nil, "", validPayload)
		expectedErr := "msg type is required"

		require.ErrorContains(t, err, expectedErr)
	})

	t.Run("err: payload is empty", func(t *testing.T) {
		_, err := o.CreateMessage(ctx, nil, validMsgType, nil)
		expectedErr := "payload is required"

		require.ErrorContains(t, err, expectedErr)
	})

	t.Run("err: createMessage fails", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.EXPECT().CreateMessage(ctx, nil, validMsgType, validPayload).Return(uuid.Nil, expectedErr)

		_, err := o.CreateMessage(ctx, nil, validMsgType, validPayload)

		require.ErrorContains(t, err, expectedErr.Error())
	})

	t.Run("happy path", func(t *testing.T) {
		mockRepo.EXPECT().CreateMessage(ctx, nil, validMsgType, validPayload).Return(uuid.Nil, nil)

		_, err := o.CreateMessage(ctx, nil, validMsgType, validPayload)

		require.NoError(t, err)
	})
}
