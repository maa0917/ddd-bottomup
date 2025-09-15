package valueobject

import (
	"errors"
	"fmt"
	"regexp"
)

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	if err := validateEmail(value); err != nil {
		return nil, err
	}
	return &Email{value: value}, nil
}

func (e *Email) Value() string {
	return e.value
}

func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

func (e *Email) String() string {
	return e.value
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	
	return nil
}