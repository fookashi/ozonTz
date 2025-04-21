package model

import (
	"time"
)

type Comment struct {
	ID        string     `json:"id"`
	User      *User      `json:"user"`
	ParentID  *string    `json:"parentId,omitempty"`
	Replies   []*Comment `json:"replies"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
}
