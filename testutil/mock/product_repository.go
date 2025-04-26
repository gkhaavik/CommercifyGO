package mock

import (
	"errors"
	"strings"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockProductRepository is a mock implementation of product repository for testing
type MockProductRepository struct {
	products map[uint]*entity.Product
	bySeller map[uint][]*entity.Product
	lastID   uint
	// Add MockSearch function field for custom search behavior in tests
	MockSearch func(query string, categoryID uint, minPrice, maxPrice float64, offset, limit int) ([]*entity.Product, error)
}

// NewMockProductRepository creates a new instance of MockProductRepository
func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[uint]*entity.Product),
		bySeller: make(map[uint][]*entity.Product),
		lastID:   0,
	}
}

// Create adds a product to the repository
func (r *MockProductRepository) Create(product *entity.Product) error {
	// Assign ID
	r.lastID++
	product.ID = r.lastID

	// Store product
	r.products[product.ID] = product

	return nil
}

// GetByID retrieves a product by ID
func (r *MockProductRepository) GetByID(id uint) (*entity.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// GetByIDWithVariants retrieves a product by ID including its variants
func (r *MockProductRepository) GetByIDWithVariants(id uint) (*entity.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	// Return a copy of the product to prevent unintended modifications
	productCopy := *product

	return &productCopy, nil
}

// Update updates a product
func (r *MockProductRepository) Update(product *entity.Product) error {
	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}

	// Update product
	r.products[product.ID] = product

	return nil
}

// Delete removes a product
func (r *MockProductRepository) Delete(id uint) error {
	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)
	return nil
}

// List retrieves products with pagination
func (r *MockProductRepository) List(offset, limit int) ([]*entity.Product, error) {
	result := make([]*entity.Product, 0)
	count := 0
	skip := offset

	for _, product := range r.products {
		if skip > 0 {
			skip--
			continue
		}

		result = append(result, product)
		count++

		if count >= limit {
			break
		}
	}

	return result, nil
}

// Search searches for products based on criteria
func (r *MockProductRepository) Search(query string, categoryID uint, minPrice, maxPrice float64, offset, limit int) ([]*entity.Product, error) {
	if r.MockSearch != nil {
		// Use the custom search function if provided
		return r.MockSearch(query, categoryID, minPrice, maxPrice, offset, limit)
	}

	result := make([]*entity.Product, 0)
	count := 0
	skip := offset

	for _, product := range r.products {
		// Apply search filters
		if query != "" && !strings.Contains(strings.ToLower(product.Name), strings.ToLower(query)) &&
			!strings.Contains(strings.ToLower(product.Description), strings.ToLower(query)) {
			continue
		}

		if categoryID > 0 && product.CategoryID != categoryID {
			continue
		}

		if minPrice > 0 && product.Price < minPrice {
			continue
		}

		if maxPrice > 0 && product.Price > maxPrice {
			continue
		}

		// Apply pagination
		if skip > 0 {
			skip--
			continue
		}

		result = append(result, product)
		count++

		if count >= limit {
			break
		}
	}

	return result, nil
}

// GetBySeller retrieves products by seller ID
func (r *MockProductRepository) GetBySeller(sellerID uint, offset, limit int) ([]*entity.Product, error) {
	result := make([]*entity.Product, 0)
	count := 0
	skip := offset

	for _, product := range r.products {
		if product.SellerID != sellerID {
			continue
		}

		if skip > 0 {
			skip--
			continue
		}

		result = append(result, product)
		count++

		if count >= limit {
			break
		}
	}

	return result, nil
}
