package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
)

// CurrencyUseCase implements currency-related use cases
type CurrencyUseCase struct {
	currencyRepo    repository.CurrencyRepository
	currencyService service.CurrencyService
}

// NewCurrencyUseCase creates a new CurrencyUseCase
func NewCurrencyUseCase(
	currencyRepo repository.CurrencyRepository,
	currencyService service.CurrencyService,
) *CurrencyUseCase {
	return &CurrencyUseCase{
		currencyRepo:    currencyRepo,
		currencyService: currencyService,
	}
}

// GetAllCurrencies returns all currencies in the system
func (uc *CurrencyUseCase) GetAllCurrencies() ([]*entity.Currency, error) {
	return uc.currencyRepo.GetAll()
}

// GetEnabledCurrencies returns all enabled currencies
func (uc *CurrencyUseCase) GetEnabledCurrencies() ([]*entity.Currency, error) {
	return uc.currencyRepo.GetEnabled()
}

// GetCurrencyByCode retrieves a specific currency by its ISO code
func (uc *CurrencyUseCase) GetCurrencyByCode(code string) (*entity.Currency, error) {
	if code == "" {
		return nil, errors.New("currency code cannot be empty")
	}
	return uc.currencyRepo.GetByCode(code)
}

// GetDefaultCurrency retrieves the store's default currency
func (uc *CurrencyUseCase) GetDefaultCurrency() (*entity.Currency, error) {
	return uc.currencyRepo.GetDefault()
}

// CreateCurrency creates a new currency
func (uc *CurrencyUseCase) CreateCurrency(
	code, name, symbol string,
	precision int,
	rate float64,
	isDefault, isEnabled bool,
) (*entity.Currency, error) {
	// Validate input parameters
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

	// Check if currency already exists
	existing, err := uc.currencyRepo.GetByCode(code)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("currency with code %s already exists", code)
	}

	// Create new currency entity
	currency, err := entity.NewCurrency(code, name, symbol, precision, rate, isDefault, isEnabled)
	if err != nil {
		return nil, err
	}

	// If this is the default currency, ensure exchange rate is 1.0
	if isDefault {
		currency.ExchangeRate = 1.0
	}

	// Save the new currency
	if err := uc.currencyRepo.Create(currency); err != nil {
		return nil, err
	}

	// If this is set as default, update any existing default currency
	if isDefault {
		if err := uc.currencyRepo.SetDefault(code); err != nil {
			return nil, err
		}
	}

	return currency, nil
}

// UpdateCurrency updates an existing currency
func (uc *CurrencyUseCase) UpdateCurrency(
	code string,
	name, symbol *string,
	precision *int,
	rate *float64,
	isEnabled *bool,
	isDefault *bool,
) (*entity.Currency, error) {
	// Get the existing currency
	currency, err := uc.currencyRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	if currency == nil {
		return nil, fmt.Errorf("currency with code %s not found", code)
	}

	// Update fields if provided
	if name != nil {
		currency.Name = *name
	}
	if symbol != nil {
		currency.Symbol = *symbol
	}
	if precision != nil {
		currency.Precision = *precision
	}

	// Handle the default currency setting first, as it affects the exchange rate
	if isDefault != nil && *isDefault {
		// If setting as default, ensure exchange rate is 1.0 regardless of provided rate
		currency.ExchangeRate = 1.0
		currency.IsDefault = true
	} else if rate != nil {
		// Only update rate if not being set as default
		currency.ExchangeRate = *rate
	}

	if isEnabled != nil {
		currency.IsEnabled = *isEnabled
	}

	// Update the timestamp
	currency.UpdatedAt = time.Now()

	// Save the changes
	if err := uc.currencyRepo.Update(currency); err != nil {
		return nil, err
	}

	// Handle default currency change if requested
	if isDefault != nil && *isDefault {
		if err := uc.currencyRepo.SetDefault(code); err != nil {
			return nil, err
		}

		// Refresh the currency to get the updated data
		return uc.currencyRepo.GetByCode(code)
	}

	return currency, nil
}

