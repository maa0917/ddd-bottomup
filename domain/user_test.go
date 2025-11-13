package domain

import (
	"testing"
)

// UserID tests
func TestNewUserID_GeneratesValidUUID(t *testing.T) {
	userID := NewUserID()

	if userID == nil {
		t.Fatal("Expected UserID, but got nil")
	}

	if len(userID.Value()) != 36 {
		t.Errorf("Expected UUID length 36, but got %d", len(userID.Value()))
	}

	// UUIDフォーマットの基本チェック（8-4-4-4-12）
	value := userID.Value()
	if value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		t.Errorf("Expected UUID format, but got %s", value)
	}
}

func TestReconstructUserID_ValidUUID_Success(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	userID, err := ReconstructUserID(validUUID)
	if err != nil {
		t.Errorf("Expected no error for valid UUID, but got: %v", err)
	}
	if userID == nil {
		t.Fatal("Expected UserID, but got nil")
	}
	if userID.Value() != validUUID {
		t.Errorf("Expected UUID %s, but got %s", validUUID, userID.Value())
	}
}

func TestReconstructUserID_EmptyString_ReturnsError(t *testing.T) {
	userID, err := ReconstructUserID("")

	if err == nil {
		t.Error("Expected error for empty string, but got none")
	}
	if userID != nil {
		t.Error("Expected nil UserID for empty string")
	}

	// エラー型の確認
	if emptyFieldErr, ok := err.(EmptyFieldError); !ok {
		t.Errorf("Expected EmptyFieldError, but got %T", err)
	} else if emptyFieldErr.Field != "user ID" {
		t.Errorf("Expected field 'user ID', but got '%s'", emptyFieldErr.Field)
	}
}

func TestReconstructUserID_InvalidUUID_ReturnsError(t *testing.T) {
	invalidUUIDs := []string{
		"invalid-uuid",
		"12345",
		"550e8400-e29b-41d4-a716-44665544000",   // 1文字少ない
		"550e8400-e29b-41d4-a716-4466554400000", // 1文字多い
	}

	for _, uuid := range invalidUUIDs {
		t.Run(uuid, func(t *testing.T) {
			userID, err := ReconstructUserID(uuid)

			if err == nil {
				t.Errorf("Expected error for invalid UUID %s, but got none", uuid)
			}
			if userID != nil {
				t.Error("Expected nil UserID for invalid UUID")
			}

			// エラー型の確認
			if invalidUserIDErr, ok := err.(InvalidUserIDError); !ok {
				t.Errorf("Expected InvalidUserIDError, but got %T", err)
			} else if invalidUserIDErr.Value != uuid {
				t.Errorf("Expected error value '%s', but got '%s'", uuid, invalidUserIDErr.Value)
			}
		})
	}
}

func TestUserID_Equals(t *testing.T) {
	userID1 := NewUserID()
	userID2, _ := ReconstructUserID(userID1.Value())
	userID3 := NewUserID()

	// 同じ値のUserID
	if !userID1.Equals(userID2) {
		t.Error("Expected equal UserIDs to return true")
	}

	// 異なる値のUserID
	if userID1.Equals(userID3) {
		t.Error("Expected different UserIDs to return false")
	}

	// nilとの比較
	if userID1.Equals(nil) {
		t.Error("Expected UserID compared to nil to return false")
	}
}

// User tests
func TestNewUser_ValidData_Success(t *testing.T) {
	name, _ := NewFullName("太郎", "田中")
	email, _ := NewEmail("taro@example.com")

	user := NewUser(name, email, true)

	if user == nil {
		t.Fatal("Expected User, but got nil")
	}
	if user.ID() == nil {
		t.Error("Expected user to have ID")
	}
	if !user.Name().Equals(name) {
		t.Error("Expected user name to match input")
	}
	if !user.Email().Equals(email) {
		t.Error("Expected user email to match input")
	}
	if !user.IsPremium() {
		t.Error("Expected user to be premium")
	}
}

