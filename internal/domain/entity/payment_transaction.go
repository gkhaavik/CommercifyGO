package entity

import (
	"time"
)

// TransactionType represents the type of payment transaction
type TransactionType string

const (
	TransactionTypeAuthorize TransactionType = "authorize"
	TransactionTypeCapture   TransactionType = "capture"
	TransactionTypeRefund    TransactionType = "refund"
	TransactionTypeCancel    TransactionType = "cancel"
)

// TransactionStatus represents the status of a payment transaction
type TransactionStatus string

const (
	TransactionStatusSuccessful TransactionStatus = "successful"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusPending    TransactionStatus = "pending"
)

// PaymentTransaction represents a payment transaction record
type PaymentTransaction struct {
	ID            uint
	OrderID       uint
	TransactionID string            // External transaction ID from payment provider
	Type          TransactionType   // Type of transaction (authorize, capture, refund, cancel)
	Status        TransactionStatus // Status of the transaction
	Amount        int64             // Amount of the transaction
	Currency      string            // Currency of the transaction
	Provider      string            // Payment provider (stripe, paypal, etc.)
	RawResponse   string            // Raw response from payment provider (JSON)
	Metadata      map[string]string // Additional metadata
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewPaymentTransaction creates a new payment transaction
func NewPaymentTransaction(
	orderID uint,
	transactionID string,
	transactionType TransactionType,
	status TransactionStatus,
	amount int64,
	currency string,
	provider string,
) (*PaymentTransaction, error) {
	if orderID == 0 {
		return nil, ErrInvalidInput{Field: "OrderID", Message: "cannot be zero"}
	}
	if transactionID == "" {
		return nil, ErrInvalidInput{Field: "TransactionID", Message: "cannot be empty"}
	}
	if string(transactionType) == "" {
		return nil, ErrInvalidInput{Field: "TransactionType", Message: "cannot be empty"}
	}
	if string(status) == "" {
		return nil, ErrInvalidInput{Field: "Status", Message: "cannot be empty"}
	}
	if provider == "" {
		return nil, ErrInvalidInput{Field: "Provider", Message: "cannot be empty"}
	}
	if currency == "" {
		currency = "USD" // Default currency
	}

	now := time.Now()

	return &PaymentTransaction{
		OrderID:       orderID,
		TransactionID: transactionID,
		Type:          transactionType,
		Status:        status,
		Amount:        amount,
		Currency:      currency,
		Provider:      provider,
		Metadata:      make(map[string]string),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// AddMetadata adds metadata to the transaction
func (pt *PaymentTransaction) AddMetadata(key, value string) {
	if pt.Metadata == nil {
		pt.Metadata = make(map[string]string)
	}
	pt.Metadata[key] = value
	pt.UpdatedAt = time.Now()
}

// SetRawResponse sets the raw response from the payment provider
func (pt *PaymentTransaction) SetRawResponse(response string) {
	pt.RawResponse = response
	pt.UpdatedAt = time.Now()
}

// UpdateStatus updates the status of the transaction
func (pt *PaymentTransaction) UpdateStatus(status TransactionStatus) {
	pt.Status = status
	pt.UpdatedAt = time.Now()
}
