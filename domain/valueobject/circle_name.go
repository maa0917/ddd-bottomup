package valueobject

import (
	"errors"
	"strings"
)

type CircleName struct {
	value string
}

func NewCircleName(value string) (*CircleName, error) {
	if strings.TrimSpace(value) == "" {
		return nil, errors.New("circle name cannot be empty")
	}

	if len(value) > 50 {
		return nil, errors.New("circle name cannot exceed 50 characters")
	}

	if len(value) < 3 {
		return nil, errors.New("circle name must be at least 3 characters")
	}

	return &CircleName{value: strings.TrimSpace(value)}, nil
}

func (c *CircleName) Value() string {
	return c.value
}

func (c *CircleName) Equals(other *CircleName) bool {
	if other == nil {
		return false
	}
	return c.value == other.value
}

func (c *CircleName) String() string {
	return c.value
}