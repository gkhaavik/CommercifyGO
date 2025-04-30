package usecase

import (
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ShippingUseCase implements shipping-related use cases
type ShippingUseCase struct {
	shippingMethodRepo repository.ShippingMethodRepository
	shippingZoneRepo   repository.ShippingZoneRepository
	shippingRateRepo   repository.ShippingRateRepository
}

// NewShippingUseCase creates a new ShippingUseCase
func NewShippingUseCase(
	shippingMethodRepo repository.ShippingMethodRepository,
	shippingZoneRepo repository.ShippingZoneRepository,
	shippingRateRepo repository.ShippingRateRepository,
) *ShippingUseCase {
	return &ShippingUseCase{
		shippingMethodRepo: shippingMethodRepo,
		shippingZoneRepo:   shippingZoneRepo,
		shippingRateRepo:   shippingRateRepo,
	}
}

// CreateShippingMethodInput contains the data needed to create a shipping method
type CreateShippingMethodInput struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EstimatedDeliveryDays int    `json:"estimated_delivery_days"`
}

// CreateShippingMethod creates a new shipping method
func (uc *ShippingUseCase) CreateShippingMethod(input CreateShippingMethodInput) (*entity.ShippingMethod, error) {
	// Create shipping method
	method, err := entity.NewShippingMethod(input.Name, input.Description, input.EstimatedDeliveryDays)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.shippingMethodRepo.Create(method); err != nil {
		return nil, err
	}

	return method, nil
}

// GetShippingMethodByID retrieves a shipping method by ID
func (uc *ShippingUseCase) GetShippingMethodByID(id uint) (*entity.ShippingMethod, error) {
	return uc.shippingMethodRepo.GetByID(id)
}

// ListShippingMethods lists all shipping methods
func (uc *ShippingUseCase) ListShippingMethods(activeOnly bool) ([]*entity.ShippingMethod, error) {
	return uc.shippingMethodRepo.List(activeOnly)
}

// UpdateShippingMethodInput contains the data needed to update a shipping method
type UpdateShippingMethodInput struct {
	ID                    uint   `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EstimatedDeliveryDays int    `json:"estimated_delivery_days"`
	Active                bool   `json:"active"`
}

// UpdateShippingMethod updates a shipping method
func (uc *ShippingUseCase) UpdateShippingMethod(input UpdateShippingMethodInput) (*entity.ShippingMethod, error) {
	// Get existing shipping method
	method, err := uc.shippingMethodRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	method.Name = input.Name
	method.Description = input.Description
	method.EstimatedDeliveryDays = input.EstimatedDeliveryDays
	method.Active = input.Active
	method.UpdatedAt = time.Now()

	// Save changes
	if err := uc.shippingMethodRepo.Update(method); err != nil {
		return nil, err
	}

	return method, nil
}

// CreateShippingZoneInput contains the data needed to create a shipping zone
type CreateShippingZoneInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Countries   []string `json:"countries"`
	States      []string `json:"states"`
	ZipCodes    []string `json:"zip_codes"`
}

// CreateShippingZone creates a new shipping zone
func (uc *ShippingUseCase) CreateShippingZone(input CreateShippingZoneInput) (*entity.ShippingZone, error) {
	// Create shipping zone
	zone, err := entity.NewShippingZone(input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.shippingZoneRepo.Create(zone); err != nil {
		return nil, err
	}

	return zone, nil
}

// GetShippingZoneByID retrieves a shipping zone by ID
func (uc *ShippingUseCase) GetShippingZoneByID(id uint) (*entity.ShippingZone, error) {
	return uc.shippingZoneRepo.GetByID(id)
}

// ListShippingZones lists all shipping zones
func (uc *ShippingUseCase) ListShippingZones(activeOnly bool) ([]*entity.ShippingZone, error) {
	return uc.shippingZoneRepo.List(activeOnly)
}

// UpdateShippingZoneInput contains the data needed to update a shipping zone
type UpdateShippingZoneInput struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Countries   []string `json:"countries"`
	States      []string `json:"states"`
	ZipCodes    []string `json:"zip_codes"`
	Active      bool     `json:"active"`
}

// UpdateShippingZone updates a shipping zone
func (uc *ShippingUseCase) UpdateShippingZone(input UpdateShippingZoneInput) (*entity.ShippingZone, error) {
	// Get existing shipping zone
	zone, err := uc.shippingZoneRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	zone.Name = input.Name
	zone.Description = input.Description
	zone.Countries = input.Countries
	zone.States = input.States
	zone.ZipCodes = input.ZipCodes
	zone.Active = input.Active
	zone.UpdatedAt = time.Now()

	// Save changes
	if err := uc.shippingZoneRepo.Update(zone); err != nil {
		return nil, err
	}

	return zone, nil
}

// CreateShippingRateInput contains the data needed to create a shipping rate
type CreateShippingRateInput struct {
	ShippingMethodID      uint     `json:"shipping_method_id"`
	ShippingZoneID        uint     `json:"shipping_zone_id"`
	BaseRate              float64  `json:"base_rate"`
	MinOrderValue         float64  `json:"min_order_value"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold"`
	Active                bool     `json:"active"`
}

