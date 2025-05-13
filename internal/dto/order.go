package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/service"
)

// OrderDTO represents an order in the system
type OrderDTO struct {
	BaseDTO
	UserID           uint             `json:"user_id"`
	OrderNumber      string           `json:"order_number"`
	Items            []OrderItemDTO   `json:"items,omitempty"`
	Status           OrderStatus      `json:"status"`
	TotalAmount      float64          `json:"total_amount"`
	FinalAmount      float64          `json:"final_amount"`
	Currency         string           `json:"currency"`
	ShippingAddress  *AddressDTO      `json:"shipping_address,omitempty"`
	BillingAddress   *AddressDTO      `json:"billing_address,omitempty"`
	PaymentProvider  PaymentProvider  `json:"payment_provider"`
	PaymentID        string           `json:"payment_id"`
	ShippingMethodID uint             `json:"shipping_method_id"`
	ShippingCost     float64          `json:"shipping_cost"`
	DiscountAmount   float64          `json:"discount_amount"`
	DiscountCode     string           `json:"discount_code,omitempty"`
	Customer         *CustomerDetails `json:"customer,omitempty"`
	ActionURL        string           `json:"action_url,omitempty"`
}

type CustomerDetails struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

// OrderItemDTO represents an item in an order
type OrderItemDTO struct {
	BaseDTO
	OrderID    uint    `json:"order_id"`
	ProductID  uint    `json:"product_id"`
	VariantID  uint    `json:"variant_id,omitempty"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

// AddressDTO represents a shipping or billing address
type AddressDTO struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
}

// CreateOrderRequest represents the data needed to create a new order
type CreateOrderRequest struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	PhoneNumber      string     `json:"phone_number,omitempty"`
	ShippingAddress  AddressDTO `json:"shipping_address"`
	BillingAddress   AddressDTO `json:"billing_address"`
	ShippingMethodID uint       `json:"shipping_method_id"`
}

// CreateOrderItemRequest represents the data needed to create a new order item
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"`
	Quantity  int  `json:"quantity"`
}

// UpdateOrderRequest represents the data needed to update an existing order
type UpdateOrderRequest struct {
	Status            string     `json:"status,omitempty"`
	PaymentStatus     string     `json:"payment_status,omitempty"`
	TrackingNumber    string     `json:"tracking_number,omitempty"`
	EstimatedDelivery *time.Time `json:"estimated_delivery,omitempty"`
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	ListResponseDTO[OrderDTO]
}

// OrderSearchRequest represents the parameters for searching orders
type OrderSearchRequest struct {
	UserID        uint        `json:"user_id,omitempty"`
	Status        OrderStatus `json:"status,omitempty"`
	PaymentStatus string      `json:"payment_status,omitempty"`
	StartDate     *time.Time  `json:"start_date,omitempty"`
	EndDate       *time.Time  `json:"end_date,omitempty"`
	PaginationDTO
}

// ProcessPaymentRequest represents the data needed to process a payment
type ProcessPaymentRequest struct {
	PaymentMethod   PaymentMethod        `json:"payment_method"`
	PaymentProvider PaymentProvider      `json:"payment_provider"`
	CardDetails     *service.CardDetails `json:"card_details,omitempty"`
	PhoneNumber     string               `json:"phone_number,omitempty"`
}

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

// PaymentMethod represents the payment method used for an order
type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "credit_card"
	PaymentMethodWallet PaymentMethod = "wallet"
	PaymentMethodBank   PaymentMethod = "bank_transfer"
)

// PaymentProvider represents the payment provider used for an order
type PaymentProvider string

const (
	PaymentProviderStripe    PaymentProvider = "stripe"
	PaymentProviderMobilePay PaymentProvider = "mobilepay"
)
