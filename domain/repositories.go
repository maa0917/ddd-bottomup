package domain

type UserRepository interface {
	FindByID(id *UserID) (*User, error)
	FindByName(name *FullName) (*User, error)
	Save(user *User) error
	Delete(id *UserID) error
}

type CircleRepository interface {
	FindByID(id *CircleID) (*Circle, error)
	FindByName(name *CircleName) (*Circle, error)
	FindAll() ([]*Circle, error)
	Save(circle *Circle) error
	Delete(id *CircleID) error
}
