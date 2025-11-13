package usecase

import (
	"ddd-bottomup/domain"
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

func NewGetUserOutput(user *domain.User) *GetUserOutput {
	return &GetUserOutput{
		UserID:    user.ID().Value(),
		FirstName: user.Name().FirstName(),
		LastName:  user.Name().LastName(),
		Email:     user.Email().Value(),
	}
}

type GetUserUseCase struct {
	userRepository domain.UserRepository
}

func NewGetUserUseCase(userRepository domain.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserUseCase) Execute(input GetUserInput) (*GetUserOutput, error) {
	userID, err := domain.ReconstructUserID(input.UserID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.UserNotFoundError{ID: input.UserID}
	}

	return NewGetUserOutput(user), nil
}
