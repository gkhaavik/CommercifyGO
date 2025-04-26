package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockOrderRepository is a mock implementation of the order repository
type MockOrderRepository struct {
	orders               map[uint]*entity.Order
	ordersByUser         map[uint][]*entity.Order
	lastID               uint
	MockIsDiscountIdUsed func(id uint) (bool, error)
}

// NewMockOrderRepository creates a new instance of MockOrderRepository
func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		orders:       make(map[uint]*entity.Order),
		ordersByUser: make(map[uint][]*entity.Order),
		lastID:       0,
		MockIsDiscountIdUsed: func(id uint) (bool, error) {
			return false, nil
		},
	}
}

// Create adds an order to the repository
func (r *MockOrderRepository) Create(order *entity.Order) error {
	// Assign ID
	r.lastID++
	order.ID = r.lastID

	// Store order
	r.orders[order.ID] = order

	// Add to user's orders
	userOrders, exists := r.ordersByUser[order.UserID]
	if !exists {
		userOrders = make([]*entity.Order, 0)
	}
	userOrders = append(userOrders, order)
	r.ordersByUser[order.UserID] = userOrders

	return nil
}

// GetByID retrieves an order by ID
func (r *MockOrderRepository) GetByID(id uint) (*entity.Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

// Update updates an order
func (r *MockOrderRepository) Update(order *entity.Order) error {
	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	// Get current user's orders
	oldOrder := r.orders[order.ID]
	userOrders := r.ordersByUser[oldOrder.UserID]

	// If user ID changed, update ordersByUser mapping
	if oldOrder.UserID != order.UserID {
		// Remove from old user's orders
		for i, o := range userOrders {
			if o.ID == order.ID {
				userOrders = append(userOrders[:i], userOrders[i+1:]...)
				break
			}
		}
		r.ordersByUser[oldOrder.UserID] = userOrders

		// Add to new user's orders
		newUserOrders, exists := r.ordersByUser[order.UserID]
		if !exists {
			newUserOrders = make([]*entity.Order, 0)
		}
		newUserOrders = append(newUserOrders, order)
		r.ordersByUser[order.UserID] = newUserOrders
	}

	// Update the order
	r.orders[order.ID] = order

	return nil
}

// GetByUser retrieves orders for a user with pagination
func (r *MockOrderRepository) GetByUser(userID uint, offset, limit int) ([]*entity.Order, error) {
	userOrders, exists := r.ordersByUser[userID]
	if !exists {
		return []*entity.Order{}, nil
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(userOrders) {
		return []*entity.Order{}, nil
	}
	if end > len(userOrders) {
		end = len(userOrders)
	}

	return userOrders[start:end], nil
}

// ListByStatus retrieves orders by status with pagination
func (r *MockOrderRepository) ListByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	statusOrders := make([]*entity.Order, 0)

	for _, order := range r.orders {
		if entity.OrderStatus(order.Status) == status {
			statusOrders = append(statusOrders, order)
		}
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(statusOrders) {
		return []*entity.Order{}, nil
	}
	if end > len(statusOrders) {
		end = len(statusOrders)
	}

	return statusOrders[start:end], nil
}

// IsDiscountIdUsed checks if a discount ID is used by any order
func (r *MockOrderRepository) IsDiscountIdUsed(id uint) (bool, error) {
	if r.MockIsDiscountIdUsed != nil {
		return r.MockIsDiscountIdUsed(id)
	}

	// Default implementation
	for _, order := range r.orders {
		if order.AppliedDiscount != nil && order.AppliedDiscount.DiscountID == id {
			return true, nil
		}
	}

	return false, nil
}
