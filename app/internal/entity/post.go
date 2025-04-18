package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"userId"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	IsCommentable bool      `json:"isCommentable"`
	CreatedAt     time.Time `json:"createdAt"`
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
