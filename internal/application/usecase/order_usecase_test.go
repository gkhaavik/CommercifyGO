package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
	"github.com/zenfulcode/commercify/testutil/mock"
)

const (
	// Define payment method constants for testing
	PaymentMethodCard   = "card"
	PaymentMethodBank   = "bank"
	PaymentMethodPayPal = "paypal"
)

func TestOrderUseCase_CreateOrderFromCart(t *testing.T) {
	t.Run("Create order successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Simple mock payment service that always succeeds
		paymentSvc := payment.NewMockPaymentService()
		// paymentSvc := &mockPaymentService{
		// 	availableProviders: []service.PaymentProvider{
		// 		{
		// 			Type:    service.PaymentProviderStripe,
		// 			Name:    "Stripe",
		// 			Enabled: true,
		// 		},
		// 	},
		// }

		// Simple mock email service
		emailSvc := &mockEmailService{}

		// Create a test user
		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Password: "hashed_password",
		}
		userRepo.Create(user)

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			CategoryID:  1,
			SellerID:    2,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: false,
		}
		productRepo.Create(product)

		// Create a test cart with one item
		cart := &entity.Cart{
			ID:     1,
			UserID: 1,
			Items: []entity.CartItem{
				{
					ID:        1,
					ProductID: 1,
					Quantity:  2,
					CartID:    1,
				},
			},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Create order input
		input := usecase.CreateOrderInput{
			UserID: 1,
			ShippingAddr: entity.Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
			BillingAddr: entity.Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
		}

		// Execute
		order, err := orderUseCase.CreateOrderFromCart(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, input.UserID, order.UserID)
		assert.Equal(t, input.ShippingAddr, order.ShippingAddr)
		assert.Equal(t, input.BillingAddr, order.BillingAddr)
		assert.Len(t, order.Items, 1)
		assert.Equal(t, uint(1), order.Items[0].ProductID)
		assert.Equal(t, 2, order.Items[0].Quantity)
		assert.Equal(t, 99.99, order.Items[0].Price)
		assert.Equal(t, 99.99*2, order.Items[0].Subtotal)
		assert.Equal(t, 99.99*2, order.TotalAmount)
		assert.Equal(t, 99.99*2, order.FinalAmount) // No discount applied
		assert.Equal(t, string(entity.OrderStatusPending), order.Status)

		// Verify cart is emptied
		updatedCart, _ := cartRepo.GetByUserID(1)
		assert.Len(t, updatedCart.Items, 0)

		// Verify product stock is updated
		updatedProduct, _ := productRepo.GetByID(1)
		assert.Equal(t, 98, updatedProduct.Stock)
	})

	t.Run("Create guest order successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Simple mock payment service that always succeeds
		paymentSvc := payment.NewMockPaymentService()
		// paymentSvc := &mockPaymentService{
		// 	availableProviders: []service.PaymentProvider{
		// 		{
		// 			Type:    service.PaymentProviderStripe,
		// 			Name:    "Stripe",
		// 			Enabled: true,
		// 		},
		// 	},
		// }

		// Simple mock email service
		emailSvc := &mockEmailService{}

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			CategoryID:  1,
			SellerID:    2,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: false,
		}
		productRepo.Create(product)

		// Create a test guest cart with one item
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:        1,
					ProductID: 1,
					Quantity:  2,
					CartID:    1,
				},
			},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Create order input for guest
		input := usecase.CreateOrderInput{
			SessionID:   sessionID,
			Email:       "guest@example.com",
			FullName:    "Guest User",
			PhoneNumber: "555-1234",
			ShippingAddr: entity.Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
			BillingAddr: entity.Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
		}

		// Execute
		order, err := orderUseCase.CreateOrderFromCart(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, uint(0), order.UserID) // Guest order has UserID = 0
		assert.True(t, order.IsGuestOrder)
		assert.Equal(t, input.Email, order.GuestEmail)
		assert.Equal(t, input.FullName, order.GuestFullName)
		assert.Equal(t, input.PhoneNumber, order.GuestPhone)
		assert.Equal(t, input.ShippingAddr, order.ShippingAddr)
		assert.Equal(t, input.BillingAddr, order.BillingAddr)
		assert.Len(t, order.Items, 1)
		assert.Equal(t, uint(1), order.Items[0].ProductID)
		assert.Equal(t, 2, order.Items[0].Quantity)
		assert.Equal(t, 99.99, order.Items[0].Price)
		assert.Equal(t, 99.99*2, order.Items[0].Subtotal)
		assert.Equal(t, 99.99*2, order.TotalAmount)
		assert.Equal(t, 99.99*2, order.FinalAmount) // No discount applied
		assert.Equal(t, string(entity.OrderStatusPending), order.Status)

		// Verify cart is emptied
		updatedCart, _ := cartRepo.GetBySessionID(sessionID)
		assert.Len(t, updatedCart.Items, 0)

		// Verify product stock is updated
		updatedProduct, _ := productRepo.GetByID(1)
		assert.Equal(t, 98, updatedProduct.Stock)
	})

	t.Run("Create order with empty cart", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create a test user
		user := &entity.User{
			ID:    1,
			Email: "test@example.com",
		}
		userRepo.Create(user)

		// Create an empty cart
		cart := &entity.Cart{
			ID:     1,
			UserID: 1,
			Items:  []entity.CartItem{},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Create order input
		input := usecase.CreateOrderInput{
			UserID:       1,
			ShippingAddr: entity.Address{Street: "123 Main St"},
			BillingAddr:  entity.Address{Street: "123 Main St"},
		}

		// Execute
		order, err := orderUseCase.CreateOrderFromCart(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Contains(t, err.Error(), "cart is empty")
	})

	t.Run("Create order with insufficient stock", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create a test user
		user := &entity.User{
			ID:    1,
			Email: "test@example.com",
		}
		userRepo.Create(user)

		// Create a test product with low stock
		product := &entity.Product{
			ID:     1,
			Name:   "Test Product",
			Price:  99.99,
			Stock:  5,
			Images: []string{"image1.jpg"},
		}
		productRepo.Create(product)

		// Create a cart with quantity exceeding stock
		cart := &entity.Cart{
			ID:     1,
			UserID: 1,
			Items: []entity.CartItem{
				{
					ID:        1,
					ProductID: 1,
					Quantity:  10, // More than available stock
					CartID:    1,
				},
			},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Create order input
		input := usecase.CreateOrderInput{
			UserID:       1,
			ShippingAddr: entity.Address{Street: "123 Main St"},
			BillingAddr:  entity.Address{Street: "123 Main St"},
		}

		// Execute
		order, err := orderUseCase.CreateOrderFromCart(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Contains(t, err.Error(), "insufficient stock")
	})
}

func TestOrderUseCase_ProcessPayment(t *testing.T) {
	t.Run("Process payment successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Create custom mock payment service that succeeds
		paymentSvc := &mockPaymentService{
			shouldSucceed: true,
			transactionID: "txn_12345",
			availableProviders: []service.PaymentProvider{
				{
					Type:    service.PaymentProviderStripe,
					Name:    "Stripe",
					Enabled: true,
				},
			},
		}

		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98, // No discount
			Status:      string(entity.OrderStatusPending),
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Process payment input
		input := usecase.ProcessPaymentInput{
			OrderID:         1,
			PaymentMethod:   service.PaymentMethodCreditCard,
			PaymentProvider: service.PaymentProviderStripe,
			CardDetails: &service.CardDetails{
				CardNumber:     "4242424242424242",
				ExpiryMonth:    12,
				ExpiryYear:     2030,
				CVV:            "123",
				CardholderName: "Test User",
			},
			CustomerEmail: "test@example.com",
		}

		// Execute
		updatedOrder, err := orderUseCase.ProcessPayment(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, string(entity.OrderStatusPaid), updatedOrder.Status)
		assert.Equal(t, paymentSvc.transactionID, updatedOrder.PaymentID)
		assert.Equal(t, string(service.PaymentProviderStripe), updatedOrder.PaymentProvider)
	})

	t.Run("Process payment with action required", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Create custom mock payment service that requires action
		paymentSvc := &mockPaymentService{
			shouldSucceed:  true,
			requiresAction: true,
			actionURL:      "https://example.com/3dsecure",
			transactionID:  "txn_3ds_12345",
			availableProviders: []service.PaymentProvider{
				{
					Type:    service.PaymentProviderStripe,
					Name:    "Stripe",
					Enabled: true,
				},
			},
		}

		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98, // No discount
			Status:      string(entity.OrderStatusPending),
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Process payment input
		input := usecase.ProcessPaymentInput{
			OrderID:         1,
			PaymentMethod:   service.PaymentMethodCreditCard,
			PaymentProvider: service.PaymentProviderStripe,
			CardDetails: &service.CardDetails{
				CardNumber:     "4000002500003155", // Example 3DS card
				ExpiryMonth:    12,
				ExpiryYear:     2030,
				CVV:            "123",
				CardholderName: "Test User",
			},
			CustomerEmail: "test@example.com",
		}

		// Execute
		updatedOrder, err := orderUseCase.ProcessPayment(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, string(entity.OrderStatusPendingAction), updatedOrder.Status)
		assert.Equal(t, "txn_3ds_12345", updatedOrder.PaymentID)
		assert.Equal(t, "https://example.com/3dsecure", updatedOrder.ActionURL)
	})

	t.Run("Process payment with unavailable provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Create mock payment service with only Stripe available
		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98, // No discount
			Status:      string(entity.OrderStatusPending),
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Process payment input with unsupported provider
		input := usecase.ProcessPaymentInput{
			OrderID:         1,
			PaymentMethod:   PaymentMethodCard,
			PaymentProvider: service.PaymentProviderPayPal, // Not available in mock
			CustomerEmail:   "test@example.com",
		}

		// Execute
		updatedOrder, err := orderUseCase.ProcessPayment(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedOrder)
		assert.Contains(t, err.Error(), "payment provider not available")
	})

	t.Run("Process payment for already paid order", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create a test order that's already paid
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98,
			Status:      string(entity.OrderStatusPaid),
			PaymentID:   "txn_12345",
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Process payment input
		input := usecase.ProcessPaymentInput{
			OrderID:         1,
			PaymentMethod:   PaymentMethodCard,
			PaymentProvider: service.PaymentProviderStripe,
			CustomerEmail:   "test@example.com",
		}

		// Execute
		updatedOrder, err := orderUseCase.ProcessPayment(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedOrder)
		assert.Contains(t, err.Error(), "already paid")
	})
}

