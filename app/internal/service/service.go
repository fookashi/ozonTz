package service

import (
	"app/graph/model"
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/golang/mock/mockgen -source=service.go -destination=mocks/service.go

type User interface {
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	CreateUser(ctx context.Context, username string) (*model.User, error)
}

type Post interface {
	GetPostById(ctx context.Context, id uuid.UUID) (*model.Post, error)
	GetPosts(ctx context.Context, limit, offset int, sortBy *model.SortBy) ([]*model.Post, error)
	CreatePost(ctx context.Context, userId uuid.UUID, title string, content string, isCommentable bool) (*model.Post, error)
	TogglePostComments(ctx context.Context, postId uuid.UUID, editorId uuid.UUID, enabled bool) error
}

type Comment interface {
	CreateComment(ctx context.Context, userId uuid.UUID, postId uuid.UUID, parentId *uuid.UUID, content string) (*model.Comment, error)
	GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*model.Comment, error)
	GetCommentReplies(ctx context.Context, parentId uuid.UUID, limit, offset int) ([]*model.Comment, error)
}

type Services struct {
	Comment
	Post
	User
}
