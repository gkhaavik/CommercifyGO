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
	t.Run("Create order successfully with shipping", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()
		shippingUseCase := mock.NewMockShippingUseCase()

		// Simple mock payment service that always succeeds
		paymentSvc := payment.NewMockPaymentService()

		// Simple mock email service
		emailSvc := &mockEmailService{}

		// Create a test user
		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Password: "hashed_password",
		}
		userRepo.Create(user)

		// Create a test product with weight
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
			Weight:      0.5, // 0.5 kg
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
			paymentTxnRepo,
			shippingUseCase,
		)

		// Create order input with shipping method
		input := usecase.CreateOrderInput{
			UserID: 1,
			ShippingAddr: entity.Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "US", // Use US for domestic shipping
			},
			BillingAddr: entity.Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "US",
			},
			ShippingMethodID: 1, // Standard Shipping
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

		// Check that shipping was applied correctly
		assert.Equal(t, input.ShippingMethodID, order.ShippingMethodID)
		assert.Equal(t, 5.99, order.ShippingCost)        // From mock shipping usecase
		assert.Equal(t, 1.0, order.TotalWeight)          // 2 items * 0.5kg = 1kg
		assert.Equal(t, 99.99*2+5.99, order.FinalAmount) // Subtotal + shipping cost

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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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

func TestOrderUseCase_GetShippingOptions(t *testing.T) {
	t.Run("Get shipping options for user cart", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()
		shippingUseCase := mock.NewMockShippingUseCase()

		// Create a test product with weight
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			Weight:      0.5, // 0.5 kg
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
			nil, // payment service not needed
			nil, // email service not needed
			paymentTxnRepo,
			shippingUseCase,
		)

		// Test address
		shippingAddr := entity.Address{
			Street:     "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "US", // US for domestic shipping
		}

		// Get shipping options
		options, err := orderUseCase.GetShippingOptions(1, "", shippingAddr)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, options)
		assert.Len(t, options.Options, 2) // We should get the 2 options from our mock

		// Check the options returned
		assert.Equal(t, uint(1), options.Options[0].ShippingMethodID) // Standard shipping
		assert.Equal(t, "Standard Shipping", options.Options[0].MethodName)
		assert.Equal(t, 5.99, options.Options[0].Cost)

		assert.Equal(t, uint(2), options.Options[1].ShippingMethodID) // Express shipping
		assert.Equal(t, "Express Shipping", options.Options[1].MethodName)
		assert.Equal(t, 15.99, options.Options[1].Cost)
	})

	t.Run("Get shipping options for guest cart", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()
		shippingUseCase := mock.NewMockShippingUseCase()

		// Create a test product with weight
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			Weight:      0.5, // 0.5 kg
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
			nil, // payment service not needed
			nil, // email service not needed
			paymentTxnRepo,
			shippingUseCase,
		)

		// Test address
		shippingAddr := entity.Address{
			Street:     "123 International St",
			City:       "Foreign City",
			State:      "FC",
			PostalCode: "12345",
			Country:    "FR", // Non-US for international shipping
		}

		// Get shipping options
		options, err := orderUseCase.GetShippingOptions(0, sessionID, shippingAddr)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, options)
		// We'll get shipping options that match this country or have wildcard "*"
		assert.GreaterOrEqual(t, len(options.Options), 1)
	})

	t.Run("Get shipping options with empty cart", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()
		shippingUseCase := mock.NewMockShippingUseCase()

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
			nil,
			nil,
			paymentTxnRepo,
			shippingUseCase,
		)

		// Test address
		shippingAddr := entity.Address{
			Street:     "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "US",
		}

		// Get shipping options
		options, err := orderUseCase.GetShippingOptions(1, "", shippingAddr)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, options)
		assert.Contains(t, err.Error(), "cart is empty")
	})
}

