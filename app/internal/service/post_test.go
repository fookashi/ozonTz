package service_test

import (
	"app/internal/entity"
	"app/internal/repository"
	mock_repository "app/internal/repository/mocks"
	"app/internal/service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostService_GetPostById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo, UserRepo: mockUserRepo}
	postService := &service.PostService{RepoHolder: repoHolder}

	ctx := context.Background()
	postId := uuid.New()
	userId := uuid.New()

	t.Run("success", func(t *testing.T) {
		post := &entity.Post{
			Id:            postId,
			UserId:        userId,
			Title:         "Test Post",
			Content:       "Content",
			IsCommentable: true,
			CreatedAt:     time.Now(),
		}
		user := &entity.User{
			Id:       userId,
			Username: "testuser",
		}

		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(post, nil)
		mockUserRepo.EXPECT().GetOneById(ctx, userId).Return(user, nil)

		result, err := postService.GetPostById(ctx, postId)
		require.NoError(t, err)
		assert.Equal(t, post.Id.String(), result.ID)
		assert.Equal(t, post.Title, result.Title)
		assert.Equal(t, user.Username, result.User.Username)
	})

	t.Run("post not found", func(t *testing.T) {
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(nil, repository.ErrNotFound)

		_, err := postService.GetPostById(ctx, postId)
		assert.ErrorIs(t, err, service.ErrPostNotFound)
	})

	t.Run("user not found", func(t *testing.T) {
		post := &entity.Post{
			Id:     postId,
			UserId: userId,
		}
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(post, nil)
		mockUserRepo.EXPECT().GetOneById(ctx, userId).Return(nil, repository.ErrNotFound)

		_, err := postService.GetPostById(ctx, postId)
		assert.ErrorIs(t, err, service.ErrUserNotFound)
	})
}

