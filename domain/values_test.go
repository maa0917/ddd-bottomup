package domain

import (
	"testing"
)

// Email tests
func TestNewEmail_ValidEmail_Success(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"Standard email", "test@example.com"},
		{"Email with subdomain", "user@mail.example.com"},
		{"Email with numbers", "user123@example123.com"},
		{"Email with dash", "test-user@example.com"},
		{"Email with plus", "test+tag@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)
			if err != nil {
				t.Errorf("Expected no error for valid email %q, but got: %v", tt.email, err)
			}
			if email == nil {
				t.Errorf("Expected email object, but got nil")
			}
			if email.Value() != tt.email {
				t.Errorf("Expected email value %q, but got %q", tt.email, email.Value())
			}
		})
	}
}

func TestNewEmail_InvalidEmail_ReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"Empty email", ""},
		{"Missing @", "testexample.com"},
		{"Missing domain", "test@"},
		{"Missing local part", "@example.com"},
		{"Invalid characters", "test@exam ple.com"},
		{"Multiple @", "test@@example.com"},
		{"Missing TLD", "test@example"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)
			if err == nil {
				t.Errorf("Expected error for invalid email %q, but got none", tt.email)
			}
			if email != nil {
				t.Errorf("Expected nil email for invalid input, but got: %v", email)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")
	email3, _ := NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Expected equal emails to return true")
	}
	if email1.Equals(email3) {
		t.Error("Expected different emails to return false")
	}
	if email1.Equals(nil) {
		t.Error("Expected email compared to nil to return false")
	}
}

// FullName tests
func TestNewFullName_ValidNames_Success(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
	}{
		{"Japanese names", "太郎", "田中"},
		{"English names", "John", "Doe"},
		{"Names with spaces", " Alice ", " Smith "},
		{"Single character", "A", "B"},
		{"Max length names", "12345678901234567890123456789012345678901234567890", "12345678901234567890123456789012345678901234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullName, err := NewFullName(tt.firstName, tt.lastName)
			if err != nil {
				t.Errorf("Expected no error for valid names %q %q, but got: %v", tt.firstName, tt.lastName, err)
			}
			if fullName == nil {
				t.Error("Expected FullName object, but got nil")
			}
		})
	}
}

func TestNewFullName_InvalidNames_ReturnsError(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
	}{
		{"Empty first name", "", "Smith"},
		{"Empty last name", "John", ""},
		{"Whitespace only first name", "   ", "Smith"},
		{"Whitespace only last name", "John", "   "},
		{"Too long first name", "123456789012345678901234567890123456789012345678901", "Smith"},
		{"Too long last name", "John", "123456789012345678901234567890123456789012345678901"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullName, err := NewFullName(tt.firstName, tt.lastName)
			if err == nil {
				t.Errorf("Expected error for invalid names %q %q, but got none", tt.firstName, tt.lastName)
			}
			if fullName != nil {
				t.Errorf("Expected nil FullName for invalid input, but got: %v", fullName)
			}
		})
	}
}

func TestFullName_AccessorMethods(t *testing.T) {
	fullName, _ := NewFullName("太郎", "田中")

	if fullName.FirstName() != "太郎" {
		t.Errorf("Expected first name '太郎', but got %q", fullName.FirstName())
	}
	if fullName.LastName() != "田中" {
		t.Errorf("Expected last name '田中', but got %q", fullName.LastName())
	}
	if fullName.String() != "太郎 田中" {
		t.Errorf("Expected full name '太郎 田中', but got %q", fullName.String())
	}
}

func TestFullName_Equals(t *testing.T) {
	name1, _ := NewFullName("太郎", "田中")
	name2, _ := NewFullName("太郎", "田中")
	name3, _ := NewFullName("花子", "田中")

	if !name1.Equals(name2) {
		t.Error("Expected equal names to return true")
	}
	if name1.Equals(name3) {
		t.Error("Expected different names to return false")
	}
	if name1.Equals(nil) {
		t.Error("Expected name compared to nil to return false")
	}
}

