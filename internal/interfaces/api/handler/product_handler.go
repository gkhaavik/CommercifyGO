package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	logger         logger.Logger
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(productUseCase *usecase.ProductUseCase, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var input usecase.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set seller ID from authenticated user
	input.SellerID = userID

	// Create product
	product, err := h.productUseCase.CreateProduct(input)
	if err != nil {
		h.logger.Error("Failed to create product: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created product
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProduct handles getting a product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get product
	product, err := h.productUseCase.GetProductByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Return product
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct handles updating a product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input usecase.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update product
	product, err := h.productUseCase.UpdateProduct(uint(id), userID, input)
	if err != nil {
		h.logger.Error("Failed to update product: %v", err)
		if err.Error() == "unauthorized: not the seller of this product" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Return updated product
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct handles deleting a product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Delete product
	if err := h.productUseCase.DeleteProduct(uint(id), userID); err != nil {
		h.logger.Error("Failed to delete product: %v", err)
		if err.Error() == "unauthorized: not the seller of this product" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// AddVariant handles adding a variant to a product
func (h *ProductHandler) AddVariant(w http.ResponseWriter, r *http.Request) {
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
	var input usecase.AddVariantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set product ID from URL
	input.ProductID = uint(productID)

	// Add variant
	variant, err := h.productUseCase.AddVariant(userID, input)
	if err != nil {
		h.logger.Error("Failed to add variant: %v", err)
		if err.Error() == "unauthorized: not the seller of this product" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Return created variant
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(variant)
}

// UpdateVariant handles updating a product variant
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get product and variant IDs from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	variantID, err := strconv.ParseUint(vars["variantId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid variant ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input usecase.UpdateVariantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update variant
	variant, err := h.productUseCase.UpdateVariant(uint(productID), uint(variantID), userID, input)
	if err != nil {
		h.logger.Error("Failed to update variant: %v", err)
		if err.Error() == "unauthorized: not the seller of this product" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Return updated variant
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variant)
}

// DeleteVariant handles deleting a product variant
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get product and variant IDs from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	variantID, err := strconv.ParseUint(vars["variantId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid variant ID", http.StatusBadRequest)
		return
	}

	// Delete variant
	if err := h.productUseCase.DeleteVariant(uint(productID), uint(variantID), userID); err != nil {
		h.logger.Error("Failed to delete variant: %v", err)
		if err.Error() == "unauthorized: not the seller of this product" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// ListProducts handles listing products with pagination
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get products
	products, err := h.productUseCase.ListProducts(offset, limit)
	if err != nil {
		h.logger.Error("Failed to list products: %v", err)
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}

	// Return products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// SearchProducts handles searching for products
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse search parameters
	query := r.URL.Query().Get("q")
	categoryIDStr := r.URL.Query().Get("category")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	// Convert parameters to appropriate types
	var categoryID uint
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			categoryID = uint(id)
		}
	}

	var minPrice, maxPrice float64
	if minPriceStr != "" {
		minPrice, _ = strconv.ParseFloat(minPriceStr, 64)
	}
	if maxPriceStr != "" {
		maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
	}

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Search products
	input := usecase.SearchProductsInput{
		Query:      query,
		CategoryID: categoryID,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Offset:     offset,
		Limit:      limit,
	}

	products, err := h.productUseCase.SearchProducts(input)
	if err != nil {
		h.logger.Error("Failed to search products: %v", err)
		http.Error(w, "Failed to search products", http.StatusInternalServerError)
		return
	}

	// Return products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// ListSellerProducts handles listing products for a seller
func (h *ProductHandler) ListSellerProducts(w http.ResponseWriter, r *http.Request) {
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

	// Get seller's products
	products, err := h.productUseCase.ListProductsBySeller(userID, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list seller products: %v", err)
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}

	// Return products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// ListCategories handles listing all categories
func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	// Get categories
	categories, err := h.productUseCase.ListCategories()
	if err != nil {
		h.logger.Error("Failed to list categories: %v", err)
		http.Error(w, "Failed to list categories", http.StatusInternalServerError)
		return
	}

	// Return categories
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
