package inmemory_test

import (
	"app/internal/entity"
	"app/internal/repository"
	"app/internal/repository/inmemory"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommentRepo(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()
	parentCommentID := uuid.New()
	nonExistentID := uuid.New()

	baseComment := entity.Comment{
		Id:        uuid.New(),
		PostId:    postID,
		UserId:    userID,
		Content:   "Base comment",
		CreatedAt: time.Now(),
	}

	replyComment := entity.Comment{
		Id:        uuid.New(),
		PostId:    postID,
		UserId:    userID,
		ParentId:  &parentCommentID,
		Content:   "Reply comment",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name string
		run  func(t *testing.T, repo *inmemory.CommentRepo)
	}{
		{
			name: "Create/success",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				err := repo.Create(context.Background(), &baseComment)
				assert.NoError(t, err)
			},
		},
		{
			name: "Create/canceled context",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				err := repo.Create(ctx, &baseComment)
				assert.ErrorIs(t, err, repository.ErrContextCanceled)
			},
		},
		{
			name: "GetOneByID/success",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				_ = repo.Create(context.Background(), &baseComment)
				result, err := repo.GetOneById(context.Background(), baseComment.Id)
				assert.NoError(t, err)
				assert.Equal(t, baseComment, *result)
			},
		},
		{
			name: "GetOneByID/not found",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				_, err := repo.GetOneById(context.Background(), nonExistentID)
				assert.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
		{
			name: "GetOneByID/canceled context",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				_ = repo.Create(context.Background(), &baseComment)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := repo.GetOneById(ctx, baseComment.Id)
				assert.ErrorIs(t, err, repository.ErrContextCanceled)
			},
		},
		{
			name: "GetCommentReplies/success",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				parentComment := baseComment
				parentComment.Id = parentCommentID
				_ = repo.Create(context.Background(), &parentComment)

				_ = repo.Create(context.Background(), &replyComment)

				result, err := repo.GetCommentReplies(context.Background(), parentCommentID, 10, 0)
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, replyComment, result[0])
			},
		},
		{
			name: "GetCommentReplies/parent not found",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				_, err := repo.GetCommentReplies(context.Background(), nonExistentID, 10, 0)
				assert.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
		{
			name: "GetCommentReplies/no replies",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				parentComment := baseComment
				parentComment.Id = parentCommentID
				_ = repo.Create(context.Background(), &parentComment)

				result, err := repo.GetCommentReplies(context.Background(), parentCommentID, 10, 0)
				assert.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name: "GetCommentReplies/pagination",
			run: func(t *testing.T, repo *inmemory.CommentRepo) {
				parentComment := baseComment
				parentComment.Id = parentCommentID
				_ = repo.Create(context.Background(), &parentComment)

				for i := 0; i < 5; i++ {
					reply := replyComment
					reply.Id = uuid.New()
					_ = repo.Create(context.Background(), &reply)
				}

				result, err := repo.GetCommentReplies(context.Background(), parentCommentID, 2, 1)
				assert.NoError(t, err)
				assert.Len(t, result, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewCommentRepo(10)
			tt.run(t, repo)
		})
	}

	t.Run("Concurrency", func(t *testing.T) {
		repo := inmemory.NewCommentRepo(100)
		const numWorkers = 10
		done := make(chan struct{})

		for i := 0; i < numWorkers; i++ {
			go func(i int) {
				for j := 0; j < 10; j++ {
					comment := baseComment
					comment.Id = uuid.New()
					comment.Content = fmt.Sprintf("Comment %d-%d", i, j)
					_ = repo.Create(context.Background(), &comment)
				}
				done <- struct{}{}
			}(i)
		}

		for i := 0; i < numWorkers; i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						_, _ = repo.GetOneById(context.Background(), baseComment.Id)
						_, _ = repo.GetByPost(context.Background(), postID, 10, 0)
					}
				}
			}()
		}

		for range numWorkers {
			<-done
		}
		close(done)
	})
}
