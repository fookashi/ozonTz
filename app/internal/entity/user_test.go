package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserEntity(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		username := "User"
		expectedUser, err := NewUser(username)

		assert.NoError(t, err)
		assert.Equal(t, username, expectedUser.Username)
	})

	t.Run("Empty username", func(t *testing.T) {
		username := ""
		expectedUser, err := NewUser(username)
		assert.Nil(t, expectedUser)
		assert.ErrorIs(t, err, ErrEmptyUsername)
	})
}
