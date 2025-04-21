package model

import (
	"time"
)

type Post struct {
	ID            string     `json:"id"`
	User          *User      `json:"user"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	IsCommentable bool       `json:"isCommentable"`
	Comments      []*Comment `json:"comments"`
	CreatedAt     time.Time  `json:"createdAt"`
}
