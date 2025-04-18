package service

import (
	"app/graph/model"
	"app/internal/entity"
	"app/internal/repository"
	"context"

	"github.com/google/uuid"
)

const maxSymbolsLength int = 2000

type CommentService struct {
	RepoHolder *repository.RepoHolder
}

func (s *CommentService) CreateComment(ctx context.Context, userId uuid.UUID, postId uuid.UUID, parentId *uuid.UUID, content string) (*model.Comment, error) {
	post, err := s.RepoHolder.PostRepo.GetOneById(ctx, postId)

	if len([]rune(content)) > maxSymbolsLength {
		return nil, ErrTooManySymbols
	}

	if err != nil {
		return nil, ErrPostNotFound
	}

	if !post.IsCommentable {
		return nil, ErrPostIsNotCommentable
	}

	user, err := s.RepoHolder.UserRepo.GetOneById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if parentId != nil {
		_, err := s.RepoHolder.CommentRepo.GetOneByID(ctx, *parentId)
		if err != nil {
			return nil, ErrParentCommentNotFound
		}
	}

	newComment, err := entity.NewComment(userId, postId, parentId, content)

	if err != nil {
		return nil, err
	}

	if err := s.RepoHolder.CommentRepo.Create(ctx, *newComment); err != nil {
		return nil, err
	}

	return &model.Comment{
		ID: newComment.Id.String(),
		User: &model.User{
			ID:       userId.String(),
			Username: user.Username,
		},
		Content: content,
	}, nil
}
