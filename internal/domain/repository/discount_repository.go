package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// DiscountRepository defines the interface for discount data access
type DiscountRepository interface {
	Create(discount *entity.Discount) error
	GetByID(id uint) (*entity.Discount, error)
	GetByCode(code string) (*entity.Discount, error)
	Update(discount *entity.Discount) error
	Delete(id uint) error
	List(offset, limit int) ([]*entity.Discount, error)
	ListActive(offset, limit int) ([]*entity.Discount, error)
	IncrementUsage(id uint) error
}
