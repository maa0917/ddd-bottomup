package repository

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/valueobject"
)

type UserRepository interface {
	FindByID(id *entity.UserID) (*entity.User, error)
	FindByName(name *valueobject.FullName) (*entity.User, error)
	Save(user *entity.User) error
	Delete(id *entity.UserID) error
}
