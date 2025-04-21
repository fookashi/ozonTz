package resolver

import (
	"app/internal/pubsub"
	"app/internal/service"
)

type Resolver struct {
	UserService    service.User
	PostService    service.Post
	CommentService service.Comment
	PubSubClient   pubsub.PubSubClient
}
