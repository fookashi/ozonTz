package service

import (
	"app/graph/model"
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type IPostService interface {
	GetPostById(ctx context.Context, id uuid.UUID) (*model.Post, error)
	GetPosts(ctx context.Context, limit, offset int, sortBy *model.SortBy) ([]*model.Post, error)
	CreatePost(ctx context.Context, userId uuid.UUID, title string, content string, isCommentable bool) (*model.Post, error)
	TogglePostComments(ctx context.Context, postId uuid.UUID, editorId uuid.UUID, enabled bool) error
}

type PostService struct {
	RepoHolder *repository.RepoHolder
}

func (s *PostService) GetPostById(ctx context.Context, id uuid.UUID) (*model.Post, error) {

	newPost, err := s.RepoHolder.PostRepo.GetOneById(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	user, err := s.RepoHolder.UserRepo.GetOneById(ctx, newPost.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &model.Post{
		ID:            newPost.Id.String(),
		Title:         newPost.Title,
		Content:       newPost.Content,
		IsCommentable: newPost.IsCommentable,
		CreatedAt:     newPost.CreatedAt,
		User: &model.User{
			ID:       user.Id.String(),
			Username: user.Username,
		},
	}, nil
}

func (s *PostService) GetPosts(ctx context.Context, limit, offset int, sortBy *model.SortBy) ([]*model.Post, error) {
	rSortBy := repository.SortByNewest
	if sortBy != nil {
		rSortBy = repository.SortBy(*sortBy)
	}

	postEntities, err := s.RepoHolder.PostRepo.GetMany(
		ctx,
		limit,
		offset,
		rSortBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	userIDs := make([]uuid.UUID, 0, len(postEntities))
	for _, post := range postEntities {
		userIDs = append(userIDs, post.UserId)
	}

	users, err := s.RepoHolder.UserRepo.GetManyByIds(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	usersMap := make(map[uuid.UUID]entity.User, len(users))
	for _, user := range users {
		usersMap[user.Id] = user
	}

	posts := make([]*model.Post, 0, len(postEntities))
	for _, postEntity := range postEntities {
		userEntity, exists := usersMap[postEntity.UserId]
		if !exists {
			return nil, fmt.Errorf("user not found for post %s", postEntity.Id)
		}

		posts = append(posts, &model.Post{
			ID:            postEntity.Id.String(),
			Title:         postEntity.Title,
			Content:       postEntity.Content,
			IsCommentable: postEntity.IsCommentable,
			CreatedAt:     postEntity.CreatedAt,
			User: &model.User{
				ID:       userEntity.Id.String(),
				Username: userEntity.Username,
			},
		})
	}

	return posts, nil
}

func (s *PostService) CreatePost(ctx context.Context, userId uuid.UUID, title string, content string, isCommentable bool) (*model.Post, error) {
	newPost, err := entity.NewPost(userId, title, content, isCommentable)

	if err != nil {
		return nil, err
	}

	user, err := s.RepoHolder.UserRepo.GetOneById(ctx, userId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := s.RepoHolder.PostRepo.Create(ctx, *newPost); err != nil {
		return nil, err
	}

	return &model.Post{
		ID:            newPost.Id.String(),
		Title:         newPost.Title,
		Content:       newPost.Content,
		IsCommentable: newPost.IsCommentable,
		CreatedAt:     newPost.CreatedAt,
		User: &model.User{
			ID:       user.Id.String(),
			Username: user.Username,
		},
	}, nil
}

func (s *PostService) TogglePostComments(ctx context.Context, postId uuid.UUID, editorId uuid.UUID, enabled bool) error {
	post, err := s.RepoHolder.PostRepo.GetOneById(ctx, postId)
	if err != nil {
		return ErrPostNotFound
	}
	if post.UserId != editorId {
		return ErrNoPermissionForToggle
	}

	post.IsCommentable = enabled
	if err := s.RepoHolder.PostRepo.Update(ctx, post); err != nil {
		return err
	}

	return nil
}
