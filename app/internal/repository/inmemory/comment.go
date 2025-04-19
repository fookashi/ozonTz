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

func (r *CommentRepo) Create(ctx context.Context, comment entity.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.comments[comment.Id] = comment

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

func (r *CommentRepo) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commentIDs, exists := r.postIndex[postId]
	if !exists {
		return nil, repository.ErrNotFound
	}

	start := offset
	end := start + limit
	if end > len(commentIDs) {
		end = len(commentIDs)
	}

	var result []*entity.Comment
	for _, id := range commentIDs[start:end] {
		comment := r.comments[id]
		result = append(result, &comment)
	}

	return result, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, parentId uuid.UUID, limit, offset int) ([]*entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit == 0 {
		return make([]*entity.Comment, 0), nil
	}

	_, exists := r.comments[parentId]
	if !exists {
		return nil, repository.ErrNotFound
	}

	replyIds, exists := r.repliesIndex[parentId]
	if !exists {
		return nil, repository.ErrNotFound
	}

	start := offset
	end := start + limit
	if end > len(replyIds) {
		end = len(replyIds)
	}

	var result = make([]*entity.Comment, 0, end-start+1)
	for _, id := range replyIds[start:end] {
		comment := r.comments[id]
		result = append(result, &comment)
	}

	return result, nil
}
