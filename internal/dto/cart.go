package dto

// CartDTO represents a shopping cart in the system
type CartDTO struct {
	BaseDTO
	UserID      uint          `json:"user_id"`
	SessionID   string        `json:"session_id"`
	Items       []CartItemDTO `json:"items"`
	TotalAmount float64       `json:"total_amount"`
	Currency    string        `json:"currency"`
}

// CartItemDTO represents an item in a shopping cart
type CartItemDTO struct {
	BaseDTO
	ProductID  uint    `json:"product_id"`
	VariantID  uint    `json:"variant_id,omitempty"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

// AddToCartRequest represents the data needed to add an item to the cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	VariantID uint `json:"variant_id,omitempty"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

// UpdateCartItemRequest represents the data needed to update a cart item
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
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
