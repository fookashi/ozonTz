package repository

import (
	"app/internal/entity"
	"context"

	"github.com/google/uuid"
)

type SortBy string

const (
	SortByNewest SortBy = "NEWEST"
	SortByOldest SortBy = "OLDEST"
	SortByTop    SortBy = "TOP"
)

//go:generate go run github.com/golang/mock/mockgen -source=repository.go -destination=mocks/repository.go

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetOneById(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetManyByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entity.User, error)
	GetOneByUsername(ctx context.Context, username string) (*entity.User, error)
}

type PostRepo interface {
	Create(ctx context.Context, post *entity.Post) error
	Update(ctx context.Context, post *entity.Post) error
	GetOneById(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	GetMany(ctx context.Context, limit, offset int, sortBy SortBy) ([]entity.Post, error)
}

type CommentRepo interface {
	Create(ctx context.Context, comment *entity.Comment) error
	GetOneByID(ctx context.Context, commentId uuid.UUID) (*entity.Comment, error)

	GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]entity.Comment, error)
	GetReplies(ctx context.Context, parentId []uuid.UUID, limit, offset int) (map[uuid.UUID][]entity.Comment, error)
}

type RepoHolder struct {
	UserRepo
	PostRepo
	CommentRepo
}