func TestOrderUseCase_UpdateOrderStatus(t *testing.T) {
	t.Run("Update order status successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98,
			Status:      string(entity.OrderStatusPaid),
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Update status input
		input := usecase.UpdateOrderStatusInput{
			OrderID: 1,
			Status:  entity.OrderStatusCaptured,
		}

		// Execute
		updatedOrder, err := orderUseCase.UpdateOrderStatus(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, string(entity.OrderStatusCaptured), updatedOrder.Status)
	})

	t.Run("Update order with invalid status transition", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create a test order that's pending
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98,
			Status:      string(entity.OrderStatusPending),
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Try to update to shipped without payment
		input := usecase.UpdateOrderStatusInput{
			OrderID: 1,
			Status:  entity.OrderStatusShipped,
		}

		// Execute
		updatedOrder, err := orderUseCase.UpdateOrderStatus(input)

		// Assert - this should fail because we're trying to ship an unpaid order
		assert.Error(t, err)
		assert.Nil(t, updatedOrder)
		assert.Contains(t, err.Error(), "invalid status transition")
	})
}

func TestOrderUseCase_GetOrderByID(t *testing.T) {
	t.Run("Get existing order", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			FinalAmount: 199.98,
			Status:      string(entity.OrderStatusPaid),
			Items: []entity.OrderItem{
				{
					ProductID: 1,
					Quantity:  2,
					Price:     99.99,
					Subtotal:  199.98,
				},
			},
		}
		orderRepo.Create(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		result, err := orderUseCase.GetOrderByID(1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, order.ID, result.ID)
		assert.Equal(t, order.UserID, result.UserID)
		assert.Equal(t, order.TotalAmount, result.TotalAmount)
		assert.Equal(t, order.Status, result.Status)
		assert.Equal(t, order.Items[0].ProductID, result.Items[0].ProductID)
	})

	t.Run("Get non-existent order", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute with non-existent ID
		result, err := orderUseCase.GetOrderByID(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestOrderUseCase_GetUserOrders(t *testing.T) {
	t.Run("Get user orders successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create test orders for user 1
		order1 := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			Status:      string(entity.OrderStatusPaid),
		}
		orderRepo.Create(order1)

		order2 := &entity.Order{
			ID:          2,
			UserID:      1,
			TotalAmount: 299.97,
			Status:      string(entity.OrderStatusShipped),
		}
		orderRepo.Create(order2)

		// Create an order for a different user
		order3 := &entity.Order{
			ID:          3,
			UserID:      2,
			TotalAmount: 149.99,
			Status:      string(entity.OrderStatusPending),
		}
		orderRepo.Create(order3)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		results, err := orderUseCase.GetUserOrders(1, 0, 10)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, uint(1), results[0].ID)
		assert.Equal(t, uint(2), results[1].ID)
	})
}

