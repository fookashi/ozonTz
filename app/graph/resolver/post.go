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

func (r *postResolver) Comments(ctx context.Context, obj *model.Post, limit int32, offset int32) ([]*model.Comment, error) {
	start := time.Now()
	log.Printf("Resolving Comments for post %s with limit %d, offset %d", obj.ID, limit, offset)

	postID, err := uuid.Parse(obj.ID)
	if err != nil {
		log.Printf("Invalid post ID format: %s, error: %v", obj.ID, err)
		return nil, fmt.Errorf("invalid post id: %w", err)
	}

	comments, err := r.CommentService.GetByPost(ctx, postID, int(limit), int(offset))
	if err != nil {
		log.Printf("Error fetching comments for post %s: %v", obj.ID, err)
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	log.Printf("Successfully fetched %d comments for post %s in %v", len(comments), obj.ID, time.Since(start))
	return comments, nil
}

func (r *Resolver) Post() graph.PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }
