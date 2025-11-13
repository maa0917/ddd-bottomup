package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/service"
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
	// おすすめサークルサービスを作成
	recommendationService := service.NewCircleRecommendationService(time.Now())

	// すべてのサークルを取得してフィルタリング
	allCircles, err := uc.circleRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var filteredCircles []*entity.Circle
	for _, circle := range allCircles {
		if recommendationService.IsRecommended(circle) {
			filteredCircles = append(filteredCircles, circle)
		}
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
