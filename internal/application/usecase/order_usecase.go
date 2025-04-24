package usecase

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
)

// OrderUseCase implements order-related use cases
type OrderUseCase struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	paymentSvc  service.PaymentService
	emailSvc    service.EmailService
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	paymentSvc service.PaymentService,
	emailSvc service.EmailService,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		paymentSvc:  paymentSvc,
		emailSvc:    emailSvc,
	}
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (uc *OrderUseCase) GetAvailablePaymentProviders() []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProviders()
}

// CreateOrderInput contains the data needed to create an order
type CreateOrderInput struct {
	UserID       uint           `json:"user_id"`
	ShippingAddr entity.Address `json:"shipping_address"`
	BillingAddr  entity.Address `json:"billing_address"`
}

// CreateOrderFromCart creates an order from a user's cart
func (uc *OrderUseCase) CreateOrderFromCart(input CreateOrderInput) (*entity.Order, error) {
	// Get user's cart
	cart, err := uc.cartRepo.GetByUserID(input.UserID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Get user for email notifications
	user, err := uc.userRepo.GetByID(input.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Convert cart items to order items
	orderItems := make([]entity.OrderItem, 0, len(cart.Items))
	for _, cartItem := range cart.Items {
		// Get product to get current price
		product, err := uc.productRepo.GetByID(cartItem.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		// Check stock availability
		if !product.IsAvailable(cartItem.Quantity) {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Create order item
		orderItems = append(orderItems, entity.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
			Subtotal:  float64(cartItem.Quantity) * product.Price,
		})

		// Update product stock
		if err := product.UpdateStock(-cartItem.Quantity); err != nil {
			return nil, err
		}
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	// Create order
	order, err := entity.NewOrder(input.UserID, orderItems, input.ShippingAddr, input.BillingAddr)
	if err != nil {
		return nil, err
	}

	// Save order
	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Clear cart after successful order creation
	cart.Clear()
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	// Send order confirmation email to customer
	if uc.emailSvc != nil {
		go uc.emailSvc.SendOrderConfirmation(order, user)
	}

	// Send order notification email to admin
	if uc.emailSvc != nil {
		go uc.emailSvc.SendOrderNotification(order, user)
	}

	return order, nil
}

// ProcessPaymentInput contains the data needed to process a payment
type ProcessPaymentInput struct {
	OrderID         uint                        `json:"order_id"`
	PaymentMethod   service.PaymentMethod       `json:"payment_method"`
	PaymentProvider service.PaymentProviderType `json:"payment_provider"`
	CardDetails     *service.CardDetails        `json:"card_details,omitempty"`
	PayPalDetails   *service.PayPalDetails      `json:"paypal_details,omitempty"`
	BankDetails     *service.BankDetails        `json:"bank_details,omitempty"`
	CustomerEmail   string                      `json:"customer_email,omitempty"`
	PhoneNumber     string                      `json:"phone_number,omitempty"`
}

// ProcessPayment processes payment for an order
func (uc *OrderUseCase) ProcessPayment(input ProcessPaymentInput) (*entity.Order, error) {
	// Get order
	order, err := uc.orderRepo.GetByID(input.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if order is already paid
	if order.Status == string(entity.OrderStatusPaid) ||
		order.Status == string(entity.OrderStatusShipped) ||
		order.Status == string(entity.OrderStatusDelivered) {
		return nil, errors.New("order is already paid")
	}

	// Validate payment provider
	availableProviders := uc.GetAvailablePaymentProviders()
	providerValid := false
	for _, p := range availableProviders {
		if p.Type == input.PaymentProvider && p.Enabled {
			providerValid = true
			break
		}
	}
	if !providerValid {
		return nil, errors.New("payment provider not available")
	}

	// Process payment
	paymentResult, err := uc.paymentSvc.ProcessPayment(service.PaymentRequest{
		OrderID:         order.ID,
		Amount:          order.TotalAmount,
		Currency:        "USD",
		PaymentMethod:   input.PaymentMethod,
		PaymentProvider: input.PaymentProvider,
		CardDetails:     input.CardDetails,
		PayPalDetails:   input.PayPalDetails,
		BankDetails:     input.BankDetails,
		CustomerEmail:   input.CustomerEmail,
		PhoneNumber:     input.PhoneNumber,
	})

	if err != nil {
		return nil, err
	}

	// Update order with payment ID, provider, and status
	if err := order.SetPaymentID(paymentResult.TransactionID); err != nil {
		return nil, err
	}
	if err := order.SetPaymentProvider(string(paymentResult.Provider)); err != nil {
		return nil, err
	}
	if err := order.UpdateStatus(entity.OrderStatusPending); err != nil {
		return nil, err
	}

	// Update order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatusInput contains the data needed to update an order status
type UpdateOrderStatusInput struct {
	OrderID uint               `json:"order_id"`
	Status  entity.OrderStatus `json:"status"`
}

// UpdateOrderStatus updates the status of an order
func (uc *OrderUseCase) UpdateOrderStatus(input UpdateOrderStatusInput) (*entity.Order, error) {
	// Get order
	order, err := uc.orderRepo.GetByID(input.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Update status
	if err := order.UpdateStatus(input.Status); err != nil {
		return nil, err
	}

	// Update order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (uc *OrderUseCase) GetOrderByID(id uint) (*entity.Order, error) {
	return uc.orderRepo.GetByID(id)
}

// GetUserOrders retrieves orders for a user
func (uc *OrderUseCase) GetUserOrders(userID uint, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.GetByUser(userID, offset, limit)
}

func (uc *OrderUseCase) ListOrdersByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.ListByStatus(status, offset, limit)
}
