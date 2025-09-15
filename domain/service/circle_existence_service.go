package service

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
)

type CircleExistenceService struct {
	circleRepository repository.CircleRepository
}

func NewCircleExistenceService(circleRepository repository.CircleRepository) *CircleExistenceService {
	return &CircleExistenceService{
		circleRepository: circleRepository,
	}
}

func (s *CircleExistenceService) Exists(circle *entity.Circle) (bool, error) {
	if circle == nil {
		return false, nil
	}

	found, err := s.circleRepository.FindByName(circle.Name())
	if err != nil {
		return false, err
	}

	if found == nil {
		return false, nil
	}

	// 同じIDのサークルの場合は重複ではない（更新時のチェック用）
	if found.ID().Equals(circle.ID()) {
		return false, nil
	}

	return true, nil
}