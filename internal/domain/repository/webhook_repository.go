package repository

import (
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// WebhookRepository defines the interface for webhook operations
type WebhookRepository interface {
	// Create creates a new webhook
	Create(webhook *entity.Webhook) error

	// Update updates an existing webhook
	Update(webhook *entity.Webhook) error

	// Delete deletes a webhook
	Delete(id uint) error

	// GetByID returns a webhook by ID
	GetByID(id uint) (*entity.Webhook, error)

	// GetByProvider returns all webhooks for a specific provider
	GetByProvider(provider string) ([]*entity.Webhook, error)

	// GetActive returns all active webhooks
	GetActive() ([]*entity.Webhook, error)

	// GetByExternalID returns a webhook by external ID
	GetByExternalID(provider string, externalID string) (*entity.Webhook, error)
}
