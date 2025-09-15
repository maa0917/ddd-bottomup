package repository

import (
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/specification"
	"ddd-bottomup/domain/valueobject"
)

type CircleRepository interface {
	FindByID(id *entity.CircleID) (*entity.Circle, error)
	FindByName(name *valueobject.CircleName) (*entity.Circle, error)
	FindAll() ([]*entity.Circle, error)
	FindBySpecification(spec specification.CircleSpecification) ([]*entity.Circle, error)
	Save(circle *entity.Circle) error
	Delete(id *entity.CircleID) error
}