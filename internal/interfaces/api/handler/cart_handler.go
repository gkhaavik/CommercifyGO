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

// CartHandler handles cart-related HTTP requests
type CartHandler struct {
	cartUseCase *usecase.CartUseCase
	logger      logger.Logger
}

// NewCartHandler creates a new CartHandler
func NewCartHandler(cartUseCase *usecase.CartUseCase, logger logger.Logger) *CartHandler {
	return &CartHandler{
		cartUseCase: cartUseCase,
		logger:      logger,
	}
}

// getSessionID gets or creates a session ID for guest users
func (h *CartHandler) getSessionID(w http.ResponseWriter, r *http.Request) string {
	// Check if session cookie exists
	cookie, err := r.Cookie(common.SessionCookieName)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Create new session ID if none exists
	sessionID := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     common.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   common.SessionCookieAge,
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set secure flag if connection is HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return sessionID
}

// GetCart handles getting the user's cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (if authenticated)
	userID, ok := r.Context().Value("user_id").(uint)

	var cart *entity.Cart
	var err error

	if ok && userID > 0 {
		// User is authenticated, get or create user cart
		cart, err = h.cartUseCase.GetOrCreateCart(userID)
	} else {
		// User is a guest, get or create guest cart
		sessionID := h.getSessionID(w, r)
		cart, err = h.cartUseCase.GetOrCreateGuestCart(sessionID)
	}

	if err != nil {
		h.logger.Error("Failed to get cart: %v", err)
		http.Error(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	// Convert entity to DTO
	cartDTO := convertToCartDTO(cart)

	// Return cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartDTO)
}

// AddToCart handles adding an item to the cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input dto.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.ProductID == 0 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	if input.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	cartInput := usecase.AddToCartInput{
		ProductID: input.ProductID,
		VariantID: input.VariantID,
		Quantity:  input.Quantity,
	}

	// Get user ID from context (if authenticated)
	userID, ok := r.Context().Value("user_id").(uint)

	var cart *entity.Cart
	var err error

	if ok && userID > 0 {
		// User is authenticated, add to user cart
		cart, err = h.cartUseCase.AddToCart(userID, cartInput)
	} else {
		// User is a guest, add to guest cart
		sessionID := h.getSessionID(w, r)
		cart, err = h.cartUseCase.AddToGuestCart(sessionID, cartInput)
	}

	if err != nil {
		h.logger.Error("Failed to add to cart: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	cartDTO := convertToCartDTO(cart)

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartDTO)
}

// UpdateCartItem handles updating the quantity of an item in the cart
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input dto.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	updateInput := usecase.UpdateCartItemInput{
		ProductID: uint(productID),
		VariantID: input.VariantID,
		Quantity:  input.Quantity,
	}

	// Get user ID from context (if authenticated)
	userID, ok := r.Context().Value("user_id").(uint)

	var cart *entity.Cart

	if ok && userID > 0 {
		// User is authenticated, update user cart
		cart, err = h.cartUseCase.UpdateCartItem(userID, updateInput)
	} else {
		// User is a guest, update guest cart
		sessionID := h.getSessionID(w, r)
		cart, err = h.cartUseCase.UpdateGuestCartItem(sessionID, updateInput)
	}

	if err != nil {
		h.logger.Error("Failed to update cart item: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	cartDTO := convertToCartDTO(cart)

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartDTO)
}

// RemoveFromCart handles removing an item from the cart
func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get variant ID from query parameters
	variantID, _ := strconv.ParseUint(r.URL.Query().Get("variantId"), 10, 32)

	// Get user ID from context (if authenticated)
	userID, ok := r.Context().Value("user_id").(uint)

	var cart *entity.Cart

	if ok && userID > 0 {
		// User is authenticated, remove from user cart
		cart, err = h.cartUseCase.RemoveFromCart(userID, uint(productID), uint(variantID))
	} else {
		// User is a guest, remove from guest cart
		sessionID := h.getSessionID(w, r)
		cart, err = h.cartUseCase.RemoveFromGuestCart(sessionID, uint(productID), uint(variantID))
	}

	if err != nil {
		h.logger.Error("Failed to remove from cart: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to DTO
	cartDTO := convertToCartDTO(cart)

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartDTO)
}

// ClearCart handles emptying the cart
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (if authenticated)
	userID, ok := r.Context().Value("user_id").(uint)

	var err error
	var cart *entity.Cart

	if ok && userID > 0 {
		// User is authenticated, clear user cart
		err = h.cartUseCase.ClearCart(userID)
		if err == nil {
			// Get updated cart
			cart, err = h.cartUseCase.GetOrCreateCart(userID)
		}
	} else {
		// User is a guest, clear guest cart
		sessionID := h.getSessionID(w, r)
		err = h.cartUseCase.ClearGuestCart(sessionID)
		if err == nil {
			// Get updated cart
			cart, err = h.cartUseCase.GetOrCreateGuestCart(sessionID)
		}
	}

	if err != nil {
		h.logger.Error("Failed to clear cart: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert entity to DTO
	cartDTO := convertToCartDTO(cart)

	// Return empty cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartDTO)
}

// ConvertGuestCartToUserCart converts a guest cart to a user cart
func (h *CartHandler) ConvertGuestCartToUserCart(w http.ResponseWriter, r *http.Request) {
	// Must be authenticated to convert a cart
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session ID from cookie
	cookie, err := r.Cookie(common.SessionCookieName)
	if err != nil || cookie.Value == "" {
		// No guest cart to convert
		cart, err := h.cartUseCase.GetOrCreateCart(userID)
		if err != nil {
			h.logger.Error("Failed to get cart: %v", err)
			http.Error(w, "Failed to get cart", http.StatusInternalServerError)
			return
		}

		// Convert entity to DTO
		cartDTO := convertToCartDTO(cart)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cartDTO)
		return
	}

	// Convert guest cart to user cart
	cart, err := h.cartUseCase.ConvertGuestCartToUserCart(cookie.Value, userID)
	if err != nil {
		h.logger.Error("Failed to convert guest cart to user cart: %v", err)
		http.Error(w, "Failed to convert cart", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     common.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	// Convert entity to DTO
	cartDTO := convertToCartDTO(cart)

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartDTO)
}

// convertToCartDTO converts a cart entity to a DTO
func convertToCartDTO(cart *entity.Cart) dto.CartDTO {
	// Convert cart items to DTOs
	items := make([]dto.CartItemDTO, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = dto.CartItemDTO{
			BaseDTO: dto.BaseDTO{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			},
			ProductID: uint(item.ProductID),
			Quantity:  item.Quantity,
		}
		if item.ProductVariantID > 0 {
			items[i].VariantID = uint(item.ProductVariantID)
		}
	}

	// Create cart DTO
	cartDTO := dto.CartDTO{
		BaseDTO: dto.BaseDTO{
			ID:        cart.ID,
			CreatedAt: cart.CreatedAt,
			UpdatedAt: cart.UpdatedAt,
		},
		Items: items,
	}

	// Set user ID if it exists
	if cart.UserID > 0 {
		cartDTO.UserID = cart.UserID
	}

	// Set session ID if it exists
	if cart.SessionID != "" {
		cartDTO.SessionID = cart.SessionID
	}

	return cartDTO
}
