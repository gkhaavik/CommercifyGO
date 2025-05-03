package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
)

// DefaultCurrencyService provides currency conversion and exchange rate functionality
type DefaultCurrencyService struct {
	currencyRepo  repository.CurrencyRepository
	apiKey        string
	apiURL        string
	httpClient    *http.Client
	cachedRates   map[string]float64
	lastUpdated   time.Time
	cacheDuration time.Duration
	mu            sync.RWMutex
}

// NewDefaultCurrencyService creates a new currency service
func NewDefaultCurrencyService(
	currencyRepo repository.CurrencyRepository,
	apiKey string,
	apiURL string,
) service.CurrencyService {
	return &DefaultCurrencyService{
		currencyRepo: currencyRepo,
		apiKey:       apiKey,
		apiURL:       apiURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cachedRates:   make(map[string]float64),
		lastUpdated:   time.Time{},   // Zero time to force initial update
		cacheDuration: 1 * time.Hour, // Cache rates for 1 hour
		mu:            sync.RWMutex{},
	}
}

// GetSupportedCurrencies returns a list of all supported currencies
func (s *DefaultCurrencyService) GetSupportedCurrencies() ([]*entity.Currency, error) {
	return s.currencyRepo.GetEnabled()
}

// UpdateExchangeRates refreshes the exchange rates from the provider
func (s *DefaultCurrencyService) UpdateExchangeRates() error {
	// Get the default currency to use as base
	baseCurrency, err := s.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	// Fetch the latest exchange rates from the API
	url := fmt.Sprintf("%s/latest?base=%s&access_key=%s", s.apiURL, baseCurrency.Code, s.apiKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exchange rate API returned status code %d", resp.StatusCode)
	}

	// Parse the response
	var response struct {
		Success   bool               `json:"success"`
		Timestamp int64              `json:"timestamp"`
		Base      string             `json:"base"`
		Date      string             `json:"date"`
		Rates     map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode exchange rates response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("exchange rate API returned unsuccessful response")
	}

	// Update our cached rates
	s.mu.Lock()
	s.cachedRates = response.Rates
	s.cachedRates[baseCurrency.Code] = 1.0 // Ensure base currency has rate of 1.0
	s.lastUpdated = time.Now()
	s.mu.Unlock()

	// Update the exchange rates in the database
	if err := s.currencyRepo.UpdateExchangeRates(baseCurrency.Code, response.Rates); err != nil {
		return fmt.Errorf("failed to update exchange rates in database: %w", err)
	}

	// Record each exchange rate history
	for targetCurrency, rate := range response.Rates {
		history := &entity.ExchangeRateHistory{
			BaseCurrency:   baseCurrency.Code,
			TargetCurrency: targetCurrency,
			Rate:           rate,
			Date:           time.Now(),
		}

		if err := s.currencyRepo.AddExchangeRate(history); err != nil {
			// Log the error but continue with other currencies
			fmt.Printf("Error recording exchange rate history for %s: %v\n", targetCurrency, err)
		}
	}

	return nil
}

// refreshRatesIfNeeded updates rates if cache is expired
func (s *DefaultCurrencyService) refreshRatesIfNeeded() error {
	s.mu.RLock()
	shouldUpdate := len(s.cachedRates) == 0 || time.Since(s.lastUpdated) > s.cacheDuration
	s.mu.RUnlock()

	if shouldUpdate {
		return s.UpdateExchangeRates()
	}
	return nil
}

// Convert converts money from one currency to another
func (s *DefaultCurrencyService) Convert(sourceMoney *money.Money, targetCurrency string) (*money.Money, error) {
	// Ensure we have up-to-date exchange rates
	if err := s.refreshRatesIfNeeded(); err != nil {
		return nil, fmt.Errorf("failed to refresh exchange rates: %w", err)
	}

	// Get the source and target currencies
	sourceCode := sourceMoney.Currency

	// If source and target are the same, no conversion needed
	if sourceCode == targetCurrency {
		return sourceMoney, nil
	}

	// Get the exchange rate between the two currencies
	rate, err := s.GetExchangeRate(sourceCode, targetCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Get the target currency to determine precision
	_, err = s.currencyRepo.GetByCode(targetCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get target currency: %w", err)
	}

	// Get the source amount as a float
	sourceAmount, err := sourceMoney.Float()
	if err != nil {
		return nil, fmt.Errorf("failed to convert source amount to float: %w", err)
	}

	// Calculate the converted amount
	convertedAmount := sourceAmount * rate

	// Create a new Money object with the converted amount
	return money.NewMoneyFromFloat(convertedAmount, targetCurrency)
}

// GetExchangeRate returns the exchange rate between two currencies
func (s *DefaultCurrencyService) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	// Ensure we have up-to-date exchange rates
	if err := s.refreshRatesIfNeeded(); err != nil {
		return 0, fmt.Errorf("failed to refresh exchange rates: %w", err)
	}

	// If the currencies are the same, return 1.0
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	// First try to get from cached rates
	s.mu.RLock()
	fromRate, fromOk := s.cachedRates[fromCurrency]
	toRate, toOk := s.cachedRates[toCurrency]
	s.mu.RUnlock()

	if fromOk && toOk {
		// Calculate the cross rate
		return toRate / fromRate, nil
	}

	// If not in cache, try to get from database
	history, err := s.currencyRepo.GetLatestExchangeRate(fromCurrency, toCurrency)
	if err == nil && history != nil {
		return history.Rate, nil
	}

	// If we don't have a direct exchange rate, we need to calculate it through the default currency
	sourceCurrency, err := s.currencyRepo.GetByCode(fromCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to get source currency: %w", err)
	}

	targetCurrency, err := s.currencyRepo.GetByCode(toCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to get target currency: %w", err)
	}

	// If we have exchange rates for both currencies, we can calculate the cross rate
	// Rate = (Target Rate / Source Rate)
	return targetCurrency.ExchangeRate / sourceCurrency.ExchangeRate, nil
}
