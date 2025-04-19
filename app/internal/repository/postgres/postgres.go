package postgres

import (
	"app/internal/repository"

	"github.com/jmoiron/sqlx"
)

func NewRepoHolder(db *sqlx.DB) *repository.RepoHolder {
	return &repository.RepoHolder{
		UserRepo:    NewUserRepo(db),
		PostRepo:    NewPostRepo(db),
		CommentRepo: NewCommentRepo(db),
	}
}
