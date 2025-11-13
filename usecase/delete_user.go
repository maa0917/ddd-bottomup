package usecase

import (
	"ddd-bottomup/domain"
)

type DeleteUserInput struct {
	UserID string
}

type DeleteUserUseCase struct {
	userRepository domain.UserRepository
}

func NewDeleteUserUseCase(userRepository domain.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(input DeleteUserInput) error {
	userID, err := domain.ReconstructUserID(input.UserID)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.UserNotFoundError{ID: input.UserID}
	}

	return uc.userRepository.Delete(userID)
}
