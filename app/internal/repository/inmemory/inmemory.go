package inmemory

import (
	"app/internal/repository"
)

func NewRepoHolder(initSize int) *repository.RepoHolder {
	return &repository.RepoHolder{
		UserRepo:    NewUserRepo(initSize),
		PostRepo:    NewPostRepo(initSize),
		CommentRepo: NewCommentRepo(initSize),
	}
}
