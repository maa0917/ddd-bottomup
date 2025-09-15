package specification

import "ddd-bottomup/domain/entity"

const (
	BasicMemberLimit       = 30
	PremiumMemberLimit     = 50
	PremiumMemberThreshold = 10
)

type CircleMemberLimitSpecification struct{}

func NewCircleMemberLimitSpecification() *CircleMemberLimitSpecification {
	return &CircleMemberLimitSpecification{}
}

func (spec *CircleMemberLimitSpecification) GetMaxLimit(circleMembers *entity.CircleMembers) int {
	premiumCount := circleMembers.CountPremiumMembers()
	if premiumCount >= PremiumMemberThreshold {
		return PremiumMemberLimit
	}
	return BasicMemberLimit
}

func (spec *CircleMemberLimitSpecification) IsSatisfiedBy(circleMembers *entity.CircleMembers) bool {
	currentParticipants := circleMembers.GetTotalParticipants()
	maxLimit := spec.GetMaxLimit(circleMembers)
	return currentParticipants < maxLimit
}

func (spec *CircleMemberLimitSpecification) GetAvailableSlots(circleMembers *entity.CircleMembers) int {
	maxLimit := spec.GetMaxLimit(circleMembers)
	totalParticipants := circleMembers.GetTotalParticipants()
	return maxLimit - totalParticipants
}