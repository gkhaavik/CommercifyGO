package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
)

// WebhookHandler handles webhook requests from payment providers
type WebhookHandler struct {
	cfg            *config.Config
	orderUseCase   *usecase.OrderUseCase
	webhookUseCase *usecase.WebhookUseCase
	webhookService *payment.WebhookService
	logger         logger.Logger
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(
	cfg *config.Config,
	orderUseCase *usecase.OrderUseCase,
	webhookUseCase *usecase.WebhookUseCase,
	webhookService *payment.WebhookService,
	logger logger.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		cfg:            cfg,
		orderUseCase:   orderUseCase,
		webhookUseCase: webhookUseCase,
		webhookService: webhookService,
		logger:         logger,
	}
}

// MobilePayWebhookEvent represents the structure of a MobilePay webhook event
type MobilePayWebhookEvent struct {
	MSN          string `json:"msn"`
	Reference    string `json:"reference"`
	PspReference string `json:"pspReference"`
	Name         string `json:"name"` // CREATED, ABORTED, EXPIRED, CANCELLED, CAPTURED, REFUNDED, AUTHORIZED, TERMINATED
	Amount       struct {
		Currency string `json:"currency"`
		Value    int64  `json:"value"`
	} `json:"amount"`
	Timestamp      string `json:"timestamp"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
	Success        bool   `json:"success"`
}

// RegisterWebhookRequest represents a request to register a webhook
type RegisterWebhookRequest struct {
	Provider string   `json:"provider"`
	URL      string   `json:"url"`
	Events   []string `json:"events"`
}

// RegisterMobilePayWebhook handles registering a webhook with MobilePay
func (h *WebhookHandler) RegisterMobilePayWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input RegisterWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if len(input.Events) == 0 {
		http.Error(w, "At least one event is required", http.StatusBadRequest)
		return
	}

	// Set provider to mobilepay
	input.Provider = "mobilepay"

	// Register webhook
	usecaseInput := usecase.RegisterWebhookInput{
		Provider: input.Provider,
		URL:      input.URL,
		Events:   input.Events,
	}

	h.logger.Debug("Registering MobilePay webhook: %v", usecaseInput)

	webhook, err := h.webhookUseCase.RegisterMobilePayWebhook(usecaseInput)
	if err != nil {
		h.logger.Error("Failed to register MobilePay webhook: %v", err)
		http.Error(w, fmt.Sprintf("Failed to register webhook: %v", err), http.StatusInternalServerError)
		return
	}

	// Return created webhook
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(webhook)
}

// ListMobilePayWebhooks handles listing all MobilePay webhooks
func (h *WebhookHandler) ListMobilePayWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.webhookUseCase.GetMobilePayWebhooks()
	if err != nil {
		h.logger.Error("Failed to list MobilePay webhooks: %v", err)
		http.Error(w, "Failed to list webhooks", http.StatusInternalServerError)
		return
	}

	// Return webhooks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
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

// HandleMobilePayWebhook handles webhook events from MobilePay
func (h *WebhookHandler) HandleMobilePayWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read MobilePay webhook body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log the raw webhook data for debugging
	h.logger.Debug("Received MobilePay webhook: %s", string(body))

	// Parse the webhook event
	var event MobilePayWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("Failed to parse MobilePay webhook: %v", err)
		http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
		return
	}

	// Log the event information
	h.logger.Info("Received MobilePay webhook event: %s for reference: %s", event.Name, event.Reference)

	// Verify the webhook signature if provided
	signature := r.Header.Get("Mobilepay-Signature")
	if signature != "" {
		// Get webhooks for MobilePay
		webhooks, err := h.webhookUseCase.GetMobilePayWebhooks()
		if err != nil || len(webhooks) == 0 {
			h.logger.Error("Failed to get MobilePay webhooks for signature verification: %v", err)
		} else {
			// Use the first active webhook's secret
			for _, webhook := range webhooks {
				if webhook.IsActive && webhook.Secret != "" {
					if !h.webhookService.VerifyMobilePayWebhookSignature(body, signature, webhook.Secret) {
						h.logger.Error("Invalid MobilePay webhook signature")
						http.Error(w, "Invalid webhook signature", http.StatusBadRequest)
						return
					}
					break
				}
			}
		}
	}

	// Validate the merchant serial number
	if event.MSN != h.cfg.MobilePay.MerchantSerialNumber {
		h.logger.Error("Invalid merchant serial number in webhook: %s", event.MSN)
		http.Error(w, "Invalid merchant serial number", http.StatusBadRequest)
		return
	}

	// Extract the order ID from the reference (reference format: order-{orderID}-{uuid})
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		http.Error(w, "Invalid reference format", http.StatusBadRequest)
		return
	}

	// Handle different event types
	switch event.Name {
	case "AUTHORIZED":
		h.handleMobilePayAuthorized(orderID, event)
	case "CAPTURED":
		h.handleMobilePayCaptured(orderID, event)
	case "CANCELLED":
		h.handleMobilePayCancelled(orderID, event)
	case "REFUNDED":
		h.handleMobilePayRefunded(orderID, event)
	case "ABORTED":
		h.handleMobilePayAborted(orderID, event)
	case "EXPIRED":
		h.handleMobilePayExpired(orderID, event)
	default:
		h.logger.Info("Received unhandled MobilePay webhook event: %s", event.Name)
	}

	// Return a successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleMobilePayAuthorized handles the AUTHORIZED event
func (h *WebhookHandler) handleMobilePayAuthorized(orderID uint, event MobilePayWebhookEvent) {
	// Update the order status to paid
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusPaid,
	}

	_, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status for MobilePay payment: %v", err)
		return
	}

	h.logger.Info("MobilePay payment authorized for order %d", orderID)
}

// handleMobilePayCaptured handles the CAPTURED event
func (h *WebhookHandler) handleMobilePayCaptured(orderID uint, event MobilePayWebhookEvent) {
	h.logger.Info("MobilePay payment captured for order %d", orderID)
}

// handleMobilePayCancelled handles the CANCELLED event
func (h *WebhookHandler) handleMobilePayCancelled(orderID uint, event MobilePayWebhookEvent) {
	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	_, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to cancel order for MobilePay payment: %v", err)
		return
	}

	h.logger.Info("MobilePay payment cancelled for order %d", orderID)
}

// handleMobilePayRefunded handles the REFUNDED event
func (h *WebhookHandler) handleMobilePayRefunded(orderID uint, event MobilePayWebhookEvent) {
	// Update order status to refunded
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusRefunded,
	}

	_, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to mark order as refunded for MobilePay payment: %v", err)
		return
	}

	h.logger.Info("MobilePay payment refunded for order %d", orderID)
}

// handleMobilePayAborted handles the ABORTED event
func (h *WebhookHandler) handleMobilePayAborted(orderID uint, event MobilePayWebhookEvent) {
	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	_, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to cancel order for MobilePay aborted payment: %v", err)
		return
	}

	h.logger.Info("MobilePay payment aborted for order %d", orderID)
}

// handleMobilePayExpired handles the EXPIRED event
func (h *WebhookHandler) handleMobilePayExpired(orderID uint, event MobilePayWebhookEvent) {
	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	_, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to cancel order for MobilePay expired payment: %v", err)
		return
	}

	h.logger.Info("MobilePay payment expired for order %d", orderID)
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
