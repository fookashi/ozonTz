package service

import (
	"app/internal/entity"
	"app/internal/repository"
	mock_repository "app/internal/repository/mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{UserRepo: mockUserRepo}
	service := &UserService{RepoHolder: repoHolder}

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := entity.User{
			Id:       userID,
			Username: "testuser",
		}

		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userID).
			Return(expectedUser, nil)

		result, err := service.GetUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, userID.String(), result.ID)
		assert.Equal(t, "testuser", result.Username)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New()

		mockUserRepo.EXPECT().
			GetOneById(gomock.Any(), userID).
			Return(entity.User{}, repository.ErrNotFound)

		result, err := service.GetUser(context.Background(), userID)

		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, result)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
	repoHolder := &repository.RepoHolder{UserRepo: mockUserRepo}
	service := &UserService{RepoHolder: repoHolder}

	t.Run("success", func(t *testing.T) {
		username := "newuser"

		mockUserRepo.EXPECT().
			UsernameExists(gomock.Any(), username).
			Return(false, nil)

		mockUserRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, user entity.User) {
				assert.Equal(t, username, user.Username)
			}).
			Return(nil)

		result, err := service.CreateUser(context.Background(), username)

		assert.NoError(t, err)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, username, result.Username)
	})

	t.Run("username exists", func(t *testing.T) {
		username := "existinguser"

		mockUserRepo.EXPECT().
			UsernameExists(gomock.Any(), username).
			Return(true, nil)

		result, err := service.CreateUser(context.Background(), username)

		assert.ErrorIs(t, err, ErrUsernameExists)
		assert.Nil(t, result)
	})
}
