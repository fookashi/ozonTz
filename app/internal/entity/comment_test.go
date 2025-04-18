package entity

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommentEntity(t *testing.T) {
	validUserId := uuid.New()
	validPostId := uuid.New()
	validContent := "This is a valid comment content"
	var nilParentId *uuid.UUID = nil

	t.Run("successful creation", func(t *testing.T) {
		comment, err := NewComment(validUserId, validPostId, nilParentId, validContent)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, validUserId, comment.UserId)
		assert.Equal(t, validPostId, comment.PostId)
		assert.Nil(t, comment.ParentId)
		assert.Equal(t, validContent, comment.Content)
		assert.False(t, comment.CreatedAt.IsZero())
	})

	t.Run("empty content", func(t *testing.T) {
		comment, err := NewComment(validUserId, validPostId, nilParentId, "")

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, ErrEmptyContent)
	})

	t.Run("comment too long", func(t *testing.T) {
		longContent := strings.Repeat("a", 2001)
		comment, err := NewComment(validUserId, validPostId, nilParentId, longContent)

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, ErrCommentTooLong)
	})

	t.Run("nil user id", func(t *testing.T) {
		comment, err := NewComment(uuid.Nil, validPostId, nilParentId, validContent)

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("nil post id", func(t *testing.T) {
		comment, err := NewComment(validUserId, uuid.Nil, nilParentId, validContent)

		assert.Nil(t, comment)
		assert.ErrorIs(t, err, ErrInvalidPostID)
	})

	t.Run("with parent id", func(t *testing.T) {
		parentId := uuid.New()
		comment, err := NewComment(validUserId, validPostId, &parentId, validContent)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, parentId, *comment.ParentId)
	})
}
