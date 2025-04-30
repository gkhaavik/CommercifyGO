package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ShippingMethodRepository defines the interface for shipping method data access
type ShippingMethodRepository interface {
	Create(method *entity.ShippingMethod) error
	GetByID(id uint) (*entity.ShippingMethod, error)
	List(active bool) ([]*entity.ShippingMethod, error)
	Update(method *entity.ShippingMethod) error
	Delete(id uint) error
}

// ShippingZoneRepository defines the interface for shipping zone data access
type ShippingZoneRepository interface {
	Create(zone *entity.ShippingZone) error
	GetByID(id uint) (*entity.ShippingZone, error)
	List(active bool) ([]*entity.ShippingZone, error)
	Update(zone *entity.ShippingZone) error
	Delete(id uint) error
}

// ShippingRateRepository defines the interface for shipping rate data access
type ShippingRateRepository interface {
	Create(rate *entity.ShippingRate) error
	GetByID(id uint) (*entity.ShippingRate, error)
	GetByMethodID(methodID uint) ([]*entity.ShippingRate, error)
	GetByZoneID(zoneID uint) ([]*entity.ShippingRate, error)
	GetAvailableRatesForAddress(address entity.Address, orderValue float64) ([]*entity.ShippingRate, error)
	CreateWeightBasedRate(weightRate *entity.WeightBasedRate) error
	CreateValueBasedRate(valueRate *entity.ValueBasedRate) error
	GetWeightBasedRates(rateID uint) ([]entity.WeightBasedRate, error)
	GetValueBasedRates(rateID uint) ([]entity.ValueBasedRate, error)
	Update(rate *entity.ShippingRate) error
	Delete(id uint) error
}
