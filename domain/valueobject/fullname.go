package valueobject

import (
	"errors"
	"strings"
)

type FullName struct {
	firstName string
	lastName  string
}

/*
値オブジェクトにすべきかどうかの判断基準（以下を満たす）
1. そこにルールが存在しているか
2. それ単体で取り扱いたいか

*/

func NewFullName(firstName, lastName string) (*FullName, error) {
	if err := validateName(firstName, "first name"); err != nil {
		return nil, err
	}
	if err := validateName(lastName, "last name"); err != nil {
		return nil, err
	}

	return &FullName{
		firstName: strings.TrimSpace(firstName),
		lastName:  strings.TrimSpace(lastName),
	}, nil
}

func (f *FullName) FirstName() string {
	return f.firstName
}

func (f *FullName) LastName() string {
	return f.lastName
}

func (f *FullName) String() string {
	return f.firstName + " " + f.lastName
}

func (f *FullName) Equals(other *FullName) bool {
	if other == nil {
		return false
	}
	return f.firstName == other.firstName && f.lastName == other.lastName
}

func validateName(name string, nameType string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.New(nameType + " cannot be empty")
	}

	if len(trimmed) > 50 {
		return errors.New(nameType + " cannot exceed 50 characters")
	}

	return nil
}
