package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
)

// OrderUseCase implements order-related use cases
type OrderUseCase struct {
	orderRepo       repository.OrderRepository
	cartRepo        repository.CartRepository
	productRepo     repository.ProductRepository
	userRepo        repository.UserRepository
	paymentSvc      service.PaymentService
	emailSvc        service.EmailService
	paymentTxnRepo  repository.PaymentTransactionRepository
	shippingUseCase *ShippingUseCase
	currencyRepo    repository.CurrencyRepository
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	paymentSvc service.PaymentService,
	emailSvc service.EmailService,
	paymentTxnRepo repository.PaymentTransactionRepository,
	shippingUseCase *ShippingUseCase,
	currencyRepo repository.CurrencyRepository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:       orderRepo,
		cartRepo:        cartRepo,
		productRepo:     productRepo,
		userRepo:        userRepo,
		paymentSvc:      paymentSvc,
		emailSvc:        emailSvc,
		paymentTxnRepo:  paymentTxnRepo,
		shippingUseCase: shippingUseCase,
		currencyRepo:    currencyRepo,
	}
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (uc *OrderUseCase) GetAvailablePaymentProviders() []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProviders()
}

// CreateOrderInput contains the data needed to create an order
type CreateOrderInput struct {
	UserID           uint           `json:"user_id"`
	SessionID        string         `json:"session_id"`
	ShippingAddr     entity.Address `json:"shipping_address"`
	BillingAddr      entity.Address `json:"billing_address"`
	Email            string         `json:"email,omitempty"`
	PhoneNumber      string         `json:"phone_number,omitempty"`
	FullName         string         `json:"full_name,omitempty"`
	ShippingMethodID uint           `json:"shipping_method_id"`
}

// CreateOrderFromCart creates an order from a user's cart
func (uc *OrderUseCase) CreateOrderFromCart(input CreateOrderInput) (*entity.Order, error) {
	// Check if this is a guest checkout or a user checkout
	if input.UserID > 0 {
		// Authenticated user checkout
		return uc.createOrderFromUserCart(input)
	} else if input.SessionID != "" {
		// Guest checkout
		return uc.createOrderFromGuestCart(input)
	}

	return nil, errors.New("either user ID or session ID must be provided")
}

