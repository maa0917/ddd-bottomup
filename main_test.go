package main

import (
	"ddd-bottomup/usecase"
	"testing"
)

func TestSetupApplication_Success(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "正常にアプリケーションがセットアップされること"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			app, err := setupApplication()

			// Assert
			if err != nil {
				t.Errorf("setupApplication()でエラーが発生しました: %v", err)
			}

			if app == nil {
				t.Fatal("アプリケーションがnilです")
			}

			if app.CreateUserUseCase == nil {
				t.Error("CreateUserUseCaseがnilです")
			}

			if app.GetUserUseCase == nil {
				t.Error("GetUserUseCaseがnilです")
			}

			if app.UpdateUserUseCase == nil {
				t.Error("UpdateUserUseCaseがnilです")
			}

			if app.DeleteUserUseCase == nil {
				t.Error("DeleteUserUseCaseがnilです")
			}
		})
	}
}

func TestTestApplication_UserLifecycle(t *testing.T) {
	tests := []struct {
		name              string
		createInput       usecase.CreateUserInput
		updateFirstName   string
		updateLastName    string
		updateEmail       string
		expectCreateError bool
	}{
		{
			name: "ユーザーの作成・取得・更新・削除が正常に行えること",
			createInput: usecase.CreateUserInput{
				FirstName: "太郎",
				LastName:  "田中",
				Email:     "taro@example.com",
				IsPremium: false,
			},
			updateFirstName:   "次郎",
			updateLastName:    "佐藤",
			updateEmail:       "jiro@example.com",
			expectCreateError: false,
		},
		{
			name: "プレミアムユーザーのライフサイクルが正常に動作すること",
			createInput: usecase.CreateUserInput{
				FirstName: "花子",
				LastName:  "佐藤",
				Email:     "hanako@example.com",
				IsPremium: true,
			},
			updateFirstName:   "美咲",
			updateLastName:    "山田",
			updateEmail:       "misaki@example.com",
			expectCreateError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			app, err := setupApplication()
			if err != nil {
				t.Fatalf("アプリケーションのセットアップに失敗しました: %v", err)
			}

			// Test 1: ユーザー作成
			createOutput, err := app.CreateUserUseCase.Execute(tt.createInput)
			if tt.expectCreateError {
				if err == nil {
					t.Error("ユーザー作成でエラーが期待されましたが、エラーが発生しませんでした")
				}
				return
			}
			if err != nil {
				t.Errorf("ユーザー作成でエラーが発生しました: %v", err)
				return
			}
			if createOutput == nil {
				t.Fatal("作成結果がnilです")
			}
			if createOutput.UserID == "" {
				t.Error("UserIDが空です")
			}

			userID := createOutput.UserID

			// Test 2: ユーザー取得
			getInput := usecase.GetUserInput{UserID: userID}
			getOutput, err := app.GetUserUseCase.Execute(getInput)
			if err != nil {
				t.Errorf("ユーザー取得でエラーが発生しました: %v", err)
				return
			}
			if getOutput == nil {
				t.Fatal("取得結果がnilです")
			}
			if getOutput.FirstName != tt.createInput.FirstName {
				t.Errorf("FirstNameが一致しません: expected=%s, got=%s",
					tt.createInput.FirstName, getOutput.FirstName)
			}
			if getOutput.LastName != tt.createInput.LastName {
				t.Errorf("LastNameが一致しません: expected=%s, got=%s",
					tt.createInput.LastName, getOutput.LastName)
			}
			if getOutput.Email != tt.createInput.Email {
				t.Errorf("Emailが一致しません: expected=%s, got=%s",
					tt.createInput.Email, getOutput.Email)
			}

			// Test 3: ユーザー更新
			updateInput := usecase.UpdateUserInput{
				UserID:    userID,
				FirstName: &tt.updateFirstName,
				LastName:  &tt.updateLastName,
				Email:     &tt.updateEmail,
			}
			updateOutput, err := app.UpdateUserUseCase.Execute(updateInput)
			if err != nil {
				t.Errorf("ユーザー更新でエラーが発生しました: %v", err)
				return
			}
			if updateOutput == nil {
				t.Fatal("更新結果がnilです")
			}
			if updateOutput.FirstName != tt.updateFirstName {
				t.Errorf("更新後のFirstNameが一致しません: expected=%s, got=%s",
					tt.updateFirstName, updateOutput.FirstName)
			}
			if updateOutput.LastName != tt.updateLastName {
				t.Errorf("更新後のLastNameが一致しません: expected=%s, got=%s",
					tt.updateLastName, updateOutput.LastName)
			}
			if updateOutput.Email != tt.updateEmail {
				t.Errorf("更新後のEmailが一致しません: expected=%s, got=%s",
					tt.updateEmail, updateOutput.Email)
			}

			// Test 4: 重複名チェック（更新後の名前と同じ名前で新規作成試行）
			duplicateInput := usecase.CreateUserInput{
				FirstName: tt.updateFirstName,
				LastName:  tt.updateLastName,
				Email:     "duplicate@example.com",
				IsPremium: false,
			}

			_, err = app.CreateUserUseCase.Execute(duplicateInput)
			if err == nil {
				t.Error("重複ユーザー作成でエラーが期待されましたが、エラーが発生しませんでした")
			}

			// Test 5: ユーザー削除
			deleteInput := usecase.DeleteUserInput{UserID: userID}
			err = app.DeleteUserUseCase.Execute(deleteInput)
			if err != nil {
				t.Errorf("ユーザー削除でエラーが発生しました: %v", err)
				return
			}

			// Test 6: 削除後の取得確認
			_, err = app.GetUserUseCase.Execute(getInput)
			if err == nil {
				t.Error("削除後のユーザーが取得できました（取得できないはずです）")
			}
		})
	}
}

func TestTestApplication_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		createInput usecase.CreateUserInput
		expectError bool
	}{
		{
			name: "空のメールアドレスでエラーが発生すること",
			createInput: usecase.CreateUserInput{
				FirstName: "太郎",
				LastName:  "田中",
				Email:     "",
				IsPremium: false,
			},
			expectError: true,
		},
		{
			name: "無効なメールアドレスでエラーが発生すること",
			createInput: usecase.CreateUserInput{
				FirstName: "太郎",
				LastName:  "田中",
				Email:     "invalid-email",
				IsPremium: false,
			},
			expectError: true,
		},
		{
			name: "空の名前でエラーが発生すること",
			createInput: usecase.CreateUserInput{
				FirstName: "",
				LastName:  "田中",
				Email:     "test@example.com",
				IsPremium: false,
			},
			expectError: true,
		},
		{
			name: "空の姓でエラーが発生すること",
			createInput: usecase.CreateUserInput{
				FirstName: "太郎",
				LastName:  "",
				Email:     "test@example.com",
				IsPremium: false,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			app, err := setupApplication()
			if err != nil {
				t.Fatalf("アプリケーションのセットアップに失敗しました: %v", err)
			}

			// Act
			_, err = app.CreateUserUseCase.Execute(tt.createInput)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}
		})
	}
}