// DeleteCurrency removes a currency if it's not the default
func (uc *CurrencyUseCase) DeleteCurrency(code string) error {
	// Get the currency to check if it's the default
	currency, err := uc.currencyRepo.GetByCode(code)
	if err != nil {
		return err
	}
	if currency == nil {
		return fmt.Errorf("currency with code %s not found", code)
	}

	// Prevent deletion of default currency
	if currency.IsDefault {
		return errors.New("cannot delete the default currency")
	}

	return uc.currencyRepo.Delete(code)
}

// SetDefaultCurrency sets a currency as the default
func (uc *CurrencyUseCase) SetDefaultCurrency(code string) error {
	// Check if the currency exists
	currency, err := uc.currencyRepo.GetByCode(code)
	if err != nil {
		return err
	}
	if currency == nil {
		return fmt.Errorf("currency with code %s not found", code)
	}

	// Set as default
	return uc.currencyRepo.SetDefault(code)
}

// UpdateExchangeRates updates all currency exchange rates using the external service
func (uc *CurrencyUseCase) UpdateExchangeRates() error {
	return uc.currencyService.UpdateExchangeRates()
}

// GetExchangeRateHistory returns historical exchange rates between two currencies
func (uc *CurrencyUseCase) GetExchangeRateHistory(
	baseCurrency, targetCurrency string,
	limit int,
) ([]*entity.ExchangeRateHistory, error) {
	// Validate currencies
	if baseCurrency == "" || targetCurrency == "" {
		return nil, errors.New("base and target currencies cannot be empty")
	}

	// Ensure currencies exist
	_, err := uc.currencyRepo.GetByCode(baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("base currency not found: %w", err)
	}

	_, err = uc.currencyRepo.GetByCode(targetCurrency)
	if err != nil {
		return nil, fmt.Errorf("target currency not found: %w", err)
	}

	// Default limit if not specified
	if limit <= 0 {
		limit = 30 // Default to 30 days
	}

	return uc.currencyRepo.GetExchangeRateHistory(baseCurrency, targetCurrency, limit)
}

// GetExchangeRate returns the current exchange rate between two currencies
func (uc *CurrencyUseCase) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	// Validate currencies
	if fromCurrency == "" || toCurrency == "" {
		return 0, errors.New("source and target currencies cannot be empty")
	}

	// Ensure currencies exist
	_, err := uc.currencyRepo.GetByCode(fromCurrency)
	if err != nil {
		return 0, fmt.Errorf("source currency not found: %w", err)
	}

	_, err = uc.currencyRepo.GetByCode(toCurrency)
	if err != nil {
		return 0, fmt.Errorf("target currency not found: %w", err)
	}

	// Get the exchange rate from the currency service
	return uc.currencyService.GetExchangeRate(fromCurrency, toCurrency)
}

// ConvertMoney converts an amount from one currency to another
func (uc *CurrencyUseCase) ConvertMoney(
	amount float64,
	fromCurrency, toCurrency string,
) (float64, float64, error) {
	// Create a money object for the source amount
	sourceMoney, err := money.NewMoneyFromFloat(amount, fromCurrency)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create money object: %w", err)
	}

	// Convert the amount
	convertedMoney, err := uc.currencyService.Convert(sourceMoney, toCurrency)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert currency: %w", err)
	}

	// Get the exchange rate
	exchangeRate, err := uc.currencyService.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Convert to float values
	convertedAmount, err := convertedMoney.Float()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert to float: %w", err)
	}

	return convertedAmount, exchangeRate, nil
}

// GetFormattedMoney returns a formatted money string with the appropriate currency symbol
func (uc *CurrencyUseCase) GetFormattedMoney(amount float64, currencyCode string) (string, error) {
	// Get the currency to access its symbol and precision
	_, err := uc.currencyRepo.GetByCode(currencyCode)
	if err != nil {
		return "", fmt.Errorf("failed to get currency: %w", err)
	}

	// Create a money instance
	moneyObj, err := money.NewMoneyFromFloat(amount, currencyCode)
	if err != nil {
		return "", fmt.Errorf("failed to create money object: %w", err)
	}

	// Format the money amount
	formatted, err := moneyObj.Format()
	if err != nil {
		return "", fmt.Errorf("failed to format money: %w", err)
	}

	return formatted, nil
}
