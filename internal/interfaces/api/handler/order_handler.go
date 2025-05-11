package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
	logger       logger.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderUseCase *usecase.OrderUseCase, logger logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

// CreateOrder handles order creation for both authenticated users and guests
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input usecase.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate shipping method ID
	if input.ShippingMethodID == 0 {
		http.Error(w, "Shipping method ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is authenticated
	userID, ok := r.Context().Value("user_id").(uint)

	var order *entity.Order
	var err error

	if ok && userID > 0 {
		// Authenticated user checkout
		input.UserID = userID
		order, err = h.orderUseCase.CreateOrderFromCart(input)
	} else {
		// Guest checkout
		// Get session ID from cookie
		cookie, cookieErr := r.Cookie(common.SessionCookieName)
		if cookieErr != nil || cookie.Value == "" {
			http.Error(w, "No cart session found", http.StatusBadRequest)
			return
		}

		// Validate required guest fields
		if input.Email == "" || input.FullName == "" {
			http.Error(w, "Email and full name are required for guest checkout", http.StatusBadRequest)
			return
		}

		input.SessionID = cookie.Value
		order, err = h.orderUseCase.CreateOrderFromCart(input)
	}

	if err != nil {
		h.logger.Error("Failed to create order: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrder handles getting an order by ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get order
	order, err := h.orderUseCase.GetOrderByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Check if the user is authorized to view this order
	if order.UserID != userID {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}

	// Return order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// ListOrders handles listing orders for a user
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get orders
	orders, err := h.orderUseCase.GetUserOrders(userID, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		http.Error(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// ProcessPayment handles payment processing for an order
func (h *OrderHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get the order
	order, err := h.orderUseCase.GetOrderByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var paymentInput struct {
		PaymentMethod   string                 `json:"payment_method"`
		PaymentProvider string                 `json:"payment_provider"`
		CardDetails     *service.CardDetails   `json:"card_details,omitempty"`
		PayPalDetails   *service.PayPalDetails `json:"paypal_details,omitempty"`
		BankDetails     *service.BankDetails   `json:"bank_details,omitempty"`
		PhoneNumber     string                 `json:"phone_number,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&paymentInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For registered users, verify authorization
	if order.UserID > 0 {
		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok || userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user is authorized to process payment for this order
		if order.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	} else {
		// For guest orders, check the session cookie
		if !order.IsGuestOrder {
			http.Error(w, "Invalid order type", http.StatusBadRequest)
			return
		}

		// Only allow payment processing for guest orders if they have a valid cookie
		cookie, cookieErr := r.Cookie(common.SessionCookieName)
		if cookieErr != nil || cookie.Value == "" {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// We could add additional validation here if needed
		// For example, match email in request with the email stored in the order
	}

	// Set up payment method based on input
	var paymentMethod service.PaymentMethod
	switch paymentInput.PaymentMethod {
	case "credit_card":
		paymentMethod = service.PaymentMethodCreditCard
	case "paypal":
		paymentMethod = service.PaymentMethodPayPal
	case "bank_transfer":
		paymentMethod = service.PaymentMethodBankTransfer
	case "wallet":
		paymentMethod = service.PaymentMethodWallet
	default:
		http.Error(w, "Invalid payment method", http.StatusBadRequest)
		return
	}

	// Set up payment provider based on input
	var paymentProvider service.PaymentProviderType
	switch paymentInput.PaymentProvider {
	case "stripe":
		paymentProvider = service.PaymentProviderStripe
	case "paypal":
		paymentProvider = service.PaymentProviderPayPal
	case "mock":
		paymentProvider = service.PaymentProviderMock
	case "mobilepay":
		paymentProvider = service.PaymentProviderMobilePay
	default:
		http.Error(w, "Invalid payment provider", http.StatusBadRequest)
		return
	}

	// Process payment
	input := usecase.ProcessPaymentInput{
		OrderID:         uint(id),
		PaymentMethod:   paymentMethod,
		PaymentProvider: paymentProvider,
		CardDetails:     paymentInput.CardDetails,
		PayPalDetails:   paymentInput.PayPalDetails,
		BankDetails:     paymentInput.BankDetails,
		CustomerEmail:   order.GuestEmail,
		PhoneNumber:     paymentInput.PhoneNumber,
	}

	updatedOrder, err := h.orderUseCase.ProcessPayment(input)
	if err != nil {
		h.logger.Error("Failed to process payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

// ListAllOrders handles listing all orders (admin only)
func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get orders by status if provided
	var orders []*entity.Order
	var err error

	if status != "" {
		orders, err = h.orderUseCase.ListOrdersByStatus(entity.OrderStatus(status), offset, limit)
	} else {
		// Get all orders (this would typically be implemented in OrderRepository)
		// For now, just return an empty list
		orders = []*entity.Order{}
	}

	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		http.Error(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus handles updating an order's status (admin only)
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var statusInput struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update order status
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(id),
		Status:  entity.OrderStatus(statusInput.Status),
	}

	updatedOrder, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}