func TestReconstructUser_ValidData_Success(t *testing.T) {
	userID, _ := ReconstructUserID("550e8400-e29b-41d4-a716-446655440000")
	name, _ := NewFullName("花子", "佐藤")
	email, _ := NewEmail("hanako@example.com")

	user := ReconstructUser(userID, name, email, false)

	if user == nil {
		t.Fatal("Expected User, but got nil")
	}
	if !user.ID().Equals(userID) {
		t.Error("Expected user ID to match input")
	}
	if !user.Name().Equals(name) {
		t.Error("Expected user name to match input")
	}
	if !user.Email().Equals(email) {
		t.Error("Expected user email to match input")
	}
	if user.IsPremium() {
		t.Error("Expected user to not be premium")
	}
}

func TestUser_ChangeName_Success(t *testing.T) {
	name, _ := NewFullName("太郎", "田中")
	email, _ := NewEmail("taro@example.com")
	user := NewUser(name, email, false)

	newName, _ := NewFullName("次郎", "田中")
	user.ChangeName(newName)

	if !user.Name().Equals(newName) {
		t.Error("Expected user name to be changed")
	}
}

func TestUser_ChangeEmail_Success(t *testing.T) {
	name, _ := NewFullName("太郎", "田中")
	email, _ := NewEmail("taro@example.com")
	user := NewUser(name, email, false)

	newEmail, _ := NewEmail("taro.tanaka@example.com")
	user.ChangeEmail(newEmail)

	if !user.Email().Equals(newEmail) {
		t.Error("Expected user email to be changed")
	}
}

func TestUser_Equals(t *testing.T) {
	userID, _ := ReconstructUserID("550e8400-e29b-41d4-a716-446655440000")
	name, _ := NewFullName("太郎", "田中")
	email, _ := NewEmail("taro@example.com")

	user1 := ReconstructUser(userID, name, email, false)
	user2 := ReconstructUser(userID, name, email, true) // 異なるpremium状態

	differentName, _ := NewFullName("花子", "田中")
	user3 := ReconstructUser(userID, differentName, email, false) // 異なる名前

	differentUserID, _ := ReconstructUserID("550e8400-e29b-41d4-a716-446655440001")
	user4 := ReconstructUser(differentUserID, name, email, false) // 異なるID

	// 同じIDのユーザー（他の属性が異なっても同じと判定される）
	if !user1.Equals(user2) {
		t.Error("Expected users with same ID to be equal")
	}
	if !user1.Equals(user3) {
		t.Error("Expected users with same ID to be equal regardless of other attributes")
	}

	// 異なるIDのユーザー
	if user1.Equals(user4) {
		t.Error("Expected users with different IDs to not be equal")
	}

	// nilとの比較
	if user1.Equals(nil) {
		t.Error("Expected user compared to nil to return false")
	}
}

// UserExistenceService tests
func TestUserExistenceService_Exists_UserFound(t *testing.T) {
	// メモリリポジトリを使用
	repo := &mockUserRepository{
		users: make(map[string]*User),
	}

	// テストユーザーを作成してリポジトリに保存
	name, _ := NewFullName("太郎", "田中")
	email, _ := NewEmail("taro@example.com")
	user := NewUser(name, email, false)
	repo.users[user.Name().String()] = user

	service := NewUserExistenceService(repo)

	exists, err := service.Exists(user)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if !exists {
		t.Error("Expected user to exist")
	}
}

func TestUserExistenceService_Exists_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		users: make(map[string]*User),
	}

	name, _ := NewFullName("花子", "佐藤")
	email, _ := NewEmail("hanako@example.com")
	user := NewUser(name, email, false)

	service := NewUserExistenceService(repo)

	exists, err := service.Exists(user)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if exists {
		t.Error("Expected user to not exist")
	}
}

// Mock repository for testing
type mockUserRepository struct {
	users map[string]*User
}

func (r *mockUserRepository) FindByID(id *UserID) (*User, error) {
	for _, user := range r.users {
		if user.ID().Equals(id) {
			return user, nil
		}
	}
	return nil, nil
}

func (r *mockUserRepository) FindByName(name *FullName) (*User, error) {
	if user, exists := r.users[name.String()]; exists {
		return user, nil
	}
	return nil, nil
}

func (r *mockUserRepository) Save(user *User) error {
	r.users[user.Name().String()] = user
	return nil
}

func (r *mockUserRepository) Delete(id *UserID) error {
	for key, user := range r.users {
		if user.ID().Equals(id) {
			delete(r.users, key)
			return nil
		}
	}
	return nil
}