func TestPostService_GetPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo, UserRepo: mockUserRepo}
	postService := &service.PostService{RepoHolder: repoHolder}

	ctx := context.Background()
	limit := 10
	offset := 0
	postId1 := uuid.New()
	postId2 := uuid.New()
	userId1 := uuid.New()
	userId2 := uuid.New()

	t.Run("success", func(t *testing.T) {
		posts := []entity.Post{
			{
				Id:            postId1,
				UserId:        userId1,
				Title:         "Post 1",
				Content:       "Content 1",
				IsCommentable: true,
				CreatedAt:     time.Now(),
			},
			{
				Id:            postId2,
				UserId:        userId2,
				Title:         "Post 2",
				Content:       "Content 2",
				IsCommentable: false,
				CreatedAt:     time.Now().Add(-time.Hour),
			},
		}
		users := map[uuid.UUID]entity.User{
			userId1: {Id: userId1, Username: "user1"},
			userId2: {Id: userId2, Username: "user2"},
		}

		mockPostRepo.EXPECT().GetMany(ctx, limit, offset, repository.SortByNewest).Return(posts, nil)
		mockUserRepo.EXPECT().GetManyByIds(ctx, []uuid.UUID{userId1, userId2}).Return(users, nil)

		result, err := postService.GetPosts(ctx, limit, offset, nil)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, postId1.String(), result[0].ID)
		assert.Equal(t, "user1", result[0].User.Username)
		assert.Equal(t, postId2.String(), result[1].ID)
		assert.Equal(t, "user2", result[1].User.Username)
	})

	t.Run("empty result", func(t *testing.T) {
		mockPostRepo.EXPECT().
			GetMany(ctx, limit, offset, repository.SortByNewest).
			Return([]entity.Post{}, nil)

		mockUserRepo.EXPECT().
			GetManyByIds(ctx, []uuid.UUID{}).
			Return(map[uuid.UUID]entity.User{}, nil)

		result, err := postService.GetPosts(ctx, limit, offset, nil)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
	t.Run("post repo error", func(t *testing.T) {
		expectedErr := errors.New("post repo error")
		mockPostRepo.EXPECT().GetMany(ctx, limit, offset, repository.SortByNewest).Return(nil, expectedErr)

		_, err := postService.GetPosts(ctx, limit, offset, nil)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("user repo error", func(t *testing.T) {
		posts := []entity.Post{{Id: postId1, UserId: userId1}}
		mockPostRepo.EXPECT().GetMany(ctx, limit, offset, repository.SortByNewest).Return(posts, nil)
		mockUserRepo.EXPECT().GetManyByIds(ctx, []uuid.UUID{userId1}).Return(nil, errors.New("user repo error"))

		_, err := postService.GetPosts(ctx, limit, offset, nil)
		assert.Error(t, err)
	})

	t.Run("missing user", func(t *testing.T) {
		posts := []entity.Post{{Id: postId1, UserId: userId1}}
		mockPostRepo.EXPECT().GetMany(ctx, limit, offset, repository.SortByNewest).Return(posts, nil)
		mockUserRepo.EXPECT().GetManyByIds(ctx, []uuid.UUID{userId1}).Return(map[uuid.UUID]entity.User{}, nil)

		_, err := postService.GetPosts(ctx, limit, offset, nil)
		assert.ErrorIs(t, err, service.ErrUserNotFound)
	})
}

func TestPostService_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo, UserRepo: mockUserRepo}
	postService := &service.PostService{RepoHolder: repoHolder}

	ctx := context.Background()
	userId := uuid.New()
	title := "Test Post"
	content := "Test Content"
	isCommentable := true

	t.Run("success", func(t *testing.T) {
		user := &entity.User{Id: userId, Username: "testuser"}
		mockUserRepo.EXPECT().GetOneById(ctx, userId).Return(user, nil)
		mockPostRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
			func(_ context.Context, post *entity.Post) error {
				assert.Equal(t, userId, post.UserId)
				assert.Equal(t, title, post.Title)
				assert.Equal(t, content, post.Content)
				assert.Equal(t, isCommentable, post.IsCommentable)
				return nil
			})

		result, err := postService.CreatePost(ctx, userId, title, content, isCommentable)
		require.NoError(t, err)
		assert.Equal(t, title, result.Title)
		assert.Equal(t, user.Username, result.User.Username)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetOneById(ctx, userId).Return(nil, repository.ErrNotFound)

		_, err := postService.CreatePost(ctx, userId, title, content, isCommentable)
		assert.ErrorIs(t, err, service.ErrUserNotFound)
	})

	t.Run("invalid post data", func(t *testing.T) {
		_, err := postService.CreatePost(ctx, userId, "", content, isCommentable)
		assert.ErrorIs(t, err, entity.ErrEmptyTitle)

		_, err = postService.CreatePost(ctx, userId, title, "", isCommentable)
		assert.ErrorIs(t, err, entity.ErrEmptyContent)
	})

	t.Run("post creation error", func(t *testing.T) {
		user := &entity.User{Id: userId}
		expectedErr := errors.New("creation error")
		mockUserRepo.EXPECT().GetOneById(ctx, userId).Return(user, nil)
		mockPostRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedErr)

		_, err := postService.CreatePost(ctx, userId, title, content, isCommentable)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestPostService_TogglePostComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
	repoHolder := &repository.RepoHolder{PostRepo: mockPostRepo}
	postService := &service.PostService{RepoHolder: repoHolder}

	ctx := context.Background()
	postId := uuid.New()
	ownerId := uuid.New()
	otherUserId := uuid.New()

	t.Run("success enable", func(t *testing.T) {
		post := &entity.Post{
			Id:            postId,
			UserId:        ownerId,
			IsCommentable: false,
		}
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(post, nil)
		mockPostRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(
			func(_ context.Context, p *entity.Post) error {
				assert.True(t, p.IsCommentable)
				return nil
			})

		err := postService.TogglePostComments(ctx, postId, ownerId, true)
		assert.NoError(t, err)
	})

	t.Run("success disable", func(t *testing.T) {
		post := &entity.Post{
			Id:            postId,
			UserId:        ownerId,
			IsCommentable: true,
		}
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(post, nil)
		mockPostRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(
			func(_ context.Context, p *entity.Post) error {
				assert.False(t, p.IsCommentable)
				return nil
			})

		err := postService.TogglePostComments(ctx, postId, ownerId, false)
		assert.NoError(t, err)
	})

	t.Run("post not found", func(t *testing.T) {
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(nil, repository.ErrNotFound)

		err := postService.TogglePostComments(ctx, postId, ownerId, true)
		assert.ErrorIs(t, err, service.ErrPostNotFound)
	})

	t.Run("no permission", func(t *testing.T) {
		post := &entity.Post{
			Id:     postId,
			UserId: ownerId,
		}
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(post, nil)

		err := postService.TogglePostComments(ctx, postId, otherUserId, true)
		assert.ErrorIs(t, err, service.ErrNoPermissionForToggle)
	})

	t.Run("update error", func(t *testing.T) {
		post := &entity.Post{
			Id:     postId,
			UserId: ownerId,
		}
		expectedErr := errors.New("update error")
		mockPostRepo.EXPECT().GetOneById(ctx, postId).Return(post, nil)
		mockPostRepo.EXPECT().Update(ctx, gomock.Any()).Return(expectedErr)

		err := postService.TogglePostComments(ctx, postId, ownerId, true)
		assert.ErrorIs(t, err, expectedErr)
	})
}