// createOrderFromUserCart creates an order from an authenticated user's cart
func (uc *OrderUseCase) createOrderFromUserCart(input CreateOrderInput) (*entity.Order, error) {
	// Validate shipping method ID
	if input.ShippingMethodID == 0 {
		return nil, errors.New("shipping method ID is required")
	}

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
	totalWeight := 0.0

	for _, cartItem := range cart.Items {
		// Get product to get current price
		product, err := uc.productRepo.GetByID(cartItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: ProductID=%d", cartItem.ProductID)
		}

		// Check stock availability
		if !product.IsAvailable(cartItem.Quantity) {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Create order item with weight
		orderItem := entity.OrderItem{
			ProductID:   cartItem.ProductID,
			Quantity:    cartItem.Quantity,
			Price:       product.Price,
			Subtotal:    int64(cartItem.Quantity) * product.Price,
			Weight:      product.Weight,
			ProductName: product.Name,
		}

		// TODO: Check for variant and assign variant ID
		// If this is a variant, store the variant ID
		variant := product.GetVariantByID(cartItem.ProductVariantID)
		if variant != nil {
			orderItem.SKU = variant.SKU
		}

		orderItems = append(orderItems, orderItem)
		totalWeight += product.Weight * float64(cartItem.Quantity)

		// Update product stock
		if err := product.UpdateStock(-cartItem.Quantity); err != nil {
			return nil, err
		}
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	// Create order
	order, err := entity.NewOrder(input.UserID, orderItems, input.ShippingAddr, input.BillingAddr, entity.CustomerDetails{
		Email:    input.Email,
		Phone:    input.PhoneNumber,
		FullName: input.FullName,
	})
	if err != nil {
		return nil, err
	}

	// Set the total weight
	order.TotalWeight = totalWeight

	// Set shipping method ID
	order.ShippingMethodID = input.ShippingMethodID

	// Apply shipping method if specified
	if uc.shippingUseCase != nil {
		shippingMethod, err := uc.shippingUseCase.GetShippingMethodByID(input.ShippingMethodID)
		if err != nil {
			return nil, errors.New("shipping method not found")
		}

		// Calculate shipping cost
		shippingCost, err := uc.shippingUseCase.GetShippingCost(input.ShippingMethodID, order.TotalAmount, order.TotalWeight)
		if err != nil {
			return nil, fmt.Errorf("error calculating shipping cost: %v", err)
		}

		// Apply shipping method and cost to order
		if err := order.SetShippingMethod(shippingMethod, shippingCost); err != nil {
			return nil, err
		}
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

// createOrderFromGuestCart creates an order from a guest's cart
func (uc *OrderUseCase) createOrderFromGuestCart(input CreateOrderInput) (*entity.Order, error) {
	// Validate and sanitize guest information
	input.Email = sanitizeEmail(input.Email)
	if input.Email == "" {
		return nil, errors.New("valid email is required for guest checkout")
	}

	input.FullName = sanitizeString(input.FullName)
	if input.FullName == "" {
		return nil, errors.New("valid full name is required for guest checkout")
	}

	// Validate shipping method ID
	if input.ShippingMethodID == 0 {
		return nil, errors.New("shipping method ID is required")
	}

	// Get guest's cart
	cart, err := uc.cartRepo.GetBySessionID(input.SessionID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Convert cart items to order items
	orderItems := make([]entity.OrderItem, 0, len(cart.Items))
	totalWeight := 0.0

	for _, cartItem := range cart.Items {
		// Get product to get current price
		product, err := uc.productRepo.GetByID(cartItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: ProductID=%d", cartItem.ProductID)
		}

		// Check stock availability
		if !product.IsAvailable(cartItem.Quantity) {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Calculate item weight
		itemWeight := product.Weight

		// Create order item with weight
		orderItem := entity.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
			Subtotal:  int64(cartItem.Quantity) * product.Price,
			Weight:    itemWeight,
		}

		// If this is a variant, store the variant ID
		orderItem.ProductID = cartItem.ProductID

		orderItems = append(orderItems, orderItem)
		totalWeight += itemWeight * float64(cartItem.Quantity)

		// Update product stock
		if err := product.UpdateStock(-cartItem.Quantity); err != nil {
			return nil, err
		}
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	// Create guest order (0 as UserID indicates a guest order)
	order, err := entity.NewGuestOrder(orderItems, input.ShippingAddr, input.BillingAddr, entity.CustomerDetails{
		Email:    input.Email,
		Phone:    input.PhoneNumber,
		FullName: input.FullName,
	})
	if err != nil {
		return nil, err
	}

	// Set the total weight
	order.TotalWeight = totalWeight

	// Set shipping method ID
	order.ShippingMethodID = input.ShippingMethodID

	// Apply shipping method if specified
	if uc.shippingUseCase != nil {
		shippingMethod, err := uc.shippingUseCase.GetShippingMethodByID(input.ShippingMethodID)
		if err != nil {
			return nil, errors.New("shipping method not found")
		}

		// Calculate shipping cost
		shippingCost, err := uc.shippingUseCase.GetShippingCost(input.ShippingMethodID, order.TotalAmount, order.TotalWeight)
		if err != nil {
			return nil, fmt.Errorf("error calculating shipping cost: %v", err)
		}

		// Apply shipping method and cost to order
		if err := order.SetShippingMethod(shippingMethod, shippingCost); err != nil {
			return nil, err
		}
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

	// Send order confirmation email to guest
	if uc.emailSvc != nil {
		// Create a temporary user object for the email
		guestUser := &entity.User{
			Email:     input.Email,
			FirstName: input.FullName,
		}
		go uc.emailSvc.SendOrderConfirmation(order, guestUser)
	}

	// Send order notification email to admin
	if uc.emailSvc != nil {
		// Create a temporary user object for the email
		guestUser := &entity.User{
			Email:     input.Email,
			FirstName: input.FullName,
		}
		go uc.emailSvc.SendOrderNotification(order, guestUser)
	}

	return order, nil
}

// ProcessPaymentInput contains the data needed to process a payment
type ProcessPaymentInput struct {
	OrderID         uint
	PaymentMethod   service.PaymentMethod
	PaymentProvider service.PaymentProviderType
	CardDetails     *service.CardDetails
	PayPalDetails   *service.PayPalDetails
	BankDetails     *service.BankDetails
	CustomerEmail   string
	PhoneNumber     string
}

// ProcessPayment processes payment for an order
func (uc *OrderUseCase) ProcessPayment(input ProcessPaymentInput) (*entity.Order, error) {
	// Get order
	order, err := uc.orderRepo.GetByID(input.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if order is already paid
	if order.Status == entity.OrderStatusPaid ||
		order.Status == entity.OrderStatusShipped ||
		order.Status == entity.OrderStatusDelivered {
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

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to get default currency: %w", err)
	}

	// Process payment
	paymentResult, err := uc.paymentSvc.ProcessPayment(service.PaymentRequest{
		OrderID:         order.ID,
		Amount:          order.FinalAmount, // Use final amount (after discounts)
		Currency:        defaultCurrency.Code,
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

	// Handle payment results that require additional action (like redirects)
	if paymentResult.RequiresAction && paymentResult.ActionURL != "" {
		// Update order with payment ID, provider, and status
		if err := order.SetPaymentID(paymentResult.TransactionID); err != nil {
			return nil, err
		}
		if err := order.SetPaymentProvider(string(paymentResult.Provider)); err != nil {
			return nil, err
		}
		if err := order.SetActionURL(paymentResult.ActionURL); err != nil {
			return nil, err
		}
		if err := order.UpdateStatus(entity.OrderStatusPendingAction); err != nil {
			return nil, err
		}

		// Update order in repository
		if err := uc.orderRepo.Update(order); err != nil {
			return nil, err
		}

		// Record the pending authorization transaction
		txn, err := entity.NewPaymentTransaction(
			order.ID,
			paymentResult.TransactionID,
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			order.FinalAmount,
			defaultCurrency.Code,
			string(paymentResult.Provider),
		)
		if err != nil {
			// Log the error but don't fail the payment process
			log.Printf("Failed to create payment transaction record: %v", err)
		} else {
			// Add metadata
			txn.AddMetadata("payment_method", string(input.PaymentMethod))
			txn.AddMetadata("requires_action", "true")
			txn.AddMetadata("action_url", paymentResult.ActionURL)

			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				// Log error but don't fail the payment process
				log.Printf("Failed to save payment transaction: %v\n", err)
			}
		}

		return order, nil
	}

	if !paymentResult.Success {
		// Record the failed transaction
		txn, err := entity.NewPaymentTransaction(
			order.ID,
			paymentResult.TransactionID,
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusFailed,
			order.FinalAmount,
			defaultCurrency.Code,
			string(paymentResult.Provider),
		)
		if err == nil {
			txn.AddMetadata("payment_method", string(input.PaymentMethod))
			txn.AddMetadata("error_message", paymentResult.ErrorMessage)

			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				// Log error but don't fail the process
				log.Printf("Failed to save failed payment transaction: %v\n", err)
			}
		}

		return nil, errors.New(paymentResult.ErrorMessage)
	}

	// Update order with payment ID, provider, and status
	if err := order.SetPaymentID(paymentResult.TransactionID); err != nil {
		return nil, err
	}
	if err := order.SetPaymentProvider(string(paymentResult.Provider)); err != nil {
		return nil, err
	}
	if err := order.SetPaymentMethod(string(input.PaymentMethod)); err != nil {
		return nil, err
	}
	if err := order.UpdateStatus(entity.OrderStatusPaid); err != nil {
		return nil, err
	}

	// Update order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return nil, err
	}

	// Record the successful authorization transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		paymentResult.TransactionID,
		entity.TransactionTypeAuthorize,
		entity.TransactionStatusSuccessful,
		order.FinalAmount,
		defaultCurrency.Code,
		string(paymentResult.Provider),
	)
	if err == nil {
		txn.AddMetadata("payment_method", string(input.PaymentMethod))

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save payment transaction: %v\n", err)
		}
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
	if id == 0 {
		return nil, errors.New("order ID cannot be 0")
	}

	order, err := uc.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	return order, nil
}

// GetOrderByPaymentID retrieves an order by its payment ID
func (uc *OrderUseCase) GetOrderByPaymentID(paymentID string) (*entity.Order, error) {
	if paymentID == "" {
		return nil, errors.New("payment ID cannot be empty")
	}

	// Delegate to the order repository which has this functionality
	order, err := uc.orderRepo.GetByPaymentID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by payment ID: %w", err)
	}

	return order, nil
}

// GetUserOrders retrieves orders for a user
func (uc *OrderUseCase) GetUserOrders(userID uint, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.GetByUser(userID, offset, limit)
}

func (uc *OrderUseCase) ListOrdersByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.ListByStatus(status, offset, limit)
}

// CapturePayment captures an authorized payment
func (uc *OrderUseCase) CapturePayment(transactionID string, amount int64) error {
	// Find the order with this payment ID
	order, err := uc.orderRepo.GetByPaymentID(transactionID)
	if err != nil {
		return errors.New("order not found for payment ID")
	}

	// Check if the order is already captured
	if order.Status == entity.OrderStatusCaptured {
		return errors.New("payment already captured")
	}
	// Check if the order is in a state that allows capture
	if order.Status != entity.OrderStatusPaid {
		return errors.New("payment capture not allowed in current order status")
	}

	// Check if the amount is valid
	if amount <= 0 {
		return errors.New("capture amount must be greater than zero")
	}

	// Check if amount is greater than the order amount
	if amount > order.FinalAmount {
		return errors.New("capture amount cannot exceed the original payment amount")
	}

	providerType := service.PaymentProviderType(order.PaymentProvider)

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	// Call payment service to capture payment
	err = uc.paymentSvc.CapturePayment(transactionID, amount, providerType)
	if err != nil {
		// Record failed capture attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			entity.TransactionTypeCapture,
			entity.TransactionStatusFailed,
			amount,
			defaultCurrency.Code,
			string(providerType),
		)

		if txErr == nil {
			txn.AddMetadata("error", err.Error())
			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				log.Printf("Failed to save capture transaction: %v\n", err)
			}
		}

		return fmt.Errorf("failed to capture payment: %v", err)
	}

	if err := order.UpdateStatus(entity.OrderStatusCaptured); err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// Save the updated order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save order status: %v", err)
	}

	// Record successful capture transaction
	// Track if this is a full or partial capture
	isFullCapture := amount >= order.FinalAmount

	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		entity.TransactionTypeCapture,
		entity.TransactionStatusSuccessful,
		amount,
		defaultCurrency.Code,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("full_capture", fmt.Sprintf("%t", isFullCapture))

		// Record total authorized amount
		if isFullCapture {
			txn.AddMetadata("remaining_amount", "0")
		} else {
			remainingAmount := order.FinalAmount - amount
			txn.AddMetadata("remaining_amount", fmt.Sprintf("%.2f", money.FromCents(remainingAmount)))
		}

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save capture transaction: %v\n", err)
		}
	}

	return nil
}