func TestOrderUseCase_ListOrdersByStatus(t *testing.T) {
	t.Run("List orders by status successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}

		// Create test orders with different statuses
		order1 := &entity.Order{
			ID:          1,
			UserID:      1,
			TotalAmount: 199.98,
			Status:      string(entity.OrderStatusPending),
		}
		orderRepo.Create(order1)

		order2 := &entity.Order{
			ID:          2,
			UserID:      2,
			TotalAmount: 299.97,
			Status:      string(entity.OrderStatusPaid),
		}
		orderRepo.Create(order2)

		order3 := &entity.Order{
			ID:          3,
			UserID:      3,
			TotalAmount: 149.99,
			Status:      string(entity.OrderStatusPaid),
		}
		orderRepo.Create(order3)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		results, err := orderUseCase.ListOrdersByStatus(entity.OrderStatusPaid, 0, 10)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, string(entity.OrderStatusPaid), results[0].Status)
		assert.Equal(t, string(entity.OrderStatusPaid), results[1].Status)
	})
}

func TestOrderUseCase_CapturePayment(t *testing.T) {
	t.Run("Capture payment successfully with MobilePay", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Create a mock payment service with MobilePay support
		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order with MobilePay payment
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98,
			Status:          string(entity.OrderStatusPaid),
			PaymentID:       "mp_payment_12345",
			PaymentProvider: string(service.PaymentProviderMobilePay),
		}
		orderRepo.Create(order)

		// Add mock implementation for GetByPaymentID to the order repository
		orderRepo.AddMockGetByPaymentID(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		err := orderUseCase.CapturePayment("mp_payment_12345", 199.98)

		// Assert
		assert.NoError(t, err)
		capturedOrder, _ := orderRepo.GetByID(1)
		assert.Equal(t, string(entity.OrderStatusPaid), capturedOrder.Status) // Status remains as the capture doesn't change it
	})

	t.Run("Capture payment with invalid payment ID", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute with non-existent payment ID
		err := orderUseCase.CapturePayment("non_existent_payment", 100)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found")
	})

	t.Run("Capture payment with unsupported provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order with Stripe payment (which we don't support for capture in this test)
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98,
			Status:          string(entity.OrderStatusCaptured),
			PaymentID:       "stripe_payment_12345",
			PaymentProvider: string(service.PaymentProviderStripe),
		}
		orderRepo.Create(order)

		// Add mock implementation for GetByPaymentID to the order repository
		orderRepo.AddMockGetByPaymentID(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		err := orderUseCase.CapturePayment("stripe_payment_12345", 199.98)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment already captured")
	})
}

