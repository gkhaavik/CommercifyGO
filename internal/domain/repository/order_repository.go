package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(order *entity.Order) error
	GetByID(id uint) (*entity.Order, error)
	Update(order *entity.Order) error
	GetByUser(userID uint, offset, limit int) ([]*entity.Order, error)
	ListByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error)
	IsDiscountIdUsed(discountID uint) (bool, error)
}
