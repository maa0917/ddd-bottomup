package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/service"
	"ddd-bottomup/domain/valueobject"
	"ddd-bottomup/infrastructure/repository"
	"testing"
)

func TestGetUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	
	// テスト用のユーザーを作成・保存
	fullName, _ := valueobject.NewFullName("太郎", "田中")
	user := entity.NewUser(fullName)
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}
	
	useCase := NewGetUserUseCase(repo)
	input := GetUserInput{UserID: user.ID().Value()}

	// Act
	output, err := useCase.Execute(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, but got nil")
	}

	if output.UserID != user.ID().Value() {
		t.Errorf("Expected UserID %s, but got %s", user.ID().Value(), output.UserID)
	}

	if output.FirstName != "太郎" {
		t.Errorf("Expected FirstName '太郎', but got '%s'", output.FirstName)
	}

	if output.LastName != "田中" {
		t.Errorf("Expected LastName '田中', but got '%s'", output.LastName)
	}
}

func TestGetUserUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	useCase := NewGetUserUseCase(repo)
	
	// 存在しないUserIDを使用
	nonExistentID := entity.NewUserID()
	input := GetUserInput{UserID: nonExistentID.Value()}

	// Act
	output, err := useCase.Execute(input)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent user, but got nil")
	}

	if output != nil {
		t.Error("Expected no output for non-existent user, but got output")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, but got '%s'", err.Error())
	}
}

func TestGetUserUseCase_Execute_InvalidUserID(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	useCase := NewGetUserUseCase(repo)

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
			input := GetUserInput{UserID: tc.userID}

			// Act
			output, err := useCase.Execute(input)

			// Assert
			if err == nil {
				t.Error("Expected error for invalid UserID, but got nil")
			}

			if output != nil {
				t.Error("Expected no output for invalid UserID, but got output")
			}
		})
	}
}

func TestGetUserUseCase_Execute_MultipleUsers(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepositoryMemory()
	userExistenceService := service.NewUserExistenceService(repo)
	createUseCase := NewCreateUserUseCase(repo, userExistenceService)
	getUserUseCase := NewGetUserUseCase(repo)

	// 複数ユーザーを作成
	users := []CreateUserInput{
		{FirstName: "太郎", LastName: "田中"},
		{FirstName: "花子", LastName: "佐藤"},
		{FirstName: "次郎", LastName: "山田"},
	}

	var createdUserIDs []string
	for _, userInput := range users {
		output, err := createUseCase.Execute(userInput)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		createdUserIDs = append(createdUserIDs, output.UserID)
	}

	// Act & Assert - 各ユーザーを取得
	for i, userID := range createdUserIDs {
		input := GetUserInput{UserID: userID}
		output, err := getUserUseCase.Execute(input)

		if err != nil {
			t.Errorf("Failed to get user %d: %v", i, err)
			continue
		}

		if output.UserID != userID {
			t.Errorf("Expected UserID %s, but got %s", userID, output.UserID)
		}

		expectedUser := users[i]
		if output.FirstName != expectedUser.FirstName {
			t.Errorf("Expected FirstName %s, but got %s", expectedUser.FirstName, output.FirstName)
		}

		if output.LastName != expectedUser.LastName {
			t.Errorf("Expected LastName %s, but got %s", expectedUser.LastName, output.LastName)
		}
	}
}