// Package money provides utilities for handling monetary values.
// Money is stored as integers (cents) in the database to avoid floating-point precision issues.
package money

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Currency represents a monetary currency using ISO 4217 standards
type Currency struct {
	Code      string  // Three-letter ISO 4217 code (e.g., USD, EUR, GBP)
	Name      string  // Full name of the currency
	Symbol    string  // Currency symbol (e.g., $, €, £)
	Precision int     // Number of decimal places (usually 2)
	Rate      float64 // Exchange rate relative to base currency
	IsDefault bool    // Whether this is the default currency
}

// Money represents a monetary amount in a specific currency
type Money struct {
	Amount   int64  // Amount in smallest currency unit (e.g., cents)
	Currency string // Three-letter ISO 4217 code
}

// Common errors
var (
	ErrInvalidCurrencyCode    = errors.New("invalid currency code")
	ErrCurrencyNotFound       = errors.New("currency not found")
	ErrInvalidConversionRate  = errors.New("invalid conversion rate")
	ErrCurrenciesMustMatch    = errors.New("currencies must match for this operation")
	ErrAmountCannotBeNegative = errors.New("amount cannot be negative")
)

// Common currencies
var (
	USD = Currency{Code: "USD", Name: "US Dollar", Symbol: "$", Precision: 2, Rate: 1.0, IsDefault: true}
	EUR = Currency{Code: "EUR", Name: "Euro", Symbol: "€", Precision: 2, Rate: 1.0}
	GBP = Currency{Code: "GBP", Name: "British Pound", Symbol: "£", Precision: 2, Rate: 1.0}
	JPY = Currency{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", Precision: 0, Rate: 1.0}
	CAD = Currency{Code: "CAD", Name: "Canadian Dollar", Symbol: "$", Precision: 2, Rate: 1.0}
)

// supportedCurrencies is a map of supported currencies
var supportedCurrencies = map[string]Currency{
	"USD": USD,
	"EUR": EUR,
	"GBP": GBP,
	"JPY": JPY,
	"CAD": CAD,
}

// NewMoney creates a new Money instance with the specified amount and currency
func NewMoney(amount int64, currencyCode string) (*Money, error) {
	if amount < 0 {
		return nil, ErrAmountCannotBeNegative
	}

	code := strings.ToUpper(currencyCode)
	if _, ok := supportedCurrencies[code]; !ok {
		return nil, ErrInvalidCurrencyCode
	}

	return &Money{
		Amount:   amount,
		Currency: code,
	}, nil
}

// NewMoneyFromFloat creates a new Money instance from a floating point value
func NewMoneyFromFloat(amount float64, currencyCode string) (*Money, error) {
	if amount < 0 {
		return nil, ErrAmountCannotBeNegative
	}

	code := strings.ToUpper(currencyCode)
	if currency, ok := supportedCurrencies[code]; !ok {
		return nil, ErrInvalidCurrencyCode
	} else {
		// Convert to smallest unit based on currency precision
		scale := math.Pow(10, float64(currency.Precision))
		units := int64(math.Round(amount * scale))

		return &Money{
			Amount:   units,
			Currency: code,
		}, nil
	}
}

// ToCents converts a dollar amount (float64) to cents (int64)
// This avoids floating-point precision issues when storing money values
func ToCents(dollars float64) int64 {
	// Round to nearest cent to avoid floating point issues
	// Multiply by 100 to convert dollars to cents
	return int64(math.Round(dollars * 100))
}

// FromCents converts a cent amount (int64) to dollars (float64)
func FromCents(cents int64) float64 {
	// Divide by 100 to convert cents to dollars
	return float64(cents) / 100
}

// ConvertNullableToCents converts a nullable dollar amount (*float64) to nullable cents (*int64)
func ConvertNullableToCents(dollars *float64) *int64 {
	if dollars == nil {
		return nil
	}
	cents := ToCents(*dollars)
	return &cents
}

// ConvertNullableFromCents converts a nullable cent amount (*int64) to nullable dollars (*float64)
func ConvertNullableFromCents(cents *int64) *float64 {
	if cents == nil {
		return nil
	}
	dollars := FromCents(*cents)
	return &dollars
}

// FormatCurrency formats a cents value (int64) as a currency string
func FormatCurrency(cents int64, symbol string) string {
	// Convert to dollars
	dollars := FromCents(cents)
	return symbol + formatDollars(dollars)
}

// FormatDollars formats a dollars value as a currency string
func FormatDollars(dollars float64, symbol string) string {
	return symbol + formatDollars(dollars)
}

// formatDollars is a helper function to format dollar values with 2 decimal places
func formatDollars(dollars float64) string {
	return fmt.Sprintf("%.2f", dollars)
}

// ApplyPercentage applies a percentage to a cents value
func ApplyPercentage(cents int64, percentage float64) int64 {
	return ToCents(FromCents(cents) * percentage / 100)
}

// Format returns the formatted string representation of the Money
func (m *Money) Format() (string, error) {
	currency, ok := supportedCurrencies[m.Currency]
	if !ok {
		return "", ErrCurrencyNotFound
	}

	value := float64(m.Amount) / math.Pow(10, float64(currency.Precision))
	format := "%s%." + fmt.Sprintf("%d", currency.Precision) + "f"

	return fmt.Sprintf(format, currency.Symbol, value), nil
}

// Add adds two Money values and returns a new Money result
func (m *Money) Add(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, ErrCurrenciesMustMatch
	}

	return &Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract subtracts the other Money value from this one and returns a new Money result
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, ErrCurrenciesMustMatch
	}

	result := m.Amount - other.Amount
	if result < 0 {
		return nil, ErrAmountCannotBeNegative
	}

	return &Money{
		Amount:   result,
		Currency: m.Currency,
	}, nil
}

// Multiply multiplies the Money value by a factor and returns a new Money result
func (m *Money) Multiply(factor float64) (*Money, error) {
	if factor < 0 {
		return nil, ErrAmountCannotBeNegative
	}

	result := int64(math.Round(float64(m.Amount) * factor))

	return &Money{
		Amount:   result,
		Currency: m.Currency,
	}, nil
}

// Equal checks if two Money values are equal
func (m *Money) Equal(other *Money) bool {
	return m.Currency == other.Currency && m.Amount == other.Amount
}

// Float returns the Money value as a float according to the currency's precision
func (m *Money) Float() (float64, error) {
	currency, ok := supportedCurrencies[m.Currency]
	if !ok {
		return 0, ErrCurrencyNotFound
	}

	return float64(m.Amount) / math.Pow(10, float64(currency.Precision)), nil
}

// GetSupportedCurrencies returns a slice of all supported currencies
func GetSupportedCurrencies() []Currency {
	currencies := make([]Currency, 0, len(supportedCurrencies))
	for _, currency := range supportedCurrencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

// IsSupported checks if a currency code is supported
func IsSupported(code string) bool {
	_, ok := supportedCurrencies[strings.ToUpper(code)]
	return ok
}

// GetCurrency returns the Currency for a given currency code
func GetCurrency(code string) (Currency, error) {
	currency, ok := supportedCurrencies[strings.ToUpper(code)]
	if !ok {
		return Currency{}, ErrCurrencyNotFound
	}
	return currency, nil
}

// GetDefaultCurrency returns the default currency
func GetDefaultCurrency() Currency {
	for _, currency := range supportedCurrencies {
		if currency.IsDefault {
			return currency
		}
	}
	// Fallback to USD if no default is set
	return USD
}
