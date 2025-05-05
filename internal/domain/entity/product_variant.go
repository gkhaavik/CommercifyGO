package entity

import (
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/money"
)

// VariantAttribute represents a single attribute of a product variant
type VariantAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	ID           uint               `json:"id"`
	ProductID    uint               `json:"product_id"`
	SKU          string             `json:"sku"`
	Price        int64              `json:"price"`                   // Stored as cents
	ComparePrice int64              `json:"compare_price,omitempty"` // Stored as cents
	Stock        int                `json:"stock"`
	Attributes   []VariantAttribute `json:"attributes"`
	Images       []string           `json:"images,omitempty"`
	IsDefault    bool               `json:"is_default"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// NewProductVariant creates a new product variant
func NewProductVariant(productID uint, sku string, price int64, stock int, weight float64, attributes []VariantAttribute, images []string, isDefault bool) (*ProductVariant, error) {
	if productID == 0 {
		return nil, errors.New("product ID cannot be empty")
	}
	if sku == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	if price <= 0 { // Check cents
		return nil, errors.New("price must be greater than zero")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if weight < 0 {
		return nil, errors.New("weight cannot be negative")
	}
	if len(attributes) == 0 {
		return nil, errors.New("variant must have at least one attribute")
	}

	now := time.Now()
	return &ProductVariant{
		ProductID:  productID,
		SKU:        sku,
		Price:      price, // Already in cents
		Stock:      stock,
		Attributes: attributes,
		Images:     images,
		IsDefault:  isDefault,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// SetComparePrice sets the compare price for the variant (input in cents)
func (v *ProductVariant) SetComparePrice(comparePrice int64) error {
	if comparePrice <= 0 { // Check cents
		return errors.New("compare price must be greater than zero")
	}

	v.ComparePrice = comparePrice
	v.UpdatedAt = time.Now()
	return nil
}

// UpdateStock updates the variant's stock
func (v *ProductVariant) UpdateStock(quantity int) error {
	newStock := v.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	v.Stock = newStock
	v.UpdatedAt = time.Now()
	return nil
}

// IsAvailable checks if the variant is available in the requested quantity
func (v *ProductVariant) IsAvailable(quantity int) bool {
	return v.Stock >= quantity
}

// GetTotalWeight calculates the total weight for a quantity of this variant
// func (v *ProductVariant) GetTotalWeight(quantity int) float64 {
// 	if quantity <= 0 {
// 		return 0
// 	}
// 	return v.Weight * float64(quantity)
// }

// GetPriceDollars returns the price in dollars
func (v *ProductVariant) GetPriceDollars() float64 {
	return money.FromCents(v.Price)
}

// GetComparePriceDollars returns the compare price in dollars
func (v *ProductVariant) GetComparePriceDollars() float64 {
	return money.FromCents(v.ComparePrice)
}
