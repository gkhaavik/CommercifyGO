package dto

import "time"

// CartDTO represents a shopping cart in the system
type CartDTO struct {
	ID        uint          `json:"id"`
	UserID    uint          `json:"user_id"`
	SessionID string        `json:"session_id"`
	Items     []CartItemDTO `json:"items"`
	Currency  string        `json:"currency"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// CartItemDTO represents an item in a shopping cart
type CartItemDTO struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	VariantID uint      `json:"variant_id,omitempty"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddToCartRequest represents the data needed to add an item to the cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"`
	Quantity  int  `json:"quantity"`
}

// UpdateCartItemRequest represents the data needed to update a cart item
type UpdateCartItemRequest struct {
	Quantity  int  `json:"quantity"`
	VariantID uint `json:"variant_id,omitempty"`
}

// CartListResponse represents a paginated list of carts
type CartListResponse struct {
	ListResponseDTO[CartDTO]
}

// CartSearchRequest represents the parameters for searching carts
type CartSearchRequest struct {
	UserID uint `json:"user_id,omitempty"`
	PaginationDTO
}