// CircleName tests
func TestNewCircleName_ValidNames_Success(t *testing.T) {
	tests := []string{
		"ABC",
		"プログラミング研究会",
		"Software Development Circle",
		"12345678901234567890123456789012345678901234567890", // 50文字
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			circleName, err := NewCircleName(name)
			if err != nil {
				t.Errorf("Expected no error for valid circle name %q, but got: %v", name, err)
			}
			if circleName == nil {
				t.Error("Expected CircleName object, but got nil")
			}
		})
	}
}

func TestNewCircleName_InvalidNames_ReturnsError(t *testing.T) {
	tests := []struct {
		name        string
		circleName  string
		description string
	}{
		{"Empty name", "", "empty circle name"},
		{"Whitespace only", "   ", "whitespace only circle name"},
		{"Too short", "AB", "too short circle name"},
		{"Too long", "123456789012345678901234567890123456789012345678901", "too long circle name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			circleName, err := NewCircleName(tt.circleName)
			if err == nil {
				t.Errorf("Expected error for %s %q, but got none", tt.description, tt.circleName)
			}
			if circleName != nil {
				t.Errorf("Expected nil CircleName for invalid input, but got: %v", circleName)
			}
		})
	}
}

// Money tests
func TestNewMoney_ValidCurrency_Success(t *testing.T) {
	tests := []struct {
		amount   int64
		currency string
	}{
		{1000, "JPY"},
		{2500, "USD"},
		{1500, "EUR"},
		{0, "GBP"},
		{-100, "AUD"},
	}

	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)
			if err != nil {
				t.Errorf("Expected no error for valid currency %q, but got: %v", tt.currency, err)
			}
			if money == nil {
				t.Error("Expected Money object, but got nil")
			}
			if money.Amount() != tt.amount {
				t.Errorf("Expected amount %d, but got %d", tt.amount, money.Amount())
			}
			if money.Currency() != tt.currency {
				t.Errorf("Expected currency %q, but got %q", tt.currency, money.Currency())
			}
		})
	}
}

func TestNewMoney_InvalidCurrency_ReturnsError(t *testing.T) {
	tests := []string{
		"",
		"INVALID",
		"123",
		"jpy", // 小文字（内部で大文字変換されるが、バリデーション用）
	}

	for _, currency := range tests {
		t.Run(currency, func(t *testing.T) {
			money, err := NewMoney(1000, currency)
			if currency == "" || currency == "INVALID" || currency == "123" {
				if err == nil {
					t.Errorf("Expected error for invalid currency %q, but got none", currency)
				}
				if money != nil {
					t.Errorf("Expected nil Money for invalid currency, but got: %v", money)
				}
			}
		})
	}
}

func TestMoney_Add_Success(t *testing.T) {
	money1, _ := NewMoney(1000, "JPY")
	money2, _ := NewMoney(500, "JPY")

	result, err := money1.Add(money2)
	if err != nil {
		t.Errorf("Expected no error when adding same currency, but got: %v", err)
	}
	if result.Amount() != 1500 {
		t.Errorf("Expected amount 1500, but got %d", result.Amount())
	}
	if result.Currency() != "JPY" {
		t.Errorf("Expected currency JPY, but got %q", result.Currency())
	}
}

func TestMoney_Add_DifferentCurrency_ReturnsError(t *testing.T) {
	money1, _ := NewMoney(1000, "JPY")
	money2, _ := NewMoney(500, "USD")

	result, err := money1.Add(money2)
	if err == nil {
		t.Error("Expected error when adding different currencies, but got none")
	}
	if result != nil {
		t.Errorf("Expected nil result for different currencies, but got: %v", result)
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tests := []struct {
		amount   int64
		expected bool
	}{
		{1000, true},
		{1, true},
		{0, false},
		{-1, false},
		{-1000, false},
	}

	for _, tt := range tests {
		money, _ := NewMoney(tt.amount, "JPY")
		if money.IsPositive() != tt.expected {
			t.Errorf("Expected IsPositive() to return %t for amount %d, but got %t",
				tt.expected, tt.amount, money.IsPositive())
		}
	}
}
