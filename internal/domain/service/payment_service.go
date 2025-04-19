package service

// PaymentProviderType represents a payment provider type
type PaymentProviderType string

const (
	PaymentProviderStripe PaymentProviderType = "stripe"
	PaymentProviderPayPal PaymentProviderType = "paypal"
	PaymentProviderMock   PaymentProviderType = "mock"
)

// PaymentMethod represents a payment method type
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodPayPal       PaymentMethod = "paypal"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// PaymentProvider represents information about a payment provider
type PaymentProvider struct {
	Type        PaymentProviderType `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	IconURL     string              `json:"icon_url,omitempty"`
	Methods     []PaymentMethod     `json:"methods"`
	Enabled     bool                `json:"enabled"`
}

// PaymentRequest represents a request to process a payment
type PaymentRequest struct {
	OrderID         uint
	Amount          float64
	Currency        string
	PaymentMethod   PaymentMethod
	PaymentProvider PaymentProviderType
	CardDetails     *CardDetails
	PayPalDetails   *PayPalDetails
	BankDetails     *BankDetails
	CustomerEmail   string
}

// CardDetails represents credit card payment details
type CardDetails struct {
	CardNumber     string
	ExpiryMonth    int
	ExpiryYear     int
	CVV            string
	CardholderName string
	Token          string
}

// PayPalDetails represents PayPal payment details
type PayPalDetails struct {
	Email string
	Token string
}

// BankDetails represents bank transfer details
type BankDetails struct {
	AccountNumber string
	BankCode      string
	AccountName   string
}

// PaymentResult represents the result of a payment processing
type PaymentResult struct {
	Success        bool
	TransactionID  string
	ErrorMessage   string
	RequiresAction bool
	ActionURL      string
	Provider       PaymentProviderType
}

// PaymentService defines the interface for payment processing
type PaymentService interface {
	// GetAvailableProviders returns a list of available payment providers
	GetAvailableProviders() []PaymentProvider

	// ProcessPayment processes a payment request
	ProcessPayment(request PaymentRequest) (*PaymentResult, error)

	// VerifyPayment verifies a payment
	VerifyPayment(transactionID string, provider PaymentProviderType) (bool, error)

	// RefundPayment refunds a payment
	RefundPayment(transactionID string, amount float64, provider PaymentProviderType) error
}
