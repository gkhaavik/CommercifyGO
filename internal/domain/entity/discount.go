package entity

import (
	"errors"
	"slices"
	"time"
)

// DiscountType represents the type of discount
type DiscountType string

const (
	// DiscountTypeBasket applies to the entire order
	DiscountTypeBasket DiscountType = "basket"
	// DiscountTypeProduct applies to specific products
	DiscountTypeProduct DiscountType = "product"
)

// DiscountMethod represents how the discount is calculated
type DiscountMethod string

const (
	// DiscountMethodFixed is a fixed amount discount
	DiscountMethodFixed DiscountMethod = "fixed"
	// DiscountMethodPercentage is a percentage discount
	DiscountMethodPercentage DiscountMethod = "percentage"
)

// Discount represents a discount in the system
type Discount struct {
	ID               uint           `json:"id"`
	Code             string         `json:"code"`
	Type             DiscountType   `json:"type"`
	Method           DiscountMethod `json:"method"`
	Value            float64        `json:"value"`
	MinOrderValue    float64        `json:"min_order_value"`
	MaxDiscountValue float64        `json:"max_discount_value"`
	ProductIDs       []uint         `json:"product_ids,omitempty"`
	CategoryIDs      []uint         `json:"category_ids,omitempty"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"end_date"`
	UsageLimit       int            `json:"usage_limit"`
	CurrentUsage     int            `json:"current_usage"`
	Active           bool           `json:"active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// NewDiscount creates a new discount
func NewDiscount(
	code string,
	discountType DiscountType,
	method DiscountMethod,
	value float64,
	minOrderValue float64,
	maxDiscountValue float64,
	productIDs []uint,
	categoryIDs []uint,
	startDate time.Time,
	endDate time.Time,
	usageLimit int,
) (*Discount, error) {
	if code == "" {
		return nil, errors.New("discount code cannot be empty")
	}

	if value <= 0 {
		return nil, errors.New("discount value must be greater than zero")
	}

	if method == DiscountMethodPercentage && value > 100 {
		return nil, errors.New("percentage discount cannot exceed 100%")
	}

	if discountType == DiscountTypeProduct && len(productIDs) == 0 && len(categoryIDs) == 0 {
		return nil, errors.New("product discount must specify at least one product or category")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	now := time.Now()
	return &Discount{
		Code:             code,
		Type:             discountType,
		Method:           method,
		Value:            value,
		MinOrderValue:    minOrderValue,
		MaxDiscountValue: maxDiscountValue,
		ProductIDs:       productIDs,
		CategoryIDs:      categoryIDs,
		StartDate:        startDate,
		EndDate:          endDate,
		UsageLimit:       usageLimit,
		CurrentUsage:     0,
		Active:           true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// IsValid checks if the discount is valid for the current time and usage
func (d *Discount) IsValid() bool {
	now := time.Now()
	return d.Active &&
		now.After(d.StartDate) &&
		now.Before(d.EndDate) &&
		(d.UsageLimit == 0 || d.CurrentUsage < d.UsageLimit)
}

// IsApplicableToOrder checks if the discount is applicable to the given order
func (d *Discount) IsApplicableToOrder(order *Order) bool {
	if !d.IsValid() {
		return false
	}

	// Check minimum order value
	if d.MinOrderValue > 0 && order.TotalAmount < d.MinOrderValue {
		return false
	}

	switch d.Type {
	case DiscountTypeBasket:
		return true
	case DiscountTypeProduct:
		for _, item := range order.Items {
			// Check if the product is directly included
			if slices.Contains(d.ProductIDs, item.ProductID) {
				return true
			}
			// Note: Category check is handled separately in the CalculateDiscount method
			// since we need product details from the repository
		}
		// If we have category IDs but no direct product matches,
		// we still need to check if any product belongs to those categories
		// This is handled in the use case layer
		if len(d.CategoryIDs) > 0 {
			return true
		}
		return false
	}

	return false
}

// CalculateDiscount calculates the discount amount for an order
func (d *Discount) CalculateDiscount(order *Order) float64 {
	if !d.IsApplicableToOrder(order) {
		return 0
	}

	var discountAmount float64

	if d.Type == DiscountTypeBasket {
		// Calculate discount for the entire order
		if d.Method == DiscountMethodFixed {
			discountAmount = d.Value
		} else if d.Method == DiscountMethodPercentage {
			discountAmount = order.TotalAmount * (d.Value / 100)
		}
	} else if d.Type == DiscountTypeProduct {
		// Calculate discount for eligible products only
		for _, item := range order.Items {
			isEligible := slices.Contains(d.ProductIDs, item.ProductID)

			if isEligible {
				itemTotal := float64(item.Quantity) * item.Price
				if d.Method == DiscountMethodFixed {
					// For fixed discount, apply once per item (not per quantity)
					// This matches with the current implementation in ApplyDiscountToOrder
					itemDiscount := min(d.Value, itemTotal)
					discountAmount += itemDiscount
				} else if d.Method == DiscountMethodPercentage {
					itemDiscount := itemTotal * (d.Value / 100)
					discountAmount += itemDiscount
				}
			}
		}
	}

	// Apply maximum discount cap if specified
	if d.MaxDiscountValue > 0 && discountAmount > d.MaxDiscountValue {
		discountAmount = d.MaxDiscountValue
	}

	// Ensure discount doesn't exceed order total
	if discountAmount > order.TotalAmount {
		discountAmount = order.TotalAmount
	}

	return discountAmount
}

// IncrementUsage increments the usage count of the discount
func (d *Discount) IncrementUsage() {
	d.CurrentUsage++
	d.UpdatedAt = time.Now()
}

// AppliedDiscount represents a discount applied to an order
type AppliedDiscount struct {
	DiscountID     uint    `json:"discount_id"`
	DiscountCode   string  `json:"discount_code"`
	DiscountAmount float64 `json:"discount_amount"`
}
