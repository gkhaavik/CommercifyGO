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
	orderRepo      repository.OrderRepository
	productRepo    repository.ProductRepository
	userRepo       repository.UserRepository
	paymentSvc     service.PaymentService
	emailSvc       service.EmailService
	paymentTxnRepo repository.PaymentTransactionRepository
	currencyRepo   repository.CurrencyRepository
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	paymentSvc service.PaymentService,
	emailSvc service.EmailService,
	paymentTxnRepo repository.PaymentTransactionRepository,
	currencyRepo repository.CurrencyRepository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		userRepo:       userRepo,
		paymentSvc:     paymentSvc,
		emailSvc:       emailSvc,
		paymentTxnRepo: paymentTxnRepo,
		currencyRepo:   currencyRepo,
	}
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (uc *OrderUseCase) GetAvailablePaymentProviders() []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProviders()
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
