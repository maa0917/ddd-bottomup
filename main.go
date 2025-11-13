package main

import (
	"ddd-bottomup/domain/service"
	"ddd-bottomup/infrastructure"
	"ddd-bottomup/presentation/router"
	"ddd-bottomup/usecase"
	"log"
	"net/http"
)

type Application struct {
	CreateUserUseCase *usecase.CreateUserUseCase
	GetUserUseCase    *usecase.GetUserUseCase
	UpdateUserUseCase *usecase.UpdateUserUseCase
	DeleteUserUseCase *usecase.DeleteUserUseCase
}

func main() {
	log.Println("Starting DDD Bottom-Up HTTP Server...")

	app, err := setupApplication()
	if err != nil {
		log.Fatalf("Failed to setup application: %v", err)
	}

	log.Println("Application setup completed successfully!")

	// HTTPルーターの設定
	mux := router.NewRouter(
		app.CreateUserUseCase,
		app.GetUserUseCase,
		app.UpdateUserUseCase,
		app.DeleteUserUseCase,
	)

	// HTTPサーバー起動
	port := ":8080"
	log.Printf("Starting HTTP server on port %s", port)
	log.Println("Available endpoints:")
	log.Println("  POST   /users      - Create user")
	log.Println("  GET    /users/{id} - Get user")
	log.Println("  PUT    /users/{id} - Update user")
	log.Println("  DELETE /users/{id} - Delete user")
	log.Println("  GET    /health     - Health check")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func setupApplication() (*Application, error) {
	log.Println("Setting up application dependencies...")

	// 1. リポジトリ層の初期化
	log.Println("Initializing repositories...")
	userRepo := infrastructure.NewMemoryUserRepository()

	// 2. ドメインサービス層の初期化
	log.Println("Initializing domain services...")
	userExistenceService := service.NewUserExistenceService(userRepo)

	// 3. ユースケース層の初期化
	log.Println("Initializing use cases...")
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo, userExistenceService)
	getUserUseCase := usecase.NewGetUserUseCase(userRepo)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userRepo, userExistenceService)
	deleteUserUseCase := usecase.NewDeleteUserUseCase(userRepo)

	return &Application{
		CreateUserUseCase: createUserUseCase,
		GetUserUseCase:    getUserUseCase,
		UpdateUserUseCase: updateUserUseCase,
		DeleteUserUseCase: deleteUserUseCase,
	}, nil
}

func testApplication(app *Application) error {
	log.Println("Running application tests...")

	// テスト1: ユーザー作成
	log.Println("Test 1: Creating user...")
	createInput := usecase.CreateUserInput{
		FirstName: "太郎",
		LastName:  "田中",
		Email:     "taro@example.com",
	}

	createOutput, err := app.CreateUserUseCase.Execute(createInput)
	if err != nil {
		return err
	}

	userID := createOutput.UserID
	log.Printf("✓ User created successfully: ID=%s", userID)

	// テスト2: ユーザー取得
	log.Println("Test 2: Getting user...")
	getInput := usecase.GetUserInput{UserID: userID}
	getOutput, err := app.GetUserUseCase.Execute(getInput)
	if err != nil {
		return err
	}

	log.Printf("✓ User retrieved: %s %s (%s)", 
		getOutput.FirstName, getOutput.LastName, getOutput.Email)

	// テスト3: ユーザー更新
	log.Println("Test 3: Updating user...")
	firstName := "次郎"
	email := "jiro@example.com"
	updateInput := usecase.UpdateUserInput{
		UserID:    userID,
		FirstName: &firstName,
		Email:     &email,
	}

	updateOutput, err := app.UpdateUserUseCase.Execute(updateInput)
	if err != nil {
		return err
	}

	log.Printf("✓ User updated: %s %s (%s)", 
		updateOutput.FirstName, updateOutput.LastName, updateOutput.Email)

	// テスト4: 重複チェック
	log.Println("Test 4: Testing duplicate name check...")
	duplicateInput := usecase.CreateUserInput{
		FirstName: "次郎",
		LastName:  "田中",
		Email:     "another@example.com",
	}

	_, err = app.CreateUserUseCase.Execute(duplicateInput)
	if err == nil {
		log.Println("⚠ Warning: Duplicate name check might not be working")
	} else {
		log.Printf("✓ Duplicate check working: %v", err)
	}

	// テスト5: ユーザー削除
	log.Println("Test 5: Deleting user...")
	deleteInput := usecase.DeleteUserInput{UserID: userID}
	err = app.DeleteUserUseCase.Execute(deleteInput)
	if err != nil {
		return err
	}

	log.Printf("✓ User deleted successfully")

	// テスト6: 削除後の取得確認
	log.Println("Test 6: Confirming user deletion...")
	_, err = app.GetUserUseCase.Execute(getInput)
	if err == nil {
		log.Println("⚠ Warning: User should not exist after deletion")
	} else {
		log.Printf("✓ User not found after deletion: %v", err)
	}

	return nil
}