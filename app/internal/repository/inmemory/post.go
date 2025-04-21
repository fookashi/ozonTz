package inmemory

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type PostRepo struct {
	posts map[uuid.UUID]entity.Post
	mu    sync.RWMutex
}

func NewPostRepo(initSize int) *PostRepo {
	return &PostRepo{
		posts: make(map[uuid.UUID]entity.Post, initSize),
	}
}

func (r *PostRepo) Create(ctx context.Context, post *entity.Post) error {
	if err := ctx.Err(); err != nil {
		return repository.ErrContextCanceled
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.Id] = *post
	return nil
}

func (r *PostRepo) GetOneById(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	if err := ctx.Err(); err != nil {
		return nil, repository.ErrContextCanceled
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return &post, nil
}

func (r *PostRepo) GetMany(ctx context.Context, limit, offset int, sortBy repository.SortBy) ([]entity.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit == 0 {
		return []entity.Post{}, nil
	}
	allPosts := make([]entity.Post, 0, len(r.posts))

	if err := ctx.Err(); err != nil {
		return allPosts, repository.ErrContextCanceled
	}

	for _, post := range r.posts {
		allPosts = append(allPosts, post)
	}

	switch sortBy {
	case repository.SortByNewest:
		sort.Slice(allPosts, func(i, j int) bool {
			return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
		})
	case repository.SortByOldest:
		sort.Slice(allPosts, func(i, j int) bool {
			return allPosts[i].CreatedAt.Before(allPosts[j].CreatedAt)
		})
	default:
		sort.Slice(allPosts, func(i, j int) bool {
			return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
		})
	}

	if offset >= len(allPosts) {
		return []entity.Post{}, nil
	}

	end := offset + limit
	if end > len(allPosts) {
		end = len(allPosts)
	}

	return allPosts[offset:end], nil
}

func (r *PostRepo) Update(ctx context.Context, post *entity.Post) error {
	if err := ctx.Err(); err != nil {
		return repository.ErrContextCanceled
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.Id] = *post
	return nil
}
