package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CurrencyRepository implements the domain's CurrencyRepository interface
type CurrencyRepository struct {
	db *sql.DB
}

// NewCurrencyRepository creates a new PostgreSQL-backed currency repository
func NewCurrencyRepository(db *sql.DB) repository.CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

// GetAll returns all currencies in the system
func (r *CurrencyRepository) GetAll() ([]*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		ORDER BY is_default DESC, code ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query currencies: %w", err)
	}
	defer rows.Close()

	var currencies []*entity.Currency
	for rows.Next() {
		var c entity.Currency
		if err := rows.Scan(
			&c.Code,
			&c.Name,
			&c.Symbol,
			&c.Precision,
			&c.ExchangeRate,
			&c.IsDefault,
			&c.IsEnabled,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan currency row: %w", err)
		}
		currencies = append(currencies, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating currency rows: %w", err)
	}

	return currencies, nil
}

// GetByCode returns a currency by its ISO code
func (r *CurrencyRepository) GetByCode(code string) (*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		WHERE code = $1
	`

	var c entity.Currency
	err := r.db.QueryRow(query, code).Scan(
		&c.Code,
		&c.Name,
		&c.Symbol,
		&c.Precision,
		&c.ExchangeRate,
		&c.IsDefault,
		&c.IsEnabled,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("currency with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to get currency by code: %w", err)
	}

	return &c, nil
}

// GetEnabled returns all enabled currencies
func (r *CurrencyRepository) GetEnabled() ([]*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		WHERE is_enabled = true
		ORDER BY is_default DESC, code ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled currencies: %w", err)
	}
	defer rows.Close()

	var currencies []*entity.Currency
	for rows.Next() {
		var c entity.Currency
		if err := rows.Scan(
			&c.Code,
			&c.Name,
			&c.Symbol,
			&c.Precision,
			&c.ExchangeRate,
			&c.IsDefault,
			&c.IsEnabled,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan currency row: %w", err)
		}
		currencies = append(currencies, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating enabled currency rows: %w", err)
	}

	return currencies, nil
}

// GetDefault returns the default currency
func (r *CurrencyRepository) GetDefault() (*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		WHERE is_default = true
	`

	var c entity.Currency
	err := r.db.QueryRow(query).Scan(
		&c.Code,
		&c.Name,
		&c.Symbol,
		&c.Precision,
		&c.ExchangeRate,
		&c.IsDefault,
		&c.IsEnabled,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no default currency found")
		}
		return nil, fmt.Errorf("failed to get default currency: %w", err)
	}

	return &c, nil
}

