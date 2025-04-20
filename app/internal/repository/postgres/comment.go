package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	query := `
		INSERT INTO comments (id, user_id, post_id, parent_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		comment.Id, comment.UserId, comment.PostId, comment.ParentId, comment.Content, comment.CreatedAt)
	return err
}

func (r *CommentRepo) GetOneByID(ctx context.Context, commentId uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	query := `
		SELECT id, user_id, post_id, parent_id, content, created_at
		FROM comments
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &comment, query, commentId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &comment, err
}

func (r *CommentRepo) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	query := `
        SELECT id, user_id, post_id, parent_id, content, created_at
        FROM comments
        WHERE post_id = $1 AND parent_id IS NULL
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	var comments []entity.Comment
	err := r.db.SelectContext(ctx, &comments, query, postId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by posts: %w", err)
	}

	return comments, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, parentIds []uuid.UUID, limit, offset int) (map[uuid.UUID][]entity.Comment, error) {
	if len(parentIds) == 0 {
		return map[uuid.UUID][]entity.Comment{}, nil
	}

	query := `
        SELECT id, user_id, post_id, parent_id, content, created_at
        FROM comments
        WHERE parent_id = ANY($1)
        ORDER BY created_at ASC
        LIMIT $2 OFFSET $3
    `

	var comments []entity.Comment
	err := r.db.SelectContext(ctx, &comments, query, pq.Array(parentIds), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	replies := make(map[uuid.UUID][]entity.Comment)
	for _, comment := range comments {
		replies[*comment.ParentId] = append(replies[*comment.ParentId], comment)
	}

	return replies, nil
}
