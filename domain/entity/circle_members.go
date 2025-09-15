package entity

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