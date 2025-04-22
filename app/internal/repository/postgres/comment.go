package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type CommentRepo struct {
	db Database
}

func NewCommentRepo(db Database) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) GetOneById(ctx context.Context, commentId uuid.UUID) (*entity.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var comment entity.Comment
	query := `
        SELECT id, user_id, post_id, parent_id, content, created_at
        FROM comments
        WHERE id = $1
    `

	err := r.db.QueryRow(ctx, query, commentId).Scan(
		&comment.Id,
		&comment.UserId,
		&comment.PostId,
		&comment.ParentId,
		&comment.Content,
		&comment.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	query := `
		INSERT INTO comments (id, user_id, post_id, parent_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, comment.Id, comment.UserId, comment.PostId, comment.ParentId, comment.Content, comment.CreatedAt)
	return err
}

func (r *CommentRepo) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	query := `
        SELECT id, user_id, post_id, parent_id, content, created_at
        FROM comments
        WHERE post_id = $1 AND parent_id IS NULL
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, postId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(
			&comment.Id, &comment.UserId, &comment.PostId, &comment.ParentId, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *CommentRepo) GetCommentReplies(ctx context.Context, parentId uuid.UUID, limit, offset int) ([]entity.Comment, error) {
	query := `
        SELECT id, user_id, post_id, parent_id, content, created_at
        FROM comments
        WHERE parent_id = $1
        ORDER BY created_at ASC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, parentId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(
			&comment.Id, &comment.UserId, &comment.PostId, &comment.ParentId, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}
