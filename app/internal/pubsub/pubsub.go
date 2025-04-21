package pubsub

import (
	"app/graph/model"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/golang/mock/mockgen -source=pubsub.go -destination=mocks/pubsub.go

type PubSubClient interface {
	PublishComment(ctx context.Context, postId uuid.UUID, comment *model.Comment) error

	SubscribeOnComments(ctx context.Context, postId uuid.UUID) (<-chan *model.Comment, error)

	Close() error
}