// CancelPayment cancels a payment
func (uc *OrderUseCase) CancelPayment(transactionID string) error {
	// Find the order with this payment ID
	order, err := uc.orderRepo.GetByPaymentID(transactionID)
	if err != nil {
		return errors.New("order not found for payment ID")
	}

	// Check if the order is already canceled
	if order.Status == entity.OrderStatusCancelled {
		return errors.New("payment already canceled")
	}
	// Check if the order is in a state that allows cancellation
	if order.Status != entity.OrderStatusPendingAction {
		return errors.New("payment cancellation not allowed in current order status")
	}
	// Check if the transaction ID is valid
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	providerType := service.PaymentProviderType(order.PaymentProvider)

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	err = uc.paymentSvc.CancelPayment(transactionID, providerType)
	if err != nil {
		// Record failed cancellation attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			entity.TransactionTypeCancel,
			entity.TransactionStatusFailed,
			0, // No amount for cancellation
			defaultCurrency.Code,
			string(providerType),
		)
		if txErr == nil {
			txn.AddMetadata("error", err.Error())
			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				log.Printf("Failed to save cancel transaction: %v\n", err)
			}
		}

		return fmt.Errorf("failed to cancel payment: %v", err)
	}

	// Update the order status to cancelled after successful payment cancellation
	if err := order.UpdateStatus(entity.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// Save the updated order in the repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save order status: %v", err)
	}

	// Record successful cancellation transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		entity.TransactionTypeCancel,
		entity.TransactionStatusSuccessful,
		0, // No amount for cancellation
		defaultCurrency.Code,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("previous_status", string(entity.OrderStatusPendingAction))

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save cancel transaction: %v\n", err)
		}
	}

	return nil
}

