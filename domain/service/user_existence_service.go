package service

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
)

type UserExistenceService struct {
	userRepository repository.UserRepository
}

func NewUserExistenceService(userRepository repository.UserRepository) *UserExistenceService {
	return &UserExistenceService{
		userRepository: userRepository,
	}
}

func (s *UserExistenceService) Exists(user *entity.User) (bool, error) {
	existingUser, err := s.userRepository.FindByName(user.Name())
	if err != nil {
		return false, err
	}
	return existingUser != nil, nil
}