package service

import "ddd-bottomup/domain/entity"

const (
	BasicMemberLimit       = 30
	PremiumMemberLimit     = 50
	PremiumMemberThreshold = 10
)

type CircleMemberService struct{}

func NewCircleMemberService() *CircleMemberService {
	return &CircleMemberService{}
}

func (s *CircleMemberService) GetMaxLimit(circleMembers *entity.CircleMembers) int {
	premiumCount := circleMembers.CountPremiumMembers()
	if premiumCount >= PremiumMemberThreshold {
		return PremiumMemberLimit
	}
	return BasicMemberLimit
}

func (s *CircleMemberService) CanAddMember(circleMembers *entity.CircleMembers) bool {
	currentParticipants := circleMembers.GetTotalParticipants()
	maxLimit := s.GetMaxLimit(circleMembers)
	return currentParticipants < maxLimit
}

func (s *CircleMemberService) GetAvailableSlots(circleMembers *entity.CircleMembers) int {
	maxLimit := s.GetMaxLimit(circleMembers)
	totalParticipants := circleMembers.GetTotalParticipants()
	return maxLimit - totalParticipants
}