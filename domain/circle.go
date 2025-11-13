package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type CircleID struct {
	value string
}

func NewCircleID() *CircleID {
	return &CircleID{value: uuid.New().String()}
}

func ReconstructCircleID(value string) (*CircleID, error) {
	if value == "" {
		return nil, errors.New("circle ID cannot be empty")
	}
	if _, err := uuid.Parse(value); err != nil {
		return nil, errors.New("invalid circle ID format")
	}
	return &CircleID{value: value}, nil
}

func (c *CircleID) Value() string {
	return c.value
}

func (c *CircleID) Equals(other *CircleID) bool {
	if other == nil {
		return false
	}
	return c.value == other.value
}

func (c *CircleID) String() string {
	return c.value
}

type Circle struct {
	id        *CircleID
	name      *CircleName
	ownerID   *UserID
	memberIDs []*UserID
	createdAt time.Time
}

func NewCircle(name *CircleName, ownerID *UserID) *Circle {
	return &Circle{
		id:        NewCircleID(),
		name:      name,
		ownerID:   ownerID,
		memberIDs: []*UserID{},
		createdAt: time.Now(),
	}
}

func ReconstructCircle(id *CircleID, name *CircleName, ownerID *UserID, memberIDs []*UserID, createdAt time.Time) *Circle {
	return &Circle{
		id:        id,
		name:      name,
		ownerID:   ownerID,
		memberIDs: memberIDs,
		createdAt: createdAt,
	}
}

func (c *Circle) ID() *CircleID {
	return c.id
}

func (c *Circle) Name() *CircleName {
	return c.name
}

func (c *Circle) OwnerID() *UserID {
	return c.ownerID
}

func (c *Circle) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Circle) GetMemberIDs() []*UserID {
	// 防御的コピーを返す
	memberIDs := make([]*UserID, len(c.memberIDs))
	copy(memberIDs, c.memberIDs)
	return memberIDs
}

func (c *Circle) GetMemberCount() int {
	return len(c.memberIDs)
}

func (c *Circle) GetTotalParticipants() int {
	return 1 + len(c.memberIDs) // オーナー1名 + メンバー数
}

func (c *Circle) ChangeName(name *CircleName) {
	c.name = name
}

func (c *Circle) AddMember(userID *UserID) {
	c.memberIDs = append(c.memberIDs, userID)
}

func (c *Circle) RemoveMember(userID *UserID) {
	for i, memberID := range c.memberIDs {
		if memberID.Equals(userID) {
			c.memberIDs = append(c.memberIDs[:i], c.memberIDs[i+1:]...)
			break
		}
	}
}

func (c *Circle) IsMember(userID *UserID) bool {
	for _, memberID := range c.memberIDs {
		if memberID.Equals(userID) {
			return true
		}
	}
	return false
}

func (c *Circle) IsOwner(userID *UserID) bool {
	if userID == nil || c.ownerID == nil {
		return false
	}
	return c.ownerID.Equals(userID)
}

func (c *Circle) CanAddMember() bool {
	return c.GetTotalParticipants() < 30
}

func (c *Circle) IsFull() bool {
	return !c.CanAddMember()
}

func (c *Circle) GetAvailableSlots() int {
	return 30 - c.GetTotalParticipants()
}

func (c *Circle) Equals(other *Circle) bool {
	if other == nil {
		return false
	}
	return c.id.Equals(other.id)
}

// CircleMembers - サークルメンバー集合エンティティ
type CircleMembers struct {
	owner   *User
	members []*User
}

func NewCircleMembers(owner *User, members []*User) *CircleMembers {
	return &CircleMembers{
		owner:   owner,
		members: members,
	}
}

func (cm *CircleMembers) CountPremiumMembers() int {
	count := 0

	// オーナーのプレミアム判定
	if cm.owner != nil && cm.owner.IsPremium() {
		count++
	}

	// メンバーのプレミアム判定
	for _, member := range cm.members {
		if member != nil && member.IsPremium() {
			count++
		}
	}

	return count
}

func (cm *CircleMembers) GetTotalParticipants() int {
	return 1 + len(cm.members) // オーナー1名 + メンバー数
}

func (cm *CircleMembers) GetMemberCount() int {
	return len(cm.members)
}

// CircleMemberService - サークルメンバー管理サービス
const (
	BasicMemberLimit       = 30
	PremiumMemberLimit     = 50
	PremiumMemberThreshold = 10
)

type CircleMemberService struct{}

func NewCircleMemberService() *CircleMemberService {
	return &CircleMemberService{}
}

func (s *CircleMemberService) GetMaxLimit(circleMembers *CircleMembers) int {
	premiumCount := circleMembers.CountPremiumMembers()
	if premiumCount >= PremiumMemberThreshold {
		return PremiumMemberLimit
	}
	return BasicMemberLimit
}

func (s *CircleMemberService) CanAddMember(circleMembers *CircleMembers) bool {
	currentParticipants := circleMembers.GetTotalParticipants()
	maxLimit := s.GetMaxLimit(circleMembers)
	return currentParticipants < maxLimit
}

func (s *CircleMemberService) GetAvailableSlots(circleMembers *CircleMembers) int {
	maxLimit := s.GetMaxLimit(circleMembers)
	totalParticipants := circleMembers.GetTotalParticipants()
	return maxLimit - totalParticipants
}

// CircleExistenceService - サークル存在確認サービス
type CircleExistenceService struct {
	circleRepository CircleRepository
}

func NewCircleExistenceService(circleRepository CircleRepository) *CircleExistenceService {
	return &CircleExistenceService{
		circleRepository: circleRepository,
	}
}

func (s *CircleExistenceService) Exists(circle *Circle) (bool, error) {
	if circle == nil {
		return false, nil
	}

	found, err := s.circleRepository.FindByName(circle.Name())
	if err != nil {
		return false, err
	}

	if found == nil {
		return false, nil
	}

	// 同じIDのサークルの場合は重複ではない（更新時のチェック用）
	if found.ID().Equals(circle.ID()) {
		return false, nil
	}

	return true, nil
}

// CircleRecommendationService - サークル推薦サービス
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

func (s *CircleRecommendationService) IsRecommended(circle *Circle) bool {
	return s.isRecentlyCreated(circle) && s.hasEnoughMembers(circle)
}

func (s *CircleRecommendationService) isRecentlyCreated(circle *Circle) bool {
	oneMonthAgo := s.baseTime.AddDate(0, -1, 0)
	return circle.CreatedAt().After(oneMonthAgo)
}

func (s *CircleRecommendationService) hasEnoughMembers(circle *Circle) bool {
	return circle.GetTotalParticipants() >= MinMembersForRecommendation
}
