package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(ctx context.Context, post *entity.Post) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		INSERT INTO posts (id, user_id, title, content, is_commentable, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctxWithTimeout, query,
		post.Id, post.UserId, post.Title, post.Content, post.IsCommentable, post.CreatedAt)
	return err
}

func (r *PostRepo) Update(ctx context.Context, post *entity.Post) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		UPDATE posts
		SET title = $2, content = $3, is_commentable = $4
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctxWithTimeout, query,
		post.Id, post.Title, post.Content, post.IsCommentable)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *PostRepo) GetOneById(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var post entity.Post
	query := `
		SELECT id, user_id, title, content, is_commentable, created_at
		FROM posts
		WHERE id = $1
	`
	err := r.db.GetContext(ctxWithTimeout, &post, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &post, err
}

func (r *PostRepo) GetMany(ctx context.Context, limit, offset int, sortBy repository.SortBy) ([]entity.Post, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var posts []entity.Post
	builder := strings.Builder{}
	builder.WriteString("SELECT id, user_id, title, content, is_commentable, created_at FROM posts")

	switch sortBy {
	case repository.SortByNewest:
		builder.WriteString(" ORDER BY created_at DESC")
	case repository.SortByOldest:
		builder.WriteString(" ORDER BY created_at ASC")
	default:
		builder.WriteString(" ORDER BY created_at DESC")
	}

	builder.WriteString(" LIMIT $1 OFFSET $2")

	err := r.db.SelectContext(ctxWithTimeout, &posts, builder.String(), limit, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
