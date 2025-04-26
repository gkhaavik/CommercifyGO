package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockCategoryRepository is a mock implementation of the category repository
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

// Delete deletes a category
func (r *MockCategoryRepository) Delete(id uint) error {
	if _, exists := r.categories[id]; !exists {
		return errors.New("category not found")
	}

	delete(r.categories, id)
	return nil
}

// List retrieves all categories
func (r *MockCategoryRepository) List() ([]*entity.Category, error) {
	categories := make([]*entity.Category, 0, len(r.categories))
	for _, category := range r.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

// GetByParent retrieves categories by parent ID
func (r *MockCategoryRepository) GetByParent(parentID uint) ([]*entity.Category, error) {
	categories := make([]*entity.Category, 0)
	for _, category := range r.categories {
		if category.ParentID != nil && *category.ParentID == parentID {
			categories = append(categories, category)
		}
	}
	return categories, nil
}

// GetChildren recursively retrieves all child categories for a category
func (r *MockCategoryRepository) GetChildren(categoryID uint) ([]*entity.Category, error) {
	result := make([]*entity.Category, 0)

	// First, get direct children
	directChildren, err := r.GetByParent(categoryID)
	if err != nil {
		return nil, err
	}

	result = append(result, directChildren...)

	// Then recursively get children of children
	for _, child := range directChildren {
		childrenOfChild, err := r.GetChildren(child.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, childrenOfChild...)
	}

	return result, nil
}
