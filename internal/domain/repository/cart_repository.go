package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// CartRepository defines the interface for cart data access
type CartRepository interface {
	Create(cart *entity.Cart) error
	GetByUserID(userID uint) (*entity.Cart, error)
	Update(cart *entity.Cart) error
	Delete(id uint) error
}
