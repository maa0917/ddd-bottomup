package domain

import (
	"net/http"

	"github.com/google/uuid"
)

type UserID struct {
	value string
}

func NewUserID() *UserID {
	return &UserID{value: uuid.New().String()}
}

func ReconstructUserID(value string) (*UserID, error) {
	if value == "" {
		return nil, EmptyFieldError{Field: "user ID"}
	}
	if _, err := uuid.Parse(value); err != nil {
		return nil, InvalidUserIDError{Value: value}
	}
	return &UserID{value: value}, nil
}

func (u *UserID) Value() string {
	return u.value
}

func (u *UserID) Equals(other *UserID) bool {
	if other == nil {
		return false
	}
	return u.value == other.value
}

func (u *UserID) String() string {
	return u.value
}

type User struct {
	id        *UserID
	name      *FullName
	email     *Email
	isPremium bool
}

func NewUser(name *FullName, email *Email, isPremium bool) *User {
	return &User{
		id:        NewUserID(),
		name:      name,
		email:     email,
		isPremium: isPremium,
	}
}

func ReconstructUser(id *UserID, name *FullName, email *Email, isPremium bool) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		isPremium: isPremium,
	}
}

func (u *User) ID() *UserID {
	return u.id
}

func (u *User) Name() *FullName {
	return u.name
}

func (u *User) Email() *Email {
	return u.email
}

func (u *User) ChangeName(name *FullName) {
	u.name = name
}

func (u *User) ChangeEmail(email *Email) {
	u.email = email
}

func (u *User) IsPremium() bool {
	return u.isPremium
}

func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id.Equals(other.id)
}

// UserExistenceService - ユーザー存在確認サービス
type UserExistenceService struct {
	userRepository UserRepository
}

func NewUserExistenceService(userRepository UserRepository) *UserExistenceService {
	return &UserExistenceService{
		userRepository: userRepository,
	}
}

func (s *UserExistenceService) Exists(user *User) (bool, error) {
	existingUser, err := s.userRepository.FindByName(user.Name())
	if err != nil {
		return false, err
	}
	// 同じユーザーIDの場合は重複ではない
	if existingUser != nil && existingUser.ID().Equals(user.ID()) {
		return false, nil
	}
	return existingUser != nil, nil
}

// User related errors
type UserNotFoundError struct {
	ID string
}

func (e UserNotFoundError) Error() string {
	return "user not found: " + e.ID
}

func (e UserNotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

type UserAlreadyExistsError struct {
	Email string
}

func (e UserAlreadyExistsError) Error() string {
	return "user with email already exists: " + e.Email
}

func (e UserAlreadyExistsError) HTTPStatus() int {
	return http.StatusBadRequest
}

type DuplicateUserNameError struct {
	Name string
}

func (e DuplicateUserNameError) Error() string {
	return "user with name already exists: " + e.Name
}

func (e DuplicateUserNameError) HTTPStatus() int {
	return http.StatusBadRequest
}

type InvalidUserIDError struct {
	Value string
}

func (e InvalidUserIDError) Error() string {
	return "invalid user ID: " + e.Value
}

func (e InvalidUserIDError) HTTPStatus() int {
	return http.StatusBadRequest
}
