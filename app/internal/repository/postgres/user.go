package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `INSERT INTO users (id, username, roles) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctxWithTimeout, query, user.Id, user.Username, pq.Array(user.Roles))
	return err
}

func (r *UserRepo) GetOneById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user entity.User
	query := `SELECT id, username, roles FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctxWithTimeout, query, id).Scan(&user.Id, &user.Username, pq.Array(&user.Roles))
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("user not found: %v", user)
		return nil, repository.ErrNotFound
	}
	return &user, err
}

func (r *UserRepo) GetManyByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entity.User, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]entity.User{}, nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, username, roles FROM users WHERE id = ANY($1)`
	rows, err := r.db.QueryContext(ctxWithTimeout, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to get users by ids: %w", err)
	}
	defer rows.Close()

	users := make(map[uuid.UUID]entity.User, len(ids))
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Username, pq.Array(user.Roles))
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users[user.Id] = user
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (r *UserRepo) GetOneByUsername(ctx context.Context, username string) (*entity.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var user entity.User
	query := `SELECT id, username, roles FROM users WHERE username = $1`
	err := r.db.QueryRowContext(ctxWithTimeout, query, username).Scan(&user.Id, &user.Username, pq.Array(&user.Roles))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &user, err
}
