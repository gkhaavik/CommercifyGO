package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// WebhookHandler handles webhook requests from payment providers
type WebhookHandler struct {
	cfg          *config.Config
	orderUseCase *usecase.OrderUseCase
	logger       logger.Logger
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(cfg *config.Config, orderUseCase *usecase.OrderUseCase, logger logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		cfg:          cfg,
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

// HandleStripeWebhook handles webhook events from Stripe
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify the webhook signature
	webhookSecret := h.cfg.Stripe.WebhookSecret
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		h.logger.Error("Failed to verify webhook signature: %v", err)
		http.Error(w, "Failed to verify webhook signature", http.StatusBadRequest)
		return
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		h.handlePaymentSucceeded(event)
	case "payment_intent.payment_failed":
		h.handlePaymentFailed(event)
	case "charge.refunded":
		h.handleRefund(event)
	default:
		h.logger.Info("Received unhandled webhook event: %s", event.Type)
	}

	// Return a successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handlePaymentSucceeded handles the payment_intent.succeeded event
func (h *WebhookHandler) handlePaymentSucceeded(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	// Update the order status to paid
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(orderID),
		Status:  entity.OrderStatusPaid,
	}

	_, err = h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		return
	}

	h.logger.Info("Payment succeeded for order %d", orderID)
}

// handlePaymentFailed handles the payment_intent.payment_failed event
func (h *WebhookHandler) handlePaymentFailed(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	// Log the payment failure
	h.logger.Info("Payment failed for order %d: %s", orderID, paymentIntent.LastPaymentError.Msg)
}

// handleRefund handles the charge.refunded event
func (h *WebhookHandler) handleRefund(event stripe.Event) {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		h.logger.Error("Failed to parse charge: %v", err)
		return
	}

	// Get payment intent ID
	paymentIntentID := charge.PaymentIntent.ID

	// Log the refund
	h.logger.Info("Refund processed for payment %s", paymentIntentID)
}
