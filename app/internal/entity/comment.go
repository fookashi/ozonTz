package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID  `db:"id"`
	UserId    uuid.UUID  `db:"user_id"`
	PostId    uuid.UUID  `db:"post_id"`
	ParentId  *uuid.UUID `db:"parent_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
}

func NewComment(userId, postId uuid.UUID, parentId *uuid.UUID, content string) (*Comment, error) {
	comment := &Comment{
		Id:        uuid.New(),
		UserId:    userId,
		PostId:    postId,
		ParentId:  parentId,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := comment.Validate(); err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *Comment) Validate() error {
	if c.UserId == uuid.Nil {
		return ErrInvalidUserID
	}
	if c.PostId == uuid.Nil {
		return ErrInvalidPostID
	}
	if c.Content == "" {
		return ErrEmptyContent
	}
	if len(c.Content) > 2000 {
		return ErrCommentTooLong
	}
	return nil
}
