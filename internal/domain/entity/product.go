package entity

import (
	"errors"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	CategoryID  uint              `json:"category_id"`
	SellerID    uint              `json:"seller_id"`
	Images      []string          `json:"images"`
	HasVariants bool              `json:"has_variants"`
	Variants    []*ProductVariant `json:"variants,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// NewProduct creates a new product with the given details
func NewProduct(name, description string, price float64, stock int, categoryID, sellerID uint, images []string) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	now := time.Now()
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CategoryID:  categoryID,
		SellerID:    sellerID,
		Images:      images,
		HasVariants: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateStock updates the product's stock
func (p *Product) UpdateStock(quantity int) error {
	// If product has variants, stock should be managed at variant level
	if p.HasVariants {
		return errors.New("product has variants, stock should be updated at variant level")
	}

	newStock := p.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}
	p.Stock = newStock
	p.UpdatedAt = time.Now()
	return nil
}

// IsAvailable checks if the product is available in the requested quantity
func (p *Product) IsAvailable(quantity int) bool {
	// If product has variants, availability should be checked at variant level
	if p.HasVariants {
		return false
	}
	return p.Stock >= quantity
}

// AddVariant adds a variant to the product
func (p *Product) AddVariant(variant *ProductVariant) error {
	if variant == nil {
		return errors.New("variant cannot be nil")
	}

	// Ensure variant belongs to this product
	if variant.ProductID != p.ID {
		return errors.New("variant does not belong to this product")
	}

	// Set product to have variants
	p.HasVariants = true

	// If this is the first variant and it's the default, set product price to match
	if len(p.Variants) == 0 && variant.IsDefault {
		p.Price = variant.Price
	}

	// Add variant to product
	p.Variants = append(p.Variants, variant)
	p.UpdatedAt = time.Now()

	return nil
}

// GetDefaultVariant returns the default variant of the product
func (p *Product) GetDefaultVariant() *ProductVariant {
	if !p.HasVariants || len(p.Variants) == 0 {
		return nil
	}

	for _, variant := range p.Variants {
		if variant.IsDefault {
			return variant
		}
	}

	// If no default is set, return the first variant
	return p.Variants[0]
}

// GetVariantByID returns a variant by its ID
func (p *Product) GetVariantByID(variantID uint) *ProductVariant {
	if !p.HasVariants || len(p.Variants) == 0 {
		return nil
	}

	for _, variant := range p.Variants {
		if variant.ID == variantID {
			return variant
		}
	}

	return nil
}

// GetVariantBySKU returns a variant by its SKU
func (p *Product) GetVariantBySKU(sku string) *ProductVariant {
	if !p.HasVariants || len(p.Variants) == 0 || sku == "" {
		return nil
	}

	for _, variant := range p.Variants {
		if variant.SKU == sku {
			return variant
		}
	}

	return nil
}

// Category represents a product category
type Category struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCategory creates a new category
func NewCategory(name, description string, parentID *uint) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	now := time.Now()
	return &Category{
		Name:        name,
		Description: description,
		ParentID:    parentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
