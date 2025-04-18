package inmemory

import (
	"app/internal/repository"
)

func NewInMemoryRepoHolder(initSize int) *repository.RepoHolder {
	return &repository.RepoHolder{
		UserRepo:    NewInMemoryUserRepo(initSize),
		PostRepo:    NewInMemoryPostRepo(initSize),
		CommentRepo: NewInMemoryCommentRepo(initSize),
	}
}
