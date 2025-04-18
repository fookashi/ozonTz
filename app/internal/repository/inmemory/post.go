package inmemory

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type inMemoryPostRepo struct {
	posts map[uuid.UUID]entity.Post
	mu    sync.RWMutex
}

func NewInMemoryPostRepo(initSize int) *inMemoryPostRepo {
	return &inMemoryPostRepo{
		posts: make(map[uuid.UUID]entity.Post, initSize),
	}
}

func (r *inMemoryPostRepo) Create(ctx context.Context, post entity.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.Id] = post
	return nil
}

func (r *inMemoryPostRepo) GetOneById(ctx context.Context, id uuid.UUID) (entity.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return entity.Post{}, repository.ErrNotFound
	}
	return post, nil
}

func (r *inMemoryPostRepo) GetMany(ctx context.Context, limit, offset int, sortBy repository.SortBy) ([]entity.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allPosts := make([]entity.Post, 0, len(r.posts))

	if limit == 0 {
		return allPosts, nil
	}

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
	case repository.SortByTop:
		sort.Slice(allPosts, func(i, j int) bool {
			return len(allPosts[i].Content) > len(allPosts[j].Content)
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

func (r *inMemoryPostRepo) Update(ctx context.Context, post entity.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.Id] = post
	return nil
}
