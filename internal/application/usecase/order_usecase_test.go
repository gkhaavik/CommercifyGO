package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
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
		paymentSvc := &mockPaymentService{
			availableProviders: []service.PaymentProvider{
				{
					Type:    service.PaymentProviderStripe,
					Name:    "Stripe",
					Enabled: true,
				},
			},
		}

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

	t.Run("Create order with empty cart", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := &mockPaymentService{}
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
		paymentSvc := &mockPaymentService{}
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

		// Create mock payment service that succeeds
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
			PaymentMethod:   PaymentMethodCard,
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

		// Create mock payment service that requires action
		paymentSvc := &mockPaymentService{
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
			PaymentMethod:   PaymentMethodCard,
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
		assert.Equal(t, paymentSvc.transactionID, updatedOrder.PaymentID)
		assert.Equal(t, paymentSvc.actionURL, updatedOrder.ActionURL)
	})

	t.Run("Process payment with unavailable provider", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()

		// Create mock payment service with only Stripe available
		paymentSvc := &mockPaymentService{
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
		paymentSvc := &mockPaymentService{}
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
		paymentSvc := &mockPaymentService{}
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
			Status:  entity.OrderStatusShipped,
		}

		// Execute
		updatedOrder, err := orderUseCase.UpdateOrderStatus(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, string(entity.OrderStatusShipped), updatedOrder.Status)
	})

	t.Run("Update order with invalid status transition", func(t *testing.T) {
		// Setup mocks
		orderRepo := mock.NewMockOrderRepository()
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()
		userRepo := mock.NewMockUserRepository()
		paymentSvc := &mockPaymentService{}
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
		paymentSvc := &mockPaymentService{}
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
		paymentSvc := &mockPaymentService{}
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
		paymentSvc := &mockPaymentService{}
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
		paymentSvc := &mockPaymentService{}
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

// Mock implementations for payment and email services

type mockPaymentService struct {
	shouldSucceed      bool
	requiresAction     bool
	actionURL          string
	transactionID      string
	availableProviders []service.PaymentProvider
}

func (m *mockPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	result := &service.PaymentResult{
		Success:        m.shouldSucceed,
		TransactionID:  m.transactionID,
		Provider:       request.PaymentProvider,
		RequiresAction: m.requiresAction,
		ActionURL:      m.actionURL,
	}

	if !m.shouldSucceed && !m.requiresAction {
		result.ErrorMessage = "Payment failed"
	}

	return result, nil
}

func (m *mockPaymentService) RefundPayment(transactionID string, amount float64, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return fmt.Errorf("transaction ID is required")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	// Simulate refund
	if m.shouldSucceed {
		return nil
	}
	return fmt.Errorf("refund failed")
}

func (m *mockPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	if transactionID == "" {
		return false, fmt.Errorf("transaction ID is required")
	}

	// Simulate verification
	if m.shouldSucceed {
		return true, nil
	}
	return false, fmt.Errorf("verification failed")
}

func (m *mockPaymentService) GetAvailableProviders() []service.PaymentProvider {
	return m.availableProviders
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
