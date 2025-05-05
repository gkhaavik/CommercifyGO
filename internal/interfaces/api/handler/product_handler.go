package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
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

// --- Response Structs --- //

// ProductVariantResponse is the API representation of a product variant (prices in dollars)
type ProductVariantResponse struct {
	ID           uint                      `json:"id"`
	ProductID    uint                      `json:"product_id"`
	SKU          string                    `json:"sku"`
	Price        float64                   `json:"price"`
	ComparePrice float64                   `json:"compare_price,omitempty"`
	Stock        int                       `json:"stock"`
	Attributes   []entity.VariantAttribute `json:"attributes"`
	Images       []string                  `json:"images,omitempty"`
	IsDefault    bool                      `json:"is_default"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// ProductResponse is the API representation of a product (prices in dollars)
type ProductResponse struct {
	ID            uint                      `json:"id"`
	ProductNumber string                    `json:"product_number"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Price         float64                   `json:"price"`
	Stock         int                       `json:"stock"`
	Weight        float64                   `json:"weight"`
	CategoryID    uint                      `json:"category_id"`
	SellerID      uint                      `json:"seller_id"`
	Images        []string                  `json:"images"`
	HasVariants   bool                      `json:"has_variants"`
	Variants      []*ProductVariantResponse `json:"variants,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

// --- Helper Functions --- //

func toProductVariantResponse(variant *entity.ProductVariant) *ProductVariantResponse {
	if variant == nil {
		return nil
	}
	return &ProductVariantResponse{
		ID:           variant.ID,
		ProductID:    variant.ProductID,
		SKU:          variant.SKU,
		Price:        money.FromCents(variant.Price),
		ComparePrice: money.FromCents(variant.ComparePrice),
		Stock:        variant.Stock,
		Attributes:   variant.Attributes,
		Images:       variant.Images,
		IsDefault:    variant.IsDefault,
		CreatedAt:    variant.CreatedAt, // Assign directly
		UpdatedAt:    variant.UpdatedAt, // Assign directly
	}
}

func toProductResponse(product *entity.Product) *ProductResponse {
	if product == nil {
		return nil
	}
	variantsResponse := make([]*ProductVariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		variantsResponse[i] = toProductVariantResponse(v)
	}

	return &ProductResponse{
		ID:            product.ID,
		ProductNumber: product.ProductNumber,
		Name:          product.Name,
		Description:   product.Description,
		Price:         money.FromCents(product.Price),
		Stock:         product.Stock,
		Weight:        product.Weight,
		CategoryID:    product.CategoryID,
		SellerID:      product.SellerID,
		Images:        product.Images,
		HasVariants:   product.HasVariants,
		Variants:      variantsResponse,
		CreatedAt:     product.CreatedAt, // Assign directly
		UpdatedAt:     product.UpdatedAt, // Assign directly
	}
}

func toProductListResponse(products []*entity.Product) []*ProductResponse {
	list := make([]*ProductResponse, len(products))
	for i, p := range products {
		list[i] = toProductResponse(p)
	}
	return list
}

// --- Handlers --- //

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body (expects float64 for prices)
	var input usecase.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set seller ID from authenticated user
	input.SellerID = userID

	// Create product (use case handles conversion to cents)
	product, err := h.productUseCase.CreateProduct(input)
	if err != nil {
		h.logger.Error("Failed to create product: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert entity to response struct (converts cents to dollars)
	response := toProductResponse(product)

	// Return created product response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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

	// Get product (use case returns entity with cents)
	currencyCode := vars["currency"]

	var product *entity.Product

	if currencyCode != "" {
		// Get product with specific currency prices
		product, err = h.productUseCase.GetProductByCurrency(uint(id), currencyCode)
	} else {
		// Get product with default currency prices
		product, err = h.productUseCase.GetProductByID(uint(id))
	}

	if err != nil {
		h.logger.Error("Failed to get product: %v", err)
		if err.Error() == "product not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Convert entity to response struct (converts cents to dollars)
	response := toProductResponse(product)

	// Return product response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Parse request body (expects float64 for prices)
	var input usecase.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update product (use case handles conversion to cents)
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

	// Convert entity to response struct (converts cents to dollars)
	response := toProductResponse(product)

	// Return updated product response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Parse request body (expects float64 for prices)
	var input usecase.AddVariantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set product ID from URL
	input.ProductID = uint(productID)

	// Add variant (use case handles conversion to cents)
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

	// Convert entity to response struct (converts cents to dollars)
	response := toProductVariantResponse(variant)

	// Return created variant response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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

	// Parse request body (expects float64 for prices)
	var input usecase.UpdateVariantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update variant (use case handles conversion to cents)
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

	// Convert entity to response struct (converts cents to dollars)
	response := toProductVariantResponse(variant)

	// Return updated variant response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Get products (use case returns entities with cents)
	products, err := h.productUseCase.ListProducts(offset, limit)
	if err != nil {
		h.logger.Error("Failed to list products: %v", err)
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}

	// Convert entities to response structs (converts cents to dollars)
	response := toProductListResponse(products)

	// Return products response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchProducts handles searching for products
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse search parameters (prices are float64 dollars)
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

	// Search products (use case expects float64 dollars)
	input := usecase.SearchProductsInput{
		Query:      query,
		CategoryID: categoryID,
		MinPrice:   minPrice, // Pass dollars
		MaxPrice:   maxPrice, // Pass dollars
		Offset:     offset,
		Limit:      limit,
	}

	products, err := h.productUseCase.SearchProducts(input)
	if err != nil {
		h.logger.Error("Failed to search products: %v", err)
		http.Error(w, "Failed to search products", http.StatusInternalServerError)
		return
	}

	// Convert entities to response structs (converts cents to dollars)
	response := toProductListResponse(products)

	// Return products response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Get seller's products (use case returns entities with cents)
	products, err := h.productUseCase.ListProductsBySeller(userID, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list seller products: %v", err)
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}

	// Convert entities to response structs (converts cents to dollars)
	response := toProductListResponse(products)

	// Return products response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
