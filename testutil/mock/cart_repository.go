package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockCartRepository is a mock implementation of cart repository for testing
type MockCartRepository struct {
	carts      map[uint]*entity.Cart
	cartByUser map[uint]*entity.Cart
	lastID     uint
}

// NewMockCartRepository creates a new instance of MockCartRepository
func NewMockCartRepository() *MockCartRepository {
	return &MockCartRepository{
		carts:      make(map[uint]*entity.Cart),
		cartByUser: make(map[uint]*entity.Cart),
		lastID:     0,
	}
}

// Create adds a cart to the repository
func (r *MockCartRepository) Create(cart *entity.Cart) error {
	// Assign ID
	r.lastID++
	cart.ID = r.lastID

	// Store cart
	r.carts[cart.ID] = cart
	r.cartByUser[cart.UserID] = cart

	return nil
}

// CreateWithID adds a cart to the repository with a specific ID (for testing)
func (r *MockCartRepository) CreateWithID(cart *entity.Cart) error {
	// Store cart
	r.carts[cart.ID] = cart
	r.cartByUser[cart.UserID] = cart

	// Update lastID if necessary
	if cart.ID > r.lastID {
		r.lastID = cart.ID
	}

	return nil
}

// GetByID retrieves a cart by ID
func (r *MockCartRepository) GetByID(id uint) (*entity.Cart, error) {
	cart, exists := r.carts[id]
	if !exists {
		return nil, errors.New("cart not found")
	}
	return cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *MockCartRepository) GetByUserID(userID uint) (*entity.Cart, error) {
	cart, exists := r.cartByUser[userID]
	if !exists {
		return nil, errors.New("cart not found")
	}
	return cart, nil
}

// Update updates a cart
func (r *MockCartRepository) Update(cart *entity.Cart) error {
	if _, exists := r.carts[cart.ID]; !exists {
		return errors.New("cart not found")
	}

	// Update cart
	r.carts[cart.ID] = cart
	r.cartByUser[cart.UserID] = cart

	return nil
}

// Delete deletes a cart
func (r *MockCartRepository) Delete(id uint) error {
	cart, exists := r.carts[id]
	if !exists {
		return errors.New("cart not found")
	}

	// Remove cart from both maps
	delete(r.cartByUser, cart.UserID)
	delete(r.carts, id)

	return nil
}