// Create adds a new currency
func (r *CurrencyRepository) Create(currency *entity.Currency) error {
	// First check if the code already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)", currency.Code).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if currency exists: %w", err)
	}

	if exists {
		return fmt.Errorf("currency with code %s already exists", currency.Code)
	}

	// Start a transaction if this is the default currency
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// If this currency is being set as default, unset all others
	if currency.IsDefault {
		_, err = tx.Exec("UPDATE currencies SET is_default = FALSE")
		if err != nil {
			return fmt.Errorf("failed to unset default currencies: %w", err)
		}
	}

	// Insert the new currency
	query := `
		INSERT INTO currencies (code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = tx.Exec(
		query,
		currency.Code,
		currency.Name,
		currency.Symbol,
		currency.Precision,
		currency.ExchangeRate,
		currency.IsDefault,
		currency.IsEnabled,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert currency: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Update modifies an existing currency
func (r *CurrencyRepository) Update(currency *entity.Currency) error {
	query := `
		UPDATE currencies
		SET name = $1, 
			symbol = $2, 
			precision = $3, 
			exchange_rate = $4, 
			is_default = $5, 
			is_enabled = $6, 
			updated_at = NOW()
		WHERE code = $7
	`

	result, err := r.db.Exec(
		query,
		currency.Name,
		currency.Symbol,
		currency.Precision,
		currency.ExchangeRate,
		currency.IsDefault,
		currency.IsEnabled,
		currency.Code,
	)

	if err != nil {
		return fmt.Errorf("failed to update currency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("currency with code %s not found", currency.Code)
	}

	return nil
}

// Delete removes a currency
func (r *CurrencyRepository) Delete(code string) error {
	// First check if this is the default currency
	var isDefault bool
	err := r.db.QueryRow("SELECT is_default FROM currencies WHERE code = $1", code).Scan(&isDefault)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("currency with code %s not found", code)
		}
		return fmt.Errorf("failed to check if currency is default: %w", err)
	}

	if isDefault {
		return fmt.Errorf("cannot delete the default currency")
	}

	// Then try to delete the currency
	result, err := r.db.Exec("DELETE FROM currencies WHERE code = $1", code)
	if err != nil {
		return fmt.Errorf("failed to delete currency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("currency with code %s not found", code)
	}

	return nil
}

// SetDefault sets a currency as the default
func (r *CurrencyRepository) SetDefault(code string) error {
	// Start a transaction to ensure atomicity when changing the default currency
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// First check if the currency exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)", code).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if currency exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("currency with code %s not found", code)
	}

	// Unset all currencies as default
	_, err = tx.Exec("UPDATE currencies SET is_default = FALSE")
	if err != nil {
		return fmt.Errorf("failed to unset default currencies: %w", err)
	}

	// Set the new default currency
	_, err = tx.Exec("UPDATE currencies SET is_default = TRUE WHERE code = $1", code)
	if err != nil {
		return fmt.Errorf("failed to set default currency: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddExchangeRate records a new exchange rate
func (r *CurrencyRepository) AddExchangeRate(history *entity.ExchangeRateHistory) error {
	query := `
		INSERT INTO exchange_rate_history (base_currency, target_currency, rate, date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (base_currency, target_currency, date) DO UPDATE
		SET rate = EXCLUDED.rate
	`

	_, err := r.db.Exec(
		query,
		history.BaseCurrency,
		history.TargetCurrency,
		history.Rate,
		history.Date,
	)

	if err != nil {
		return fmt.Errorf("failed to add exchange rate: %w", err)
	}

	return nil
}

// GetExchangeRateHistory gets historical exchange rates for a currency pair
func (r *CurrencyRepository) GetExchangeRateHistory(baseCurrency, targetCurrency string, limit int) ([]*entity.ExchangeRateHistory, error) {
	query := `
		SELECT id, base_currency, target_currency, rate, date
		FROM exchange_rate_history
		WHERE base_currency = $1 AND target_currency = $2
		ORDER BY date DESC
		LIMIT $3
	`

	rows, err := r.db.Query(query, baseCurrency, targetCurrency, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query exchange rate history: %w", err)
	}
	defer rows.Close()

	var history []*entity.ExchangeRateHistory
	for rows.Next() {
		var h entity.ExchangeRateHistory
		if err := rows.Scan(&h.ID, &h.BaseCurrency, &h.TargetCurrency, &h.Rate, &h.Date); err != nil {
			return nil, fmt.Errorf("failed to scan exchange rate history row: %w", err)
		}
		history = append(history, &h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through exchange rate history rows: %w", err)
	}

	return history, nil
}

// GetLatestExchangeRate gets the most recent exchange rate for a currency pair
func (r *CurrencyRepository) GetLatestExchangeRate(baseCurrency, targetCurrency string) (*entity.ExchangeRateHistory, error) {
	query := `
		SELECT id, base_currency, target_currency, rate, date
		FROM exchange_rate_history
		WHERE base_currency = $1 AND target_currency = $2
		ORDER BY date DESC
		LIMIT 1
	`

	var history entity.ExchangeRateHistory
	err := r.db.QueryRow(query, baseCurrency, targetCurrency).Scan(
		&history.ID,
		&history.BaseCurrency,
		&history.TargetCurrency,
		&history.Rate,
		&history.Date,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no exchange rate found for %s to %s", baseCurrency, targetCurrency)
		}
		return nil, fmt.Errorf("failed to get latest exchange rate: %w", err)
	}

	return &history, nil
}

// UpdateExchangeRates updates exchange rates for all currencies
func (r *CurrencyRepository) UpdateExchangeRates(baseCurrency string, rates map[string]float64) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update the exchange rates in currencies table
	for currencyCode, rate := range rates {
		// Skip the base currency itself (should always be 1.0)
		if currencyCode == baseCurrency {
			continue
		}

		_, err = tx.Exec(
			"UPDATE currencies SET exchange_rate = $1, updated_at = NOW() WHERE code = $2",
			rate,
			currencyCode,
		)

		if err != nil {
			return fmt.Errorf("failed to update exchange rate for %s: %w", currencyCode, err)
		}
	}

	// Also need to add records to the exchange_rate_history table
	// Get current timestamp for consistent date across all entries
	now := time.Now()

	// Insert history records
	stmt, err := tx.Prepare(`
		INSERT INTO exchange_rate_history (base_currency, target_currency, rate, date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (base_currency, target_currency, date) DO UPDATE
		SET rate = EXCLUDED.rate
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for currencyCode, rate := range rates {
		// Skip the base currency itself
		if currencyCode == baseCurrency {
			continue
		}

		_, err = stmt.Exec(baseCurrency, currencyCode, rate, now)
		if err != nil {
			return fmt.Errorf("failed to insert exchange rate history for %s: %w", currencyCode, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
