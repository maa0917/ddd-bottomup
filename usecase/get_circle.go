package usecase

import (
	"ddd-bottomup/domain"
	"errors"
)

type GetCircleInput struct {
	CircleID string
}

type GetCircleOutput struct {
	CircleID       string
	CircleName     string
	OwnerID        string
	MemberIDs      []string
	TotalMembers   int
	AvailableSlots int
}

type GetCircleUseCase struct {
	circleRepository domain.CircleRepository
	userRepository   domain.UserRepository
}

func NewGetCircleUseCase(circleRepository domain.CircleRepository, userRepository domain.UserRepository) *GetCircleUseCase {
	return &GetCircleUseCase{
		circleRepository: circleRepository,
		userRepository:   userRepository,
	}
}

func (uc *GetCircleUseCase) Execute(input GetCircleInput) (*GetCircleOutput, error) {
	// CircleIDを再構成
	circleID, err := domain.ReconstructCircleID(input.CircleID)
	if err != nil {
		return nil, err
	}

	// リポジトリからエンティティを取得
	circle, err := uc.circleRepository.FindByID(circleID)
	if err != nil {
		return nil, err
	}
	if circle == nil {
		return nil, errors.New("circle not found")
	}

	// オーナーを取得
	owner, err := uc.userRepository.FindByID(circle.OwnerID())
	if err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, errors.New("owner not found")
	}

	// メンバーリストを構築
	var members []*domain.User
	for _, memberID := range circle.GetMemberIDs() {
		member, err := uc.userRepository.FindByID(memberID)
		if err != nil {
			return nil, err
		}
		if member != nil {
			members = append(members, member)
		}
	}

	// プレミアム制限を考慮した利用可能枠を計算
	circleMembers := domain.NewCircleMembers(owner, members)
	memberService := domain.NewCircleMemberService()

	// アウトプットに変換
	return &GetCircleOutput{
		CircleID:       circle.ID().Value(),
		CircleName:     circle.Name().Value(),
		OwnerID:        circle.OwnerID().Value(),
		MemberIDs:      convertUserIDsToStrings(circle.GetMemberIDs()),
		TotalMembers:   circle.GetTotalParticipants(),
		AvailableSlots: memberService.GetAvailableSlots(circleMembers),
	}, nil
}

func convertUserIDsToStrings(userIDs []*domain.UserID) []string {
	result := make([]string, len(userIDs))
	for i, userID := range userIDs {
		result[i] = userID.Value()
	}
	return result
}
