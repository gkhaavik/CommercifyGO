package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// MockCartRepository is a mock implementation of the CartRepository interface
type MockCartRepository struct {
	mutex sync.Mutex
	carts map[uint]*entity.Cart
	// Map to store carts by session ID for guest carts
	guestCarts map[string]*entity.Cart
	nextID     uint
}

// NewMockCartRepository creates a new mock cart repository
func NewMockCartRepository() repository.CartRepository {
	return &MockCartRepository{
		carts:      make(map[uint]*entity.Cart),
		guestCarts: make(map[string]*entity.Cart),
		nextID:     1,
	}
}

// Create adds a cart to the repository
func (r *MockCartRepository) Create(cart *entity.Cart) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cart.ID = r.nextID
	r.nextID++

	// Store cart based on whether it's a user cart or guest cart
	if cart.SessionID != "" {
		r.guestCarts[cart.SessionID] = cart
	} else {
		r.carts[cart.UserID] = cart
	}

	return nil
}

// CreateWithID adds a cart with the specified ID to the repository
func (r *MockCartRepository) CreateWithID(cart *entity.Cart) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update nextID if the cart's ID is greater
	if cart.ID >= r.nextID {
		r.nextID = cart.ID + 1
	}

	// Store cart based on whether it's a user cart or guest cart
	if cart.SessionID != "" {
		r.guestCarts[cart.SessionID] = cart
	} else {
		r.carts[cart.UserID] = cart
	}

	return nil
}

// GetByUserID retrieves a cart by user ID
func (r *MockCartRepository) GetByUserID(userID uint) (*entity.Cart, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cart, ok := r.carts[userID]
	if !ok {
		return nil, errors.New("cart not found")
	}

	return cart, nil
}

// GetBySessionID retrieves a cart by session ID
func (r *MockCartRepository) GetBySessionID(sessionID string) (*entity.Cart, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cart, ok := r.guestCarts[sessionID]
	if !ok {
		return nil, errors.New("cart not found")
	}

	return cart, nil
}

// Update updates a cart in the repository
func (r *MockCartRepository) Update(cart *entity.Cart) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update based on whether it's a user cart or guest cart
	if cart.SessionID != "" {
		_, ok := r.guestCarts[cart.SessionID]
		if !ok {
			return errors.New("cart not found")
		}
		r.guestCarts[cart.SessionID] = cart
	} else {
		_, ok := r.carts[cart.UserID]
		if !ok {
			return errors.New("cart not found")
		}
		r.carts[cart.UserID] = cart
	}

	return nil
}

// Delete deletes a cart from the repository
func (r *MockCartRepository) Delete(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Find and delete the cart with the given ID
	for userID, cart := range r.carts {
		if cart.ID == id {
			delete(r.carts, userID)
			return nil
		}
	}

	for sessionID, cart := range r.guestCarts {
		if cart.ID == id {
			delete(r.guestCarts, sessionID)
			return nil
		}
	}

	return errors.New("cart not found")
}

// ConvertGuestCartToUserCart converts a guest cart to a user cart
func (r *MockCartRepository) ConvertGuestCartToUserCart(sessionID string, userID uint) (*entity.Cart, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Get guest cart
	guestCart, ok := r.guestCarts[sessionID]
	if !ok {
		return nil, errors.New("cart not found")
	}

	// Check if user already has a cart
	existingCart, userCartExists := r.carts[userID]

	// If user already has a cart, merge items
	if userCartExists {
		// Add items from guest cart to user cart
		for _, item := range guestCart.Items {
			found := false
			for i, userItem := range existingCart.Items {
				// Check both product ID and variant ID to determine if it's the same item
				if userItem.ProductID == item.ProductID && userItem.ProductVariantID == item.ProductVariantID {
					// Update quantity if product and variant already exist
					existingCart.Items[i].Quantity += item.Quantity
					found = true
					break
				}
			}
			if !found {
				// Add new item if product/variant doesn't exist in user cart
				existingCart.Items = append(existingCart.Items, entity.CartItem{
					ID:               r.nextID,
					CartID:           existingCart.ID,
					ProductID:        item.ProductID,
					ProductVariantID: item.ProductVariantID,
					Quantity:         item.Quantity,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				r.nextID++
			}
		}

		// Delete guest cart
		delete(r.guestCarts, sessionID)

		return existingCart, nil
	}

	// Otherwise convert guest cart to user cart
	guestCart.UserID = userID
	guestCart.SessionID = ""

	// Move from guestCarts to carts
	r.carts[userID] = guestCart
	delete(r.guestCarts, sessionID)

	return guestCart, nil
}
