package specification

import (
	"ddd-bottomup/domain/entity"
	"time"
)

const (
	MinMembersForRecommendation = 10
)

type RecommendedCircleSpecification struct{
	baseTime time.Time
}

// CircleSpecificationインターフェースの実装確認
var _ CircleSpecification = (*RecommendedCircleSpecification)(nil)

func NewRecommendedCircleSpecification(baseTime time.Time) *RecommendedCircleSpecification {
	return &RecommendedCircleSpecification{
		baseTime: baseTime,
	}
}

func (spec *RecommendedCircleSpecification) IsSatisfiedBy(circle *entity.Circle) bool {
	return spec.isRecentlyCreated(circle) && spec.hasEnoughMembers(circle)
}

func (spec *RecommendedCircleSpecification) isRecentlyCreated(circle *entity.Circle) bool {
	oneMonthAgo := spec.baseTime.AddDate(0, -1, 0)
	return circle.CreatedAt().After(oneMonthAgo)
}

func (spec *RecommendedCircleSpecification) hasEnoughMembers(circle *entity.Circle) bool {
	return circle.GetTotalParticipants() >= MinMembersForRecommendation
}

