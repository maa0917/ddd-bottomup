package entity

import (
	"ddd-bottomup/domain/valueobject"
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
	name      *valueobject.CircleName
	ownerID   *UserID
	memberIDs []*UserID
	createdAt time.Time
}

func NewCircle(name *valueobject.CircleName, ownerID *UserID) *Circle {
	return &Circle{
		id:        NewCircleID(),
		name:      name,
		ownerID:   ownerID,
		memberIDs: []*UserID{},
		createdAt: time.Now(),
	}
}

func ReconstructCircle(id *CircleID, name *valueobject.CircleName, ownerID *UserID, memberIDs []*UserID, createdAt time.Time) *Circle {
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

func (c *Circle) Name() *valueobject.CircleName {
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




func (c *Circle) ChangeName(name *valueobject.CircleName) {
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
