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

	if len([]rune(content)) > maxSymbolsLength {
		return nil, ErrTooManySymbols
	}
	post, err := s.RepoHolder.PostRepo.GetOneById(ctx, postId)

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
		_, err := s.RepoHolder.CommentRepo.GetOneById(ctx, *parentId)
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

	var ParentIdString *string
	if parentId != nil {
		parentStr := parentId.String()
		ParentIdString = &parentStr
	}

	return &model.Comment{
		ID:       newComment.Id.String(),
		ParentID: ParentIdString,
		User: &model.User{
			ID:       userId.String(),
			Username: user.Username,
		},
		Content:   content,
		CreatedAt: newComment.CreatedAt,
	}, nil
}

func (s *CommentService) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	commentEntities, err := s.RepoHolder.CommentRepo.GetByPost(ctx, postId, limit, offset)
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
			ParentID:  nil,
			User: &model.User{
				ID:       userEntity.Id.String(),
				Username: userEntity.Username,
			},
		})
	}

	return comments, nil
}

func (s *CommentService) GetCommentReplies(ctx context.Context, parentId uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	replyEntities, err := s.RepoHolder.CommentRepo.GetCommentReplies(ctx, parentId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	if len(replyEntities) == 0 {
		return []*model.Comment{}, nil
	}

	userIds := make([]uuid.UUID, 0, len(replyEntities))
	for _, reply := range replyEntities {
		userIds = append(userIds, reply.UserId)
	}

	users, err := s.RepoHolder.UserRepo.GetManyByIds(ctx, userIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for replies: %w", err)
	}

	usersMap := make(map[uuid.UUID]entity.User, len(users))
	for _, user := range users {
		usersMap[user.Id] = user
	}

	replies := make([]*model.Comment, 0, len(replyEntities))
	for _, replyEntity := range replyEntities {
		userEntity, exists := usersMap[replyEntity.UserId]
		if !exists {
			return nil, fmt.Errorf("user not found for reply %s", replyEntity.Id)
		}
		var replyEntityParent *string
		if replyEntity.ParentId != nil {
			parentStr := replyEntity.ParentId.String()
			replyEntityParent = &parentStr
		}
		replies = append(replies, &model.Comment{
			ID:        replyEntity.Id.String(),
			Content:   replyEntity.Content,
			ParentID:  replyEntityParent,
			CreatedAt: replyEntity.CreatedAt,
			User: &model.User{
				ID:       userEntity.Id.String(),
				Username: userEntity.Username,
			},
		})
	}

	return replies, nil
}
