package pubsub_redis

import (
	"app/graph/model"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisPubSub(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	comment := &model.Comment{
		ID:      uuid.New().String(),
		Content: "Test comment",
	}

	t.Run("PublishComment", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			pubsub := NewRedisPubSub(client)

			payload, err := json.Marshal(comment)
			require.NoError(t, err)

			channel := pubsub.getChannel(postID)
			mock.ExpectPublish(channel, payload).SetVal(1)

			err = pubsub.PublishComment(ctx, postID, comment)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("publish error", func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			pubsub := NewRedisPubSub(client)

			payload, err := json.Marshal(comment)
			require.NoError(t, err)

			channel := pubsub.getChannel(postID)
			mock.ExpectPublish(channel, payload).SetErr(errors.New("publish failed"))

			err = pubsub.PublishComment(ctx, postID, comment)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to publish comment")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	})
	t.Run("GetChannelName", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			pubsub := NewRedisPubSub(nil)
			postID := uuid.New()
			expected := "comments:" + postID.String()
			assert.Equal(t, expected, pubsub.getChannel(postID))
		})
	})
}
