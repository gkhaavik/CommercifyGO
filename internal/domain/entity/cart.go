package entity

import (
	"errors"
	"slices"
	"time"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id,omitempty"`
	SessionID string     `json:"session_id,omitempty"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID               uint      `json:"id"`
	CartID           uint      `json:"cart_id"`
	ProductID        uint      `json:"product_id"`
	ProductVariantID uint      `json:"product_variant_id,omitempty"` // Added field for variant ID
	Quantity         int       `json:"quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewCart creates a new shopping cart for a user
func NewCart(userID uint) (*Cart, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be empty")
	}

	now := time.Now()
	return &Cart{
		UserID:    userID,
		Items:     []CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NewGuestCart creates a new cart for a guest user
func NewGuestCart(sessionID string) (*Cart, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}

	now := time.Now()
	return &Cart{
		SessionID: sessionID,
		Items:     []CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddItem adds a product to the cart
func (c *Cart) AddItem(productID uint, variantID uint, quantity int) error {
	if productID == 0 {
		return errors.New("product ID cannot be empty")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	// Check if the product is already in the cart
	for i, item := range c.Items {
		// Match by both product ID and variant ID (if variant ID is provided)
		if item.ProductID == productID &&
			(variantID == 0 || item.ProductVariantID == variantID) {
			// Update quantity if product already exists
			c.Items[i].Quantity += quantity
			c.Items[i].UpdatedAt = time.Now()
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	// Add new item if product doesn't exist in cart
	now := time.Now()
	c.Items = append(c.Items, CartItem{
		ProductID:        productID,
		ProductVariantID: variantID, // Store variant ID
		Quantity:         quantity,
		CreatedAt:        now,
		UpdatedAt:        now,
	})
	c.UpdatedAt = now

	return nil
}

// UpdateItem updates the quantity of a product in the cart
func (c *Cart) UpdateItem(productID uint, variantID uint, quantity int) error {
	if productID == 0 {
		return errors.New("product ID cannot be empty")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	for i, item := range c.Items {
		// Match by both product ID and variant ID (if variant ID is provided)
		if item.ProductID == productID &&
			(variantID == 0 || item.ProductVariantID == variantID) {
			c.Items[i].Quantity = quantity
			c.Items[i].UpdatedAt = time.Now()
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("product not found in cart")
}

// RemoveItem removes a product from the cart
func (c *Cart) RemoveItem(productID uint, variantID uint) error {
	if productID == 0 {
		return errors.New("product ID cannot be empty")
	}

	for i, item := range c.Items {
		// Match by both product ID and variant ID (if variant ID is provided)
		if item.ProductID == productID &&
			(variantID == 0 || item.ProductVariantID == variantID) {
			// Remove item from slice
			c.Items = slices.Delete(c.Items, i, i+1)
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("product not found in cart")
}

// Clear empties the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
	c.UpdatedAt = time.Now()
}

// TotalItems returns the total number of items in the cart
func (c *Cart) TotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}
