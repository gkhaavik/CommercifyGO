// Package currency provides implementations for currency conversion services
package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ExchangeRateProvider defines the interface for external exchange rate providers
type ExchangeRateProvider interface {
	// FetchLatestRates fetches the latest exchange rates for a base currency
	FetchLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error)

	// GetSupportedCurrencies returns a list of supported currency codes
	GetSupportedCurrencies() []string
}

// ExchangeRatesAPIProvider implements the ExchangeRateProvider interface using ExchangeRatesAPI
type ExchangeRatesAPIProvider struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// ExchangeRateResponse represents the response format from ExchangeRatesAPI
type ExchangeRateResponse struct {
	Base      string             `json:"base"`
	Timestamp int64              `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
	Success   bool               `json:"success"`
}

// NewExchangeRatesAPIProvider creates a new ExchangeRatesAPI provider
func NewExchangeRatesAPIProvider(apiKey string) *ExchangeRatesAPIProvider {
	return &ExchangeRatesAPIProvider{
		APIKey:  apiKey,
		BaseURL: "https://api.exchangeratesapi.io/v1",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchLatestRates fetches the latest exchange rates for the base currency
func (p *ExchangeRatesAPIProvider) FetchLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	// Build request URL
	url := fmt.Sprintf("%s/latest?access_key=%s&base=%s", p.BaseURL, p.APIKey, baseCurrency)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var rateResponse ExchangeRateResponse
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API error
	if !rateResponse.Success {
		return nil, fmt.Errorf("API returned unsuccessful response")
	}

	return rateResponse.Rates, nil
}

// MockExchangeRateProvider provides mock exchange rates for development and testing
type MockExchangeRateProvider struct {
	// Static rates against USD
	rates map[string]float64
}

// NewMockExchangeRateProvider creates a new mock exchange rate provider
func NewMockExchangeRateProvider() *MockExchangeRateProvider {
	return &MockExchangeRateProvider{
		rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.85,
			"GBP": 0.75,
			"JPY": 110.0,
			"CAD": 1.25,
			"AUD": 1.35,
			"CHF": 0.92,
			"CNY": 6.45,
			"INR": 74.5,
			"BRL": 5.2,
			"SEK": 8.6,
			"NOK": 8.8,
			"DKK": 6.3,
		},
	}
}

// FetchLatestRates implements ExchangeRateProvider.FetchLatestRates
func (p *MockExchangeRateProvider) FetchLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	// If the base currency is not USD, we need to convert the rates
	if baseCurrency == "USD" {
		return p.rates, nil
	}

	// Get the rate for the base currency against USD
	baseRate, exists := p.rates[baseCurrency]
	if !exists {
		return nil, fmt.Errorf("rates not available for base currency: %s", baseCurrency)
	}

	// Calculate rates with the new base currency
	result := make(map[string]float64)
	for curr, rate := range p.rates {
		result[curr] = rate / baseRate
	}

	return result, nil
}

// GetSupportedCurrencies implements ExchangeRateProvider.GetSupportedCurrencies
func (p *MockExchangeRateProvider) GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(p.rates))
	for curr := range p.rates {
		currencies = append(currencies, curr)
	}
	return currencies
}
