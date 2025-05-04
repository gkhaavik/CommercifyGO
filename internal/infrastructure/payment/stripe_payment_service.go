package payment

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// StripePaymentService implements a Stripe payment service
type StripePaymentService struct {
	config config.StripeConfig
	logger logger.Logger
}

// NewStripePaymentService creates a new StripePaymentService
func NewStripePaymentService(config config.StripeConfig, logger logger.Logger) *StripePaymentService {
	// Initialize Stripe with the API key
	stripe.Key = config.SecretKey

	return &StripePaymentService{
		config: config,
		logger: logger,
	}
}

// GetAvailableProviders returns a list of available payment providers
func (s *StripePaymentService) GetAvailableProviders() []service.PaymentProvider {
	return []service.PaymentProvider{
		{
			Type:        service.PaymentProviderStripe,
			Name:        "Stripe",
			Description: "Pay with credit or debit card",
			IconURL:     "/assets/images/stripe-logo.png",
			Methods:     []service.PaymentMethod{service.PaymentMethodCreditCard},
			Enabled:     true,
		},
	}
}

// createPaymentMethodFromCard creates a payment method from card details
func (s *StripePaymentService) createPaymentMethodFromCard(cardDetails *service.CardDetails) (string, error) {
	if cardDetails == nil {
		return "", errors.New("card details are required")
	}

	// If a token was provided, use it directly
	if cardDetails.Token != "" {
		return cardDetails.Token, nil
	}

	// Otherwise create a payment method from the card details
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(cardDetails.CardNumber),
			ExpMonth: stripe.String(string(cardDetails.ExpiryMonth)),
			ExpYear:  stripe.String(string(cardDetails.ExpiryYear)),
			CVC:      stripe.String(cardDetails.CVV),
		},
		Type: stripe.String("card"),
	}

	if cardDetails.CardholderName != "" {
		params.BillingDetails = &stripe.BillingDetailsParams{
			Name: stripe.String(cardDetails.CardholderName),
		}
	}

	// Create the payment method
	pm, err := paymentmethod.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create payment method: %w", err)
	}

	return pm.ID, nil
}

// createCustomer creates a customer in Stripe
func (s *StripePaymentService) createCustomer(email string, name string) (string, error) {
	if email == "" {
		return "", errors.New("email is required to create customer")
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	if name != "" {
		params.Name = stripe.String(name)
	}

	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create customer: %w", err)
	}

	return c.ID, nil
}

// ProcessPayment processes a payment request using Stripe
func (s *StripePaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	// Convert amount to cents (Stripe requires amounts in the smallest currency unit)
	amountInCents := int64(request.Amount)

	// Set up payment method based on the payment method type
	var paymentMethodID string
	var paymentMethodType string
	var err error

	switch request.PaymentMethod {
	case service.PaymentMethodCreditCard:
		if request.CardDetails == nil {
			return &service.PaymentResult{
				Success:      false,
				ErrorMessage: "card details are required for credit card payment",
				Provider:     service.PaymentProviderStripe,
			}, nil
		}
		paymentMethodType = "card"

		// Create payment method from card details or use token
		paymentMethodID, err = s.createPaymentMethodFromCard(request.CardDetails)
		if err != nil {
			s.logger.Error("Failed to create payment method: %v", err)
			return &service.PaymentResult{
				Success:      false,
				ErrorMessage: "failed to create payment method: " + err.Error(),
				Provider:     service.PaymentProviderStripe,
			}, nil
		}

	case service.PaymentMethodPayPal:
		// Stripe supports PayPal through payment methods API
		if request.PayPalDetails == nil {
			return &service.PaymentResult{
				Success:      false,
				ErrorMessage: "PayPal details are required for PayPal payment",
				Provider:     service.PaymentProviderStripe,
			}, nil
		}
		paymentMethodType = "paypal"
		paymentMethodID = request.PayPalDetails.Token

	default:
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "unsupported payment method for Stripe",
			Provider:     service.PaymentProviderStripe,
		}, nil
	}

	// If no payment method ID is provided, return an error
	if paymentMethodID == "" {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "payment method token is required",
			Provider:     service.PaymentProviderStripe,
		}, nil
	}

	// Create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amountInCents),
		Currency:      stripe.String(s.getCurrencyCode(request.Currency)),
		PaymentMethod: stripe.String(paymentMethodID),
		Description:   stripe.String(s.config.PaymentDescription),
		Confirm:       stripe.Bool(true), // Confirm the payment intent immediately
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": fmt.Sprint(request.OrderID),
				"method":   paymentMethodType,
			},
		},
	}

	// Create a customer if email is provided
	if request.CustomerEmail != "" {
		// First, attach email to receipt
		params.ReceiptEmail = stripe.String(request.CustomerEmail)

		// Then, create customer and attach to payment
		customerName := ""
		if request.CardDetails != nil && request.CardDetails.CardholderName != "" {
			customerName = request.CardDetails.CardholderName
		}

		customerID, err := s.createCustomer(request.CustomerEmail, customerName)
		if err != nil {
			s.logger.Warn("Failed to create customer, proceeding with payment: %v", err)
			// Continue with payment, just without customer association
		} else {
			// Associate payment with customer
			params.Customer = stripe.String(customerID)

			// Save payment method for future use if it's a card
			if paymentMethodType == "card" {
				params.SetupFutureUsage = stripe.String("off_session")
			}
		}
	}

	// Create and confirm the payment intent
	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe payment intent: %v", err)
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "failed to process payment: " + err.Error(),
			Provider:     service.PaymentProviderStripe,
		}, nil
	}

	// Check payment intent status
	switch paymentIntent.Status {
	case stripe.PaymentIntentStatusSucceeded:
		// Payment succeeded
		return &service.PaymentResult{
			Success:       true,
			TransactionID: paymentIntent.ID,
			Provider:      service.PaymentProviderStripe,
		}, nil

	case stripe.PaymentIntentStatusRequiresAction:
		// Payment requires additional action (e.g., 3D Secure)
		return &service.PaymentResult{
			Success:        false,
			TransactionID:  paymentIntent.ID,
			ErrorMessage:   "payment requires additional action",
			RequiresAction: true,
			ActionURL:      paymentIntent.NextAction.RedirectToURL.URL,
			Provider:       service.PaymentProviderStripe,
		}, nil

	default:
		// Payment failed or is in another state
		return &service.PaymentResult{
			Success:       false,
			TransactionID: paymentIntent.ID,
			ErrorMessage:  fmt.Sprintf("payment status: %s", paymentIntent.Status),
			Provider:      service.PaymentProviderStripe,
		}, nil
	}
}

