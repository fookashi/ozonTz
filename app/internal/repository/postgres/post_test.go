package postgres_test

import (
	"app/internal/entity"
	"app/internal/repository"
	"app/internal/repository/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepo(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := postgres.NewPostRepo(mock)

	t.Run("Create", func(t *testing.T) {
		post := &entity.Post{
			Id:            uuid.New(),
			UserId:        uuid.New(),
			Title:         "Test Post",
			Content:       "Test Content",
			IsCommentable: true,
			CreatedAt:     time.Now(),
		}

		mock.ExpectExec("INSERT INTO posts").
			WithArgs(post.Id, post.UserId, post.Title, post.Content, post.IsCommentable, post.CreatedAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(context.Background(), post)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update", func(t *testing.T) {
		post := &entity.Post{
			Id:            uuid.New(),
			Title:         "Updated Post",
			Content:       "Updated Content",
			IsCommentable: false,
		}

		mock.ExpectExec("UPDATE posts").
			WithArgs(post.Id, post.Title, post.Content, post.IsCommentable).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(context.Background(), post)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update not found", func(t *testing.T) {
		post := &entity.Post{
			Id:            uuid.New(),
			Title:         "Updated Post",
			Content:       "Updated Content",
			IsCommentable: false,
		}

		mock.ExpectExec("UPDATE posts").
			WithArgs(post.Id, post.Title, post.Content, post.IsCommentable).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := repo.Update(context.Background(), post)
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOneById", func(t *testing.T) {
		expectedPost := &entity.Post{
			Id:            uuid.New(),
			UserId:        uuid.New(),
			Title:         "Test Post",
			Content:       "Test Content",
			IsCommentable: true,
			CreatedAt:     time.Now(),
		}

		mock.ExpectQuery("SELECT id, user_id, title, content, is_commentable, created_at FROM posts").
			WithArgs(expectedPost.Id).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "is_commentable", "created_at"}).
				AddRow(expectedPost.Id, expectedPost.UserId, expectedPost.Title, expectedPost.Content,
					expectedPost.IsCommentable, expectedPost.CreatedAt))

		post, err := repo.GetOneById(context.Background(), expectedPost.Id)
		assert.NoError(t, err)
		assert.Equal(t, expectedPost, post)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetMany newest first", func(t *testing.T) {
		posts := []entity.Post{
			{
				Id:            uuid.New(),
				UserId:        uuid.New(),
				Title:         "Post 1",
				Content:       "Content 1",
				IsCommentable: true,
				CreatedAt:     time.Now(),
			},
			{
				Id:            uuid.New(),
				UserId:        uuid.New(),
				Title:         "Post 2",
				Content:       "Content 2",
				IsCommentable: false,
				CreatedAt:     time.Now().Add(-time.Hour),
			},
		}

		mock.ExpectQuery("SELECT id, user_id, title, content, is_commentable, created_at FROM posts ORDER BY created_at DESC").
			WithArgs(10, 0).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "is_commentable", "created_at"}).
				AddRow(posts[0].Id, posts[0].UserId, posts[0].Title, posts[0].Content,
					posts[0].IsCommentable, posts[0].CreatedAt).
				AddRow(posts[1].Id, posts[1].UserId, posts[1].Title, posts[1].Content,
					posts[1].IsCommentable, posts[1].CreatedAt))

		result, err := repo.GetMany(context.Background(), 10, 0, repository.SortByNewest)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, posts[0].Title, result[0].Title)
		assert.Equal(t, posts[1].Title, result[1].Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
