package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(ctx context.Context, post entity.Post) error {
	query := `
		INSERT INTO posts (id, user_id, title, content, is_commentable, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		post.Id, post.UserId, post.Title, post.Content, post.IsCommentable, post.CreatedAt)
	return err
}

func (r *PostRepo) Update(ctx context.Context, post entity.Post) error {
	query := `
		UPDATE posts
		SET title = $2, content = $3, is_commentable = $4
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
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

func (r *PostRepo) GetOneById(ctx context.Context, id uuid.UUID) (entity.Post, error) {
	var post entity.Post
	query := `
		SELECT id, user_id, title, content, is_commentable, created_at
		FROM posts
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &post, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Post{}, repository.ErrNotFound
	}
	return post, err
}

func (r *PostRepo) GetMany(ctx context.Context, limit, offset int, sortBy repository.SortBy) ([]entity.Post, error) {
	var posts []entity.Post

	baseQuery := `SELECT id, user_id, title, content, is_commentable, created_at FROM posts`

	switch sortBy {
	case repository.SortByNewest:
		baseQuery += " ORDER BY created_at DESC"
	case repository.SortByOldest:
		baseQuery += " ORDER BY created_at ASC"
	default:
		baseQuery += " ORDER BY created_at DESC"
	}

	baseQuery += " LIMIT $1 OFFSET $2"

	err := r.db.SelectContext(ctx, &posts, baseQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
