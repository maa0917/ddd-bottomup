package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"errors"
)

type GetUserInput struct {
	UserID string
}

type GetUserOutput struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
}

func NewGetUserOutput(user *entity.User) *GetUserOutput {
	return &GetUserOutput{
		UserID:    user.ID().Value(),
		FirstName: user.Name().FirstName(),
		LastName:  user.Name().LastName(),
		Email:     user.Email().Value(),
	}
}

type GetUserUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUserUseCase(userRepository repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserUseCase) Execute(input GetUserInput) (*GetUserOutput, error) {
	userID, err := entity.ReconstructUserID(input.UserID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return NewGetUserOutput(user), nil
}
