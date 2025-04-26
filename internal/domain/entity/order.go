package entity

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusPendingAction OrderStatus = "pending_action" // Requires user action (e.g., redirect to payment provider)
	OrderStatusPaid          OrderStatus = "paid"
	OrderStatusCaptured      OrderStatus = "captured" // Payment captured
	OrderStatusShipped       OrderStatus = "shipped"
	OrderStatusDelivered     OrderStatus = "delivered"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusRefunded      OrderStatus = "refunded"
)

// Order represents an order entity
type Order struct {
	ID              uint
	OrderNumber     string
	UserID          uint // 0 for guest orders
	Items           []OrderItem
	TotalAmount     float64
	Status          string
	ShippingAddr    Address
	BillingAddr     Address
	PaymentID       string
	PaymentProvider string
	TrackingCode    string
	ActionURL       string // URL for redirect to payment provider
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CompletedAt     *time.Time

	// Guest information (only used for guest orders where UserID is 0)
	GuestEmail    string
	GuestPhone    string
	GuestFullName string
	IsGuestOrder  bool

	// Discount-related fields
	DiscountAmount  float64
	FinalAmount     float64
	AppliedDiscount *AppliedDiscount
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

// Address represents a shipping or billing address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// NewOrder creates a new order
func NewOrder(userID uint, items []OrderItem, shippingAddr, billingAddr Address) (*Order, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be empty")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	totalAmount := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than zero")
		}
		if item.Price <= 0 {
			return nil, errors.New("item price must be greater than zero")
		}
		item.Subtotal = float64(item.Quantity) * item.Price
		totalAmount += item.Subtotal
	}

	now := time.Now()

	// Generate a friendly order number (will be replaced with actual ID after creation)
	// Format: ORD-YYYYMMDD-TEMP
	orderNumber := fmt.Sprintf("ORD-%s-TEMP", now.Format("20060102"))

	return &Order{
		UserID:         userID,
		OrderNumber:    orderNumber,
		Items:          items,
		TotalAmount:    totalAmount,
		DiscountAmount: 0,
		FinalAmount:    totalAmount, // Initially same as total amount
		Status:         string(OrderStatusPending),
		ShippingAddr:   shippingAddr,
		BillingAddr:    billingAddr,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// NewGuestOrder creates a new order for a guest user
func NewGuestOrder(items []OrderItem, shippingAddr, billingAddr Address, email, phoneNumber, fullName string) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	if email == "" {
		return nil, errors.New("guest email cannot be empty")
	}

	if fullName == "" {
		return nil, errors.New("guest full name cannot be empty")
	}

	totalAmount := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than zero")
		}
		if item.Price <= 0 {
			return nil, errors.New("item price must be greater than zero")
		}
		item.Subtotal = float64(item.Quantity) * item.Price
		totalAmount += item.Subtotal
	}

	now := time.Now()

	// Generate a friendly order number (will be replaced with actual ID after creation)
	// Format: GS-YYYYMMDD-TEMP (GS prefix for guest orders)
	orderNumber := fmt.Sprintf("GS-%s-TEMP", now.Format("20060102"))

	return &Order{
		UserID:         0, // Using 0 to indicate it should be NULL in the database
		OrderNumber:    orderNumber,
		Items:          items,
		TotalAmount:    totalAmount,
		DiscountAmount: 0,
		FinalAmount:    totalAmount, // Initially same as total amount
		Status:         string(OrderStatusPending),
		ShippingAddr:   shippingAddr,
		BillingAddr:    billingAddr,
		CreatedAt:      now,
		UpdatedAt:      now,

		// Guest-specific information
		GuestEmail:    email,
		GuestPhone:    phoneNumber,
		GuestFullName: fullName,
		IsGuestOrder:  true,
	}, nil
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) error {
	if status == "" {
		return errors.New("status cannot be empty")
	}

	// Validate status transitions
	currentStatus := OrderStatus(o.Status)
	if !isValidStatusTransition(currentStatus, status) {
		return fmt.Errorf("invalid status transition: %s -> %s", currentStatus, status)
	}

	o.Status = string(status)
	o.UpdatedAt = time.Now()

	if status == OrderStatusDelivered {
		now := time.Now()
		o.CompletedAt = &now
	}

	return nil
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(from, to OrderStatus) bool {
	// Define valid transitions
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {
			OrderStatusPaid,
			OrderStatusCancelled,
			OrderStatusPendingAction,
		},
		OrderStatusPendingAction: {
			OrderStatusPaid,
			OrderStatusCancelled,
		},
		OrderStatusPaid: {
			OrderStatusCaptured,
			OrderStatusRefunded,
		},
		OrderStatusShipped: {
			OrderStatusDelivered,
		},
		OrderStatusDelivered: {
			OrderStatusRefunded,
		},
		OrderStatusCaptured: {
			OrderStatusShipped,
			OrderStatusRefunded,
		},
		OrderStatusCancelled: {},
		OrderStatusRefunded:  {},
	}

	// If it's the same status, always allow
	if from == to {
		return true
	}

	// Check if the transition is valid
	validNextStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	return slices.Contains(validNextStates, to)
}

// SetPaymentID sets the payment ID for the order
func (o *Order) SetPaymentID(paymentID string) error {
	if paymentID == "" {
		return errors.New("payment ID cannot be empty")
	}

	o.PaymentID = paymentID
	o.UpdatedAt = time.Now()
	return nil
}

// SetPaymentProvider sets the payment provider for the order
func (o *Order) SetPaymentProvider(provider string) error {
	if provider == "" {
		return errors.New("payment provider cannot be empty")
	}

	o.PaymentProvider = provider
	o.UpdatedAt = time.Now()
	return nil
}

// SetTrackingCode sets the tracking code for the order
func (o *Order) SetTrackingCode(trackingCode string) error {
	if trackingCode == "" {
		return errors.New("tracking code cannot be empty")
	}

	o.TrackingCode = trackingCode
	o.UpdatedAt = time.Now()
	return nil
}

// SetOrderNumber sets the order number
func (o *Order) SetOrderNumber(id uint) {
	// Format: ORD-YYYYMMDD-000001
	o.OrderNumber = fmt.Sprintf("ORD-%s-%06d", o.CreatedAt.Format("20060102"), id)
}

// ApplyDiscount applies a discount to the order
func (o *Order) ApplyDiscount(discount *Discount) error {
	if discount == nil {
		return errors.New("discount cannot be nil")
	}

	if !discount.IsValid() {
		return errors.New("discount is not valid")
	}

	if !discount.IsApplicableToOrder(o) {
		return errors.New("discount is not applicable to this order")
	}

	// Calculate discount amount
	discountAmount := discount.CalculateDiscount(o)
	if discountAmount <= 0 {
		return errors.New("discount amount must be greater than zero")
	}

	// Apply the discount
	o.DiscountAmount = discountAmount
	o.FinalAmount = o.TotalAmount - o.DiscountAmount
	o.AppliedDiscount = &AppliedDiscount{
		DiscountID:     discount.ID,
		DiscountCode:   discount.Code,
		DiscountAmount: discountAmount,
	}
	o.UpdatedAt = time.Now()

	return nil
}

// RemoveDiscount removes any applied discount from the order
func (o *Order) RemoveDiscount() {
	o.DiscountAmount = 0
	o.FinalAmount = o.TotalAmount
	o.AppliedDiscount = nil
	o.UpdatedAt = time.Now()
}

// SetActionURL sets the action URL for the order
func (o *Order) SetActionURL(actionURL string) error {
	o.ActionURL = actionURL
	o.UpdatedAt = time.Now()
	return nil
}
