package repository

import (
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// CurrencyRepository defines the interface for currency data operations
type CurrencyRepository interface {
	// GetAll returns all currencies in the system
	GetAll() ([]*entity.Currency, error)

	// GetByCode returns a currency by its ISO code
	GetByCode(code string) (*entity.Currency, error)

	// GetDefault returns the default currency
	GetDefault() (*entity.Currency, error)

	// GetEnabled returns all enabled currencies
	GetEnabled() ([]*entity.Currency, error)

	// Create adds a new currency
	Create(currency *entity.Currency) error

	// Update modifies an existing currency
	Update(currency *entity.Currency) error

	// Delete removes a currency
	Delete(code string) error

	// SetDefault sets a currency as the default
	SetDefault(code string) error

	// AddExchangeRate records a new exchange rate
	AddExchangeRate(history *entity.ExchangeRateHistory) error

	// GetExchangeRateHistory gets historical exchange rates for a currency pair
	GetExchangeRateHistory(baseCurrency, targetCurrency string, limit int) ([]*entity.ExchangeRateHistory, error)

	// GetLatestExchangeRate gets the most recent exchange rate for a currency pair
	GetLatestExchangeRate(baseCurrency, targetCurrency string) (*entity.ExchangeRateHistory, error)

	// UpdateExchangeRates updates exchange rates for all currencies
	UpdateExchangeRates(baseCurrency string, rates map[string]float64) error
}
