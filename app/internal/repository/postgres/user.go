package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (id, username, roles) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, user.Id, user.Username, user.Roles)
	return err
}

func (r *UserRepo) GetOneById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user entity.User
	query := `SELECT id, username, roles FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.Id, &user.Username, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &user, err
}

func (r *UserRepo) GetManyByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entity.User, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]entity.User{}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, username, roles FROM users WHERE id = ANY($1)`
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[uuid.UUID]entity.User)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Roles); err != nil {
			return nil, err
		}
		users[user.Id] = user
	}

	return users, rows.Err()
}

func (r *UserRepo) GetOneByUsername(ctx context.Context, username string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user entity.User
	query := `SELECT id, username, roles FROM users WHERE username = $1`
	err := r.pool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &user, err
}