// RefundPayment refunds a payment
func (uc *OrderUseCase) RefundPayment(transactionID string, amount int64) error {
	// Find the order with this payment ID
	order, err := uc.orderRepo.GetByPaymentID(transactionID)
	if err != nil {
		return errors.New("order not found for payment ID")
	}

	// Check if the order is already refunded
	if order.Status == entity.OrderStatusRefunded {
		return errors.New("payment already refunded")
	}
	// Check if the order is in a state that allows refund
	if order.Status != entity.OrderStatusPaid && order.Status != entity.OrderStatusCaptured {
		return errors.New("payment refund not allowed in current order status")
	}
	// Check if the amount is valid
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Check if the refund amount exceeds the original amount
	if amount > order.FinalAmount {
		return errors.New("refund amount cannot exceed the original payment amount")
	}

	providerType := service.PaymentProviderType(order.PaymentProvider)

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	// Get total refunded amount so far (if any)
	var totalRefundedSoFar int64 = 0
	totalRefundedSoFar, _ = uc.paymentTxnRepo.SumAmountByOrderIDAndType(order.ID, entity.TransactionTypeRefund)

	// Check if we're trying to refund more than the original amount when combining with previous refunds
	if totalRefundedSoFar+amount > order.FinalAmount {
		return errors.New("total refund amount would exceed the original payment amount")
	}

	err = uc.paymentSvc.RefundPayment(transactionID, amount, providerType)
	if err != nil {
		// Record failed refund attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			entity.TransactionTypeRefund,
			entity.TransactionStatusFailed,
			amount,
			defaultCurrency.Code,
			string(providerType),
		)
		if txErr == nil {
			txn.AddMetadata("error", err.Error())
			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				log.Printf("Failed to save refund transaction: %v\n", err)
			}
		}

		return fmt.Errorf("failed to refund payment: %v", err)
	}

	// Calculate if this is a full refund
	isFullRefund := false
	if amount >= order.FinalAmount || (totalRefundedSoFar+amount) >= order.FinalAmount {
		isFullRefund = true
	}

	// Only update the order status to refunded if it's a full refund
	if isFullRefund {
		if err := order.UpdateStatus(entity.OrderStatusRefunded); err != nil {
			return fmt.Errorf("failed to update order status: %v", err)
		}

		// Save the updated order in the repository
		if err := uc.orderRepo.Update(order); err != nil {
			return fmt.Errorf("failed to save order status: %v", err)
		}
	}

	// Record successful refund transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		entity.TransactionTypeRefund,
		entity.TransactionStatusSuccessful,
		amount,
		defaultCurrency.Code,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("full_refund", fmt.Sprintf("%t", isFullRefund))
		txn.AddMetadata("previous_status", string(order.Status))

		// Record total refunded amount including this transaction
		totalRefunded := totalRefundedSoFar + amount
		txn.AddMetadata("total_refunded", fmt.Sprintf("%.2f", money.FromCents(totalRefunded)))

		// Record remaining amount still available for refund
		remainingAmount := max(order.FinalAmount-totalRefunded, 0)
		txn.AddMetadata("remaining_available", fmt.Sprintf("%.2f", money.FromCents(remainingAmount)))

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save refund transaction: %v\n", err)
		}
	}

	return nil
}

