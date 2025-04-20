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
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[commentID]
	if !exists {
		return nil, repository.ErrNotFound
	}

	return &comment, nil
}

func (r *CommentRepo) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]entity.Comment, error) {
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

func (r *CommentRepo) GetReplies(ctx context.Context, parentIds []uuid.UUID, limit, offset int) (map[uuid.UUID][]entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(parentIds) == 0 {
		return map[uuid.UUID][]entity.Comment{}, nil
	}

	result := make(map[uuid.UUID][]entity.Comment)

	for _, parentId := range parentIds {
		if _, exists := r.comments[parentId]; !exists {
			continue
		}

		replyIds, exists := r.repliesIndex[parentId]
		if !exists || len(replyIds) == 0 {
			result[parentId] = []entity.Comment{}
			continue
		}

		start := offset
		if start > len(replyIds) {
			start = len(replyIds)
		}
		end := start + limit
		if end > len(replyIds) {
			end = len(replyIds)
		}

		replies := make([]entity.Comment, 0, end-start)
		for _, id := range replyIds[start:end] {
			if comment, exists := r.comments[id]; exists {
				replies = append(replies, comment)
			}
		}

		if len(replies) > 0 {
			result[parentId] = replies
		} else {
			result[parentId] = []entity.Comment{}
		}
	}

	if len(result) != len(parentIds) {
		return nil, repository.ErrPartiallyFound
	}

	return result, nil
}