func TestOrderUseCase_CancelPayment(t *testing.T) {
	t.Run("Cancel payment successfully with MobilePay", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order with MobilePay payment
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98,
			Status:          string(entity.OrderStatusPendingAction),
			PaymentID:       "mp_payment_12345",
			PaymentProvider: string(service.PaymentProviderMobilePay),
		}
		orderRepo.Create(order)

		// Add mock implementation for GetByPaymentID to the order repository
		orderRepo.AddMockGetByPaymentID(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		err := orderUseCase.CancelPayment("mp_payment_12345")

		// Assert
		assert.NoError(t, err)
		cancelledOrder, _ := orderRepo.GetByID(1)
		assert.Equal(t, string(entity.OrderStatusCancelled), cancelledOrder.Status)
	})

	t.Run("Cancel payment with invalid payment ID", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute with non-existent payment ID
		err := orderUseCase.CancelPayment("non_existent_payment")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found")
	})

	t.Run("Cancel payment with unsupported provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order with Stripe payment (which we don't support for cancellation in this test)
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98,
			Status:          string(entity.OrderStatusPaid),
			PaymentID:       "stripe_payment_12345",
			PaymentProvider: string(service.PaymentProviderStripe),
		}
		orderRepo.Create(order)

		// Add mock implementation for GetByPaymentID to the order repository
		orderRepo.AddMockGetByPaymentID(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		err := orderUseCase.CancelPayment("stripe_payment_12345")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment cancellation not allowed in current order status")
	})
}

