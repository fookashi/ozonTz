package postgres_test

import (
	"app/internal/entity"
	"app/internal/repository/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentRepo(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := postgres.NewCommentRepo(mock)

	t.Run("GetOneById", func(t *testing.T) {
		expectedComment := &entity.Comment{
			Id:        uuid.New(),
			UserId:    uuid.New(),
			PostId:    uuid.New(),
			ParentId:  nil,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		mock.ExpectQuery("SELECT id, user_id, post_id, parent_id, content, created_at FROM comments").
			WithArgs(expectedComment.Id).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "post_id", "parent_id", "content", "created_at"}).
				AddRow(expectedComment.Id, expectedComment.UserId, expectedComment.PostId,
					expectedComment.ParentId, expectedComment.Content, expectedComment.CreatedAt))

		comment, err := repo.GetOneById(context.Background(), expectedComment.Id)
		assert.NoError(t, err)
		assert.Equal(t, expectedComment, comment)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create", func(t *testing.T) {
		comment := &entity.Comment{
			Id:        uuid.New(),
			UserId:    uuid.New(),
			PostId:    uuid.New(),
			ParentId:  nil,
			Content:   "New comment",
			CreatedAt: time.Now(),
		}

		mock.ExpectExec("INSERT INTO comments").
			WithArgs(comment.Id, comment.UserId, comment.PostId, comment.ParentId,
				comment.Content, comment.CreatedAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(context.Background(), comment)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetByPost", func(t *testing.T) {
		postId := uuid.New()
		comments := []entity.Comment{
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				PostId:    postId,
				ParentId:  nil,
				Content:   "Comment 1",
				CreatedAt: time.Now(),
			},
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				PostId:    postId,
				ParentId:  nil,
				Content:   "Comment 2",
				CreatedAt: time.Now().Add(-time.Hour),
			},
		}

		mock.ExpectQuery("SELECT id, user_id, post_id, parent_id, content, created_at FROM comments").
			WithArgs(postId, 10, 0).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "post_id", "parent_id", "content", "created_at"}).
				AddRow(comments[0].Id, comments[0].UserId, comments[0].PostId,
					comments[0].ParentId, comments[0].Content, comments[0].CreatedAt).
				AddRow(comments[1].Id, comments[1].UserId, comments[1].PostId,
					comments[1].ParentId, comments[1].Content, comments[1].CreatedAt))

		result, err := repo.GetByPost(context.Background(), postId, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, comments[0].Content, result[0].Content)
		assert.Equal(t, comments[1].Content, result[1].Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetCommentReplies", func(t *testing.T) {
		parentId := uuid.New()
		replies := []entity.Comment{
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				PostId:    uuid.New(),
				ParentId:  &parentId,
				Content:   "Reply 1",
				CreatedAt: time.Now(),
			},
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				PostId:    uuid.New(),
				ParentId:  &parentId,
				Content:   "Reply 2",
				CreatedAt: time.Now().Add(-time.Hour),
			},
		}

		mock.ExpectQuery("SELECT id, user_id, post_id, parent_id, content, created_at FROM comments").
			WithArgs(parentId, 10, 0).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "post_id", "parent_id", "content", "created_at"}).
				AddRow(replies[0].Id, replies[0].UserId, replies[0].PostId,
					replies[0].ParentId, replies[0].Content, replies[0].CreatedAt).
				AddRow(replies[1].Id, replies[1].UserId, replies[1].PostId,
					replies[1].ParentId, replies[1].Content, replies[1].CreatedAt))

		result, err := repo.GetCommentReplies(context.Background(), parentId, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, replies[0].Content, result[0].Content)
		assert.Equal(t, replies[1].Content, result[1].Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
