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

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user entity.User) error {
	query := `INSERT INTO users (id, username, roles) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, user.Id, user.Username, user.Roles)
	return err
}

func (r *UserRepo) GetOneById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	var user entity.User
	query := `SELECT id, username, roles FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, repository.ErrNotFound
	}
	return user, err
}

func (r *UserRepo) GetManyByIds(ctx context.Context, ids []uuid.UUID) ([]entity.User, error) {
	if len(ids) == 0 {
		return []entity.User{}, nil
	}

	query, args, err := sqlx.In(`SELECT id, username, roles FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	var users []entity.User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}
	if len(users) < len(ids) {
		return users, repository.ErrNotFound
	}

	return users, nil
}

func (r *UserRepo) GetOneByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, username, roles FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, &user, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, repository.ErrNotFound
	}
	return user, err
}

func (r *UserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.GetContext(ctx, &exists, query, username)
	return exists, err
}