func TestOrderUseCase_ProcessPayment(t *testing.T) {
	t.Run("Process payment successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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

		// Verify that a payment transaction was created
		transactions, err := paymentTxnRepo.GetByOrderID(order.ID)
		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, entity.TransactionTypeAuthorize, transactions[0].Type)
		assert.Equal(t, entity.TransactionStatusSuccessful, transactions[0].Status)
		assert.Equal(t, order.FinalAmount, transactions[0].Amount)
		assert.Equal(t, paymentSvc.transactionID, transactions[0].TransactionID)
	})

	t.Run("Process payment with action required", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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

		// Verify that a pending payment transaction was created
		transactions, err := paymentTxnRepo.GetByOrderID(order.ID)
		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, entity.TransactionTypeAuthorize, transactions[0].Type)
		assert.Equal(t, entity.TransactionStatusPending, transactions[0].Status)
		assert.Equal(t, order.FinalAmount, transactions[0].Amount)
	})

	t.Run("Process payment with unavailable provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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

		// Verify that no payment transaction was created
		assert.True(t, paymentTxnRepo.IsEmpty())
	})

	t.Run("Process payment for already paid order", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := payment.NewMockPaymentService()
		emailSvc := &mockEmailService{}
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

		// Create use case with mocks
		orderUseCase := usecase.NewOrderUseCase(
			orderRepo,
			cartRepo,
			productRepo,
			userRepo,
			paymentSvc,
			emailSvc,
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute
		err := orderUseCase.CapturePayment("mp_payment_12345", 199.98)

		// Assert
		assert.NoError(t, err)
		capturedOrder, _ := orderRepo.GetByID(1)
		assert.Equal(t, string(entity.OrderStatusCaptured), capturedOrder.Status) // Status remains as the capture doesn't change it
	})

	t.Run("Capture payment with invalid payment ID", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute with non-existent payment ID
		err := orderUseCase.CapturePayment("non_existent_payment", 100)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found")

		// Verify that no transaction was created
		assert.True(t, paymentTxnRepo.IsEmpty())
	})

	t.Run("Capture payment with unsupported provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

		// Create a mock payment service that fails on capture for Stripe
		paymentSvc := &mockPaymentService{
			shouldSucceed: false, // Set this to false to simulate the failure
		}

		emailSvc := &mockEmailService{}

		// Create a test order with Stripe payment
		order := &entity.Order{
			ID:              1,
			UserID:          1,
			TotalAmount:     199.98,
			FinalAmount:     199.98, // No discount
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
			paymentTxnRepo,
		)

		// Execute
		err := orderUseCase.CapturePayment("stripe_payment_12345", 199.98)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to capture payment")
	})
}

func TestOrderUseCase_CancelPayment(t *testing.T) {
	t.Run("Cancel payment successfully with MobilePay", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute
		err := orderUseCase.CancelPayment("mp_payment_12345")

		// Assert
		assert.NoError(t, err)
		cancelledOrder, _ := orderRepo.GetByID(1)
		assert.Equal(t, string(entity.OrderStatusCancelled), cancelledOrder.Status)

		// Verify that a cancel transaction was created
		transactions, err := paymentTxnRepo.GetByOrderID(order.ID)
		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, entity.TransactionTypeCancel, transactions[0].Type)
		assert.Equal(t, entity.TransactionStatusSuccessful, transactions[0].Status)
		assert.Equal(t, float64(0), transactions[0].Amount) // Cancel transactions have amount 0
	})

	t.Run("Cancel payment with invalid payment ID", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute with non-existent payment ID
		err := orderUseCase.CancelPayment("non_existent_payment")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found")

		// Verify that no transaction was created
		assert.True(t, paymentTxnRepo.IsEmpty())
	})

	t.Run("Cancel payment with unsupported provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
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
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute - full refund
		err := orderUseCase.RefundPayment("payment_12345", 199.98)

		// Assert
		assert.NoError(t, err)
		refundedOrder, _ := orderRepo.GetByID(1)
		assert.Equal(t, string(entity.OrderStatusRefunded), refundedOrder.Status)

		// Verify that a refund transaction was created
		transactions, err := paymentTxnRepo.GetByOrderID(order.ID)
		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, entity.TransactionTypeRefund, transactions[0].Type)
		assert.Equal(t, entity.TransactionStatusSuccessful, transactions[0].Status)
		assert.Equal(t, order.FinalAmount, transactions[0].Amount)
		// Check metadata
		assert.Equal(t, "true", transactions[0].Metadata["full_refund"])
	})

	t.Run("Partial refund payment successfully", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute - partial refund
		err := orderUseCase.RefundPayment("payment_12345", 50.00)

		// Assert
		assert.NoError(t, err)
		refundedOrder, _ := orderRepo.GetByID(1)
		// Status should remain paid for partial refunds
		assert.Equal(t, string(entity.OrderStatusPaid), refundedOrder.Status)

		// Verify that a refund transaction was created
		transactions, err := paymentTxnRepo.GetByOrderID(order.ID)
		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, entity.TransactionTypeRefund, transactions[0].Type)
		assert.Equal(t, entity.TransactionStatusSuccessful, transactions[0].Status)
		assert.Equal(t, 50.0, transactions[0].Amount)
		// Check metadata
		assert.Equal(t, "false", transactions[0].Metadata["full_refund"])
	})

	t.Run("Refund payment with invalid amount", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute with invalid amount
		err := orderUseCase.RefundPayment("payment_12345", -10.00)

		// Assert
		assert.Error(t, err)

		// Verify that no transaction was created
		assert.True(t, paymentTxnRepo.IsEmpty())
	})

	t.Run("Refund payment with failed payment service", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentTxnRepo := mock.NewMockPaymentTransactionRepository()

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
			paymentTxnRepo,
		)

		// Execute
		err := orderUseCase.RefundPayment("payment_12345", 100.00)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found for payment ID")

		// Verify that no transaction was created
		assert.True(t, paymentTxnRepo.IsEmpty())
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
