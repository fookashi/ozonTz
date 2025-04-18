package repository

import (
	"app/internal/entity"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetOneById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) GetManyByIds(ctx context.Context, ids []uuid.UUID) ([]entity.User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserRepo) GetOneByUsername(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Create(ctx context.Context, post entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) Update(ctx context.Context, post entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) GetOneById(ctx context.Context, id uuid.UUID) (entity.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Post), args.Error(1)
}

func (m *MockPostRepo) GetMany(ctx context.Context, limit, offset int, sortBy SortBy) ([]entity.Post, error) {
	args := m.Called(ctx, limit, offset, sortBy)
	return args.Get(0).([]entity.Post), args.Error(1)
}

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) Create(ctx context.Context, comment entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepo) GetOneByID(ctx context.Context, commentId uuid.UUID) (*entity.Comment, error) {
	args := m.Called(ctx, commentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentRepo) GetByPost(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*entity.Comment, error) {
	args := m.Called(ctx, postId, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *MockCommentRepo) GetReplies(ctx context.Context, parentId uuid.UUID, limit, offset int) ([]*entity.Comment, error) {
	args := m.Called(ctx, parentId, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

type MockRepoHolder struct {
	UserRepo    *MockUserRepo
	PostRepo    *MockPostRepo
	CommentRepo *MockCommentRepo
}

func NewMockRepoHolder() *MockRepoHolder {
	return &MockRepoHolder{
		UserRepo:    &MockUserRepo{},
		PostRepo:    &MockPostRepo{},
		CommentRepo: &MockCommentRepo{},
	}
}

func (m *MockRepoHolder) GetHolder() *RepoHolder {
	return &RepoHolder{
		UserRepo:    m.UserRepo,
		PostRepo:    m.PostRepo,
		CommentRepo: m.CommentRepo,
	}
}
