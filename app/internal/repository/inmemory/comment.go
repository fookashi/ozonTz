package inmemory

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"sync"

	"github.com/google/uuid"
)

type CommentRepo struct {
	comments     map[uuid.UUID]entity.Comment
	postIndex    map[uuid.UUID][]uuid.UUID
	repliesIndex map[uuid.UUID][]uuid.UUID

	mu sync.RWMutex
}

func NewCommentRepo(initSize int) *CommentRepo {
	return &CommentRepo{
		comments:     make(map[uuid.UUID]entity.Comment, initSize),
		postIndex:    make(map[uuid.UUID][]uuid.UUID, initSize),
		repliesIndex: make(map[uuid.UUID][]uuid.UUID, initSize),
	}
}

func (r *CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	if err := ctx.Err(); err != nil {
		return repository.ErrContextCanceled
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.comments[comment.Id] = *comment

	r.postIndex[comment.PostId] = append(r.postIndex[comment.PostId], comment.Id)

	if comment.ParentId != nil {
		r.repliesIndex[*comment.ParentId] = append(r.repliesIndex[*comment.ParentId], comment.Id)
	}

	return nil
}

func (r *CommentRepo) GetOneByID(ctx context.Context, commentID uuid.UUID) (*entity.Comment, error) {
	if err := ctx.Err(); err != nil {
		return nil, repository.ErrContextCanceled
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[commentID]
	if !exists {
		return nil, repository.ErrNotFound
	}

	return &comment, nil
}

func (r *CommentRepo) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	if err := ctx.Err(); err != nil {
		return []entity.Comment{}, repository.ErrContextCanceled
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	commentIds, exists := r.postIndex[postId]
	if !exists {
		return []entity.Comment{}, repository.ErrNotFound
	}
	result := make([]entity.Comment, len(commentIds))
	for ind, id := range commentIds {
		if ind < offset {
			continue
		}
		if ind > limit {
			break
		}
		result = append(result, r.comments[id])
	}

	return result, nil
}

func (r *CommentRepo) GetCommentReplies(ctx context.Context, parentId uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	if err := ctx.Err(); err != nil {
		return []entity.Comment{}, repository.ErrContextCanceled
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, exists := r.comments[parentId]; !exists {
		return nil, repository.ErrNotFound
	}

	replyIDs, exists := r.repliesIndex[parentId]
	if !exists || len(replyIDs) == 0 {
		return []entity.Comment{}, nil
	}

	start := offset
	end := offset + limit
	if start > len(replyIDs) {
		return []entity.Comment{}, nil
	}
	if end > len(replyIDs) {
		end = len(replyIDs)
	}
	paginatedIDs := replyIDs[start:end]

	result := make([]entity.Comment, 0, len(paginatedIDs))
	for _, id := range paginatedIDs {
		if comment, exists := r.comments[id]; exists {
			result = append(result, comment)
		}
	}

	return result, nil
}
