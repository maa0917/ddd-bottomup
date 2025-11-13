package service

import (
	"ddd-bottomup/domain/entity"
	"time"
)

const (
	MinMembersForRecommendation = 10
)

type CircleRecommendationService struct {
	baseTime time.Time
}

func NewCircleRecommendationService(baseTime time.Time) *CircleRecommendationService {
	return &CircleRecommendationService{
		baseTime: baseTime,
	}
}

func (s *CircleRecommendationService) IsRecommended(circle *entity.Circle) bool {
	return s.isRecentlyCreated(circle) && s.hasEnoughMembers(circle)
}

func (s *CircleRecommendationService) isRecentlyCreated(circle *entity.Circle) bool {
	oneMonthAgo := s.baseTime.AddDate(0, -1, 0)
	return circle.CreatedAt().After(oneMonthAgo)
}

func (s *CircleRecommendationService) hasEnoughMembers(circle *entity.Circle) bool {
	return circle.GetTotalParticipants() >= MinMembersForRecommendation
}