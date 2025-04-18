package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostEntity(t *testing.T) {
	validUserId := uuid.New()
	validTitle := "Valid Title"
	validContent := "Valid content that is long enough"

	t.Run("success", func(t *testing.T) {
		post, err := NewPost(validUserId, validTitle, validContent, true)

		assert.NoError(t, err)
		assert.NotNil(t, post)
		assert.Equal(t, validTitle, post.Title)
		assert.Equal(t, validContent, post.Content)
		assert.Equal(t, validUserId, post.UserId)
		assert.True(t, post.IsCommentable)
		assert.False(t, post.CreatedAt.IsZero())
	})

	t.Run("empty title", func(t *testing.T) {
		post, err := NewPost(validUserId, "", validContent, true)

		assert.Nil(t, post)
		assert.ErrorIs(t, err, ErrEmptyTitle)
	})

	t.Run("empty content", func(t *testing.T) {
		post, err := NewPost(validUserId, validTitle, "", true)

		assert.Nil(t, post)
		assert.ErrorIs(t, err, ErrEmptyContent)
	})

	t.Run("nil user id", func(t *testing.T) {
		post, err := NewPost(uuid.Nil, validTitle, validContent, true)

		assert.Nil(t, post)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("validate method", func(t *testing.T) {
		post := &Post{
			UserId:  validUserId,
			Title:   validTitle,
			Content: validContent,
		}

		assert.NoError(t, post.Validate())
	})
}
