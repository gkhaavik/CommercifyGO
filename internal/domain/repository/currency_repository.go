package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// CurrencyRepository defines the contract for currency operations
type CurrencyRepository interface {
	// Currency operations
	Create(currency *entity.Currency) error
	Update(currency *entity.Currency) error
	Delete(code string) error
	GetByCode(code string) (*entity.Currency, error)
	GetDefault() (*entity.Currency, error)
	List() ([]*entity.Currency, error)
	ListEnabled() ([]*entity.Currency, error)
	SetDefault(code string) error

	// Product price operations
	GetProductPrices(productID uint) ([]entity.ProductPrice, error)
	// SetProductPrices(productID uint, prices []entity.ProductPrice) error
	DeleteProductPrice(productID uint, currencyCode string) error
	// SetProductPrice(price *entity.ProductPrice) error

	// Product variant price operations
	GetVariantPrices(variantID uint) ([]entity.ProductVariantPrice, error)
	// SetVariantPrices(variantID uint, prices []entity.ProductVariantPrice) error
	// SetVariantPrice(prices *entity.ProductVariantPrice) error
	DeleteVariantPrice(variantID uint, currencyCode string) error
}
