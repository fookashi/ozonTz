package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id            uuid.UUID `db:"id"`
	UserId        uuid.UUID `db:"user_id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	IsCommentable bool      `db:"is_commentable"`
	CreatedAt     time.Time `db:"created_at"`
}

func NewPost(userId uuid.UUID, title string, content string, isCommentable bool) (*Post, error) {
	post := &Post{
		Id:            uuid.New(),
		UserId:        userId,
		Title:         title,
		Content:       content,
		IsCommentable: isCommentable,
		CreatedAt:     time.Now(),
	}

	if err := post.Validate(); err != nil {
		return nil, err
	}

	return post, nil
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return ErrEmptyTitle
	}

	if p.Content == "" {
		return ErrEmptyContent
	}

	if p.UserId == uuid.Nil {
		return ErrInvalidUserID
	}

	return nil
}
