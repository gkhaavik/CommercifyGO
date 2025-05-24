package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// CheckoutHandler handles checkout-related HTTP requests
type CheckoutHandler struct {
	checkoutUseCase *usecase.CheckoutUseCase
	orderUseCase    *usecase.OrderUseCase
	logger          logger.Logger
}

// NewCheckoutHandler creates a new CheckoutHandler
func NewCheckoutHandler(checkoutUseCase *usecase.CheckoutUseCase, logger logger.Logger) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutUseCase: checkoutUseCase,
		logger:          logger,
	}
}

// getCheckoutSessionID gets or creates a checkout session ID
func (h *CheckoutHandler) getCheckoutSessionID(w http.ResponseWriter, r *http.Request) string {
	// Check if checkout session cookie exists
	cookie, err := r.Cookie(common.CheckoutSessionCookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Create new checkout session ID if none exists
	sessionID := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     common.CheckoutSessionCookie,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   common.CheckoutSessionMaxAge,
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set secure flag if connection is HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return sessionID
}

// GetCheckout handles getting a user's checkout
func (h *CheckoutHandler) GetCheckout(w http.ResponseWriter, r *http.Request) {
	// Always get checkout session ID, needed for all checkouts
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetOrCreateCheckoutBySessionID(checkoutSessionID)

	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// AddToCheckout handles adding an item to the checkout
func (h *CheckoutHandler) AddToCheckout(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.AddToCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Always get checkout session ID, needed for all checkouts
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	// Try to find checkout by checkout session ID first
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DTO to usecase input
	checkoutInput := usecase.CheckoutInput{
		ProductID: request.ProductID,
		VariantID: request.VariantID,
		Quantity:  request.Quantity,
	}

	// Add item to checkout
	checkout, err = h.checkoutUseCase.AddItemToCheckout(checkout.ID, checkoutInput)

	if err != nil {
		h.logger.Error("Failed to add to checkout: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// UpdateCheckoutItem handles updating an item in the checkout
func (h *CheckoutHandler) UpdateCheckoutItem(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.UpdateCheckoutItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DTO to usecase input
	updateInput := usecase.UpdateCheckoutItemInput{
		ProductID: uint(productID),
		VariantID: request.VariantID,
		Quantity:  request.Quantity,
	}

	checkout.UpdateItem(updateInput.ProductID, updateInput.VariantID, updateInput.Quantity)
	h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to update checkout item: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// RemoveFromCheckout handles removing an item from the checkout
func (h *CheckoutHandler) RemoveFromCheckout(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove the item from checkout
	err = checkout.RemoveItem(uint(productID), 0)
	if err != nil {
		h.logger.Error("Failed to remove item from checkout: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the checkout
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to remove item from checkout: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// ClearCheckout handles emptying the checkout
func (h *CheckoutHandler) ClearCheckout(w http.ResponseWriter, r *http.Request) {
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkout.Clear()
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to clear checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return empty checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// SetShippingAddress handles setting the shipping address for a checkout
func (h *CheckoutHandler) SetShippingAddress(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetShippingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	address := entity.Address{
		Street:     request.AddressLine1,
		City:       request.City,
		State:      request.State,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	checkout.SetShippingAddress(address)
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to set shipping address: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// SetBillingAddress handles setting the billing address for a checkout
func (h *CheckoutHandler) SetBillingAddress(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetBillingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DTO to address entity
	address := entity.Address{
		Street:     request.AddressLine1,
		City:       request.City,
		State:      request.State,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	checkout.SetBillingAddress(address)
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to set billing address: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// SetCustomerDetails handles setting the customer details for a checkout
func (h *CheckoutHandler) SetCustomerDetails(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetCustomerDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DTO to customer details entity
	customerDetails := entity.CustomerDetails{
		Email:    request.Email,
		Phone:    request.Phone,
		FullName: request.FullName,
	}

	checkout.SetCustomerDetails(customerDetails)
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to set customer details: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// SetShippingMethod handles setting the shipping method for a checkout
func (h *CheckoutHandler) SetShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetShippingMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkout, err = h.checkoutUseCase.SetShippingMethod(checkout, request.ShippingMethodID)

	if err != nil {
		h.logger.Error("Failed to set shipping method: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// ApplyDiscount handles applying a discount code to a checkout
func (h *CheckoutHandler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.ApplyDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkout, err = h.checkoutUseCase.ApplyDiscountCode(checkout, request.DiscountCode)

	if err != nil {
		h.logger.Error("Failed to apply discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// RemoveDiscount handles removing a discount from a checkout
func (h *CheckoutHandler) RemoveDiscount(w http.ResponseWriter, r *http.Request) {
	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.checkoutUseCase.RemoveDiscountCode(checkout)

	if err != nil {
		h.logger.Error("Failed to remove discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	checkoutDTO := convertToCheckoutDTO(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkoutDTO)
}

// CompleteOrder handles converting a checkout to an order
func (h *CheckoutHandler) CompleteCheckout(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var paymentInput dto.CompleteCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get checkout session ID
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	h.logger.Info("Converting checkout to order. Authenticated: %v, UserID: %d, CheckoutSessionID: %s",
		false, 0, checkoutSessionID)

	var order *entity.Order
	var err error

	// Try to find checkout by checkout session ID first
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If checkout exists for this session, convert it to order
	order, err = h.checkoutUseCase.CreateOrderFromCheckout(checkout.ID)
	if err != nil {
		h.logger.Error("Failed to convert checkout to order: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate payment data
	if paymentInput.PaymentData.CardDetails == nil && paymentInput.PaymentData.PhoneNumber == "" {
		http.Error(w, "Payment data is required", http.StatusBadRequest)
		return
	}
	// Process payment
	order, err = h.checkoutUseCase.ProcessPayment(order, paymentInput.PaymentData)
	if err != nil {
		h.logger.Error("Failed to process payment: %v", err)
		http.Error(w, "Failed to process payment", http.StatusBadRequest)
		return
	}

	// Return created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Create response
	response := dto.ResponseDTO[dto.OrderDTO]{
		Success: true,
		Message: "Order created successfully",
		Data:    convertToOrderDTO(order),
	}

	json.NewEncoder(w).Encode(response)
}

// ListAdminCheckouts handles listing all checkouts (admin only)
func (h *CheckoutHandler) ListAdminCheckouts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get checkouts by status if provided
	var checkouts []*entity.Checkout
	var err error

	if status != "" {
		checkouts, err = h.checkoutUseCase.GetCheckoutsByStatus(entity.CheckoutStatus(status), offset, limit)
	} else {
		checkouts, err = h.checkoutUseCase.GetAllCheckouts(offset, limit)
	}

	if err != nil {
		h.logger.Error("Failed to list checkouts: %v", err)
		http.Error(w, "Failed to list checkouts", http.StatusInternalServerError)
		return
	}

	// Convert checkouts to DTOs
	checkoutDTOs := make([]dto.CheckoutDTO, len(checkouts))
	for i, checkout := range checkouts {
		checkoutDTOs[i] = convertToCheckoutDTO(checkout)
	}

	// Create response
	response := dto.CheckoutListResponse{
		ListResponseDTO: dto.ListResponseDTO[dto.CheckoutDTO]{
			Success: true,
			Data:    checkoutDTOs,
			Pagination: dto.PaginationDTO{
				Page:     offset/limit + 1,
				PageSize: limit,
				Total:    len(checkoutDTOs),
			},
		},
	}

	// Return checkouts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAdminCheckout handles retrieving a checkout by ID for admin
func (h *CheckoutHandler) GetAdminCheckout(w http.ResponseWriter, r *http.Request) {
	// Get checkout ID from URL
	vars := mux.Vars(r)
	checkoutID, err := strconv.ParseUint(vars["checkoutId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid checkout ID", http.StatusBadRequest)
		return
	}

	// Get checkout
	checkout, err := h.checkoutUseCase.GetCheckoutByID(uint(checkoutID))
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		http.Error(w, "Failed to get checkout", http.StatusInternalServerError)
		return
	}

	if checkout == nil {
		http.Error(w, "Checkout not found", http.StatusNotFound)
		return
	}

	// Convert checkout to DTO and return response
	checkoutDTO := convertToCheckoutDTO(checkout)
	response := dto.ResponseDTO[dto.CheckoutDTO]{
		Success: true,
		Data:    checkoutDTO,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteAdminCheckout handles deleting a checkout by ID (admin only)
func (h *CheckoutHandler) DeleteAdminCheckout(w http.ResponseWriter, r *http.Request) {
	// Get checkout ID from URL
	vars := mux.Vars(r)
	checkoutID, err := strconv.ParseUint(vars["checkoutId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid checkout ID", http.StatusBadRequest)
		return
	}

	// Delete checkout
	err = h.checkoutUseCase.DeleteCheckout(uint(checkoutID))
	if err != nil {
		h.logger.Error("Failed to delete checkout: %v", err)
		http.Error(w, "Failed to delete checkout", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := dto.ResponseDTO[string]{
		Success: true,
		Message: "Checkout deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to convert checkout entity to DTO
func convertToCheckoutDTO(checkout *entity.Checkout) dto.CheckoutDTO {
	// Convert checkout items
	items := make([]dto.CheckoutItemDTO, len(checkout.Items))
	for i, item := range checkout.Items {
		items[i] = dto.CheckoutItemDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			VariantID:   item.ProductVariantID,
			ProductName: item.ProductName,
			VariantName: item.VariantName,
			SKU:         item.SKU,
			Price:       float64(item.Price) / 100, // Convert cents to currency units
			Quantity:    item.Quantity,
			Weight:      item.Weight,
			Subtotal:    float64(item.Price*int64(item.Quantity)) / 100, // Convert cents to currency units
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	// Convert shipping method if exists
	var shippingMethod *dto.ShippingMethodDTO
	if checkout.ShippingMethod != nil {
		shippingMethod = &dto.ShippingMethodDTO{
			ID:          checkout.ShippingMethod.ID,
			Name:        checkout.ShippingMethod.Name,
			Description: checkout.ShippingMethod.Description,
			Cost:        float64(checkout.ShippingCost) / 100, // Use checkout.ShippingCost instead of ShippingMethod.Cost
		}
	}

	// Convert applied discount if exists
	var appliedDiscount *dto.AppliedDiscountDTO
	if checkout.AppliedDiscount != nil {
		// Based on the code we've seen, it looks like AppliedDiscount has these fields:
		// DiscountID, DiscountCode, DiscountAmount
		appliedDiscount = &dto.AppliedDiscountDTO{
			ID:     checkout.AppliedDiscount.DiscountID,
			Code:   checkout.AppliedDiscount.DiscountCode,
			Type:   "",                                                     // We don't have this info in the AppliedDiscount entity
			Method: "",                                                     // We don't have this info in the AppliedDiscount entity
			Value:  0,                                                      // We don't have this info in the AppliedDiscount entity
			Amount: float64(checkout.AppliedDiscount.DiscountAmount) / 100, // Convert cents to currency units
		}
	}

	// Convert addresses
	shippingAddress := dto.AddressDTO{
		AddressLine1: checkout.ShippingAddr.Street,
		AddressLine2: "", // Entity doesn't have AddressLine2
		City:         checkout.ShippingAddr.City,
		State:        checkout.ShippingAddr.State,
		PostalCode:   checkout.ShippingAddr.PostalCode,
		Country:      checkout.ShippingAddr.Country,
	}

	billingAddress := dto.AddressDTO{
		AddressLine1: checkout.BillingAddr.Street,
		AddressLine2: "", // Entity doesn't have AddressLine2
		City:         checkout.BillingAddr.City,
		State:        checkout.BillingAddr.State,
		PostalCode:   checkout.BillingAddr.PostalCode,
		Country:      checkout.BillingAddr.Country,
	}

	// Convert customer details
	customerDetails := dto.CustomerDetailsDTO{
		Email:    checkout.CustomerDetails.Email,
		Phone:    checkout.CustomerDetails.Phone,
		FullName: checkout.CustomerDetails.FullName,
	}

	return dto.CheckoutDTO{
		ID:               checkout.ID,
		UserID:           checkout.UserID,
		SessionID:        checkout.SessionID,
		Items:            items,
		Status:           string(checkout.Status),
		ShippingAddress:  shippingAddress,
		BillingAddress:   billingAddress,
		ShippingMethodID: checkout.ShippingMethodID,
		ShippingMethod:   shippingMethod,
		PaymentProvider:  checkout.PaymentProvider,
		TotalAmount:      float64(checkout.TotalAmount) / 100,  // Convert cents to currency units
		ShippingCost:     float64(checkout.ShippingCost) / 100, // Convert cents to currency units
		TotalWeight:      checkout.TotalWeight,
		CustomerDetails:  customerDetails,
		Currency:         checkout.Currency,
		DiscountCode:     checkout.DiscountCode,
		DiscountAmount:   float64(checkout.DiscountAmount) / 100, // Convert cents to currency units
		FinalAmount:      float64(checkout.FinalAmount) / 100,    // Convert cents to currency units
		AppliedDiscount:  appliedDiscount,
		CreatedAt:        checkout.CreatedAt,
		UpdatedAt:        checkout.UpdatedAt,
		LastActivityAt:   checkout.LastActivityAt,
		ExpiresAt:        checkout.ExpiresAt,
		CompletedAt:      checkout.CompletedAt,
		ConvertedOrderID: checkout.ConvertedOrderID,
	}
}
