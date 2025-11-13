package usecase

import (
	"ddd-bottomup/domain"
	"ddd-bottomup/infrastructure"
	"strings"
	"testing"
)

func TestCreateUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := domain.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)

	input := CreateUserInput{
		FirstName: "太郎",
		LastName:  "田中",
		Email:     "taro@example.com",
		IsPremium: false,
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

	if output.UserID == "" {
		t.Error("Expected UserID to be set, but got empty string")
	}
}

func TestCreateUserUseCase_Execute_DuplicateName_ReturnsError(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := domain.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)

	input := CreateUserInput{
		FirstName: "太郎",
		LastName:  "田中",
		Email:     "taro@example.com",
		IsPremium: false,
	}

	// 最初のユーザーを作成
	_, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// 同じ名前で異なるメールの2番目のユーザーを作成試行
	duplicateInput := CreateUserInput{
		FirstName: "太郎",
		LastName:  "田中",
		Email:     "taro2@example.com",
		IsPremium: false,
	}

	// Act
	output, err := useCase.Execute(duplicateInput)

	// Assert
	if err == nil {
		t.Error("Expected error for duplicate name, but got nil")
	}

	if duplicateNameErr, ok := err.(domain.DuplicateUserNameError); !ok {
		t.Errorf("Expected DuplicateUserNameError, but got %T", err)
	} else if !strings.Contains(duplicateNameErr.Name, "太郎") {
		t.Errorf("Expected error message to contain '太郎', but got '%s'", duplicateNameErr.Name)
	}

	if output != nil {
		t.Error("Expected no output for duplicate name, but got output")
	}
}

func TestCreateUserUseCase_Execute_InvalidEmail_ReturnsError(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := domain.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)

	tests := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"invalid format", "invalid-email"},
		{"missing @", "testexample.com"},
		{"missing domain", "test@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CreateUserInput{
				FirstName: "太郎",
				LastName:  "田中",
				Email:     tt.email,
				IsPremium: false,
			}

			// Act
			output, err := useCase.Execute(input)

			// Assert
			if err == nil {
				t.Errorf("Expected error for invalid email %q, but got nil", tt.email)
			}

			if output != nil {
				t.Error("Expected no output for invalid email, but got output")
			}
		})
	}
}

func TestCreateUserUseCase_Execute_InvalidNames_ReturnsError(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := domain.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)

	tests := []struct {
		name      string
		firstName string
		lastName  string
	}{
		{"empty first name", "", "田中"},
		{"empty last name", "太郎", ""},
		{"both empty", "", ""},
		{"whitespace only first name", "   ", "田中"},
		{"whitespace only last name", "太郎", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CreateUserInput{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Email:     "test@example.com",
				IsPremium: false,
			}

			// Act
			output, err := useCase.Execute(input)

			// Assert
			if err == nil {
				t.Errorf("Expected error for invalid names %q %q, but got nil", tt.firstName, tt.lastName)
			}

			if _, ok := err.(domain.EmptyFieldError); ok {
				// EmptyFieldErrorが期待される
			} else {
				t.Errorf("Expected EmptyFieldError, but got %T: %v", err, err)
			}

			if output != nil {
				t.Error("Expected no output for invalid names, but got output")
			}
		})
	}
}

func TestCreateUserUseCase_Execute_PremiumUser_Success(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := domain.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)

	input := CreateUserInput{
		FirstName: "花子",
		LastName:  "佐藤",
		Email:     "hanako@example.com",
		IsPremium: true,
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

	if output.UserID == "" {
		t.Error("Expected UserID to be set, but got empty string")
	}

	// プレミアムユーザーとして保存されているか確認（リポジトリから取得して確認）
	userID, _ := domain.ReconstructUserID(output.UserID)
	savedUser, err := repo.FindByID(userID)
	if err != nil {
		t.Errorf("Failed to retrieve saved user: %v", err)
	}

	if savedUser == nil {
		t.Error("Expected saved user, but got nil")
	} else if !savedUser.IsPremium() {
		t.Error("Expected user to be premium, but was not")
	}
}
