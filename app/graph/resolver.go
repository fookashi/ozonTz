package graph

import (
	"app/internal/service"
)

type Resolver struct {
	UserService    *service.UserService
	PostService    *service.PostService
	CommentService *service.CommentService
}
