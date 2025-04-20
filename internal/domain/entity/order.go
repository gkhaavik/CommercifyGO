package entity

import (
	"errors"
	"fmt"
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID              uint        `json:"id"`
	OrderNumber     string      `json:"order_number"`
	UserID          uint        `json:"user_id"`
	Items           []OrderItem `json:"items"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `json:"status"`
	ShippingAddr    Address     `json:"shipping_address"`
	BillingAddr     Address     `json:"billing_address"`
	PaymentID       string      `json:"payment_id"`
	PaymentProvider string      `json:"payment_provider"`
	TrackingCode    string      `json:"tracking_code"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CompletedAt     *time.Time  `json:"completed_at"`
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
		UserID:       userID,
		OrderNumber:  orderNumber,
		Items:        items,
		TotalAmount:  totalAmount,
		Status:       string(OrderStatusPending),
		ShippingAddr: shippingAddr,
		BillingAddr:  billingAddr,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) error {
	if status == "" {
		return errors.New("status cannot be empty")
	}

	o.Status = string(status)
	o.UpdatedAt = time.Now()

	if status == OrderStatusDelivered {
		now := time.Now()
		o.CompletedAt = &now
	}

	return nil
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
