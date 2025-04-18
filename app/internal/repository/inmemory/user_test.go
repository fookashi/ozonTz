package inmemory

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserRepo(t *testing.T) {
	repo := NewInMemoryUserRepo(10)
	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

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

	t.Run("Create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			err := repo.Create(ctx, user1)
			assert.NoError(t, err)
		})

		t.Run("canceled context", func(t *testing.T) {
			err := repo.Create(canceledCtx, user2)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("GetOneById", func(t *testing.T) {
		// Setup
		_ = repo.Create(ctx, user1)

		t.Run("success", func(t *testing.T) {
			result, err := repo.GetOneById(ctx, user1.Id)
			assert.NoError(t, err)
			assert.Equal(t, user1, result)
		})

		t.Run("not found", func(t *testing.T) {
			_, err := repo.GetOneById(ctx, uuid.New())
			assert.ErrorIs(t, err, repository.ErrNotFound)
		})

		t.Run("canceled context", func(t *testing.T) {
			_, err := repo.GetOneById(canceledCtx, user1.Id)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("GetManyByIds", func(t *testing.T) {
		// Setup
		_ = repo.Create(ctx, user2)
		nonExistentID := uuid.New()

		t.Run("success - all found", func(t *testing.T) {
			ids := []uuid.UUID{user1.Id, user2.Id}
			result, err := repo.GetManyByIds(ctx, ids)
			assert.NoError(t, err)
			assert.Len(t, result, 2)
			assert.Contains(t, result, user1)
			assert.Contains(t, result, user2)
		})

		t.Run("partial found", func(t *testing.T) {
			ids := []uuid.UUID{user1.Id, nonExistentID}
			result, err := repo.GetManyByIds(ctx, ids)
			assert.ErrorIs(t, err, repository.ErrNotFound)
			assert.Len(t, result, 1)
			assert.Contains(t, result, user1)
		})

		t.Run("none found", func(t *testing.T) {
			ids := []uuid.UUID{nonExistentID}
			result, err := repo.GetManyByIds(ctx, ids)
			assert.ErrorIs(t, err, repository.ErrNotFound)
			assert.Empty(t, result)
		})

		t.Run("empty ids", func(t *testing.T) {
			result, err := repo.GetManyByIds(ctx, []uuid.UUID{})
			assert.NoError(t, err)
			assert.Empty(t, result)
		})

		t.Run("canceled context", func(t *testing.T) {
			_, err := repo.GetManyByIds(canceledCtx, []uuid.UUID{user1.Id})
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("GetOneByUsername", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result, err := repo.GetOneByUsername(ctx, user1.Username)
			assert.NoError(t, err)
			assert.Equal(t, user1, result)
		})

		t.Run("not found", func(t *testing.T) {
			_, err := repo.GetOneByUsername(ctx, "nonexistent")
			assert.ErrorIs(t, err, repository.ErrNotFound)
		})

		t.Run("canceled context", func(t *testing.T) {
			_, err := repo.GetOneByUsername(canceledCtx, user1.Username)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("UsernameExists", func(t *testing.T) {
		t.Run("exists", func(t *testing.T) {
			exists, err := repo.UsernameExists(ctx, user1.Username)
			assert.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("not exists", func(t *testing.T) {
			exists, err := repo.UsernameExists(ctx, "nonexistent")
			assert.NoError(t, err)
			assert.False(t, exists)
		})

		t.Run("canceled context", func(t *testing.T) {
			_, err := repo.UsernameExists(canceledCtx, user1.Username)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("Concurrency", func(t *testing.T) {
		t.Run("parallel read/write", func(t *testing.T) {
			const numWorkers = 10
			done := make(chan struct{})

			go func() {
				for i := range numWorkers {
					user := entity.User{
						Id:       uuid.New(),
						Username: "concurrent_" + string(rune(i)),
					}
					_ = repo.Create(ctx, user)
				}
				close(done)
			}()

			for range numWorkers {
				go func() {
					for {
						select {
						case <-done:
							return
						default:
							_, _ = repo.GetOneById(ctx, user1.Id)
							_, _ = repo.GetManyByIds(ctx, []uuid.UUID{user1.Id, user2.Id})
							_, _ = repo.UsernameExists(ctx, user1.Username)
						}
					}
				}()
			}

			<-done
		})
	})
}
