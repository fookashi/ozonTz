package entity

import "errors"

var (
	ErrEmptyUsername  = errors.New("username cannot be empty")
	ErrEmptyTitle     = errors.New("post title cannot be empty")
	ErrEmptyContent   = errors.New("content cannot be empty")
	ErrInvalidUserID  = errors.New("invalid user ID")
	ErrInvalidPostID  = errors.New("invalid post ID")
	ErrCommentTooLong = errors.New("comment is too long")
)

type Entity interface {
	Validate() error
}
