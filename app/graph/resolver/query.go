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

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	start := time.Now()
	log.Printf("Resolving User query with ID: %s", id)

	userId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid user ID format: %s, error: %v", id, err)
		return nil, fmt.Errorf("invalid ID format")
	}

	user, err := r.UserService.GetUser(ctx, userId)
	if err != nil {
		log.Printf("Error fetching user with ID %s: %v", id, err)
	} else {
		log.Printf("Successfully fetched user with ID %s in %v", id, time.Since(start))
	}

	return user, err
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	start := time.Now()
	log.Printf("Resolving Post query with ID: %s", id)

	postId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid post ID format: %s, error: %v", id, err)
		return nil, fmt.Errorf("invalid post ID format")
	}

	post, err := r.PostService.GetPostById(ctx, postId)
	if err != nil {
		log.Printf("Error fetching post with ID %s: %v", id, err)
	} else {
		log.Printf("Successfully fetched post with ID %s in %v", id, time.Since(start))
	}

	return post, err
}

func (r *queryResolver) Replies(ctx context.Context, commentID string, limit int32, offset int32) ([]*model.Comment, error) {
	start := time.Now()
	log.Printf("Resolving Replies for commentID: %s, limit: %d, offset: %d", commentID, limit, offset)

	parentId, err := uuid.Parse(commentID)
	if err != nil {
		log.Printf("Invalid comment ID format: %s, error: %v", commentID, err)
		return nil, fmt.Errorf("invalid comment ID format")
	}

	replies, err := r.CommentService.GetCommentReplies(ctx, parentId, int(limit), int(offset))
	if err != nil {
		log.Printf("Error fetching replies for comment %s: %v", commentID, err)
	} else {
		log.Printf("Successfully fetched %d replies for comment %s in %v", len(replies), commentID, time.Since(start))
	}

	return replies, err
}

func (r *queryResolver) Posts(ctx context.Context, limit int32, offset int32, sortBy *model.SortBy) ([]*model.Post, error) {
	start := time.Now()
	sort := "default"
	if sortBy != nil {
		sort = string(*sortBy)
	}
	log.Printf("Resolving Posts query with limit: %d, offset: %d, sortBy: %s", limit, offset, sort)

	posts, err := r.PostService.GetPosts(ctx, int(limit), int(offset), sortBy)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
	} else {
		log.Printf("Successfully fetched %d posts in %v", len(posts), time.Since(start))
	}

	return posts, err
}

func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
