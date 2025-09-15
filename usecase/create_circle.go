package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/service"
	"ddd-bottomup/domain/valueobject"
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
	circleRepository       repository.CircleRepository
	userRepository         repository.UserRepository
	circleExistenceService *service.CircleExistenceService
}

func NewCreateCircleUseCase(
	circleRepository repository.CircleRepository,
	userRepository repository.UserRepository,
	circleExistenceService *service.CircleExistenceService,
) *CreateCircleUseCase {
	return &CreateCircleUseCase{
		circleRepository:       circleRepository,
		userRepository:         userRepository,
		circleExistenceService: circleExistenceService,
	}
}

func (uc *CreateCircleUseCase) Execute(input CreateCircleInput) (*CreateCircleOutput, error) {
	// サークル名の値オブジェクト作成
	circleName, err := valueobject.NewCircleName(input.CircleName)
	if err != nil {
		return nil, err
	}

	// オーナーIDからUserIDを再構成
	ownerID, err := entity.ReconstructUserID(input.OwnerID)
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
	circle := entity.NewCircle(circleName, ownerID)

	// サークル保存
	err = uc.circleRepository.Save(circle)
	if err != nil {
		return nil, err
	}

	return &CreateCircleOutput{
		CircleID: circle.ID().Value(),
	}, nil
}

