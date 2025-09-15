package usecase

import (
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/specification"
	"time"
)

type GetRecommendedCirclesOutput struct {
	Circles []RecommendedCircleInfo
}

type RecommendedCircleInfo struct {
	CircleID     string
	CircleName   string
	OwnerID      string
	MemberCount  int
	TotalMembers int
	CreatedAt    string
}

type GetRecommendedCirclesUseCase struct {
	circleRepository repository.CircleRepository
}

func NewGetRecommendedCirclesUseCase(
	circleRepository repository.CircleRepository,
) *GetRecommendedCirclesUseCase {
	return &GetRecommendedCirclesUseCase{
		circleRepository: circleRepository,
	}
}

func (uc *GetRecommendedCirclesUseCase) Execute() (*GetRecommendedCirclesOutput, error) {
	// おすすめサークル仕様を作成
	recommendedSpec := specification.NewRecommendedCircleSpecification(time.Now())

	// リポジトリでフィルタリング済みのサークルを取得
	filteredCircles, err := uc.circleRepository.FindBySpecification(recommendedSpec)
	if err != nil {
		return nil, err
	}

	var recommendedCircles []RecommendedCircleInfo
	for _, circle := range filteredCircles {
		recommendedCircles = append(recommendedCircles, RecommendedCircleInfo{
			CircleID:     circle.ID().Value(),
			CircleName:   circle.Name().Value(),
			OwnerID:      circle.OwnerID().Value(),
			MemberCount:  circle.GetMemberCount(),
			TotalMembers: circle.GetTotalParticipants(),
			CreatedAt:    circle.CreatedAt().Format("2006-01-02"),
		})
	}

	return &GetRecommendedCirclesOutput{
		Circles: recommendedCircles,
	}, nil
}
