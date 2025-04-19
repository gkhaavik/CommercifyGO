package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

func TestNewOrder(t *testing.T) {
	// Create valid order items
	items := []entity.OrderItem{
		{
			ProductID: 1,
			Quantity:  2,
			Price:     10.0,
		},
		{
			ProductID: 2,
			Quantity:  1,
			Price:     15.0,
		},
	}

	// Create addresses
	shippingAddr := entity.Address{
		Street:     "123 Shipping St",
		City:       "Shipping City",
		State:      "SS",
		PostalCode: "12345",
		Country:    "Shipping Country",
	}

	billingAddr := entity.Address{
		Street:     "456 Billing St",
		City:       "Billing City",
		State:      "BS",
		PostalCode: "67890",
		Country:    "Billing Country",
	}

	tests := []struct {
		name         string
		userID       uint
		items        []entity.OrderItem
		shippingAddr entity.Address
		billingAddr  entity.Address
		wantErr      bool
	}{
		{
			name:         "Valid order",
			userID:       1,
			items:        items,
			shippingAddr: shippingAddr,
			billingAddr:  billingAddr,
			wantErr:      false,
		},
		{
			name:         "Zero user ID",
			userID:       0,
			items:        items,
			shippingAddr: shippingAddr,
			billingAddr:  billingAddr,
			wantErr:      true,
		},
		{
			name:         "No items",
			userID:       1,
			items:        []entity.OrderItem{},
			shippingAddr: shippingAddr,
			billingAddr:  billingAddr,
			wantErr:      true,
		},
		{
			name:   "Invalid item quantity",
			userID: 1,
			items: []entity.OrderItem{
				{
					ProductID: 1,
					Quantity:  0, // Invalid quantity
					Price:     10.0,
				},
			},
			shippingAddr: shippingAddr,
			billingAddr:  billingAddr,
			wantErr:      true,
		},
		{
			name:   "Invalid item price",
			userID: 1,
			items: []entity.OrderItem{
				{
					ProductID: 1,
					Quantity:  2,
					Price:     0, // Invalid price
				},
			},
			shippingAddr: shippingAddr,
			billingAddr:  billingAddr,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := entity.NewOrder(
				tt.userID,
				tt.items,
				tt.shippingAddr,
				tt.billingAddr,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.userID, order.UserID)
				assert.Equal(t, string(entity.OrderStatusPending), order.Status)
				assert.Equal(t, tt.shippingAddr, order.ShippingAddr)
				assert.Equal(t, tt.billingAddr, order.BillingAddr)

				// Check total amount calculation
				expectedTotal := 0.0
				for _, item := range tt.items {
					expectedTotal += float64(item.Quantity) * item.Price

					// Check each item's subtotal calculation
					for _, orderItem := range order.Items {
						if orderItem.ProductID == item.ProductID {
							assert.Equal(t, float64(item.Quantity)*item.Price, orderItem.Subtotal)
						}
					}
				}
				assert.Equal(t, expectedTotal, order.TotalAmount)
			}
		})
	}
}

func TestOrder_UpdateStatus(t *testing.T) {
	// Create a simple order for testing
	items := []entity.OrderItem{
		{
			ProductID: 1,
			Quantity:  2,
			Price:     10.0,
		},
	}

	shippingAddr := entity.Address{
		Street: "123 Test St",
		City:   "Test City",
	}

	billingAddr := entity.Address{
		Street: "123 Test St",
		City:   "Test City",
	}

	order, err := entity.NewOrder(1, items, shippingAddr, billingAddr)
	assert.NoError(t, err)
	assert.Equal(t, string(entity.OrderStatusPending), order.Status)

	// Test updating to paid status
	err = order.UpdateStatus(entity.OrderStatusPaid)
	assert.NoError(t, err)
	assert.Equal(t, string(entity.OrderStatusPaid), order.Status)
	assert.Nil(t, order.CompletedAt)

	// Test updating to delivered status (should set CompletedAt)
	err = order.UpdateStatus(entity.OrderStatusDelivered)
	assert.NoError(t, err)
	assert.Equal(t, string(entity.OrderStatusDelivered), order.Status)
	assert.NotNil(t, order.CompletedAt)

	// Test empty status
	err = order.UpdateStatus("")
	assert.Error(t, err)
}

func TestOrder_SetPaymentID(t *testing.T) {
	order, _ := entity.NewOrder(
		1,
		[]entity.OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
		entity.Address{Street: "Test St"},
		entity.Address{Street: "Test St"},
	)

	// Test setting valid payment ID
	err := order.SetPaymentID("payment-123")
	assert.NoError(t, err)
	assert.Equal(t, "payment-123", order.PaymentID)

	// Test setting empty payment ID
	err = order.SetPaymentID("")
	assert.Error(t, err)
}

func TestOrder_SetTrackingCode(t *testing.T) {
	order, _ := entity.NewOrder(
		1,
		[]entity.OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
		entity.Address{Street: "Test St"},
		entity.Address{Street: "Test St"},
	)

	// Test setting valid tracking code
	err := order.SetTrackingCode("track-123")
	assert.NoError(t, err)
	assert.Equal(t, "track-123", order.TrackingCode)

	// Test setting empty tracking code
	err = order.SetTrackingCode("")
	assert.Error(t, err)
}
