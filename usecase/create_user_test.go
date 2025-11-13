package usecase

import (
	"ddd-bottomup/domain"
	"ddd-bottomup/infrastructure"
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
		E
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

	// リポジトリに保存されているか確認
	memoryRepo := repo.(*repository.UserRepositoryMemory)
	memoryRepo := repo.(*infrastructure.MemoryUserRepository)
		t.Errorf("Expected 1 user in repository, but got %d", memoryRepo.Count())
	}
}

func TestCreateUserUseCase_Execute_DuplicateName(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := service.NewUserExistenceService(repo)
	userExistenceService := domain.NewUserExistenceService(repo)

	input := CreateUserInput{
		FirstName: "太郎",
		LastName:  "田中",
	}
		LastName:  "田

	// 最初のユーザーを作成
	_, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Act - 同じ名前で再度作成を試行
	output, err := useCase.Execute(input)

	// Assert
	if err == nil {
		t.Error("Expected error for duplicate name, but got nil")
	}

	if output != nil {
		t.Error("Expected no output for duplicate name, but got output")
	}

	// リポジトリには1人だけ保存されている
	memoryRepo := repo.(*repository.UserRepositoryMemory)
	if memoryRepo.Count() != 1 {
	memoryRepo := repo.(*infrastructure.MemoryUserRepository)
	}
}

func TestCreateUserUseCase_Execute_InvalidName(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := service.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)
	userExistenceService := domain.NewUserExistenceService(repo)
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
			input := CreateUserInput{
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

func TestCreateUserUseCase_Execute_MultipleUsers(t *testing.T) {
	// Arrange
	repo := infrastructure.NewMemoryUserRepository()
	userExistenceService := service.NewUserExistenceService(repo)
	useCase := NewCreateUserUseCase(repo, userExistenceService)
	userExistenceService := domain.NewUserExistenceService(repo)
	users := []CreateUserInput{
		{FirstName: "太郎", LastName: "田中"},
		{FirstName: "花子", LastName: "佐藤"},
		{FirstName: "次郎", LastName: "山田"},
	}

	// Act
	var outputs []*CreateUserOutput
	for _, input := range users {
		output, err := useCase.Execute(input)
		if err != nil {
			t.Errorf("Failed to create user %s %s: %v", input.FirstName, input.LastName, err)
			continue
		}
		outputs = append(outputs, output)
	}

	// Assert
	if len(outputs) != 3 {
		t.Errorf("Expected 3 users created, but got %d", len(outputs))
	}

	memoryRepo := repo.(*repository.UserRepositoryMemory)
	if memoryRepo.Count() != 3 {
	memoryRepo := repo.(*infrastructure.MemoryUserRepository)
	}

	// 各ユーザーIDが一意であることを確認
	userIDs := make(map[string]bool)
	for _, output := range outputs {
		if userIDs[output.UserID] {
			t.Errorf("Duplicate UserID found: %s", output.UserID)
		}
		userIDs[output.UserID] = true
	}
}