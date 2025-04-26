package payment

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
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

// ProcessPayment processes a payment request using Stripe
func (s *StripePaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	// Convert amount to cents (Stripe requires amounts in the smallest currency unit)
	amountInCents := int64(request.Amount * 100)

	// Set up payment method based on the payment method type
	var paymentMethodID string
	var paymentMethodType string

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
		// In a real implementation, you would create a payment method using the card details
		// For now, assume the card token is passed directly
		paymentMethodID = request.CardDetails.Token

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
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
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

	// Add receipt email if available
	if request.CustomerEmail != "" {
		params.ReceiptEmail = stripe.String(request.CustomerEmail)
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
			ErrorMessage:  fmt.Sprint(paymentIntent.Status),
			Provider:      service.PaymentProviderStripe,
		}, nil
	}
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
		return false, err
	}

	// Check if the payment intent was successful
	return paymentIntent.Status == stripe.PaymentIntentStatusSucceeded, nil
}

// RefundPayment refunds a payment
func (s *StripePaymentService) RefundPayment(transactionID string, amount float64, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Convert amount to cents
	amountInCents := int64(amount * 100)

	// Create refund params
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(transactionID),
		Amount:        stripe.Int64(amountInCents),
	}

	// Process the refund
	_, err := refund.New(params)
	if err != nil {
		s.logger.Error("Failed to process Stripe refund: %v", err)
		return err
	}

	return nil
}
