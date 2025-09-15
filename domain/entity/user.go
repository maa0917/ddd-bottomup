package entity

import (
	"ddd-bottomup/domain/valueobject"
	"errors"
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
		return nil, errors.New("user ID cannot be empty")
	}
	if _, err := uuid.Parse(value); err != nil {
		return nil, errors.New("invalid user ID format")
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
	name      *valueobject.FullName
	email     *valueobject.Email
	isPremium bool
}

func NewUser(name *valueobject.FullName, email *valueobject.Email, isPremium bool) *User {
	return &User{
		id:        NewUserID(),
		name:      name,
		email:     email,
		isPremium: isPremium,
	}
}

func ReconstructUser(id *UserID, name *valueobject.FullName, email *valueobject.Email, isPremium bool) *User {
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

func (u *User) Name() *valueobject.FullName {
	return u.name
}

func (u *User) Email() *valueobject.Email {
	return u.email
}

func (u *User) ChangeName(name *valueobject.FullName) {
	u.name = name
}

func (u *User) ChangeEmail(email *valueobject.Email) {
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
