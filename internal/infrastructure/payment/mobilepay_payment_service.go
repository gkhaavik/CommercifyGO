package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

const (
	// MobilePay API endpoints
	mobilePayAPITestBaseURL     = "https://apitest.vipps.no"
	mobilePayAPIProdBaseURL     = "https://api.vipps.no"
	mobilePayAccessTokenPath    = "/accesstoken/get"
	mobilePayPaymentsPath       = "/epayment/v1/payments"
	mobilePayPaymentDetailsPath = "/epayment/v1/payments/%s" // %s: payment reference
	mobilePayCapturePaymentPath = "/epayment/v1/payments/%s/capture"
	mobilePayRefundPaymentPath  = "/epayment/v1/payments/%s/refund"
	mobilePayCancelPaymentPath  = "/epayment/v1/payments/%s/cancel"
)

// MobilePayAccessTokenResponse represents the response from MobilePay access token API
type MobilePayAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// MobilePayPaymentService implements a MobilePay payment service
type MobilePayPaymentService struct {
	config      config.MobilePayConfig
	logger      logger.Logger
	accessToken string
	expiresAt   time.Time
}

// NewMobilePayPaymentService creates a new MobilePayPaymentService
func NewMobilePayPaymentService(config config.MobilePayConfig, logger logger.Logger) *MobilePayPaymentService {
	return &MobilePayPaymentService{
		config: config,
		logger: logger,
	}
}

// GetAvailableProviders returns a list of available payment providers
func (s *MobilePayPaymentService) GetAvailableProviders() []service.PaymentProvider {
	return []service.PaymentProvider{
		{
			Type:        service.PaymentProviderMobilePay,
			Name:        "MobilePay",
			Description: "Pay with MobilePay app",
			IconURL:     "/assets/images/mobilepay-logo.png",
			Methods:     []service.PaymentMethod{service.PaymentMethodWallet},
			Enabled:     true,
		},
	}
}

// ProcessPayment processes a payment request using MobilePay
func (s *MobilePayPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	// Get access token if needed
	if err := s.ensureAccessToken(); err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to get MobilePay access token: " + err.Error(),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Only wallet payment method is supported for MobilePay
	if request.PaymentMethod != service.PaymentMethodWallet {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "unsupported payment method for MobilePay, only wallet is supported",
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Generate a unique reference for this payment
	reference := fmt.Sprintf("order-%d-%s", request.OrderID, uuid.New().String())

	// Convert amount to smallest currency unit (øre/cents)
	amountInSmallestUnit := int64(request.Amount * 100)

	// Construct the payment request
	paymentRequest := map[string]interface{}{
		"amount": map[string]interface{}{
			"currency": s.config.Market, // NOK, DKK, EUR
			"value":    amountInSmallestUnit,
		},
		"paymentMethod": map[string]interface{}{
			"type": "WALLET",
		},
		"reference": reference,
		"returnUrl": s.config.ReturnURL, // URL to redirect to after payment
		"userFlow":  "WEB_REDIRECT",     // Default flow for browser-based payments
	}

	// Add customer phone number if available (optional)
	if request.CustomerEmail != "" {
		paymentRequest["paymentDescription"] = s.config.PaymentDescription
	}

	// Convert to JSON
	paymentJSON, err := json.Marshal(paymentRequest)
	if err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to create payment request: " + err.Error(),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", baseURL+mobilePayPaymentsPath, bytes.NewBuffer(paymentJSON))
	if err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to create HTTP request: " + err.Error(),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Generate a unique idempotency key for this request
	idempotencyKey := uuid.New().String()

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MerchantSerialNumber)
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("Vipps-System-Name", s.config.PaymentDescription)
	req.Header.Set("Vipps-System-Version", "1.0.0")
	req.Header.Set("Vipps-System-Plugin-Name", "commercify-backend")
	req.Header.Set("Vipps-System-Plugin-Version", "1.0.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to execute HTTP request: " + err.Error(),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to read response: " + err.Error(),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Check for errors
	if resp.StatusCode != http.StatusCreated {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("failed to create payment: %s (status: %d)", string(body), resp.StatusCode),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Parse response
	var paymentResponse map[string]interface{}
	if err := json.Unmarshal(body, &paymentResponse); err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to parse response: " + err.Error(),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Extract redirect URL and transaction ID
	redirectURL, ok := paymentResponse["redirectUrl"].(string)
	if !ok {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "missing redirect URL in response",
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// MobilePay requires a redirect to complete the payment
	// Return a result with action URL
	return &service.PaymentResult{
		Success:        false,
		TransactionID:  reference, // Use the reference as the transaction ID
		ErrorMessage:   "payment requires user action",
		RequiresAction: true,
		ActionURL:      redirectURL,
		Provider:       service.PaymentProviderMobilePay,
	}, nil
}

// VerifyPayment verifies a payment
func (s *MobilePayPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	if provider != service.PaymentProviderMobilePay {
		return false, errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return false, errors.New("transaction ID is required")
	}

	// Ensure we have a valid access token
	if err := s.ensureAccessToken(); err != nil {
		return false, fmt.Errorf("failed to get MobilePay access token: %v", err)
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request to get payment details
	endpoint := fmt.Sprintf(baseURL+mobilePayPaymentDetailsPath, transactionID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MerchantSerialNumber)
	req.Header.Set("Vipps-System-Name", "Commercify")
	req.Header.Set("Vipps-System-Version", "1.0.0")
	req.Header.Set("Vipps-System-Plugin-Name", "commercify-backend")
	req.Header.Set("Vipps-System-Plugin-Version", "1.0.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to get payment details: %s", string(body))
	}

	// Parse response
	var paymentDetails map[string]interface{}
	if err := json.Unmarshal(body, &paymentDetails); err != nil {
		return false, fmt.Errorf("failed to parse response: %v", err)
	}

	// Check payment state
	state, ok := paymentDetails["state"].(string)
	if !ok {
		return false, errors.New("missing state in payment details")
	}

	// Return true if payment is authorized
	return state == "AUTHORIZED", nil
}

// RefundPayment refunds a payment
func (s *MobilePayPaymentService) RefundPayment(transactionID string, amount float64, provider service.PaymentProviderType) error {
	if provider != service.PaymentProviderMobilePay {
		return errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Ensure we have a valid access token
	if err := s.ensureAccessToken(); err != nil {
		return fmt.Errorf("failed to get MobilePay access token: %v", err)
	}

	// Generate a unique idempotency key for this refund
	idempotencyKey := uuid.New().String()

	// Convert amount to smallest currency unit (øre/cents)
	amountInSmallestUnit := int64(amount * 100)

	// Prepare refund request
	refundRequest := map[string]interface{}{
		"modificationAmount": map[string]interface{}{
			"currency": s.config.Market,
			"value":    amountInSmallestUnit,
		},
	}

	// Convert to JSON
	refundJSON, err := json.Marshal(refundRequest)
	if err != nil {
		return fmt.Errorf("failed to create refund request: %v", err)
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request
	endpoint := fmt.Sprintf(baseURL+mobilePayRefundPaymentPath, transactionID)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(refundJSON))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MerchantSerialNumber)
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("Vipps-System-Name", "Commercify")
	req.Header.Set("Vipps-System-Version", "1.0.0")
	req.Header.Set("Vipps-System-Plugin-Name", "commercify-backend")
	req.Header.Set("Vipps-System-Plugin-Version", "1.0.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to refund payment (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// CapturePayment captures an authorized payment
func (s *MobilePayPaymentService) CapturePayment(transactionID string, amount float64) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	if amount <= 0 {
		return errors.New("capture amount must be greater than zero")
	}

	// Ensure we have a valid access token
	if err := s.ensureAccessToken(); err != nil {
		return fmt.Errorf("failed to get MobilePay access token: %v", err)
	}

	// Generate a unique idempotency key for this capture
	idempotencyKey := uuid.New().String()

	// Convert amount to smallest currency unit (øre/cents)
	amountInSmallestUnit := int64(amount * 100)

	// Prepare capture request
	captureRequest := map[string]interface{}{
		"modificationAmount": map[string]interface{}{
			"currency": s.config.Market,
			"value":    amountInSmallestUnit,
		},
	}

	// Convert to JSON
	captureJSON, err := json.Marshal(captureRequest)
	if err != nil {
		return fmt.Errorf("failed to create capture request: %v", err)
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request
	endpoint := fmt.Sprintf(baseURL+mobilePayCapturePaymentPath, transactionID)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(captureJSON))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MerchantSerialNumber)
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("Vipps-System-Name", "Commercify")
	req.Header.Set("Vipps-System-Version", "1.0.0")
	req.Header.Set("Vipps-System-Plugin-Name", "commercify-backend")
	req.Header.Set("Vipps-System-Plugin-Version", "1.0.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to capture payment (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// CancelPayment cancels a payment
func (s *MobilePayPaymentService) CancelPayment(transactionID string) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	// Ensure we have a valid access token
	if err := s.ensureAccessToken(); err != nil {
		return fmt.Errorf("failed to get MobilePay access token: %v", err)
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request
	endpoint := fmt.Sprintf(baseURL+mobilePayCancelPaymentPath, transactionID)
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MerchantSerialNumber)
	req.Header.Set("Vipps-System-Name", "Commercify")
	req.Header.Set("Vipps-System-Version", "1.0.0")
	req.Header.Set("Vipps-System-Plugin-Name", "commercify-backend")
	req.Header.Set("Vipps-System-Plugin-Version", "1.0.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to cancel payment (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// ensureAccessToken ensures there is a valid access token
func (s *MobilePayPaymentService) ensureAccessToken() error {
	// Check if the current token is still valid
	if s.accessToken != "" && time.Now().Before(s.expiresAt) {
		return nil // Token is still valid
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request for token
	req, err := http.NewRequest("POST", baseURL+mobilePayAccessTokenPath, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers for token request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("client_id", s.config.ClientID)
	req.Header.Set("client_secret", s.config.ClientSecret)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MerchantSerialNumber)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get access token (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResponse MobilePayAccessTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %v", err)
	}

	// Save token and expiration
	s.accessToken = tokenResponse.AccessToken
	// tokenResponse.ExpiresIn is in seconds, convert to time.Duration
	expiresIn, err := time.ParseDuration(tokenResponse.ExpiresIn + "s")

	if err != nil {
		return fmt.Errorf("failed to parse expires_in: %v", err)
	}

	// Set expiration with a safety margin (5 minutes before actual expiration)
	s.expiresAt = time.Now().Add(time.Duration(expiresIn-300) * time.Second)

	return nil
}