// CreateShippingRate creates a new shipping rate
func (uc *ShippingUseCase) CreateShippingRate(input CreateShippingRateInput) (*entity.ShippingRate, error) {
	// Validate shipping method exists
	method, err := uc.shippingMethodRepo.GetByID(input.ShippingMethodID)
	if err != nil {
		return nil, errors.New("shipping method not found")
	}

	// Validate shipping zone exists
	zone, err := uc.shippingZoneRepo.GetByID(input.ShippingZoneID)
	if err != nil {
		return nil, errors.New("shipping zone not found")
	}

	// Create shipping rate
	rate := &entity.ShippingRate{
		ShippingMethodID:      input.ShippingMethodID,
		ShippingZoneID:        input.ShippingZoneID,
		ShippingMethod:        method,
		ShippingZone:          zone,
		BaseRate:              input.BaseRate,
		MinOrderValue:         input.MinOrderValue,
		FreeShippingThreshold: input.FreeShippingThreshold,
		Active:                input.Active,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Save to repository
	if err := uc.shippingRateRepo.Create(rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// CreateWeightBasedRateInput contains the data needed to create a weight-based rate
type CreateWeightBasedRateInput struct {
	ShippingRateID uint    `json:"shipping_rate_id"`
	MinWeight      float64 `json:"min_weight"`
	MaxWeight      float64 `json:"max_weight"`
	Rate           float64 `json:"rate"`
}

// CreateWeightBasedRate creates a weight-based shipping rate
func (uc *ShippingUseCase) CreateWeightBasedRate(input CreateWeightBasedRateInput) (*entity.WeightBasedRate, error) {
	// Validate shipping rate exists
	_, err := uc.shippingRateRepo.GetByID(input.ShippingRateID)
	if err != nil {
		return nil, errors.New("shipping rate not found")
	}

	// Create weight-based rate
	weightRate := &entity.WeightBasedRate{
		ShippingRateID: input.ShippingRateID,
		MinWeight:      input.MinWeight,
		MaxWeight:      input.MaxWeight,
		Rate:           input.Rate,
	}

	// Save to repository
	if err := uc.shippingRateRepo.CreateWeightBasedRate(weightRate); err != nil {
		return nil, err
	}

	return weightRate, nil
}

// CreateValueBasedRateInput contains the data needed to create a value-based rate
type CreateValueBasedRateInput struct {
	ShippingRateID uint    `json:"shipping_rate_id"`
	MinOrderValue  float64 `json:"min_order_value"`
	MaxOrderValue  float64 `json:"max_order_value"`
	Rate           float64 `json:"rate"`
}

// CreateValueBasedRate creates a value-based shipping rate
func (uc *ShippingUseCase) CreateValueBasedRate(input CreateValueBasedRateInput) (*entity.ValueBasedRate, error) {
	// Validate shipping rate exists
	_, err := uc.shippingRateRepo.GetByID(input.ShippingRateID)
	if err != nil {
		return nil, errors.New("shipping rate not found")
	}

	// Create value-based rate
	valueRate := &entity.ValueBasedRate{
		ShippingRateID: input.ShippingRateID,
		MinOrderValue:  input.MinOrderValue,
		MaxOrderValue:  input.MaxOrderValue,
		Rate:           input.Rate,
	}

	// Save to repository
	if err := uc.shippingRateRepo.CreateValueBasedRate(valueRate); err != nil {
		return nil, err
	}

	return valueRate, nil
}

// GetShippingRateByID retrieves a shipping rate by ID
func (uc *ShippingUseCase) GetShippingRateByID(id uint) (*entity.ShippingRate, error) {
	return uc.shippingRateRepo.GetByID(id)
}

// UpdateShippingRateInput contains the data needed to update a shipping rate
type UpdateShippingRateInput struct {
	ID                    uint     `json:"id"`
	BaseRate              float64  `json:"base_rate"`
	MinOrderValue         float64  `json:"min_order_value"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold"`
	Active                bool     `json:"active"`
}

// UpdateShippingRate updates a shipping rate
func (uc *ShippingUseCase) UpdateShippingRate(input UpdateShippingRateInput) (*entity.ShippingRate, error) {
	// Get existing shipping rate
	rate, err := uc.shippingRateRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	rate.BaseRate = input.BaseRate
	rate.MinOrderValue = input.MinOrderValue
	rate.FreeShippingThreshold = input.FreeShippingThreshold
	rate.Active = input.Active
	rate.UpdatedAt = time.Now()

	// Save changes
	if err := uc.shippingRateRepo.Update(rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// ShippingOptions represents available shipping options for an order
type ShippingOptions struct {
	Options []*entity.ShippingOption `json:"options"`
}

// CalculateShippingOptions calculates available shipping options for an order
func (uc *ShippingUseCase) CalculateShippingOptions(address entity.Address, orderValue float64, orderWeight float64) (*ShippingOptions, error) {
	// Get available shipping rates for address and order value
	rates, err := uc.shippingRateRepo.GetAvailableRatesForAddress(address, orderValue)
	if err != nil {
		return nil, err
	}

	options := &ShippingOptions{
		Options: make([]*entity.ShippingOption, 0, len(rates)),
	}

	for _, rate := range rates {
		cost := rate.BaseRate

		// Check if there are weight-based rates
		if len(rate.WeightBasedRates) > 0 {
			for _, weightRate := range rate.WeightBasedRates {
				if orderWeight >= weightRate.MinWeight && (weightRate.MaxWeight == 0 || orderWeight <= weightRate.MaxWeight) {
					cost += weightRate.Rate
					break
				}
			}
		}

		// Check if there are value-based rates
		if len(rate.ValueBasedRates) > 0 {
			for _, valueRate := range rate.ValueBasedRates {
				if orderValue >= valueRate.MinOrderValue && (valueRate.MaxOrderValue == 0 || orderValue <= valueRate.MaxOrderValue) {
					cost += valueRate.Rate
					break
				}
			}
		}

		// Check if free shipping applies
		freeShipping := false
		if rate.FreeShippingThreshold != nil && orderValue >= *rate.FreeShippingThreshold {
			cost = 0
			freeShipping = true
		}

		option := &entity.ShippingOption{
			ShippingRateID:        rate.ID,
			ShippingMethodID:      rate.ShippingMethodID,
			Name:                  rate.ShippingMethod.Name,
			Description:           rate.ShippingMethod.Description,
			EstimatedDeliveryDays: rate.ShippingMethod.EstimatedDeliveryDays,
			Cost:                  cost,
			FreeShipping:          freeShipping,
		}

		options.Options = append(options.Options, option)
	}

	return options, nil
}

// GetShippingCost calculates the shipping cost for a specific shipping rate
func (uc *ShippingUseCase) GetShippingCost(rateID uint, orderValue float64, orderWeight float64) (float64, error) {
	// Get shipping rate
	rate, err := uc.shippingRateRepo.GetByID(rateID)
	if err != nil {
		return 0, err
	}

	// Start with base rate
	cost := rate.BaseRate

	// Check if there are weight-based rates
	if len(rate.WeightBasedRates) > 0 {
		for _, weightRate := range rate.WeightBasedRates {
			if orderWeight >= weightRate.MinWeight && (weightRate.MaxWeight == 0 || orderWeight <= weightRate.MaxWeight) {
				cost += weightRate.Rate
				break
			}
		}
	}

	// Check if there are value-based rates
	if len(rate.ValueBasedRates) > 0 {
		for _, valueRate := range rate.ValueBasedRates {
			if orderValue >= valueRate.MinOrderValue && (valueRate.MaxOrderValue == 0 || orderValue <= valueRate.MaxOrderValue) {
				cost += valueRate.Rate
				break
			}
		}
	}

	// Check if free shipping applies
	if rate.FreeShippingThreshold != nil && orderValue >= *rate.FreeShippingThreshold {
		cost = 0
	}

	return cost, nil
}
