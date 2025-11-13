package usecase

import (
	"ddd-bottomup/domain"
)

type UpdateUserInput struct {
	UserID    string
	FirstName *string // オプショナル
	LastName  *string // オプショナル
	Email     *string // オプショナル
}

type UpdateUserOutput struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
}

func NewUpdateUserOutput(user *domain.User) *UpdateUserOutput {
	return &UpdateUserOutput{
		UserID:    user.ID().Value(),
		FirstName: user.Name().FirstName(),
		LastName:  user.Name().LastName(),
		Email:     user.Email().Value(),
	}
}

type UpdateUserUseCase struct {
	userRepository       domain.UserRepository
	userExistenceService *domain.UserExistenceService
}

func NewUpdateUserUseCase(userRepository domain.UserRepository, userExistenceService *domain.UserExistenceService) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository:       userRepository,
		userExistenceService: userExistenceService,
	}
}

func (uc *UpdateUserUseCase) Execute(input UpdateUserInput) (*UpdateUserOutput, error) {
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

	// 名前更新（指定されている場合）
	if input.FirstName != nil && input.LastName != nil {
		newName, err := domain.NewFullName(*input.FirstName, *input.LastName)
		if err != nil {
			return nil, err
		}

		// 現在の名前と同じかチェック
		if !user.Name().Equals(newName) {
			// 名前変更前に重複チェック（変更先の名前で一時的にユーザーを作成してチェック）
			tempUser := domain.ReconstructUser(domain.NewUserID(), newName, user.Email(), user.IsPremium())
			exists, err := uc.userExistenceService.Exists(tempUser)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, domain.DuplicateUserNameError{Name: newName.String()}
			}

			user.ChangeName(newName)
		}
	}

	// メール更新（指定されている場合）
	if input.Email != nil {
		newEmail, err := domain.NewEmail(*input.Email)
		if err != nil {
			return nil, err
		}

		user.ChangeEmail(newEmail)
	}

	err = uc.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return NewUpdateUserOutput(user), nil
}
