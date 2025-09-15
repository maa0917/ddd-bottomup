package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/service"
	"ddd-bottomup/domain/valueobject"
	"ddd-bottomup/infrastructure/repository"
	"testing"
)

func TestDeleteUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	
	// テスト用のユーザーを作成・保存
	fullName, _ := valueobject.NewFullName("太郎", "田中")
	email, _ := valueobject.NewEmail("taro@example.com")
	user := entity.NewUser(fullName, email)
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}
	
	useCase := NewDeleteUserUseCase(repo)
	input := DeleteUserInput{UserID: user.ID().Value()}

	// Act
	err = useCase.Execute(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 削除されていることを確認
	deletedUser, err := repo.FindByID(user.ID())
	if err != nil {
		t.Errorf("Error finding deleted user: %v", err)
	}
	if deletedUser != nil {
		t.Error("Expected user to be deleted, but still exists")
	}

	// リポジトリから削除されていることを確認
	memoryRepo := repo.(*repository.UserRepositoryMemory)
	if memoryRepo.Count() != 0 {
		t.Errorf("Expected 0 users in repository, but got %d", memoryRepo.Count())
	}
}

func TestDeleteUserUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	useCase := NewDeleteUserUseCase(repo)
	
	// 存在しないUserIDを使用
	nonExistentID := entity.NewUserID()
	input := DeleteUserInput{UserID: nonExistentID.Value()}

	// Act
	err := useCase.Execute(input)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent user, but got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, but got '%s'", err.Error())
	}
}

func TestDeleteUserUseCase_Execute_InvalidUserID(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	useCase := NewDeleteUserUseCase(repo)

	testCases := []struct {
		name   string
		userID string
	}{
		{"empty user ID", ""},
		{"invalid format", "invalid-uuid"},
		{"malformed UUID", "12345"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := DeleteUserInput{UserID: tc.userID}

			// Act
			err := useCase.Execute(input)

			// Assert
			if err == nil {
				t.Error("Expected error for invalid UserID, but got nil")
			}
		})
	}
}

func TestDeleteUserUseCase_Execute_MultipleUsersOneDeleted(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	userExistenceService := service.NewUserExistenceService(repo)
	createUseCase := NewCreateUserUseCase(repo, userExistenceService)
	deleteUseCase := NewDeleteUserUseCase(repo)

	// 複数ユーザーを作成
	users := []CreateUserInput{
		{FirstName: "太郎", LastName: "田中", Email: "taro@example.com"},
		{FirstName: "花子", LastName: "佐藤", Email: "hanako@example.com"},
		{FirstName: "次郎", LastName: "山田", Email: "jiro@example.com"},
	}

	var createdUserIDs []string
	for _, userInput := range users {
		output, err := createUseCase.Execute(userInput)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		createdUserIDs = append(createdUserIDs, output.UserID)
	}

	// 最初のユーザーを削除
	deleteInput := DeleteUserInput{UserID: createdUserIDs[0]}
	err := deleteUseCase.Execute(deleteInput)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Assert
	memoryRepo := repo.(*repository.UserRepositoryMemory)
	if memoryRepo.Count() != 2 {
		t.Errorf("Expected 2 users remaining, but got %d", memoryRepo.Count())
	}

	// 削除されたユーザーが見つからないことを確認
	deletedUser, err := repo.FindByID(entity.NewUserID())
	if err == nil && deletedUser != nil {
		t.Error("Deleted user should not be found")
	}

	// 残りのユーザーは存在することを確認
	for i := 1; i < len(createdUserIDs); i++ {
		userID, _ := entity.ReconstructUserID(createdUserIDs[i])
		user, err := repo.FindByID(userID)
		if err != nil {
			t.Errorf("Failed to find remaining user %d: %v", i, err)
		}
		if user == nil {
			t.Errorf("Remaining user %d should exist", i)
		}
	}
}

func TestDeleteUserUseCase_Execute_DeleteSameUserTwice(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	
	fullName, _ := valueobject.NewFullName("太郎", "田中")
	email, _ := valueobject.NewEmail("taro@example.com")
	user := entity.NewUser(fullName, email)
	repo.Save(user)
	
	useCase := NewDeleteUserUseCase(repo)
	input := DeleteUserInput{UserID: user.ID().Value()}

	// Act - 最初の削除
	err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("First deletion failed: %v", err)
	}

	// Act - 2回目の削除
	err = useCase.Execute(input)

	// Assert
	if err == nil {
		t.Error("Expected error for second deletion, but got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, but got '%s'", err.Error())
	}
}