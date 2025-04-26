package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	orderUseCase *usecase.OrderUseCase
	logger       logger.Logger
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(orderUseCase *usecase.OrderUseCase, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (h *PaymentHandler) GetAvailablePaymentProviders(w http.ResponseWriter, r *http.Request) {
	// Get available payment providers
	providers := h.orderUseCase.GetAvailablePaymentProviders()

	// Return providers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}
