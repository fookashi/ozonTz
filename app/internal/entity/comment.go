package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID  `json:"id"`
	UserId    uuid.UUID  `json:"userId"`
	PostId    uuid.UUID  `json:"postId"`
	ParentId  *uuid.UUID `json:"parentId"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
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
