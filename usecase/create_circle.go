package usecase

import (
	"ddd-bottomup/domain"
	"errors"
)

type CreateCircleInput struct {
	CircleName string
	OwnerID    string
}

type CreateCircleOutput struct {
	CircleID string
}

type CreateCircleUseCase struct {
	circleRepository       domain.CircleRepository
	userRepository         domain.UserRepository
	circleExistenceService *domain.CircleExistenceService
}

func NewCreateCircleUseCase(
	circleRepository domain.CircleRepository,
	userRepository domain.UserRepository,
	circleExistenceService *domain.CircleExistenceService,
) *CreateCircleUseCase {
	return &CreateCircleUseCase{
		circleRepository:       circleRepository,
		userRepository:         userRepository,
		circleExistenceService: circleExistenceService,
	}
}

func (uc *CreateCircleUseCase) Execute(input CreateCircleInput) (*CreateCircleOutput, error) {
	// サークル名の値オブジェクト作成
	circleName, err := domain.NewCircleName(input.CircleName)
	if err != nil {
		return nil, err
	}

	// オーナーIDからUserIDを再構成
	ownerID, err := domain.ReconstructUserID(input.OwnerID)
	if err != nil {
		return nil, err
	}

	// オーナーユーザーの存在確認
	owner, err := uc.userRepository.FindByID(ownerID)
	if err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, errors.New("owner user not found")
	}

	// 同名のサークルが存在しないかチェック
	existingCircle, err := uc.circleRepository.FindByName(circleName)
	if err != nil {
		return nil, err
	}
	if existingCircle != nil {
		return nil, errors.New("circle name already exists")
	}

	// サークル作成
	circle := domain.NewCircle(circleName, ownerID)

	// サークル保存
	err = uc.circleRepository.Save(circle)
	if err != nil {
		return nil, err
	}

	return &CreateCircleOutput{
		CircleID: circle.ID().Value(),
	}, nil
}
