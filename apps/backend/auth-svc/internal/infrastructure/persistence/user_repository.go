package persistence

import (
	"context"
	"fmt"
	"sync"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/domain/entity"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		users: make(map[string]*entity.User),
	}

	adminUser, _ := entity.NewUser("admin", "admin123")
	adminUser.SetRole(entity.RoleAdmin)
	repo.users["admin"] = adminUser

	return repo
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	r.users[user.Username] = user
	return nil
}

func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", username)
	}

	return user, nil
}

func (r *InMemoryUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[username]
	return exists, nil
}
