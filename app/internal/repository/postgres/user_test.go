package postgres

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		user := &entity.User{
			Id:       uuid.New(),
			Username: "testuser",
			Roles:    []string{"user"},
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.Id, user.Username, user.Roles).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepo_GetOneById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		expectedUser := &entity.User{
			Id:       uuid.New(),
			Username: "testuser",
			Roles:    []string{"user"},
		}

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE id =").
			WithArgs(expectedUser.Id).
			WillReturnRows(pgxmock.NewRows([]string{"id", "username", "roles"}).
				AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Roles))

		user, err := repo.GetOneById(context.Background(), expectedUser.Id)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		id := uuid.New()

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE id =").
			WithArgs(id).
			WillReturnError(pgx.ErrNoRows)

		user, err := repo.GetOneById(context.Background(), id)
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		id := uuid.New()
		expectedErr := errors.New("database error")

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE id =").
			WithArgs(id).
			WillReturnError(expectedErr)

		user, err := repo.GetOneById(context.Background(), id)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepo_GetOneByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		expectedUser := &entity.User{
			Id:       uuid.New(),
			Username: "testuser",
			Roles:    []string{"user"},
		}

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE username =").
			WithArgs(expectedUser.Username).
			WillReturnRows(pgxmock.NewRows([]string{"id", "username", "roles"}).
				AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Roles))

		user, err := repo.GetOneByUsername(context.Background(), expectedUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		username := "nonexistent"

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE username =").
			WithArgs(username).
			WillReturnError(pgx.ErrNoRows)

		user, err := repo.GetOneByUsername(context.Background(), username)
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		username := "testuser"
		expectedErr := errors.New("database error")

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE username =").
			WithArgs(username).
			WillReturnError(expectedErr)

		user, err := repo.GetOneByUsername(context.Background(), username)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepo_GetManyByIds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		users := []entity.User{
			{Id: uuid.New(), Username: "user1", Roles: []string{"user"}},
			{Id: uuid.New(), Username: "user2", Roles: []string{"admin"}},
		}

		ids := make([]uuid.UUID, len(users))
		for i, u := range users {
			ids[i] = u.Id
		}

		rows := pgxmock.NewRows([]string{"id", "username", "roles"})
		for _, u := range users {
			rows.AddRow(u.Id, u.Username, u.Roles)
		}

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE id = ANY").
			WithArgs(ids).
			WillReturnRows(rows)

		result, err := repo.GetManyByIds(context.Background(), ids)
		assert.NoError(t, err)
		assert.Len(t, result, len(users))
		for _, u := range users {
			assert.Equal(t, u, result[u.Id])
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty ids", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		result, err := repo.GetManyByIds(context.Background(), []uuid.UUID{})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		ids := []uuid.UUID{uuid.New()}
		expectedErr := errors.New("database error")

		mock.ExpectQuery("SELECT id, username, roles FROM users WHERE id = ANY").
			WithArgs(ids).
			WillReturnError(expectedErr)

		result, err := repo.GetManyByIds(context.Background(), ids)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
