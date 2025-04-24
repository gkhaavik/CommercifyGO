package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// WebhookHandler handles webhook requests from payment providers
type WebhookHandler struct {
	cfg            *config.Config
	orderUseCase   *usecase.OrderUseCase
	webhookUseCase *usecase.WebhookUseCase
	logger         logger.Logger
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(
	cfg *config.Config,
	orderUseCase *usecase.OrderUseCase,
	webhookUseCase *usecase.WebhookUseCase,
	logger logger.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		cfg:            cfg,
		orderUseCase:   orderUseCase,
		webhookUseCase: webhookUseCase,
		logger:         logger,
	}
}

// RegisterWebhookRequest represents a request to register a webhook
type RegisterWebhookRequest struct {
	Provider string   `json:"provider"`
	URL      string   `json:"url"`
	Events   []string `json:"events"`
}

// ListWebhooks handles listing all webhooks
func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.webhookUseCase.GetAllWebhooks()
	if err != nil {
		h.logger.Error("Failed to list webhooks: %v", err)
		http.Error(w, "Failed to list webhooks", http.StatusInternalServerError)
		return
	}

	// Return webhooks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

// GetWebhook handles getting a webhook by ID
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	// Get webhook ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["webhookId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	// Get webhook
	webhook, err := h.webhookUseCase.GetWebhookByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get webhook: %v", err)
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	// Return webhook
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook)
}

// DeleteWebhook handles deleting a webhook
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	// Get webhook ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["webhookId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	// Delete webhook
	if err := h.webhookUseCase.DeleteWebhook(uint(id)); err != nil {
		h.logger.Error("Failed to delete webhook: %v", err)
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// HandleMobilePayAuthorized handles the AUTHORIZED event
func (h *WebhookHandler) HandleMobilePayAuthorized(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update the order status to paid
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusPaid,
	}

	_, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to update order status for MobilePay payment: %v", err2)
		return err2
	}

	h.logger.Info("MobilePay payment authorized for order %d", orderID)
	return nil
}

// HandleMobilePayCaptured handles the CAPTURED event
func (h *WebhookHandler) HandleMobilePayCaptured(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment captured for order %d", orderID)
	return nil
}

// HandleMobilePayCancelled handles the CANCELLED event
func (h *WebhookHandler) HandleMobilePayCancelled(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	_, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to cancel order for MobilePay payment: %v", err2)
		return err2
	}

	h.logger.Info("MobilePay payment cancelled for order %d", orderID)
	return nil
}

// HandleMobilePayRefunded handles the REFUNDED event
func (h *WebhookHandler) HandleMobilePayRefunded(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to refunded
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusRefunded,
	}

	_, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to mark order as refunded for MobilePay payment: %v", err2)
		return err2
	}

	h.logger.Info("MobilePay payment refunded for order %d", orderID)
	return nil
}

// HandleMobilePayAborted handles the ABORTED event
func (h *WebhookHandler) HandleMobilePayAborted(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	_, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to cancel order for MobilePay aborted payment: %v", err2)
		return err2
	}

	h.logger.Info("MobilePay payment aborted for order %d", orderID)
	return nil
}

// HandleMobilePayExpired handles the EXPIRED event
func (h *WebhookHandler) HandleMobilePayExpired(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	_, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to cancel order for MobilePay expired payment: %v", err2)
		return err2
	}

	h.logger.Info("MobilePay payment expired for order %d", orderID)
	return nil
}

// extractOrderIDFromReference extracts the order ID from the reference
// Reference format: "order-{orderID}-{uuid}"
func extractOrderIDFromReference(reference string) (uint, error) {
	var orderID uint
	_, err := fmt.Sscanf(reference, "order-%d-", &orderID)
	if err != nil {
		return 0, fmt.Errorf("invalid reference format: %v", err)
	}
	return orderID, nil
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
