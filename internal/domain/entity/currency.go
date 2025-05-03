package entity

import (
	"errors"
	"time"
)

// Currency represents a monetary currency configuration in the system
type Currency struct {
	// Code is the three-letter ISO 4217 code (e.g., USD, EUR, GBP)
	Code string `json:"code"`

	// Name is the full name of the currency
	Name string `json:"name"`

	// Symbol is the currency symbol (e.g., $, €, £)
	Symbol string `json:"symbol"`

	// Precision is the number of decimal places for the currency
	Precision int `json:"precision"`

	// ExchangeRate is the current exchange rate relative to the base currency
	ExchangeRate float64 `json:"exchange_rate"`

	// IsDefault indicates if this is the default currency for the store
	IsDefault bool `json:"is_default"`

	// IsEnabled indicates if this currency is available for customers
	IsEnabled bool `json:"is_enabled"`

	// CreatedAt is when the currency was added to the system
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the currency was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ExchangeRateHistory represents a historical exchange rate record
type ExchangeRateHistory struct {
	// ID is the unique identifier
	ID uint `json:"id"`

	// BaseCurrency is the source currency code
	BaseCurrency string `json:"base_currency"`

	// TargetCurrency is the destination currency code
	TargetCurrency string `json:"target_currency"`

	// Rate is the exchange rate from base to target
	Rate float64 `json:"rate"`

	// Date is when this exchange rate was recorded
	Date time.Time `json:"date"`
}

// NewCurrency creates a new Currency entity
func NewCurrency(code, name, symbol string, precision int, rate float64, isDefault, isEnabled bool) (*Currency, error) {
	if code == "" {
		return nil, errors.New("currency code cannot be empty")
	}

	if name == "" {
		return nil, errors.New("currency name cannot be empty")
	}

	if symbol == "" {
		return nil, errors.New("currency symbol cannot be empty")
	}

	if precision < 0 {
		return nil, errors.New("currency precision cannot be negative")
	}

	if rate <= 0 {
		return nil, errors.New("exchange rate must be positive")
	}

	now := time.Now()
	return &Currency{
		Code:         code,
		Name:         name,
		Symbol:       symbol,
		Precision:    precision,
		ExchangeRate: rate,
		IsDefault:    isDefault,
		IsEnabled:    isEnabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
