package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(product *entity.Product) error
	GetByID(id uint) (*entity.Product, error)
	GetByIDWithVariants(id uint) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint) error
	List(offset, limit int) ([]*entity.Product, error)
	// Search expects minPriceCents and maxPriceCents as int64 (cents)
	Search(query string, categoryID uint, minPriceCents, maxPriceCents int64, offset, limit int) ([]*entity.Product, error)
	GetBySeller(sellerID uint, offset, limit int) ([]*entity.Product, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(category *entity.Category) error
	GetByID(id uint) (*entity.Category, error)
	Update(category *entity.Category) error
	Delete(id uint) error
	List() ([]*entity.Category, error)
	GetChildren(parentID uint) ([]*entity.Category, error)
}
