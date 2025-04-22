package service_test

import (
	"app/internal/entity"
	"app/internal/repository"
	mock_repository "app/internal/repository/mocks"
	"app/internal/service"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommentService_CreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setup := func() (*service.CommentService, *mock_repository.MockPostRepo, *mock_repository.MockUserRepo, *mock_repository.MockCommentRepo) {
		mockPostRepo := mock_repository.NewMockPostRepo(ctrl)
		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockCommentRepo := mock_repository.NewMockCommentRepo(ctrl)

		repoHolder := &repository.RepoHolder{
			PostRepo:    mockPostRepo,
			UserRepo:    mockUserRepo,
			CommentRepo: mockCommentRepo,
		}

		return &service.CommentService{RepoHolder: repoHolder}, mockPostRepo, mockUserRepo, mockCommentRepo
	}

	userID := uuid.New()
	postID := uuid.New()
	parentID := uuid.New()
	content := "Test comment"

	t.Run("success with parent", func(t *testing.T) {
		service, mockPostRepo, mockUserRepo, mockCommentRepo := setup()

		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postID).
			Return(&entity.Post{
				Id:            postID,
				IsCommentable: true,
			}, nil)

		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userID).
			Return(&entity.User{
				Id:       userID,
				Username: "testuser",
			}, nil)

		mockCommentRepo.EXPECT().
			GetOneById(gomock.Any(), parentID).
			Return(&entity.Comment{}, nil)

		mockCommentRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil)

		result, err := service.CreateComment(context.Background(), userID, postID, &parentID, content)

		assert.NoError(t, err)
		assert.Equal(t, content, result.Content)
		assert.Equal(t, parentID.String(), *result.ParentID)
		assert.Equal(t, userID.String(), result.User.ID)
	})

	t.Run("success without parent", func(t *testing.T) {
		service, mockPostRepo, mockUserRepo, mockCommentRepo := setup()

		mockPostRepo.EXPECT().
			GetOneById(gomock.Any(), postID).
			Return(&entity.Post{
				Id:            postID,
				IsCommentable: true,
			}, nil)

		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userID).
			Return(&entity.User{
				Id:       userID,
				Username: "testuser",
			}, nil)

		mockCommentRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil)

		result, err := service.CreateComment(context.Background(), userID, postID, nil, content)

		assert.NoError(t, err)
		assert.Nil(t, result.ParentID)
	})

	t.Run("too many symbols", func(t *testing.T) {
		cService, _, _, _ := setup()

		longContent := make([]rune, 3000) // magic value
		for i := range longContent {
			longContent[i] = 'a'
		}

		result, err := cService.CreateComment(context.Background(), userID, postID, nil, string(longContent))

		assert.ErrorIs(t, err, service.ErrTooManySymbols)
		assert.Nil(t, result)
	})

}

func TestCommentService_GetByPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setup := func() (*service.CommentService, *mock_repository.MockCommentRepo, *mock_repository.MockUserRepo) {
		mockCommentRepo := mock_repository.NewMockCommentRepo(ctrl)
		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)

		repoHolder := &repository.RepoHolder{
			CommentRepo: mockCommentRepo,
			UserRepo:    mockUserRepo,
		}

		return &service.CommentService{RepoHolder: repoHolder}, mockCommentRepo, mockUserRepo
	}

	postID := uuid.New()
	userID := uuid.New()
	limit := 10
	offset := 0

	t.Run("success", func(t *testing.T) {
		service, mockCommentRepo, mockUserRepo := setup()

		commentEntities := []entity.Comment{
			{
				Id:        uuid.New(),
				UserId:    userID,
				PostId:    postID,
				Content:   "Comment 1",
				CreatedAt: time.Now(),
			},
		}

		mockCommentRepo.EXPECT().
			GetByPost(gomock.Any(), postID, limit, offset).
			Return(commentEntities, nil)

		mockUserRepo.EXPECT().
			GetManyByIds(gomock.Any(), []uuid.UUID{userID}).
			Return(map[uuid.UUID]entity.User{
				userID: {
					Id:       userID,
					Username: "testuser",
				},
			}, nil)

		result, err := service.GetByPost(context.Background(), postID, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Comment 1", result[0].Content)
	})
}

func TestCommentService_GetCommentReplies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setup := func() (*service.CommentService, *mock_repository.MockCommentRepo, *mock_repository.MockUserRepo) {
		mockCommentRepo := mock_repository.NewMockCommentRepo(ctrl)
		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)

		repoHolder := &repository.RepoHolder{
			CommentRepo: mockCommentRepo,
			UserRepo:    mockUserRepo,
		}

		return &service.CommentService{RepoHolder: repoHolder}, mockCommentRepo, mockUserRepo
	}

	parentID := uuid.New()
	userID := uuid.New()
	limit := 10
	offset := 0

	t.Run("success", func(t *testing.T) {
		service, mockCommentRepo, mockUserRepo := setup()

		replyEntities := []entity.Comment{
			{
				Id:        uuid.New(),
				UserId:    userID,
				ParentId:  &parentID,
				Content:   "Reply 1",
				CreatedAt: time.Now(),
			},
		}

		mockCommentRepo.EXPECT().
			GetCommentReplies(gomock.Any(), parentID, limit, offset).
			Return(replyEntities, nil)

		mockUserRepo.EXPECT().
			GetManyByIds(gomock.Any(), []uuid.UUID{userID}).
			Return(map[uuid.UUID]entity.User{
				userID: {
					Id:       userID,
					Username: "testuser",
				},
			}, nil)

		result, err := service.GetCommentReplies(context.Background(), parentID, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Reply 1", result[0].Content)
	})
}
