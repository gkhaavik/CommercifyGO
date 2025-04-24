package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
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
	// Prepare webhook registration request
	webhookRequest := models.WebhookRegistrationRequest{
		URL:    url,
		Events: events,
	}

	res, err := s.mobilePayService.webhookClient.Register(webhookRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to register webhook with MobilePay: %v", err)
	}

	// Create webhook record in database
	webhook := &entity.Webhook{
		Provider:   "mobilepay",
		ExternalID: res.ID,
		URL:        url,
		Events:     events,
		Secret:     res.Secret,
		IsActive:   true,
	}

	// Save webhook in database
	if err := s.webhookRepo.Create(webhook); err != nil {
		// Try to delete the webhook from MobilePay if database operation fails
		s.deleteMobilePayWebhook(res.ID)
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
	return s.mobilePayService.webhookClient.Delete(externalID)
}

// GetMobilePayWebhooks returns all registered MobilePay webhooks
func (s *WebhookService) GetMobilePayWebhooks() ([]*entity.Webhook, error) {
	// s.mobilePayService.webhookClient.GetAll()

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
