package usecase

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/service"
	"ddd-bottomup/domain/valueobject"
	"ddd-bottomup/infrastructure"
	"testing"
)

func TestUpdateUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()

	// 既存ユーザーを作成・保存
	originalName, _ := valueobject.NewFullName("太郎", "田中")
	user := entity.NewUser(originalName)
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	useCase := NewUpdateUserUseCase(repo)
	input := UpdateUserInput{
		UserID:    user.ID().Value(),
		FirstName: "次郎",
		LastName:  "佐藤",
	}

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

	if output.FirstName != "次郎" {
		t.Errorf("Expected FirstName '次郎', but got '%s'", output.FirstName)
	}

	if output.LastName != "佐藤" {
		t.Errorf("Expected LastName '佐藤', but got '%s'", output.LastName)
	}

	// リポジトリからも確認
	updatedUser, err := repo.FindByID(user.ID())
	if err != nil {
		t.Errorf("Failed to find updated user: %v", err)
	}

	if updatedUser.Name().FirstName() != "次郎" {
		t.Errorf("Expected updated FirstName '次郎', but got '%s'", updatedUser.Name().FirstName())
	}
}

func TestUpdateUserUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	useCase := NewUpdateUserUseCase(repo)

	nonExistentID := entity.NewUserID()
	input := UpdateUserInput{
		UserID:    nonExistentID.Value(),
		FirstName: "太郎",
		LastName:  "田中",
	}

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

func TestUpdateUserUseCase_Execute_DuplicateName(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()

	// 2人のユーザーを作成
	user1Name, _ := valueobject.NewFullName("太郎", "田中")
	user1 := entity.NewUser(user1Name)
	repo.Save(user1)

	user2Name, _ := valueobject.NewFullName("花子", "佐藤")
	user2 := entity.NewUser(user2Name)
	repo.Save(user2)

	useCase := NewUpdateUserUseCase(repo)

	// user2の名前をuser1と同じにしようとする
	input := UpdateUserInput{
		UserID:    user2.ID().Value(),
		FirstName: "太郎",
		LastName:  "田中",
	}

	// Act
	output, err := useCase.Execute(input)

	// Assert
	if err == nil {
		t.Error("Expected error for duplicate name, but got nil")
	}

	if output != nil {
		t.Error("Expected no output for duplicate name, but got output")
	}

	if err.Error() != "name already exists" {
		t.Errorf("Expected 'name already exists' error, but got '%s'", err.Error())
	}

	// user2の名前が変更されていないことを確認
	unchangedUser, _ := repo.FindByID(user2.ID())
	if unchangedUser.Name().FirstName() != "花子" {
		t.Errorf("Expected unchanged FirstName '花子', but got '%s'", unchangedUser.Name().FirstName())
	}
}

func TestUpdateUserUseCase_Execute_SameNameUpdate(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()

	originalName, _ := valueobject.NewFullName("太郎", "田中")
	user := entity.NewUser(originalName)
	repo.Save(user)

	useCase := NewUpdateUserUseCase(repo)

	// 同じ名前に更新（自分自身なのでOK）
	input := UpdateUserInput{
		UserID:    user.ID().Value(),
		FirstName: "太郎",
		LastName:  "田中",
	}

	// Act
	output, err := useCase.Execute(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for same name update, but got: %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, but got nil")
	}

	if output.FirstName != "太郎" || output.LastName != "田中" {
		t.Errorf("Expected name unchanged, but got %s %s", output.FirstName, output.LastName)
	}
}

func TestUpdateUserUseCase_Execute_InvalidName(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()

	originalName, _ := valueobject.NewFullName("太郎", "田中")
	user := entity.NewUser(originalName)
	repo.Save(user)

	useCase := NewUpdateUserUseCase(repo)

	testCases := []struct {
		name      string
		firstName string
		lastName  string
	}{
		{"empty first name", "", "田中"},
		{"empty last name", "太郎", ""},
		{"both empty", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := UpdateUserInput{
				UserID:    user.ID().Value(),
				FirstName: tc.firstName,
				LastName:  tc.lastName,
			}

			// Act
			output, err := useCase.Execute(input)

			// Assert
			if err == nil {
				t.Error("Expected error for invalid name, but got nil")
			}

			if output != nil {
				t.Error("Expected no output for invalid name, but got output")
			}
		})
	}
}
