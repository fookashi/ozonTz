package inmemory_test

import (
	"app/internal/entity"
	"app/internal/repository"
	"app/internal/repository/inmemory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserRepo(t *testing.T) {
	user1 := entity.User{
		Id:       uuid.New(),
		Username: "user1",
		Roles:    []string{"user"},
	}
	user2 := entity.User{
		Id:       uuid.New(),
		Username: "user2",
		Roles:    []string{"admin"},
	}
	nonExistentID := uuid.New()

	tests := []struct {
		name string
		run  func(t *testing.T, repo *inmemory.UserRepo)
	}{
		{
			name: "Create/success",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				err := repo.Create(context.Background(), &user1)
				assert.NoError(t, err)
			},
		},
		{
			name: "Create/canceled context",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				err := repo.Create(ctx, &user1)
				assert.ErrorIs(t, err, repository.ErrContextCanceled)
			},
		},
		{
			name: "GetOneById/success",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				result, err := repo.GetOneById(context.Background(), user1.Id)
				assert.NoError(t, err)
				assert.Equal(t, user1, *result)
			},
		},
		{
			name: "GetOneById/not found",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_, err := repo.GetOneById(context.Background(), nonExistentID)
				assert.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
		{
			name: "GetOneById/canceled context",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := repo.GetOneById(ctx, user1.Id)
				assert.ErrorIs(t, err, repository.ErrContextCanceled)
			},
		},
		{
			name: "GetManyByIds/success all found",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				_ = repo.Create(context.Background(), &user2)
				result, err := repo.GetManyByIds(context.Background(), []uuid.UUID{user1.Id, user2.Id})
				assert.NoError(t, err)
				assert.Len(t, result, 2)
			},
		},
		{
			name: "GetManyByIds/partial found",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				result, err := repo.GetManyByIds(context.Background(), []uuid.UUID{user1.Id, nonExistentID})
				assert.ErrorIs(t, err, repository.ErrNotFound)
				assert.Len(t, result, 1)
			},
		},
		{
			name: "GetManyByIds/none found",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				result, err := repo.GetManyByIds(context.Background(), []uuid.UUID{nonExistentID})
				assert.ErrorIs(t, err, repository.ErrNotFound)
				assert.Empty(t, result)
			},
		},
		{
			name: "GetManyByIds/empty ids",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				result, err := repo.GetManyByIds(context.Background(), []uuid.UUID{})
				assert.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name: "GetManyByIds/canceled context",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := repo.GetManyByIds(ctx, []uuid.UUID{user1.Id})
				assert.ErrorIs(t, err, repository.ErrContextCanceled)
			},
		},
		{
			name: "GetOneByUsername/success",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				result, err := repo.GetOneByUsername(context.Background(), user1.Username)
				assert.NoError(t, err)
				assert.Equal(t, user1, *result)
			},
		},
		{
			name: "GetOneByUsername/not found",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_, err := repo.GetOneByUsername(context.Background(), "nonexistent")
				assert.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
		{
			name: "GetOneByUsername/canceled context",
			run: func(t *testing.T, repo *inmemory.UserRepo) {
				_ = repo.Create(context.Background(), &user1)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := repo.GetOneByUsername(ctx, user1.Username)
				assert.ErrorIs(t, err, repository.ErrContextCanceled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewUserRepo(10)
			tt.run(t, repo)
		})
	}

	t.Run("Concurrency read/write", func(t *testing.T) {
		repo := inmemory.NewUserRepo(10)
		const numWorkers = 10
		done := make(chan struct{})

		go func() {
			for i := 0; i < numWorkers; i++ {
				user := entity.User{
					Id:       uuid.New(),
					Username: "concurrent_" + string(rune(i)),
				}
				_ = repo.Create(context.Background(), &user)
			}
			close(done)
		}()

		for i := 0; i < numWorkers; i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						_, _ = repo.GetOneById(context.Background(), user1.Id)
						_, _ = repo.GetManyByIds(context.Background(), []uuid.UUID{user1.Id, user2.Id})
					}
				}
			}()
		}

		<-done
	})
}