// getCurrencyCode returns the standardized currency code
func (s *StripePaymentService) getCurrencyCode(currency string) string {
	if currency == "" {
		return string(stripe.CurrencyUSD) // Default currency
	}
	return currency
}

// VerifyPayment verifies a payment
func (s *StripePaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	if transactionID == "" {
		return false, errors.New("transaction ID is required")
	}

	// Retrieve the payment intent from Stripe
	paymentIntent, err := paymentintent.Get(transactionID, nil)
	if err != nil {
		s.logger.Error("Failed to retrieve Stripe payment intent: %v", err)
		return false, fmt.Errorf("failed to verify payment: %w", err)
	}

	// Check if the payment intent was successful
	if paymentIntent.Status == stripe.PaymentIntentStatusSucceeded {
		return true, nil
	} else if paymentIntent.Status == stripe.PaymentIntentStatusRequiresCapture {
		// Payment is authorized but requires capture
		return true, nil
	}

	return false, nil
}

// RefundPayment refunds a payment
func (s *StripePaymentService) RefundPayment(transactionID string, amount int64, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Create refund params
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(transactionID),
		Amount:        stripe.Int64(amount),
	}

	// Process the refund
	refundResult, err := refund.New(params)
	if err != nil {
		s.logger.Error("Failed to process Stripe refund: %v", err)
		return fmt.Errorf("failed to process refund: %w", err)
	}

	if refundResult.Status != stripe.RefundStatusSucceeded {
		s.logger.Warn("Refund created with status %s", refundResult.Status)
	}

	return nil
}

// CapturePayment captures a payment
func (s *StripePaymentService) CapturePayment(transactionID string, amount int64, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}
	if amount <= 0 {
		return errors.New("capture amount must be greater than zero")
	}

	// Create capture params
	params := &stripe.PaymentIntentCaptureParams{
		AmountToCapture: stripe.Int64(amount),
	}

	// Capture the payment intent
	captureResult, err := paymentintent.Capture(transactionID, params)
	if err != nil {
		s.logger.Error("Failed to capture Stripe payment: %v", err)
		return fmt.Errorf("failed to capture payment: %w", err)
	}

	if captureResult.Status != stripe.PaymentIntentStatusSucceeded {
		return fmt.Errorf("capture resulted in unexpected status: %s", captureResult.Status)
	}

	return nil
}

// CancelPayment cancels a payment
func (s *StripePaymentService) CancelPayment(transactionID string, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	// Create cancel params
	params := &stripe.PaymentIntentCancelParams{}

	// Cancel the payment intent
	cancelResult, err := paymentintent.Cancel(transactionID, params)
	if err != nil {
		s.logger.Error("Failed to cancel Stripe payment: %v", err)
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	if cancelResult.Status != stripe.PaymentIntentStatusCanceled {
		return fmt.Errorf("cancel resulted in unexpected status: %s", cancelResult.Status)
	}

	return nil
}

// CreateSetupIntent creates a setup intent for saving a payment method without charging
func (s *StripePaymentService) CreateSetupIntent(customerEmail string) (string, string, error) {
	// This method could be used to save payment methods for future use
	// Implementation would go here
	return "", "", errors.New("not implemented")
}
