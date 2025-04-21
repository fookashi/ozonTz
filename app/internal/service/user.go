package service

import (
	"app/graph/model"
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type UserService struct {
	RepoHolder *repository.RepoHolder
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.RepoHolder.UserRepo.GetOneById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return &model.User{
		ID:       user.Id.String(),
		Username: user.Username,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, username string) (*model.User, error) {
	existingUser, err := s.RepoHolder.UserRepo.GetOneByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			// pass
		default:
			log.Fatalf("%v", err)
			return nil, err
		}
	}

	if existingUser != nil {
		return nil, ErrUsernameExists
	}

	newUser, err := entity.NewUser(username)
	if err != nil {
		return nil, fmt.Errorf("Validation error: %w", err)
	}

	if err := s.RepoHolder.UserRepo.Create(ctx, newUser); err != nil {
		log.Fatalf("%v", err)
		return nil, ErrDueUserCreation
	}

	return &model.User{
		ID:       newUser.Id.String(),
		Username: newUser.Username,
	}, nil
}
