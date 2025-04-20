package service

import (
	"app/graph/model"
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"fmt"

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

	if err := s.RepoHolder.CommentRepo.Create(ctx, newComment); err != nil {
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

func (s *CommentService) GetByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	commentEntities, err := s.RepoHolder.CommentRepo.GetByPost(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	userIds := make([]uuid.UUID, 0, len(commentEntities))
	for _, comment := range commentEntities {
		userIds = append(userIds, comment.UserId)
	}

	users, err := s.RepoHolder.UserRepo.GetManyByIds(ctx, userIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	usersMap := make(map[uuid.UUID]entity.User, len(users))
	for _, user := range users {
		usersMap[user.Id] = user
	}

	comments := make([]*model.Comment, 0, len(commentEntities))
	for _, commentEntity := range commentEntities {
		userEntity, exists := usersMap[commentEntity.UserId]
		if !exists {
			return nil, fmt.Errorf("user not found for comment %s", commentEntity.Id)
		}

		comments = append(comments, &model.Comment{
			ID:        commentEntity.Id.String(),
			Content:   commentEntity.Content,
			CreatedAt: commentEntity.CreatedAt,
			User: &model.User{
				ID:       userEntity.Id.String(),
				Username: userEntity.Username,
			},
		})
	}

	return comments, nil
}
