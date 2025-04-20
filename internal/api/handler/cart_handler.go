package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
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

// GetCart handles getting the user's cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get or create cart
	cart, err := h.cartUseCase.GetOrCreateCart(userID)
	if err != nil {
		h.logger.Error("Failed to get cart: %v", err)
		http.Error(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	// Return cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// AddToCart handles adding an item to the cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var input usecase.AddToCartInput
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

	// Add to cart
	cart, err := h.cartUseCase.AddToCart(userID, input)
	if err != nil {
		h.logger.Error("Failed to add to cart: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// UpdateCartItem handles updating the quantity of an item in the cart
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	// Update cart item
	cart, err := h.cartUseCase.UpdateCartItem(userID, usecase.UpdateCartItemInput{
		ProductID: uint(productID),
		Quantity:  input.Quantity,
	})
	if err != nil {
		h.logger.Error("Failed to update cart item: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// RemoveFromCart handles removing an item from the cart
func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Remove from cart
	cart, err := h.cartUseCase.RemoveFromCart(userID, uint(productID))
	if err != nil {
		h.logger.Error("Failed to remove from cart: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// ClearCart handles emptying the cart
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Clear cart
	if err := h.cartUseCase.ClearCart(userID); err != nil {
		h.logger.Error("Failed to clear cart: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated cart
	cart, err := h.cartUseCase.GetOrCreateCart(userID)
	if err != nil {
		h.logger.Error("Failed to get cart after clearing: %v", err)
		http.Error(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	// Return empty cart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
