package mock

import (
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockShippingUseCase is a mock implementation of shipping use case for testing
type MockShippingUseCase struct {
	ShippingMethods []*entity.ShippingMethod
	ShippingZones   []*entity.ShippingZone
	ShippingRates   []*entity.ShippingRate
}

// NewMockShippingUseCase creates a new instance of MockShippingUseCase
func NewMockShippingUseCase() *MockShippingUseCase {
	return &MockShippingUseCase{
		ShippingMethods: []*entity.ShippingMethod{
			{
				ID:          1,
				Name:        "Standard Shipping",
				Description: "Regular shipping option (3-5 business days)",
				Active:      true,
			},
			{
				ID:          2,
				Name:        "Express Shipping",
				Description: "Faster shipping option (1-2 business days)",
				Active:      true,
			},
		},
		ShippingZones: []*entity.ShippingZone{
			{
				ID:        1,
				Name:      "Domestic",
				Countries: []string{"US"},
			},
			{
				ID:        2,
				Name:      "International",
				Countries: []string{"*"},
			},
		},
		ShippingRates: []*entity.ShippingRate{
			{
				ID:               1,
				ShippingZoneID:   1,
				ShippingMethodID: 1,
				MinWeight:        0.0,
				MaxWeight:        1.0,
				MinOrderValue:    0.0,
				MaxOrderValue:    50.0,
				Cost:             5.99,
			},
			{
				ID:               2,
				ShippingZoneID:   1,
				ShippingMethodID: 1,
				MinWeight:        1.1,
				MaxWeight:        5.0,
				MinOrderValue:    0.0,
				MaxOrderValue:    100.0,
				Cost:             8.99,
			},
			{
				ID:               3,
				ShippingZoneID:   1,
				ShippingMethodID: 2,
				MinWeight:        0.0,
				MaxWeight:        5.0,
				MinOrderValue:    0.0,
				MaxOrderValue:    200.0,
				Cost:             15.99,
			},
		},
	}
}

// GetShippingMethod returns a shipping method by ID
func (m *MockShippingUseCase) GetShippingMethod(id uint) (*entity.ShippingMethod, error) {
	for _, method := range m.ShippingMethods {
		if method.ID == id {
			return method, nil
		}
	}
	return nil, entity.ErrNotFound
}

// GetShippingMethods returns all shipping methods
func (m *MockShippingUseCase) GetShippingMethods() ([]*entity.ShippingMethod, error) {
	return m.ShippingMethods, nil
}

// GetShippingOptions returns available shipping options based on the given criteria
func (m *MockShippingUseCase) GetShippingOptions(countryCode string, weight float64, orderValue float64) ([]*entity.ShippingOption, error) {
	options := make([]*entity.ShippingOption, 0)

	// Find applicable rates
	for _, rate := range m.ShippingRates {
		// Check if weight and order value are within range
		if weight >= rate.MinWeight && weight <= rate.MaxWeight &&
			orderValue >= rate.MinOrderValue && orderValue <= rate.MaxOrderValue {

			// Find the corresponding method and zone
			var method *entity.ShippingMethod
			var zone *entity.ShippingZone

			for _, m := range m.ShippingMethods {
				if m.ID == rate.ShippingMethodID {
					method = m
					break
				}
			}

			for _, z := range m.ShippingZones {
				if z.ID == rate.ShippingZoneID {
					zone = z
					break
				}
			}

			if method != nil && zone != nil && method.Active {
				// Check if country is in zone
				countryMatch := false
				for _, c := range zone.Countries {
					if c == countryCode || c == "*" {
						countryMatch = true
						break
					}
				}

				if countryMatch {
					option := &entity.ShippingOption{
						ShippingMethodID: method.ID,
						MethodName:       method.Name,
						Description:      method.Description,
						Cost:             rate.Cost,
					}
					options = append(options, option)
				}
			}
		}
	}

	return options, nil
}

// CalculateShippingCost calculates the shipping cost for a given method and order criteria
func (m *MockShippingUseCase) CalculateShippingCost(methodID uint, countryCode string, weight float64, orderValue float64) (float64, error) {
	// Find a matching shipping rate
	for _, rate := range m.ShippingRates {
		if rate.ShippingMethodID == methodID &&
			weight >= rate.MinWeight && weight <= rate.MaxWeight &&
			orderValue >= rate.MinOrderValue && orderValue <= rate.MaxOrderValue {

			// Verify the zone covers this country
			for _, zone := range m.ShippingZones {
				if zone.ID == rate.ShippingZoneID {
					for _, country := range zone.Countries {
						if country == countryCode || country == "*" {
							return rate.Cost, nil
						}
					}
				}
			}
		}
	}

	return 0, entity.ErrNotFound
}
