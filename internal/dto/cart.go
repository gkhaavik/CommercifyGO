package dto

// CartDTO represents a shopping cart in the system
type CartDTO struct {
	BaseDTO
	UserID    uint          `json:"user_id"`
	SessionID string        `json:"session_id"`
	Items     []CartItemDTO `json:"items"`
}

// CartItemDTO represents an item in a shopping cart
type CartItemDTO struct {
	BaseDTO
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"`
	Quantity  int  `json:"quantity"`
}

// AddToCartRequest represents the data needed to add an item to the cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
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