func TestOrderUseCase_RefundPayment(t *testing.T) {
	t.Run("Full refund payment successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98,
			Status:          string(entity.OrderStatusPaid),
			PaymentID:       "payment_12345",
			PaymentProvider: string(service.PaymentProviderStripe),
		}
		orderRepo.Create(order)

		// Add mock implementation for GetByPaymentID to the order repository
		orderRepo.AddMockGetByPaymentID(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute - full refund
		err := orderUseCase.RefundPayment("payment_12345", 199.98)

		// Assert
		assert.NoError(t, err)
		refundedOrder, _ := orderRepo.GetByID(1)
		assert.Equal(t, string(entity.OrderStatusRefunded), refundedOrder.Status)
	})

	t.Run("Partial refund payment successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create a test order
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98,
			Status:          string(entity.OrderStatusPaid),
			PaymentID:       "payment_12345",
			PaymentProvider: string(service.PaymentProviderStripe),
		}
		orderRepo.Create(order)

		// Add mock implementation for GetByPaymentID to the order repository
		orderRepo.AddMockGetByPaymentID(order)

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute - partial refund
		err := orderUseCase.RefundPayment("payment_12345", 50.00)

		// Assert
		assert.NoError(t, err)
		refundedOrder, _ := orderRepo.GetByID(1)
		// Status should remain paid for partial refunds
		assert.Equal(t, string(entity.OrderStatusPaid), refundedOrder.Status)
	})

	t.Run("Refund payment with invalid amount", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute with invalid amount
		err := orderUseCase.RefundPayment("payment_12345", -10.00)

		// Assert
		assert.Error(t, err)
	})

	// Test Case: Refund payment with invalid payment provider

	t.Run("Refund payment with failed payment service", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		paymentSvc := payment.NewMockPaymentService()

		emailSvc := &mockEmailService{}

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
		)

		// Execute
		err := orderUseCase.RefundPayment("payment_12345", 100.00)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found for payment ID")
	})
}

func TestOrderUseCase_NotifyShipping(t *testing.T) {
	// This function would contain tests for notifying shipping, if applicable.
}

type mockEmailService struct {
	orderConfirmationSent bool
	orderNotificationSent bool
}

func (m *mockEmailService) SendOrderConfirmation(order *entity.Order, user *entity.User) error {
	m.orderConfirmationSent = true
	return nil
}

func (m *mockEmailService) SendOrderNotification(order *entity.Order, user *entity.User) error {
	m.orderNotificationSent = true
	return nil
}

func (m *mockEmailService) SendEmail(data service.EmailData) error {
	return nil
}

type mockPaymentService struct {
	shouldSucceed      bool
	requiresAction     bool
	actionURL          string
	transactionID      string
	availableProviders []service.PaymentProvider
}

func (m *mockPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	if !m.shouldSucceed {
		return nil, errors.New("payment failed")
	}

	result := &service.PaymentResult{
		Success:        true,
		TransactionID:  m.transactionID,
		Provider:       request.PaymentProvider,
		RequiresAction: m.requiresAction,
		ActionURL:      m.actionURL,
	}

	return result, nil
}

func (m *mockPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	return m.shouldSucceed, nil
}

func (m *mockPaymentService) RefundPayment(transactionID string, amount float64, provider service.PaymentProviderType) error {
	if !m.shouldSucceed {
		return errors.New("refund failed")
	}
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

func (m *mockPaymentService) GetAvailableProviders() []service.PaymentProvider {
	return m.availableProviders
}

// CapturePayment implements the payment service interface method for capturing payments
func (m *mockPaymentService) CapturePayment(transactionID string, amount float64, provider service.PaymentProviderType) error {
	if !m.shouldSucceed {
		return errors.New("capture failed")
	}
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

// CancelPayment implements the payment service interface method for cancelling payments
func (m *mockPaymentService) CancelPayment(transactionID string, provider service.PaymentProviderType) error {
	if !m.shouldSucceed {
		return errors.New("cancel failed")
	}
	return nil
}
