package infrastructure

import (
	"ddd-bottomup/domain"
	"sync"
)

type MemoryCircleRepository struct {
	circles map[string]*domain.Circle
	mu      sync.RWMutex
}

func NewMemoryCircleRepository() domain.CircleRepository {
	return &MemoryCircleRepository{
		circles: make(map[string]*domain.Circle),
	}
}

func (r *MemoryCircleRepository) FindByID(id *domain.CircleID) (*domain.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	circle, exists := r.circles[id.Value()]
	if !exists {
		return nil, nil
	}
	return circle, nil
}

func (r *MemoryCircleRepository) FindByName(name *domain.CircleName) (*domain.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, circle := range r.circles {
		if circle.Name().Equals(name) {
			return circle, nil
		}
	}
	return nil, nil
}

func (r *MemoryCircleRepository) Save(circle *domain.Circle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.circles[circle.ID().Value()] = circle
	return nil
}

func (r *MemoryCircleRepository) Delete(id *domain.CircleID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.circles, id.Value())
	return nil
}

func (r *MemoryCircleRepository) FindAll() ([]*domain.Circle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	circles := make([]*domain.Circle, 0, len(r.circles))
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
