package resolver

import (
	"app/graph"
	"app/graph/model"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreateUser(ctx context.Context, username string) (*model.User, error) {
	start := time.Now()
	log.Printf("Creating user with username: %s", username)

	user, err := r.UserService.CreateUser(ctx, username)
	if err != nil {
		log.Printf("Error creating user %s: %v", username, err)
	} else {
		log.Printf("Successfully created user %s with ID %s in %v", username, user.ID, time.Since(start))
	}

	return user, err
}

func (r *mutationResolver) CreatePost(ctx context.Context, userID string, title string, content string, isCommentable bool) (*model.Post, error) {
	start := time.Now()
	log.Printf("Creating post for user %s, title: %s, commentable: %v", userID, title, isCommentable)

	parsedUserId, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %s, error: %v", userID, err)
		return nil, fmt.Errorf("invalid user ID format")
	}

	post, err := r.PostService.CreatePost(ctx, parsedUserId, title, content, isCommentable)
	if err != nil {
		log.Printf("Error creating post: %v", err)
	} else {
		log.Printf("Successfully created post with ID %s in %v", post.ID, time.Since(start))
	}

	return post, err
}

func (r *mutationResolver) CreateComment(ctx context.Context, userID string, postID string, parentID *string, content string) (*model.Comment, error) {
	start := time.Now()
	log.Printf("Creating comment by user %s on post %s (parent: %v)", userID, postID, parentID)

	userId, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %s, error: %v", userID, err)
		return nil, errors.New("invalid user ID format")
	}

	postId, err := uuid.Parse(postID)
	if err != nil {
		log.Printf("Invalid post ID format: %s, error: %v", postID, err)
		return nil, errors.New("invalid post ID format")
	}

	var parentId *uuid.UUID
	if parentID != nil {
		parsedParentID, err := uuid.Parse(*parentID)
		if err != nil {
			log.Printf("Invalid parent comment ID format: %s, error: %v", *parentID, err)
			return nil, errors.New("invalid parent comment ID format")
		}
		parentId = &parsedParentID
	}

	comment, err := r.CommentService.CreateComment(ctx, userId, postId, parentId, content)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		return nil, err
	}

	if err := r.PubSubClient.PublishComment(ctx, postId, comment); err != nil {
		log.Printf("Failed to publish comment: %v", err)
	} else {
		log.Printf("Successfully published comment %s to post %s", comment.ID, postID)
	}

	log.Printf("Successfully created comment %s in %v", comment.ID, time.Since(start))
	return comment, nil
}

func (r *mutationResolver) TogglePostComments(ctx context.Context, postID string, editor string, enabled bool) (string, error) {
	start := time.Now()
	action := "enable"
	if !enabled {
		action = "disable"
	}
	log.Printf("Attempting to %s comments for post %s by editor %s", action, postID, editor)

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		log.Printf("Invalid post ID format: %s, error: %v", postID, err)
		return postID, fmt.Errorf("invalid post ID format")
	}

	parsedEditorId, err := uuid.Parse(editor)
	if err != nil {
		log.Printf("Invalid editor ID format: %s, error: %v", editor, err)
		return postID, fmt.Errorf("invalid editor ID format")
	}

	err = r.PostService.TogglePostComments(ctx, parsedPostID, parsedEditorId, enabled)
	if err != nil {
		log.Printf("Error toggling comments for post %s: %v", postID, err)
	} else {
		log.Printf("Successfully %sd comments for post %s in %v", action, postID, time.Since(start))
	}

	return postID, err
}

type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }
