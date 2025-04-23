package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// WebhookService handles webhook management for payment providers
type WebhookService struct {
	config           *config.Config
	webhookRepo      repository.WebhookRepository
	logger           logger.Logger
	mobilePayService *MobilePayPaymentService
}

// MobilePayWebhookRequest represents the request body for creating a MobilePay webhook
type MobilePayWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// MobilePayWebhookResponse represents the response from creating a MobilePay webhook
type MobilePayWebhookResponse struct {
	ID     string   `json:"id"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Secret string   `json:"secret"`
}

// NewWebhookService creates a new webhook service
func NewWebhookService(
	config *config.Config,
	webhookRepo repository.WebhookRepository,
	logger logger.Logger,
	mobilePayService *MobilePayPaymentService,
) *WebhookService {
	return &WebhookService{
		config:           config,
		webhookRepo:      webhookRepo,
		logger:           logger,
		mobilePayService: mobilePayService,
	}
}

// RegisterMobilePayWebhook registers a webhook with MobilePay
func (s *WebhookService) RegisterMobilePayWebhook(url string, events []string) (*entity.Webhook, error) {
	// Ensure we have a valid access token
	if err := s.mobilePayService.ensureAccessToken(); err != nil {
		return nil, fmt.Errorf("failed to get MobilePay access token: %v", err)
	}

	// Prepare webhook registration request
	webhookRequest := MobilePayWebhookRequest{
		URL:    url,
		Events: events,
	}

	s.logger.Info("Registering MobilePay webhook: %s", webhookRequest)

	// Convert to JSON
	webhookJSON, err := json.Marshal(webhookRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook request: %v", err)
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.MobilePay.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request
	webhookURL := baseURL + "/webhooks/v1/webhooks"
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(webhookJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.mobilePayService.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.MobilePay.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MobilePay.MerchantSerialNumber)
	req.Header.Set("Idempotency-Key", uuid.New().String())

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to register webhook (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var webhookResponse MobilePayWebhookResponse
	if err := json.Unmarshal(body, &webhookResponse); err != nil {
		return nil, fmt.Errorf("failed to parse webhook response: %v", err)
	}

	// Create webhook record in database
	webhook := &entity.Webhook{
		Provider:   "mobilepay",
		ExternalID: webhookResponse.ID,
		URL:        webhookResponse.URL,
		Events:     webhookResponse.Events,
		Secret:     webhookResponse.Secret,
		IsActive:   true,
	}

	// Save webhook in database
	if err := s.webhookRepo.Create(webhook); err != nil {
		// Try to delete the webhook from MobilePay if database operation fails
		s.deleteMobilePayWebhook(webhookResponse.ID)
		return nil, fmt.Errorf("failed to save webhook: %v", err)
	}

	return webhook, nil
}

// DeleteMobilePayWebhook deletes a webhook from MobilePay
func (s *WebhookService) DeleteMobilePayWebhook(id uint) error {
	// Get webhook from database
	webhook, err := s.webhookRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("webhook not found: %v", err)
	}

	if webhook.Provider != "mobilepay" {
		return fmt.Errorf("webhook is not a MobilePay webhook")
	}

	// Delete webhook from MobilePay
	if err := s.deleteMobilePayWebhook(webhook.ExternalID); err != nil {
		return fmt.Errorf("failed to delete webhook from MobilePay: %v", err)
	}

	// Delete webhook from database
	if err := s.webhookRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete webhook from database: %v", err)
	}

	return nil
}

// deleteMobilePayWebhook deletes a webhook from MobilePay (internal method)
func (s *WebhookService) deleteMobilePayWebhook(externalID string) error {
	// Ensure we have a valid access token
	if err := s.mobilePayService.ensureAccessToken(); err != nil {
		return fmt.Errorf("failed to get MobilePay access token: %v", err)
	}

	// Determine the API base URL based on the test mode setting
	baseURL := mobilePayAPIProdBaseURL
	if s.config.MobilePay.IsTestMode {
		baseURL = mobilePayAPITestBaseURL
	}

	// Create HTTP request
	webhookURL := fmt.Sprintf("%s/webhooks/v1/webhooks/%s", baseURL, externalID)
	req, err := http.NewRequest("DELETE", webhookURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.mobilePayService.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.MobilePay.SubscriptionKey)
	req.Header.Set("Merchant-Serial-Number", s.config.MobilePay.MerchantSerialNumber)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete webhook (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetMobilePayWebhooks returns all registered MobilePay webhooks
func (s *WebhookService) GetMobilePayWebhooks() ([]*entity.Webhook, error) {
	// Get all MobilePay webhooks from database
	return s.webhookRepo.GetByProvider("mobilepay")
}

// VerifyMobilePayWebhookSignature verifies the signature of a MobilePay webhook
func (s *WebhookService) VerifyMobilePayWebhookSignature(requestBody []byte, signatureHeader string, webhookSecret string) bool {
	if signatureHeader == "" || webhookSecret == "" {
		return false
	}

	// Create HMAC signature using webhookSecret
	h := hmac.New(sha256.New, []byte(webhookSecret))
	h.Write(requestBody)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare with provided signature (case-insensitive)
	return expectedSignature == signatureHeader
}
