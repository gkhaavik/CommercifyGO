package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// SQLCurrencyRepository implements the CurrencyRepository interface using SQL database
type SQLCurrencyRepository struct {
	db *sql.DB
}

// NewSQLCurrencyRepository creates a new SQLCurrencyRepository
func NewSQLCurrencyRepository(db *sql.DB) repository.CurrencyRepository {
	return &SQLCurrencyRepository{db: db}
}

// GetAll returns all currencies in the system
func (r *SQLCurrencyRepository) GetAll() ([]*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at 
		FROM currencies
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query currencies: %w", err)
	}
	defer rows.Close()

	currencies := []*entity.Currency{}
	for rows.Next() {
		currency := &entity.Currency{}
		err := rows.Scan(
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.Precision,
			&currency.ExchangeRate,
			&currency.IsDefault,
			&currency.IsEnabled,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating currency rows: %w", err)
	}

	return currencies, nil
}

// GetByCode returns a currency by its ISO code
func (r *SQLCurrencyRepository) GetByCode(code string) (*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at 
		FROM currencies 
		WHERE code = $1
	`

	currency := &entity.Currency{}
	err := r.db.QueryRow(query, code).Scan(
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.Precision,
		&currency.ExchangeRate,
		&currency.IsDefault,
		&currency.IsEnabled,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("currency with code %s not found", code)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get currency: %w", err)
	}

	return currency, nil
}

// GetDefault returns the default currency
func (r *SQLCurrencyRepository) GetDefault() (*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at 
		FROM currencies 
		WHERE is_default = true
	`

	currency := &entity.Currency{}
	err := r.db.QueryRow(query).Scan(
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.Precision,
		&currency.ExchangeRate,
		&currency.IsDefault,
		&currency.IsEnabled,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no default currency found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get default currency: %w", err)
	}

	return currency, nil
}

// GetEnabled returns all enabled currencies
func (r *SQLCurrencyRepository) GetEnabled() ([]*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at 
		FROM currencies 
		WHERE is_enabled = true
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled currencies: %w", err)
	}
	defer rows.Close()

	currencies := []*entity.Currency{}
	for rows.Next() {
		currency := &entity.Currency{}
		err := rows.Scan(
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.Precision,
			&currency.ExchangeRate,
			&currency.IsDefault,
			&currency.IsEnabled,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating currency rows: %w", err)
	}

	return currencies, nil
}

// Create adds a new currency
func (r *SQLCurrencyRepository) Create(currency *entity.Currency) error {
	query := `
		INSERT INTO currencies (
			code, name, symbol, precision, exchange_rate, is_default, is_enabled, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.Exec(
		query,
		currency.Code,
		currency.Name,
		currency.Symbol,
		currency.Precision,
		currency.ExchangeRate,
		currency.IsDefault,
		currency.IsEnabled,
		currency.CreatedAt,
		currency.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create currency: %w", err)
	}

	// If this is the default currency, update any other currencies to not be default
	if currency.IsDefault {
		return r.SetDefault(currency.Code)
	}

	return nil
}

// Update modifies an existing currency
func (r *SQLCurrencyRepository) Update(currency *entity.Currency) error {
	query := `
		UPDATE currencies 
		SET 
			name = $1, 
			symbol = $2, 
			precision = $3, 
			exchange_rate = $4, 
			is_enabled = $5, 
			updated_at = $6
		WHERE code = $7
	`

	result, err := r.db.Exec(
		query,
		currency.Name,
		currency.Symbol,
		currency.Precision,
		currency.ExchangeRate,
		currency.IsEnabled,
		time.Now(),
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
func (r *SQLCurrencyRepository) Delete(code string) error {
	// First check if this is the default currency
	currency, err := r.GetByCode(code)
	if err != nil {
		return err
	}

	if currency.IsDefault {
		return errors.New("cannot delete the default currency")
	}

	query := "DELETE FROM currencies WHERE code = $1"
	result, err := r.db.Exec(query, code)
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
func (r *SQLCurrencyRepository) SetDefault(code string) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// First, set all currencies to not be default
	_, err = tx.Exec("UPDATE currencies SET is_default = false, updated_at = $1", time.Now())
	if err != nil {
		return fmt.Errorf("failed to unset default currencies: %w", err)
	}

	// Then, set the specified currency to be default and ensure exchange rate is 1.0
	result, err := tx.Exec(
		"UPDATE currencies SET is_default = true, exchange_rate = 1.0, updated_at = $1 WHERE code = $2",
		time.Now(),
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to set default currency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("currency with code %s not found", code)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddExchangeRate records a new exchange rate
func (r *SQLCurrencyRepository) AddExchangeRate(history *entity.ExchangeRateHistory) error {
	query := `
		INSERT INTO exchange_rate_history (
			base_currency, target_currency, rate, date
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id
	`

	err := r.db.QueryRow(
		query,
		history.BaseCurrency,
		history.TargetCurrency,
		history.Rate,
		history.Date,
	).Scan(&history.ID)

	if err != nil {
		return fmt.Errorf("failed to add exchange rate history: %w", err)
	}

	return nil
}

// GetExchangeRateHistory gets historical exchange rates for a currency pair
func (r *SQLCurrencyRepository) GetExchangeRateHistory(baseCurrency, targetCurrency string, limit int) ([]*entity.ExchangeRateHistory, error) {
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

	history := []*entity.ExchangeRateHistory{}
	for rows.Next() {
		rate := &entity.ExchangeRateHistory{}
		err := rows.Scan(
			&rate.ID,
			&rate.BaseCurrency,
			&rate.TargetCurrency,
			&rate.Rate,
			&rate.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exchange rate history: %w", err)
		}
		history = append(history, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exchange rate rows: %w", err)
	}

	return history, nil
}

// GetLatestExchangeRate gets the most recent exchange rate for a currency pair
func (r *SQLCurrencyRepository) GetLatestExchangeRate(baseCurrency, targetCurrency string) (*entity.ExchangeRateHistory, error) {
	query := `
		SELECT id, base_currency, target_currency, rate, date 
		FROM exchange_rate_history 
		WHERE base_currency = $1 AND target_currency = $2
		ORDER BY date DESC
		LIMIT 1
	`

	rate := &entity.ExchangeRateHistory{}
	err := r.db.QueryRow(query, baseCurrency, targetCurrency).Scan(
		&rate.ID,
		&rate.BaseCurrency,
		&rate.TargetCurrency,
		&rate.Rate,
		&rate.Date,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No exchange rate found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get latest exchange rate: %w", err)
	}

	return rate, nil
}

// UpdateExchangeRates updates exchange rates for all currencies
func (r *SQLCurrencyRepository) UpdateExchangeRates(baseCurrency string, rates map[string]float64) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// Update each currency's exchange rate
	stmt, err := tx.Prepare(`
		UPDATE currencies 
		SET exchange_rate = $1, updated_at = $2
		WHERE code = $3
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Set base currency exchange rate to 1.0
	_, err = stmt.Exec(1.0, now, baseCurrency)
	if err != nil {
		return fmt.Errorf("failed to update base currency exchange rate: %w", err)
	}

	// Update all other currencies
	for code, rate := range rates {
		if code == baseCurrency {
			continue // Skip base currency, we already set it to 1.0
		}

		_, err = stmt.Exec(rate, now, code)
		if err != nil {
			// Log the error but continue with other currencies
			fmt.Printf("Error updating exchange rate for %s: %v\n", code, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
