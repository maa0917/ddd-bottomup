package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/specification"
	"errors"
	"fmt"
)

type AddMemberInput struct {
	CircleID string
	UserID   string
}

type AddMemberUseCase struct {
	circleRepository repository.CircleRepository
	userRepository   repository.UserRepository
}

func NewAddMemberUseCase(
	circleRepository repository.CircleRepository,
	userRepository repository.UserRepository,
) *AddMemberUseCase {
	return &AddMemberUseCase{
		circleRepository: circleRepository,
		userRepository:   userRepository,
	}
}

func (uc *AddMemberUseCase) Execute(input AddMemberInput) error {
	// CircleIDを再構成
	circleID, err := entity.ReconstructCircleID(input.CircleID)
	if err != nil {
		return err
	}

	// UserIDを再構成
	userID, err := entity.ReconstructUserID(input.UserID)
	if err != nil {
		return err
	}

	// サークルを取得
	circle, err := uc.circleRepository.FindByID(circleID)
	if err != nil {
		return err
	}
	if circle == nil {
		return errors.New("circle not found")
	}

	// ユーザーの存在確認
	user, err := uc.userRepository.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 基本的なバリデーション
	if circle.IsOwner(userID) {
		return errors.New("owner cannot be a member")
	}
	if circle.IsMember(userID) {
		return nil // 既にメンバーの場合はエラーではない
	}

	// オーナーを取得
	owner, err := uc.userRepository.FindByID(circle.OwnerID())
	if err != nil {
		return err
	}
	if owner == nil {
		return errors.New("owner not found")
	}

	// メンバーリストを構築
	var members []*entity.User
	for _, memberID := range circle.GetMemberIDs() {
		member, err := uc.userRepository.FindByID(memberID)
		if err != nil {
			return err
		}
		if member != nil {
			members = append(members, member)
		}
	}

	// プレミアム制限をチェック
	circleMembers := entity.NewCircleMembers(owner, members)
	limitSpec := specification.NewCircleMemberLimitSpecification()
	if !limitSpec.IsSatisfiedBy(circleMembers) {
		maxLimit := limitSpec.GetMaxLimit(circleMembers)
		return fmt.Errorf("circle is full: maximum %d participants (including owner) allowed", maxLimit)
	}

	// メンバーを追加
	circle.AddMember(userID)

	// 保存
	return uc.circleRepository.Save(circle)
}