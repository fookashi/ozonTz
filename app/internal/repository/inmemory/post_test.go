package inmemory

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryPostRepo(t *testing.T) {
	// Setup
	repo := NewPostRepo(10)
	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	// Test data
	now := time.Now()
	post1 := entity.Post{
		Id:            uuid.New(),
		UserId:        uuid.New(),
		Title:         "First Post",
		Content:       "Short content",
		IsCommentable: true,
		CreatedAt:     now.Add(-2 * time.Hour),
	}
	post2 := entity.Post{
		Id:            uuid.New(),
		UserId:        uuid.New(),
		Title:         "Second Post",
		Content:       "Medium content length",
		IsCommentable: false,
		CreatedAt:     now.Add(-1 * time.Hour),
	}
	post3 := entity.Post{
		Id:            uuid.New(),
		UserId:        uuid.New(),
		Title:         "Third Post",
		Content:       "Very long content for this post",
		IsCommentable: true,
		CreatedAt:     now,
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, &post1)
		assert.NoError(t, err)

		t.Run("canceled context", func(t *testing.T) {
			err := repo.Create(canceledCtx, &post2)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("GetOneById", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result, err := repo.GetOneById(ctx, post1.Id)
			assert.NoError(t, err)
			assert.Equal(t, &post1, result)
		})

		t.Run("not found", func(t *testing.T) {
			_, err := repo.GetOneById(ctx, uuid.New())
			assert.ErrorIs(t, err, repository.ErrNotFound)
		})

		t.Run("canceled context", func(t *testing.T) {
			_, err := repo.GetOneById(canceledCtx, post1.Id)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("Update", func(t *testing.T) {
		updatedPost := post1
		updatedPost.Title = "Updated Title"

		err := repo.Update(ctx, &updatedPost)
		assert.NoError(t, err)

		t.Run("verify update", func(t *testing.T) {
			result, err := repo.GetOneById(ctx, post1.Id)
			assert.NoError(t, err)
			assert.Equal(t, "Updated Title", result.Title)
		})

		t.Run("canceled context", func(t *testing.T) {
			err := repo.Update(canceledCtx, &post1)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		_ = repo.Create(ctx, &post2)
		_ = repo.Create(ctx, &post3)

		t.Run("default sorting (newest first)", func(t *testing.T) {
			posts, err := repo.GetMany(ctx, 10, 0, repository.SortByNewest)
			assert.NoError(t, err)
			assert.Len(t, posts, 3)
			assert.Equal(t, post3.Id, posts[0].Id)
			assert.Equal(t, post2.Id, posts[1].Id)
			assert.Equal(t, post1.Id, posts[2].Id)
		})

		t.Run("oldest first", func(t *testing.T) {
			posts, err := repo.GetMany(ctx, 10, 0, repository.SortByOldest)
			assert.NoError(t, err)
			assert.Len(t, posts, 3)
			assert.Equal(t, post1.Id, posts[0].Id)
			assert.Equal(t, post2.Id, posts[1].Id)
			assert.Equal(t, post3.Id, posts[2].Id)
		})

		t.Run("top (by content length)", func(t *testing.T) {
			posts, err := repo.GetMany(ctx, 10, 0, repository.SortByTop)
			assert.NoError(t, err)
			assert.Len(t, posts, 3)
			assert.Equal(t, post3.Id, posts[0].Id)
			assert.Equal(t, post2.Id, posts[1].Id)
			assert.Equal(t, post1.Id, posts[2].Id)
		})

		t.Run("pagination", func(t *testing.T) {
			t.Run("limit", func(t *testing.T) {
				posts, err := repo.GetMany(ctx, 2, 0, repository.SortByNewest)
				assert.NoError(t, err)
				assert.Len(t, posts, 2)
			})

			t.Run("offset", func(t *testing.T) {
				posts, err := repo.GetMany(ctx, 1, 1, repository.SortByNewest)
				assert.NoError(t, err)
				assert.Len(t, posts, 1)
				assert.Equal(t, post2.Id, posts[0].Id)
			})

			t.Run("offset out of range", func(t *testing.T) {
				posts, err := repo.GetMany(ctx, 10, 10, repository.SortByNewest)
				assert.NoError(t, err)
				assert.Empty(t, posts)
			})
		})

		t.Run("limit 0", func(t *testing.T) {
			posts, err := repo.GetMany(ctx, 0, 0, repository.SortByNewest)
			assert.NoError(t, err)
			assert.Empty(t, posts)
		})

		t.Run("canceled context", func(t *testing.T) {
			_, err := repo.GetMany(canceledCtx, 10, 0, repository.SortByNewest)
			assert.ErrorIs(t, err, repository.ErrContextCanceled)
		})
	})

	t.Run("Concurrency", func(t *testing.T) {
		const numWorkers = 10
		done := make(chan struct{})

		go func() {
			for range numWorkers {
				post := entity.Post{
					Id:        uuid.New(),
					Title:     "Concurrent Post",
					CreatedAt: time.Now(),
				}
				_ = repo.Create(ctx, &post)
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
						_, _ = repo.GetOneById(ctx, post1.Id)
						_, _ = repo.GetMany(ctx, 2, 0, repository.SortByNewest)
					}
				}
			}()
		}

		<-done
	})
}