// GetShippingOptions calculates available shipping options for an order based on the cart
func (uc *OrderUseCase) GetShippingOptions(userID uint, sessionID string, shippingAddr entity.Address) (*ShippingOptions, error) {
	var cart *entity.Cart
	var err error

	// Get the appropriate cart
	if userID > 0 {
		cart, err = uc.cartRepo.GetByUserID(userID)
	} else if sessionID != "" {
		cart, err = uc.cartRepo.GetBySessionID(sessionID)
	} else {
		return nil, errors.New("either user ID or session ID must be provided")
	}

	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Calculate cart's total value and weight
	var totalValue int64
	var totalWeight float64

	for _, item := range cart.Items {
		product, err := uc.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: ProductID=%d", item.ProductID)
		}

		totalValue += int64(item.Quantity) * product.Price

		// Calculate weight based on product or product variant
		totalWeight += float64(item.Quantity) * product.Weight
	}

	// Call shipping use case to calculate options
	if uc.shippingUseCase == nil {
		return nil, errors.New("shipping use case not initialized")
	}

	return uc.shippingUseCase.CalculateShippingOptions(shippingAddr, totalValue, totalWeight)
}

// RecordPaymentTransaction records a payment transaction for an order
func (uc *OrderUseCase) RecordPaymentTransaction(transaction *entity.PaymentTransaction) error {
	if transaction == nil {
		return errors.New("payment transaction cannot be nil")
	}

	// Validate the order exists
	_, err := uc.orderRepo.GetByID(transaction.OrderID)
	if err != nil {
		return fmt.Errorf("failed to verify order existence: %w", err)
	}

	// Create transaction record
	return uc.paymentTxnRepo.Create(transaction)
}

func (uc *OrderUseCase) UpdatePaymentTransaction(transactionID string, status entity.TransactionStatus, metadata map[string]string) error {
	txn, err := uc.paymentTxnRepo.GetByTransactionID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get payment transaction: %w", err)
	}

	txn.UpdateStatus(status)

	for key, value := range metadata {
		txn.AddMetadata(key, value)
	}

	return uc.paymentTxnRepo.Update(txn)
}

// ForceApproveMobilePayPayment force approves a MobilePay payment
func (uc *OrderUseCase) ForceApproveMobilePayPayment(paymentID string, phoneNumber string) error {
	// Get the payment service
	paymentSvc, ok := uc.paymentSvc.(*payment.MultiProviderPaymentService)
	if !ok {
		return errors.New("invalid payment service")
	}

	// Force approve the payment
	return paymentSvc.ForceApprovePayment(paymentID, phoneNumber, service.PaymentProviderMobilePay)
}

// GetUserByID retrieves a user by ID
func (uc *OrderUseCase) GetUserByID(id uint) (*entity.User, error) {
	return uc.userRepo.GetByID(id)
}

// ListAllOrders lists all orders
func (uc *OrderUseCase) ListAllOrders(offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.ListAll(offset, limit)
}
