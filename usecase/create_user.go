package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/service"
	"ddd-bottomup/domain/valueobject"
	"errors"
)

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	IsPremium bool
}

type CreateUserOutput struct {
	UserID string
}

type CreateUserUseCase struct {
	userRepository       repository.UserRepository
	userExistenceService *service.UserExistenceService
}

func NewCreateUserUseCase(
	userRepository repository.UserRepository,
	userExistenceService *service.UserExistenceService,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository:       userRepository,
		userExistenceService: userExistenceService,
	}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	fullName, err := valueobject.NewFullName(input.FirstName, input.LastName)
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(fullName, email, input.IsPremium)

	exists, err := uc.userExistenceService.Exists(user)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	if err := uc.userRepository.Save(user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		UserID: user.ID().String(),
	}, nil
}
