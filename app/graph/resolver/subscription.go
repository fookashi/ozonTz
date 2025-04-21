package resolver

import (
	"app/graph"
	"app/graph/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	postId, err := uuid.Parse(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid postID format: %w", err)
	}

	return r.PubSubClient.SubscribeOnComments(ctx, postId)
}

func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
