package service

import (
	"app/internal/entity"
	"app/internal/repository"
	mock_repository "app/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostService_GetPostById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo, UserRepo: mockUserRepo}
	service := &PostService{RepoHolder: repoHolder}
	postId := uuid.New()
	userId := uuid.New()
	expectedUser := entity.User{
		Id:       userId,
		Username: "User",
	}
	expectedPost := entity.Post{
		Id:            postId,
		UserId:        userId,
		Title:         "Post",
		Content:       "Cool",
		IsCommentable: true,
	}
	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetOneById(gomock.Any(), userId).Return(expectedUser, nil)
		mockPostRepo.EXPECT().GetOneById(gomock.Any(), postId).Return(expectedPost, nil)

		result, err := service.GetPostById(context.Background(), postId)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Id.String(), result.User.ID)
		assert.Equal(t, expectedPost.Id.String(), result.ID)
	})
	t.Run("post not exists", func(t *testing.T) {
		postId := uuid.New()
		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postId).
			Return(entity.Post{}, repository.ErrNotFound)
		_, err := service.GetPostById(context.Background(), postId)
		assert.ErrorIs(t, err, ErrPostNotFound)
	})
}

func TestPostService_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo, UserRepo: mockUserRepo}
	service := &PostService{RepoHolder: repoHolder}

	userId := uuid.New()
	expectedUser := entity.User{
		Id:       userId,
		Username: "User",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userId).
			Return(expectedUser, nil)

		mockPostRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, post entity.Post) {
				assert.Equal(t, "Post", post.Title)
				assert.Equal(t, "Cool", post.Content)
				assert.Equal(t, userId, post.UserId)
				assert.True(t, post.IsCommentable)
			}).
			Return(nil)

		result, err := service.CreatePost(context.Background(), userId, "Post", "Cool", true)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.Id.String(), result.User.ID)
		assert.Equal(t, expectedUser.Username, result.User.Username)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userId).
			Return(entity.User{}, repository.ErrNotFound)

		result, err := service.CreatePost(context.Background(), userId, "Post", "Content", false)

		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, result)
	})

	t.Run("user repo error", func(t *testing.T) {

		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userId).
			Return(entity.User{}, ErrUserNotFound)

		result, err := service.CreatePost(context.Background(), userId, "Post", "Content", false)

		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, result)
	})

	t.Run("post creation error", func(t *testing.T) {
		expectedErr := errors.New("post creation failed")

		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userId).
			Return(expectedUser, nil)

		mockPostRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		result, err := service.CreatePost(context.Background(), userId, "Post", "Content", false)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, result)
	})

	t.Run("empty title", func(t *testing.T) {
		result, err := service.CreatePost(context.Background(), userId, "", "Content", false)

		assert.ErrorIs(t, err, entity.ErrEmptyTitle)
		assert.Nil(t, result)
	})

	t.Run("empty content", func(t *testing.T) {
		result, err := service.CreatePost(context.Background(), userId, "Title", "", false)

		assert.ErrorIs(t, err, entity.ErrEmptyContent)
		assert.Nil(t, result)
	})
}

func TestPostService_TogglePostComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo}
	service := &PostService{RepoHolder: repoHolder}

	postId := uuid.New()
	ownerId := uuid.New()
	editorId := uuid.New()
	otherUserId := uuid.New()

	existingPost := entity.Post{
		Id:            postId,
		UserId:        ownerId,
		Title:         "Test Post",
		Content:       "Content",
		IsCommentable: true,
	}

	t.Run("success - enable comments", func(t *testing.T) {
		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postId).
			Return(existingPost, nil)

		updatedPost := existingPost
		updatedPost.IsCommentable = true

		mockPostRepo.EXPECT().
			Update(gomock.Any(), updatedPost).
			Return(nil)

		err := service.TogglePostComments(context.Background(), postId, ownerId, true)

		assert.NoError(t, err)
	})

	t.Run("success - disable comments", func(t *testing.T) {
		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postId).
			Return(existingPost, nil)

		updatedPost := existingPost
		updatedPost.IsCommentable = false

		mockPostRepo.EXPECT().
			Update(gomock.Any(), updatedPost).
			Return(nil)

		err := service.TogglePostComments(context.Background(), postId, ownerId, false)

		assert.NoError(t, err)
	})

	t.Run("post not found", func(t *testing.T) {
		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postId).
			Return(entity.Post{}, repository.ErrNotFound)

		err := service.TogglePostComments(context.Background(), postId, editorId, true)

		assert.ErrorIs(t, err, ErrPostNotFound)
	})

	t.Run("no permission - not post owner", func(t *testing.T) {
		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postId).
			Return(existingPost, nil)

		err := service.TogglePostComments(context.Background(), postId, otherUserId, true)

		assert.ErrorIs(t, err, ErrNoPermissionForToggle)
	})

	t.Run("update error", func(t *testing.T) {
		expectedErr := errors.New("update failed")

		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postId).
			Return(existingPost, nil)

		updatedPost := existingPost
		updatedPost.IsCommentable = false

		mockPostRepo.EXPECT().
			Update(gomock.Any(), updatedPost).
			Return(expectedErr)

		err := service.TogglePostComments(context.Background(), postId, ownerId, false)

		assert.ErrorIs(t, err, expectedErr)
	})
}
