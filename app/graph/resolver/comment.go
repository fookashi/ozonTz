package resolver

import (
	"app/graph"
	"app/graph/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func (r *commentResolver) Replies(ctx context.Context, obj *model.Comment, limit int32, offset int32) ([]*model.Comment, error) {
	start := time.Now()
	log.Printf("Fetching replies for comment %s (limit: %d, offset: %d)", obj.ID, limit, offset)

	parentId, err := uuid.Parse(obj.ID)
	if err != nil {
		log.Printf("Invalid comment ID format: %s, error: %v", obj.ID, err)
		return nil, fmt.Errorf("invalid comment ID format")
	}

	replies, err := r.CommentService.GetCommentReplies(ctx, parentId, int(limit), int(offset))
	if err != nil {
		log.Printf("Error fetching replies for comment %s: %v", obj.ID, err)
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	log.Printf("Successfully fetched %d replies for comment %s in %v", len(replies), obj.ID, time.Since(start))
	return replies, nil
}

func (r *Resolver) Comment() graph.CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }
