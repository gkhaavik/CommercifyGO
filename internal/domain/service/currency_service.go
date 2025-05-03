// Package service defines interfaces for services that implement business logic
package service

import (
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

// CurrencyService defines the interface for currency-related operations
type CurrencyService interface {
	// GetSupportedCurrencies returns a list of all supported currencies
	GetSupportedCurrencies() ([]*entity.Currency, error)

	// UpdateExchangeRates refreshes the exchange rates from the provider
	UpdateExchangeRates() error

	// Convert converts money from one currency to another
	Convert(sourceMoney *money.Money, targetCurrency string) (*money.Money, error)

	// GetExchangeRate returns the exchange rate between two currencies
	GetExchangeRate(fromCurrency, toCurrency string) (float64, error)
}
