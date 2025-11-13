package infrastructure

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/valueobject"
	"sync"
)

type MemoryCircleRepository struct {
	circles map[string]*entity.Circle
	mu      sync.RWMutex
}

func NewMemoryCircleRepository() repository.CircleRepository {
	return &MemoryCircleRepository{
		circles: make(map[string]*entity.Circle),
	}
}

func (r *MemoryCircleRepository) FindByID(id *entity.CircleID) (*entity.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	circle, exists := r.circles[id.Value()]
	if !exists {
		return nil, nil
	}
	return circle, nil
}

func (r *MemoryCircleRepository) FindByName(name *valueobject.CircleName) (*entity.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, circle := range r.circles {
		if circle.Name().Equals(name) {
			return circle, nil
		}
	}
	return nil, nil
}

func (r *MemoryCircleRepository) Save(circle *entity.Circle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.circles[circle.ID().Value()] = circle
	return nil
}

func (r *MemoryCircleRepository) Delete(id *entity.CircleID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	delete(r.circles, id.Value())
	return nil
}

func (r *MemoryCircleRepository) FindAll() ([]*entity.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	circles := make([]*entity.Circle, 0, len(r.circles))
	for _, circle := range r.circles {
		circles = append(circles, circle)
	}
	return circles, nil
}

func (r *MemoryCircleRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.circles)
}