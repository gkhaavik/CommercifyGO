package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockCategoryRepository is a mock implementation of category repository for testing
type MockCategoryRepository struct {
	categories map[uint]*entity.Category
	lastID     uint
}

// NewMockCategoryRepository creates a new instance of MockCategoryRepository
func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		categories: make(map[uint]*entity.Category),
		lastID:     0,
	}
}

// Create adds a category to the repository
func (r *MockCategoryRepository) Create(category *entity.Category) error {
	// Assign ID
	r.lastID++
	category.ID = r.lastID

	// Store category
	r.categories[category.ID] = category

	return nil
}

// GetByID retrieves a category by ID
func (r *MockCategoryRepository) GetByID(id uint) (*entity.Category, error) {
	category, exists := r.categories[id]
	if !exists {
		return nil, errors.New("category not found")
	}
	return category, nil
}

// Update updates a category
func (r *MockCategoryRepository) Update(category *entity.Category) error {
	if _, exists := r.categories[category.ID]; !exists {
		return errors.New("category not found")
	}

	// Update category
	r.categories[category.ID] = category

	return nil
}

// Delete removes a category
func (r *MockCategoryRepository) Delete(id uint) error {
	if _, exists := r.categories[id]; !exists {
		return errors.New("category not found")
	}

	delete(r.categories, id)
	return nil
}

// List retrieves all categories
func (r *MockCategoryRepository) List() ([]*entity.Category, error) {
	result := make([]*entity.Category, 0, len(r.categories))

	for _, category := range r.categories {
		result = append(result, category)
	}

	return result, nil
}

// GetChildren retrieves child categories for a parent category
func (r *MockCategoryRepository) GetChildren(parentID uint) ([]*entity.Category, error) {
	result := make([]*entity.Category, 0)

	for _, category := range r.categories {
		if category.ParentID != nil && *category.ParentID == parentID {
			result = append(result, category)
		}
	}

	return result, nil
}
