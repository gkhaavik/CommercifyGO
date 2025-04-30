package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

type paymentTransactionRepository struct {
	db *sql.DB
}

// NewPaymentTransactionRepository creates a new PaymentTransactionRepository
func NewPaymentTransactionRepository(db *sql.DB) repository.PaymentTransactionRepository {
	return &paymentTransactionRepository{
		db: db,
	}
}

// Create inserts a new payment transaction into the database
func (r *paymentTransactionRepository) Create(transaction *entity.PaymentTransaction) error {
	// Convert metadata to JSON string
	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO payment_transactions 
		(order_id, transaction_id, type, status, amount, currency, provider, raw_response, metadata, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	err = r.db.QueryRow(
		query,
		transaction.OrderID,
		transaction.TransactionID,
		string(transaction.Type),
		string(transaction.Status),
		transaction.Amount,
		transaction.Currency,
		transaction.Provider,
		transaction.RawResponse,
		metadataJSON,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	).Scan(&transaction.ID)

	if err != nil {
		return fmt.Errorf("failed to create payment transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a payment transaction by ID
func (r *paymentTransactionRepository) GetByID(id uint) (*entity.PaymentTransaction, error) {
	query := `
		SELECT id, order_id, transaction_id, type, status, amount, currency, provider, raw_response, metadata, created_at, updated_at
		FROM payment_transactions
		WHERE id = $1
	`

	var metadataJSON string
	tx := &entity.PaymentTransaction{}

	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.OrderID,
		&tx.TransactionID,
		&tx.Type,
		&tx.Status,
		&tx.Amount,
		&tx.Currency,
		&tx.Provider,
		&tx.RawResponse,
		&metadataJSON,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment transaction not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get payment transaction: %w", err)
	}

	// Parse metadata JSON
	if metadataJSON != "" {
		metadata := make(map[string]string)
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		tx.Metadata = metadata
	} else {
		tx.Metadata = make(map[string]string)
	}

	return tx, nil
}

// GetByTransactionID retrieves a payment transaction by external transaction ID
func (r *paymentTransactionRepository) GetByTransactionID(transactionID string) (*entity.PaymentTransaction, error) {
	query := `
		SELECT id, order_id, transaction_id, type, status, amount, currency, provider, raw_response, metadata, created_at, updated_at
		FROM payment_transactions
		WHERE transaction_id = $1
	`

	var metadataJSON string
	tx := &entity.PaymentTransaction{}

	err := r.db.QueryRow(query, transactionID).Scan(
		&tx.ID,
		&tx.OrderID,
		&tx.TransactionID,
		&tx.Type,
		&tx.Status,
		&tx.Amount,
		&tx.Currency,
		&tx.Provider,
		&tx.RawResponse,
		&metadataJSON,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment transaction not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get payment transaction: %w", err)
	}

	// Parse metadata JSON
	if metadataJSON != "" {
		metadata := make(map[string]string)
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		tx.Metadata = metadata
	} else {
		tx.Metadata = make(map[string]string)
	}

	return tx, nil
}

// GetByOrderID retrieves all payment transactions for an order
func (r *paymentTransactionRepository) GetByOrderID(orderID uint) ([]*entity.PaymentTransaction, error) {
	query := `
		SELECT id, order_id, transaction_id, type, status, amount, currency, provider, raw_response, metadata, created_at, updated_at
		FROM payment_transactions
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entity.PaymentTransaction

	for rows.Next() {
		var metadataJSON string
		tx := &entity.PaymentTransaction{}

		err := rows.Scan(
			&tx.ID,
			&tx.OrderID,
			&tx.TransactionID,
			&tx.Type,
			&tx.Status,
			&tx.Amount,
			&tx.Currency,
			&tx.Provider,
			&tx.RawResponse,
			&metadataJSON,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment transaction: %w", err)
		}

		// Parse metadata JSON
		if metadataJSON != "" {
			metadata := make(map[string]string)
			if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			tx.Metadata = metadata
		} else {
			tx.Metadata = make(map[string]string)
		}

		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payment transactions rows: %w", err)
	}

	return transactions, nil
}

// Update updates a payment transaction
func (r *paymentTransactionRepository) Update(transaction *entity.PaymentTransaction) error {
	// Convert metadata to JSON string
	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		UPDATE payment_transactions
		SET transaction_id = $1, 
		    type = $2, 
		    status = $3, 
		    amount = $4, 
		    currency = $5, 
		    provider = $6, 
		    raw_response = $7, 
		    metadata = $8, 
		    updated_at = $9
		WHERE id = $10
	`

	result, err := r.db.Exec(
		query,
		transaction.TransactionID,
		string(transaction.Type),
		string(transaction.Status),
		transaction.Amount,
		transaction.Currency,
		transaction.Provider,
		transaction.RawResponse,
		metadataJSON,
		transaction.UpdatedAt,
		transaction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment transaction not found")
	}

	return nil
}

// Delete deletes a payment transaction
func (r *paymentTransactionRepository) Delete(id uint) error {
	query := "DELETE FROM payment_transactions WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment transaction not found")
	}

	return nil
}

// GetLatestByOrderIDAndType retrieves the latest transaction of a specific type for an order
func (r *paymentTransactionRepository) GetLatestByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (*entity.PaymentTransaction, error) {
	query := `
		SELECT id, order_id, transaction_id, type, status, amount, currency, provider, raw_response, metadata, created_at, updated_at
		FROM payment_transactions
		WHERE order_id = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var metadataJSON string
	tx := &entity.PaymentTransaction{}

	err := r.db.QueryRow(query, orderID, string(transactionType)).Scan(
		&tx.ID,
		&tx.OrderID,
		&tx.TransactionID,
		&tx.Type,
		&tx.Status,
		&tx.Amount,
		&tx.Currency,
		&tx.Provider,
		&tx.RawResponse,
		&metadataJSON,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No transaction found, not an error
		}
		return nil, fmt.Errorf("failed to get latest payment transaction: %w", err)
	}

	// Parse metadata JSON
	if metadataJSON != "" {
		metadata := make(map[string]string)
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		tx.Metadata = metadata
	} else {
		tx.Metadata = make(map[string]string)
	}

	return tx, nil
}

// CountSuccessfulByOrderIDAndType counts successful transactions of a specific type for an order
func (r *paymentTransactionRepository) CountSuccessfulByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM payment_transactions
		WHERE order_id = $1 AND type = $2 AND status = $3
	`

	var count int
	err := r.db.QueryRow(query, orderID, string(transactionType), string(entity.TransactionStatusSuccessful)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count successful transactions: %w", err)
	}

	return count, nil
}

// SumAmountByOrderIDAndType sums the amount of transactions of a specific type for an order
func (r *paymentTransactionRepository) SumAmountByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM payment_transactions
		WHERE order_id = $1 AND type = $2 AND status = $3
	`

	var total float64
	err := r.db.QueryRow(query, orderID, string(transactionType), string(entity.TransactionStatusSuccessful)).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to sum transaction amounts: %w", err)
	}

	return total, nil
}
