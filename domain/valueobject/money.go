package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

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