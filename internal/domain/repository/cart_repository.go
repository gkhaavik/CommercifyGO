package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// CartRepository defines the interface for cart data access
type CartRepository interface {
	Create(cart *entity.Cart) error
	GetByUserID(userID uint) (*entity.Cart, error)
	GetBySessionID(sessionID string) (*entity.Cart, error)
	Update(cart *entity.Cart) error
	Delete(cartID uint) error
	ConvertGuestCartToUserCart(sessionID string, userID uint) (*entity.Cart, error)
}
