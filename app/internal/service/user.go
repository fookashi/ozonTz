package service

import (
	"app/graph/model"
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"log"

	"github.com/google/uuid"
)

type UserService struct {
	RepoHolder *repository.RepoHolder
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.RepoHolder.UserRepo.GetOneById(ctx, id)
	if err != nil {
		log.Fatal(err)
		return nil, ErrUserNotFound
	}
	return &model.User{
		ID:       user.Id.String(),
		Username: user.Username,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, username string) (*model.User, error) {
	if exists, _ := s.RepoHolder.UserRepo.UsernameExists(ctx, username); exists {
		return nil, ErrUsernameExists
	}
	newUser, err := entity.NewUser(username)

	if err != nil {
		return nil, err
	}

	err = s.RepoHolder.UserRepo.Create(ctx, *newUser)
	if err != nil {
		log.Fatal(err)
		return nil, ErrDueUserCreation
	}

	return &model.User{
		ID:       newUser.Id.String(),
		Username: newUser.Username,
	}, nil
}
