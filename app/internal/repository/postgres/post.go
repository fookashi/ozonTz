package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepo struct {
	pool *pgxpool.Pool
}

func NewPostRepo(pool *pgxpool.Pool) *PostRepo {
	return &PostRepo{pool: pool}
}

func (r *PostRepo) Create(ctx context.Context, post *entity.Post) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO posts (id, user_id, title, content, is_commentable, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		post.Id, post.UserId, post.Title, post.Content, post.IsCommentable, post.CreatedAt)
	return err
}

func (r *PostRepo) Update(ctx context.Context, post *entity.Post) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE posts
		SET title = $2, content = $3, is_commentable = $4
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query,
		post.Id, post.Title, post.Content, post.IsCommentable)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *PostRepo) GetOneById(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var post entity.Post
	query := `
		SELECT id, user_id, title, content, is_commentable, created_at
		FROM posts
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.Id, &post.UserId, &post.Title, &post.Content, &post.IsCommentable, &post.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &post, err
}

func (r *PostRepo) GetMany(ctx context.Context, limit, offset int, sortBy repository.SortBy) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

	rows, err := r.pool.Query(ctx, builder.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(
			&post.Id, &post.UserId, &post.Title, &post.Content, &post.IsCommentable, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}
