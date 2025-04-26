package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockDiscountRepository is a mock implementation of the discount repository
type MockDiscountRepository struct {
	discounts      map[uint]*entity.Discount
	discountByCode map[string]*entity.Discount
	lastID         uint
}

// NewMockDiscountRepository creates a new instance of MockDiscountRepository
func NewMockDiscountRepository() *MockDiscountRepository {
	return &MockDiscountRepository{
		discounts:      make(map[uint]*entity.Discount),
		discountByCode: make(map[string]*entity.Discount),
		lastID:         0,
	}
}

// Create adds a discount to the repository
func (r *MockDiscountRepository) Create(discount *entity.Discount) error {
	// Check for duplicate code
	if _, exists := r.discountByCode[discount.Code]; exists {
		return errors.New("discount code already exists")
	}

	// Assign ID
	r.lastID++
	discount.ID = r.lastID

	// Store discount
	r.discounts[discount.ID] = discount
	r.discountByCode[discount.Code] = discount

	return nil
}

// GetByID retrieves a discount by ID
func (r *MockDiscountRepository) GetByID(id uint) (*entity.Discount, error) {
	discount, exists := r.discounts[id]
	if !exists {
		return nil, errors.New("discount not found")
	}
	return discount, nil
}

// GetByCode retrieves a discount by code
func (r *MockDiscountRepository) GetByCode(code string) (*entity.Discount, error) {
	discount, exists := r.discountByCode[code]
	if !exists {
		return nil, errors.New("discount not found")
	}
	return discount, nil
}

// Update updates a discount
func (r *MockDiscountRepository) Update(discount *entity.Discount) error {
	if _, exists := r.discounts[discount.ID]; !exists {
		return errors.New("discount not found")
	}

	// Check if updating the code and if the new code already exists
	if oldDiscount, exists := r.discounts[discount.ID]; exists {
		if oldDiscount.Code != discount.Code {
			if _, codeExists := r.discountByCode[discount.Code]; codeExists {
				return errors.New("discount code already exists")
			}
			// Remove the old code mapping
			delete(r.discountByCode, oldDiscount.Code)
		}
	}

	// Update the discount
	r.discounts[discount.ID] = discount
	r.discountByCode[discount.Code] = discount

	return nil
}

// Delete removes a discount
func (r *MockDiscountRepository) Delete(id uint) error {
	discount, exists := r.discounts[id]
	if !exists {
		return errors.New("discount not found")
	}

	// Remove discount from maps
	delete(r.discountByCode, discount.Code)
	delete(r.discounts, id)

	return nil
}

// List retrieves all discounts with pagination
func (r *MockDiscountRepository) List(offset, limit int) ([]*entity.Discount, error) {
	discounts := make([]*entity.Discount, 0, len(r.discounts))

	// Convert map to slice
	for _, discount := range r.discounts {
		discounts = append(discounts, discount)
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(discounts) {
		return []*entity.Discount{}, nil
	}
	if end > len(discounts) {
		end = len(discounts)
	}

	return discounts[start:end], nil
}

// ListActive retrieves all active discounts with pagination
func (r *MockDiscountRepository) ListActive(offset, limit int) ([]*entity.Discount, error) {
	discounts := make([]*entity.Discount, 0)

	// Filter active discounts
	for _, discount := range r.discounts {
		if discount.IsValid() {
			discounts = append(discounts, discount)
		}
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(discounts) {
		return []*entity.Discount{}, nil
	}
	if end > len(discounts) {
		end = len(discounts)
	}

	return discounts[start:end], nil
}

// IncrementUsage increments the usage count of a discount
func (r *MockDiscountRepository) IncrementUsage(id uint) error {
	discount, exists := r.discounts[id]
	if !exists {
		return errors.New("discount not found")
	}

	discount.IncrementUsage()
	return nil
}
