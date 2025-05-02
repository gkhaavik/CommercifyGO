package mock

import (
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockPaymentTransactionRepository implements a mock payment transaction repository for testing
type MockPaymentTransactionRepository struct {
	transactions    map[uint]*entity.PaymentTransaction
	nextID          uint
	byOrderID       map[uint][]*entity.PaymentTransaction
	byTransactionID map[string]*entity.PaymentTransaction
}

// NewMockPaymentTransactionRepository creates a new mock payment transaction repository
func NewMockPaymentTransactionRepository() *MockPaymentTransactionRepository {
	return &MockPaymentTransactionRepository{
		transactions:    make(map[uint]*entity.PaymentTransaction),
		byOrderID:       make(map[uint][]*entity.PaymentTransaction),
		byTransactionID: make(map[string]*entity.PaymentTransaction),
		nextID:          1,
	}
}

// Create adds a new payment transaction
func (m *MockPaymentTransactionRepository) Create(tx *entity.PaymentTransaction) error {
	if tx == nil {
		return errors.New("payment transaction cannot be nil")
	}

	tx.ID = m.nextID
	m.nextID++

	// Store transaction in our maps for quick lookup
	m.transactions[tx.ID] = tx

	// Store by order ID
	if _, ok := m.byOrderID[tx.OrderID]; !ok {
		m.byOrderID[tx.OrderID] = make([]*entity.PaymentTransaction, 0)
	}
	m.byOrderID[tx.OrderID] = append(m.byOrderID[tx.OrderID], tx)

	// Store by transaction ID
	m.byTransactionID[tx.TransactionID] = tx

	return nil
}

// GetByID retrieves a payment transaction by ID
func (m *MockPaymentTransactionRepository) GetByID(id uint) (*entity.PaymentTransaction, error) {
	tx, ok := m.transactions[id]
	if !ok {
		return nil, errors.New("payment transaction not found")
	}
	return tx, nil
}

// GetByOrderID retrieves all payment transactions for an order
func (m *MockPaymentTransactionRepository) GetByOrderID(orderID uint) ([]*entity.PaymentTransaction, error) {
	transactions, ok := m.byOrderID[orderID]
	if !ok {
		return []*entity.PaymentTransaction{}, nil
	}
	return transactions, nil
}

// GetByTransactionID retrieves a payment transaction by external transaction ID
func (m *MockPaymentTransactionRepository) GetByTransactionID(transactionID string) (*entity.PaymentTransaction, error) {
	tx, ok := m.byTransactionID[transactionID]
	if !ok {
		return nil, errors.New("payment transaction not found")
	}
	return tx, nil
}

// Update updates a payment transaction
func (m *MockPaymentTransactionRepository) Update(transaction *entity.PaymentTransaction) error {
	if transaction == nil {
		return errors.New("payment transaction cannot be nil")
	}

	_, ok := m.transactions[transaction.ID]
	if !ok {
		return errors.New("payment transaction not found")
	}

	transaction.UpdatedAt = time.Now()
	m.transactions[transaction.ID] = transaction
	m.byTransactionID[transaction.TransactionID] = transaction

	return nil
}

// Delete deletes a payment transaction
func (m *MockPaymentTransactionRepository) Delete(id uint) error {
	tx, ok := m.transactions[id]
	if !ok {
		return errors.New("payment transaction not found")
	}

	// Remove from all maps
	delete(m.transactions, id)
	delete(m.byTransactionID, tx.TransactionID)

	// Remove from byOrderID map
	if txs, ok := m.byOrderID[tx.OrderID]; ok {
		updatedTxs := make([]*entity.PaymentTransaction, 0, len(txs)-1)
		for _, t := range txs {
			if t.ID != id {
				updatedTxs = append(updatedTxs, t)
			}
		}
		if len(updatedTxs) > 0 {
			m.byOrderID[tx.OrderID] = updatedTxs
		} else {
			delete(m.byOrderID, tx.OrderID)
		}
	}

	return nil
}

// GetLatestByOrderIDAndType retrieves the latest transaction of a specific type for an order
func (m *MockPaymentTransactionRepository) GetLatestByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (*entity.PaymentTransaction, error) {
	transactions, ok := m.byOrderID[orderID]
	if !ok || len(transactions) == 0 {
		return nil, nil
	}

	var latestTx *entity.PaymentTransaction
	var latestTime time.Time

	for _, tx := range transactions {
		if tx.Type == transactionType && (latestTx == nil || tx.CreatedAt.After(latestTime)) {
			latestTx = tx
			latestTime = tx.CreatedAt
		}
	}

	if latestTx == nil {
		return nil, nil
	}

	return latestTx, nil
}

// CountSuccessfulByOrderIDAndType counts successful transactions of a specific type for an order
func (m *MockPaymentTransactionRepository) CountSuccessfulByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int, error) {
	transactions, ok := m.byOrderID[orderID]
	if !ok {
		return 0, nil
	}

	count := 0
	for _, tx := range transactions {
		if tx.Type == transactionType && tx.Status == entity.TransactionStatusSuccessful {
			count++
		}
	}

	return count, nil
}

// SumAmountByOrderIDAndType sums the amount of transactions of a specific type for an order
func (m *MockPaymentTransactionRepository) SumAmountByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int64, error) {
	transactions, ok := m.byOrderID[orderID]
	if !ok {
		return 0, nil
	}

	var total int64
	for _, tx := range transactions {
		if tx.Type == transactionType && tx.Status == entity.TransactionStatusSuccessful {
			total += tx.Amount
		}
	}

	return total, nil
}

// IsEmpty checks if the repository has any transactions
func (m *MockPaymentTransactionRepository) IsEmpty() bool {
	return len(m.transactions) == 0
}

// Count returns the total number of transactions in the repository
func (m *MockPaymentTransactionRepository) Count() int {
	return len(m.transactions)
}
