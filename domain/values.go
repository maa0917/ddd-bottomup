package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Email 値オブジェクト
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

// FullName 値オブジェクト
type FullName struct {
	firstName string
	lastName  string
}

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

// CircleName 値オブジェクト
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

// Money 値オブジェクト
type Money struct {
	amount   int64  // 最小単位で保存（例：JPYなら円、USDならcents）
	currency string // 通貨コード（ISO 4217）
}

func NewMoney(amount int64, currency string) (*Money, error) {
	if err := validateCurrency(currency); err != nil {
		return nil, err
	}

	return &Money{
		amount:   amount,
		currency: strings.ToUpper(currency),
	}, nil
}

func (m *Money) Amount() int64 {
	return m.amount
}

func (m *Money) Currency() string {
	return m.currency
}

func (m *Money) String() string {
	switch m.currency {
	case "JPY":
		return fmt.Sprintf("¥%d", m.amount)
	case "USD":
		return fmt.Sprintf("$%.2f", float64(m.amount)/100)
	case "EUR":
		return fmt.Sprintf("€%.2f", float64(m.amount)/100)
	default:
		return fmt.Sprintf("%d %s", m.amount, m.currency)
	}
}

func (m *Money) Equals(other *Money) bool {
	if other == nil {
		return false
	}
	return m.amount == other.amount && m.currency == other.currency
}

func (m *Money) Add(other *Money) (*Money, error) {
	if err := m.validateSameCurrency(other); err != nil {
		return nil, err
	}

	return &Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

func (m *Money) Subtract(other *Money) (*Money, error) {
	if err := m.validateSameCurrency(other); err != nil {
		return nil, err
	}

	return &Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

func (m *Money) Multiply(multiplier int64) *Money {
	return &Money{
		amount:   m.amount * multiplier,
		currency: m.currency,
	}
}

func (m *Money) IsPositive() bool {
	return m.amount > 0
}

func (m *Money) IsNegative() bool {
	return m.amount < 0
}

func (m *Money) IsZero() bool {
	return m.amount == 0
}

func (m *Money) validateSameCurrency(other *Money) error {
	if other == nil {
		return errors.New("cannot operate with nil money")
	}
	if m.currency != other.currency {
		return fmt.Errorf("cannot operate between different currencies: %s and %s", m.currency, other.currency)
	}
	return nil
}

func validateCurrency(currency string) error {
	if currency == "" {
		return errors.New("currency cannot be empty")
	}

	currency = strings.ToUpper(currency)
	validCurrencies := map[string]bool{
		"JPY": true,
		"USD": true,
		"EUR": true,
		"GBP": true,
		"AUD": true,
		"CAD": true,
	}

	if !validCurrencies[currency] {
		return fmt.Errorf("unsupported currency: %s", currency)
	}

	return nil
}
