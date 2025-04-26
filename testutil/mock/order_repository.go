package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// OrderRepository is a mock implementation of the order repository interface
type OrderRepository struct {
	orders               map[uint]*entity.Order
	paymentIDIndex       map[string]*entity.Order // Index to find orders by payment ID
	MockIsDiscountIdUsed func(id uint) (bool, error)
}

// NewMockOrderRepository creates a new mock order repository
func NewMockOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders:         make(map[uint]*entity.Order),
		paymentIDIndex: make(map[string]*entity.Order),
	}
}

// Create adds a new order to the mock repository
func (r *OrderRepository) Create(order *entity.Order) error {
	// If no ID provided, generate one
	if order.ID == 0 {
		maxID := uint(0)
		for id := range r.orders {
			if id > maxID {
				maxID = id
			}
		}
		order.ID = maxID + 1
	}

	// Clone the order to prevent unintended modifications
	clone := *order
	r.orders[order.ID] = &clone

	// Index by payment ID if available
	if order.PaymentID != "" {
		r.paymentIDIndex[order.PaymentID] = &clone
	}

	return nil
}

// GetByID retrieves an order by ID from the mock repository
func (r *OrderRepository) GetByID(id uint) (*entity.Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	// Return a clone to prevent unintended modifications
	clone := *order
	return &clone, nil
}

// Update updates an existing order in the mock repository
func (r *OrderRepository) Update(order *entity.Order) error {
	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	// If payment ID has changed, update the index
	existingOrder := r.orders[order.ID]
	if existingOrder.PaymentID != order.PaymentID {
		if existingOrder.PaymentID != "" {
			delete(r.paymentIDIndex, existingOrder.PaymentID)
		}
		if order.PaymentID != "" {
			r.paymentIDIndex[order.PaymentID] = order
		}
	}

	// Clone the order to prevent unintended modifications
	clone := *order
	r.orders[order.ID] = &clone

	return nil
}

// GetByUser retrieves orders for a user from the mock repository
func (r *OrderRepository) GetByUser(userID uint, offset, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			clone := *order
			orders = append(orders, &clone)
		}
	}

	// Apply offset and limit
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}
	end := min(offset+limit, len(orders))

	return orders[offset:end], nil
}

// ListByStatus retrieves orders by status from the mock repository
func (r *OrderRepository) ListByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	for _, order := range r.orders {
		if order.Status == string(status) {
			clone := *order
			orders = append(orders, &clone)
		}
	}

	// Apply offset and limit
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}
	end := min(offset+limit, len(orders))

	return orders[offset:end], nil
}

// IsDiscountIdUsed checks if a discount is used by any order in the mock repository
func (r *OrderRepository) IsDiscountIdUsed(discountID uint) (bool, error) {
	// If a mock function is provided, use it
	if r.MockIsDiscountIdUsed != nil {
		return r.MockIsDiscountIdUsed(discountID)
	}

	// Otherwise fall back to the default implementation
	for _, order := range r.orders {
		if order.AppliedDiscount != nil && order.AppliedDiscount.DiscountID == discountID {
			return true, nil
		}
	}
	return false, nil
}

// GetByPaymentID retrieves an order by payment ID from the mock repository
func (r *OrderRepository) GetByPaymentID(paymentID string) (*entity.Order, error) {
	order, exists := r.paymentIDIndex[paymentID]
	if !exists {
		return nil, errors.New("order not found for payment ID")
	}

	// Return a clone to prevent unintended modifications
	clone := *order
	return &clone, nil
}

// AddMockGetByPaymentID is a helper function to set up mock behavior for GetByPaymentID
func (r *OrderRepository) AddMockGetByPaymentID(order *entity.Order) {
	if order != nil && order.PaymentID != "" {
		r.paymentIDIndex[order.PaymentID] = order
	}
}
