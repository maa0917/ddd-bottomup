package infrastructure

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/valueobject"
	"sync"
)

type MemoryUserRepository struct {
	users map[string]*entity.User
	mutex sync.RWMutex
}

func NewMemoryUserRepository() repository.UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *MemoryUserRepository) FindByID(id *entity.UserID) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MemoryUserRepository) FindByName(name *valueobject.FullName) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Name().Equals(name) {
			return user, nil
		}
	}
	return nil, nil
}

func (r *MemoryUserRepository) Save(user *entity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.users[user.ID().Value()] = user
	return nil
}

func (r *MemoryUserRepository) Delete(id *entity.UserID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.users, id.Value())
	return nil
}

// テスト用ヘルパーメソッド
func (r *MemoryUserRepository) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.users = make(map[string]*entity.User)
}

func (r *MemoryUserRepository) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	return len(r.users)
}