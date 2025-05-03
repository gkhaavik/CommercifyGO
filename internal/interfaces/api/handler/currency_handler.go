package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// CurrencyHandler handles HTTP requests related to currencies
type CurrencyHandler struct {
	currencyUseCase *usecase.CurrencyUseCase
}

// CurrencyResponse represents a currency in API responses
type CurrencyResponse struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	Precision     int     `json:"precision"`
	ExchangeRate  float64 `json:"exchange_rate"`
	IsDefault     bool    `json:"is_default"`
	IsEnabled     bool    `json:"is_enabled"`
	FormattedName string  `json:"formatted_name"` // Format: "USD ($)"
}

// CurrencyRequest represents a currency in API requests
type CurrencyRequest struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	Precision int     `json:"precision"`
	Rate      float64 `json:"rate"`
	IsDefault bool    `json:"is_default"`
	IsEnabled bool    `json:"is_enabled"`
}

// ExchangeRateResponse represents an exchange rate in API responses
type ExchangeRateResponse struct {
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Rate           float64 `json:"rate"`
	Date           string  `json:"date"`
}

// ConversionRequest represents a money conversion request
type ConversionRequest struct {
	Amount       float64 `json:"amount"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
}

// ConversionResponse represents a money conversion response
type ConversionResponse struct {
	OriginalAmount     float64 `json:"original_amount"`
	ConvertedAmount    float64 `json:"converted_amount"`
	FromCurrency       string  `json:"from_currency"`
	ToCurrency         string  `json:"to_currency"`
	Rate               float64 `json:"rate"`
	FormattedOriginal  string  `json:"formatted_original"`
	FormattedConverted string  `json:"formatted_converted"`
}

// NewCurrencyHandler creates a new currency handler
func NewCurrencyHandler(currencyUseCase *usecase.CurrencyUseCase) *CurrencyHandler {
	return &CurrencyHandler{
		currencyUseCase: currencyUseCase,
	}
}

// currencyToResponse converts a Currency entity to a CurrencyResponse
func (h *CurrencyHandler) currencyToResponse(currency *entity.Currency) CurrencyResponse {
	return CurrencyResponse{
		Code:          currency.Code,
		Name:          currency.Name,
		Symbol:        currency.Symbol,
		Precision:     currency.Precision,
		ExchangeRate:  currency.ExchangeRate,
		IsDefault:     currency.IsDefault,
		IsEnabled:     currency.IsEnabled,
		FormattedName: currency.Name + " (" + currency.Symbol + ")",
	}
}

// GetAllCurrencies returns all currencies
func (h *CurrencyHandler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.currencyUseCase.GetAllCurrencies()
	if err != nil {
		http.Error(w, "Failed to retrieve currencies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]CurrencyResponse, len(currencies))
	for i, currency := range currencies {
		response[i] = h.currencyToResponse(currency)
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetEnabledCurrencies returns all enabled currencies
func (h *CurrencyHandler) GetEnabledCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.currencyUseCase.GetEnabledCurrencies()
	if err != nil {
		http.Error(w, "Failed to retrieve enabled currencies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]CurrencyResponse, len(currencies))
	for i, currency := range currencies {
		response[i] = h.currencyToResponse(currency)
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetCurrencyByCode retrieves a specific currency by its ISO code
func (h *CurrencyHandler) GetCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Upper case the code for consistency
	code = strings.ToUpper(code)

	currency, err := h.currencyUseCase.GetCurrencyByCode(code)
	if err != nil {
		http.Error(w, "Failed to retrieve currency: "+err.Error(), http.StatusNotFound)
		return
	}

	response := h.currencyToResponse(currency)

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetDefaultCurrency retrieves the store's default currency
func (h *CurrencyHandler) GetDefaultCurrency(w http.ResponseWriter, r *http.Request) {
	currency, err := h.currencyUseCase.GetDefaultCurrency()
	if err != nil {
		http.Error(w, "Failed to retrieve default currency: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.currencyToResponse(currency)

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateCurrency creates a new currency
func (h *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	var req CurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Upper case the code for consistency
	req.Code = strings.ToUpper(req.Code)

	// Validate request
	if req.Code == "" || len(req.Code) != 3 {
		http.Error(w, "Currency code must be a valid 3-letter ISO code", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Currency name is required", http.StatusBadRequest)
		return
	}
	if req.Symbol == "" {
		http.Error(w, "Currency symbol is required", http.StatusBadRequest)
		return
	}
	if req.Precision < 0 {
		http.Error(w, "Precision cannot be negative", http.StatusBadRequest)
		return
	}
	if req.Rate <= 0 {
		http.Error(w, "Exchange rate must be positive", http.StatusBadRequest)
		return
	}

	// Create currency
	currency, err := h.currencyUseCase.CreateCurrency(
		req.Code, req.Name, req.Symbol, req.Precision, req.Rate, req.IsDefault, req.IsEnabled,
	)
	if err != nil {
		http.Error(w, "Failed to create currency: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.currencyToResponse(currency)

	// Write JSON response with status 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateCurrency updates an existing currency
func (h *CurrencyHandler) UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Upper case the code for consistency
	code = strings.ToUpper(code)

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extract optional fields
	var name, symbol *string
	var precision *int
	var rate *float64
	var isEnabled, isDefault *bool

	if nameVal, ok := req["name"].(string); ok && nameVal != "" {
		name = &nameVal
	}

	if symbolVal, ok := req["symbol"].(string); ok && symbolVal != "" {
		symbol = &symbolVal
	}

	if precisionVal, ok := req["precision"].(float64); ok {
		precisionInt := int(precisionVal)
		precision = &precisionInt
	}

	if rateVal, ok := req["rate"].(float64); ok && rateVal > 0 {
		rate = &rateVal
	}

	if enabledVal, ok := req["is_enabled"].(bool); ok {
		isEnabled = &enabledVal
	}

	if defaultVal, ok := req["is_default"].(bool); ok {
		isDefault = &defaultVal
	}

	// Update currency
	currency, err := h.currencyUseCase.UpdateCurrency(
		code, name, symbol, precision, rate, isEnabled, isDefault,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to update currency: "+err.Error(), statusCode)
		return
	}

	response := h.currencyToResponse(currency)

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteCurrency deletes a currency
func (h *CurrencyHandler) DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Upper case the code for consistency
	code = strings.ToUpper(code)

	if err := h.currencyUseCase.DeleteCurrency(code); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "cannot delete the default") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to delete currency: "+err.Error(), statusCode)
		return
	}

	// Return success with no content
	w.WriteHeader(http.StatusNoContent)
}

// SetDefaultCurrency sets a currency as the default
func (h *CurrencyHandler) SetDefaultCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Upper case the code for consistency
	code = strings.ToUpper(code)

	if err := h.currencyUseCase.SetDefaultCurrency(code); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to set default currency: "+err.Error(), statusCode)
		return
	}

	// Get the updated currency to return in the response
	currency, err := h.currencyUseCase.GetCurrencyByCode(code)
	if err != nil {
		http.Error(w, "Currency set as default but failed to retrieve details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.currencyToResponse(currency)

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateExchangeRates updates all exchange rates
func (h *CurrencyHandler) UpdateExchangeRates(w http.ResponseWriter, r *http.Request) {
	if err := h.currencyUseCase.UpdateExchangeRates(); err != nil {
		http.Error(w, "Failed to update exchange rates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all currencies with updated rates
	currencies, err := h.currencyUseCase.GetAllCurrencies()
	if err != nil {
		http.Error(w, "Exchange rates updated but failed to retrieve currencies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]CurrencyResponse, len(currencies))
	for i, currency := range currencies {
		response[i] = h.currencyToResponse(currency)
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetExchangeRateHistory returns historical exchange rates
func (h *CurrencyHandler) GetExchangeRateHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	base := vars["base"]
	target := vars["target"]

	if base == "" || target == "" {
		http.Error(w, "Base and target currency codes are required", http.StatusBadRequest)
		return
	}

	// Upper case the codes for consistency
	base = strings.ToUpper(base)
	target = strings.ToUpper(target)

	// Parse limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 30 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	history, err := h.currencyUseCase.GetExchangeRateHistory(base, target, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve exchange rate history: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]ExchangeRateResponse, len(history))
	for i, rate := range history {
		response[i] = ExchangeRateResponse{
			BaseCurrency:   rate.BaseCurrency,
			TargetCurrency: rate.TargetCurrency,
			Rate:           rate.Rate,
			Date:           rate.Date.Format(time.RFC3339),
		}
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// ConvertMoney converts an amount from one currency to another
func (h *CurrencyHandler) ConvertMoney(w http.ResponseWriter, r *http.Request) {
	var req ConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Upper case the currency codes for consistency
	req.FromCurrency = strings.ToUpper(req.FromCurrency)
	req.ToCurrency = strings.ToUpper(req.ToCurrency)

	// Validate request
	if req.FromCurrency == "" || req.ToCurrency == "" {
		http.Error(w, "Both source and target currencies are required", http.StatusBadRequest)
		return
	}

	// Convert the amount
	convertedAmount, rate, err := h.currencyUseCase.ConvertMoney(
		req.Amount, req.FromCurrency, req.ToCurrency,
	)
	if err != nil {
		http.Error(w, "Failed to convert currency: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get formatted money strings
	formattedOriginal, err := h.currencyUseCase.GetFormattedMoney(req.Amount, req.FromCurrency)
	if err != nil {
		http.Error(w, "Failed to format original amount: "+err.Error(), http.StatusInternalServerError)
		return
	}

	formattedConverted, err := h.currencyUseCase.GetFormattedMoney(convertedAmount, req.ToCurrency)
	if err != nil {
		http.Error(w, "Failed to format converted amount: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	response := ConversionResponse{
		OriginalAmount:     req.Amount,
		ConvertedAmount:    convertedAmount,
		FromCurrency:       req.FromCurrency,
		ToCurrency:         req.ToCurrency,
		Rate:               rate,
		FormattedOriginal:  formattedOriginal,
		FormattedConverted: formattedConverted,
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
