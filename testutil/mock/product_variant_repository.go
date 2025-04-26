package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockProductVariantRepository is a mock implementation of product variant repository for testing
type MockProductVariantRepository struct {
	variants          map[uint]*entity.ProductVariant
	variantsBySKU     map[string]*entity.ProductVariant
	variantsByProduct map[uint][]*entity.ProductVariant
	lastID            uint
}

// NewMockProductVariantRepository creates a new instance of MockProductVariantRepository
func NewMockProductVariantRepository() *MockProductVariantRepository {
	return &MockProductVariantRepository{
		variants:          make(map[uint]*entity.ProductVariant),
		variantsBySKU:     make(map[string]*entity.ProductVariant),
		variantsByProduct: make(map[uint][]*entity.ProductVariant),
		lastID:            0,
	}
}

// Create adds a product variant to the repository
func (r *MockProductVariantRepository) Create(variant *entity.ProductVariant) error {
	// Check if SKU already exists
	if _, exists := r.variantsBySKU[variant.SKU]; exists {
		return errors.New("variant with this SKU already exists")
	}

	// Assign ID
	r.lastID++
	variant.ID = r.lastID

	// Store variant
	r.variants[variant.ID] = variant
	r.variantsBySKU[variant.SKU] = variant

	// Add to product's variants
	productVariants, exists := r.variantsByProduct[variant.ProductID]
	if !exists {
		productVariants = make([]*entity.ProductVariant, 0)
	}
	productVariants = append(productVariants, variant)
	r.variantsByProduct[variant.ProductID] = productVariants

	return nil
}

// GetByID retrieves a product variant by ID
func (r *MockProductVariantRepository) GetByID(id uint) (*entity.ProductVariant, error) {
	variant, exists := r.variants[id]
	if !exists {
		return nil, errors.New("product variant not found")
	}
	return variant, nil
}

// GetBySKU retrieves a product variant by SKU
func (r *MockProductVariantRepository) GetBySKU(sku string) (*entity.ProductVariant, error) {
	variant, exists := r.variantsBySKU[sku]
	if !exists {
		return nil, errors.New("product variant not found")
	}
	return variant, nil
}

// GetByProduct retrieves all variants for a product
func (r *MockProductVariantRepository) GetByProduct(productID uint) ([]*entity.ProductVariant, error) {
	variants, exists := r.variantsByProduct[productID]
	if !exists {
		return make([]*entity.ProductVariant, 0), nil
	}
	return variants, nil
}

// Update updates a product variant
func (r *MockProductVariantRepository) Update(variant *entity.ProductVariant) error {
	// Check if variant exists
	oldVariant, exists := r.variants[variant.ID]
	if !exists {
		return errors.New("product variant not found")
	}

	// If SKU changed, update variantsBySKU mapping
	if oldVariant.SKU != variant.SKU {
		delete(r.variantsBySKU, oldVariant.SKU)
		r.variantsBySKU[variant.SKU] = variant
	}

	// Update variant in maps
	r.variants[variant.ID] = variant
	r.variantsBySKU[variant.SKU] = variant

	// Update in product's variants
	productVariants, exists := r.variantsByProduct[variant.ProductID]
	if exists {
		for i, v := range productVariants {
			if v.ID == variant.ID {
				productVariants[i] = variant
				break
			}
		}
		r.variantsByProduct[variant.ProductID] = productVariants
	}

	return nil
}

// Delete deletes a product variant
func (r *MockProductVariantRepository) Delete(id uint) error {
	// Check if variant exists
	variant, exists := r.variants[id]
	if !exists {
		return errors.New("product variant not found")
	}

	// Remove from maps
	delete(r.variants, id)
	delete(r.variantsBySKU, variant.SKU)

	// Remove from product's variants
	productVariants, exists := r.variantsByProduct[variant.ProductID]
	if exists {
		for i, v := range productVariants {
			if v.ID == id {
				productVariants = append(productVariants[:i], productVariants[i+1:]...)
				break
			}
		}
		r.variantsByProduct[variant.ProductID] = productVariants
	}

	return nil
}

// BatchCreate creates multiple product variants in a single transaction
func (r *MockProductVariantRepository) BatchCreate(variants []*entity.ProductVariant) error {
	for _, variant := range variants {
		if err := r.Create(variant); err != nil {
			return err
		}
	}
	return nil
}
