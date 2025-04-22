package postgres

import (
	"app/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepoHolder(pool *pgxpool.Pool) *repository.RepoHolder {
	return &repository.RepoHolder{
		UserRepo:    NewUserRepo(pool),
		PostRepo:    NewPostRepo(pool),
		CommentRepo: NewCommentRepo(pool),
	}
}
