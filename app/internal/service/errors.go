package service

import "errors"

var (
	ErrUsernameExists        = errors.New("Username already exists")
	ErrDueUserCreation       = errors.New("Error due user creation")
	ErrDuePostCreation       = errors.New("Error due post creation")
	ErrUserNotFound          = errors.New("User not found")
	ErrPostNotFound          = errors.New("Post not found")
	ErrNoPermissionForToggle = errors.New("Only creator can toggle comments")
	ErrPostIsNotCommentable  = errors.New("Post is not commentable")
	ErrParentCommentNotFound = errors.New("Parent comment not found")
	ErrTooManySymbols        = errors.New("Too many symbols")
)
