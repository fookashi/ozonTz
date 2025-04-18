package inmemory

import (
	"app/internal/entity"
	"app/internal/repository"
	"context"
	"sync"

	"github.com/google/uuid"
)

type inMemoryUserRepo struct {
	users         map[uuid.UUID]entity.User
	usernameIndex map[string]uuid.UUID
	lock          sync.RWMutex
}

func NewInMemoryUserRepo(initSize int) *inMemoryUserRepo {
	return &inMemoryUserRepo{
		users:         make(map[uuid.UUID]entity.User, initSize),
		usernameIndex: make(map[string]uuid.UUID, initSize),
	}
}

func (repo *inMemoryUserRepo) Create(ctx context.Context, user entity.User) error {
	if err := ctx.Err(); err != nil {
		return repository.ErrContextCanceled
	}
	repo.lock.Lock()
	defer repo.lock.Unlock()
	repo.users[user.Id] = user
	repo.usernameIndex[user.Username] = user.Id
	return nil
}
func (repo *inMemoryUserRepo) GetOneById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	if err := ctx.Err(); err != nil {
		return entity.User{}, repository.ErrContextCanceled
	}
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	user, exists := repo.users[id]
	if !exists {
		return entity.User{}, repository.ErrNotFound
	}
	return user, nil
}

func (repo *inMemoryUserRepo) GetManyByIds(ctx context.Context, ids []uuid.UUID) ([]entity.User, error) {
	if err := ctx.Err(); err != nil {
		return []entity.User{}, repository.ErrContextCanceled
	}

	repo.lock.RLock()
	defer repo.lock.RUnlock()

	result := make([]entity.User, 0, len(ids))
	var notFoundIDs []uuid.UUID

	for _, id := range ids {
		user, exists := repo.users[id]
		if !exists {
			notFoundIDs = append(notFoundIDs, id)
			continue
		}
		result = append(result, user)
	}

	if len(notFoundIDs) > 0 {
		return result, repository.ErrNotFound
	}

	return result, nil
}

func (repo *inMemoryUserRepo) GetOneByUsername(ctx context.Context, username string) (entity.User, error) {
	if err := ctx.Err(); err != nil {
		return entity.User{}, repository.ErrContextCanceled
	}
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	id := repo.usernameIndex[username]
	user, exists := repo.users[id]
	if !exists {
		return entity.User{}, repository.ErrNotFound
	}
	return user, nil
}
func (repo *inMemoryUserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, repository.ErrContextCanceled
	}
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	_, exists := repo.usernameIndex[username]
	return exists, nil
}
