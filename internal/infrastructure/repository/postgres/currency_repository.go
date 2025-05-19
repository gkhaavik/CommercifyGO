package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CurrencyRepository is the PostgreSQL implementation of the currency repository
type CurrencyRepository struct {
	db *sql.DB
}

// NewCurrencyRepository creates a new currency repository
func NewCurrencyRepository(db *sql.DB) repository.CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

// Create creates a new currency
func (r *CurrencyRepository) Create(currency *entity.Currency) error {
	query := `
		INSERT INTO currencies (code, name, symbol, exchange_rate, is_default, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			symbol = EXCLUDED.symbol,
			exchange_rate = EXCLUDED.exchange_rate,
			is_default = EXCLUDED.is_default,
			is_enabled = EXCLUDED.is_enabled,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(
		query,
		currency.Code,
		currency.Name,
		currency.Symbol,
		currency.ExchangeRate,
		currency.IsDefault,
		currency.IsEnabled,
		currency.CreatedAt,
		currency.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// If this is the default currency, ensure it's the only default
	if currency.IsDefault {
		_, err = r.db.Exec(
			"UPDATE currencies SET is_default = false WHERE code != $1",
			currency.Code,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByCode retrieves a currency by its code
func (r *CurrencyRepository) GetByCode(code string) (*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		WHERE code = $1
	`

	var currency entity.Currency
	err := r.db.QueryRow(query, code).Scan(
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.ExchangeRate,
		&currency.IsDefault,
		&currency.IsEnabled,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}

	return &currency, nil
}

// GetDefault retrieves the default currency
func (r *CurrencyRepository) GetDefault() (*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		WHERE is_default = true
		LIMIT 1
	`

	var currency entity.Currency
	err := r.db.QueryRow(query).Scan(
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.ExchangeRate,
		&currency.IsDefault,
		&currency.IsEnabled,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no default currency found")
		}
		return nil, err
	}

	return &currency, nil
}

// List returns all currencies
func (r *CurrencyRepository) List() ([]*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		ORDER BY is_default DESC, code ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []*entity.Currency
	for rows.Next() {
		var currency entity.Currency
		err := rows.Scan(
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.ExchangeRate,
			&currency.IsDefault,
			&currency.IsEnabled,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, &currency)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

// ListEnabled returns all enabled currencies
func (r *CurrencyRepository) ListEnabled() ([]*entity.Currency, error) {
	query := `
		SELECT code, name, symbol, exchange_rate, is_default, is_enabled, created_at, updated_at
		FROM currencies
		WHERE is_enabled = true
		ORDER BY is_default DESC, code ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []*entity.Currency
	for rows.Next() {
		var currency entity.Currency
		err := rows.Scan(
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.ExchangeRate,
			&currency.IsDefault,
			&currency.IsEnabled,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, &currency)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

// Update updates a currency
func (r *CurrencyRepository) Update(currency *entity.Currency) error {
	query := `
		UPDATE currencies
		SET name = $2, symbol = $3, exchange_rate = $4, is_default = $5, is_enabled = $6, updated_at = $7
		WHERE code = $1
	`

	_, err := r.db.Exec(
		query,
		currency.Code,
		currency.Name,
		currency.Symbol,
		currency.ExchangeRate,
		currency.IsDefault,
		currency.IsEnabled,
		time.Now(),
	)

	if err != nil {
		return err
	}

	// If this is the default currency, ensure it's the only default
	if currency.IsDefault {
		_, err = r.db.Exec(
			"UPDATE currencies SET is_default = false WHERE code != $1",
			currency.Code,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a currency
func (r *CurrencyRepository) Delete(code string) error {
	// Check if this is the default currency
	var isDefault bool
	err := r.db.QueryRow("SELECT is_default FROM currencies WHERE code = $1", code).Scan(&isDefault)
	if err != nil {
		return err
	}

	if isDefault {
		return errors.New("cannot delete default currency")
	}

	query := "DELETE FROM currencies WHERE code = $1"
	_, err = r.db.Exec(query, code)
	return err
}

// SetDefault sets a currency as the default
func (r *CurrencyRepository) SetDefault(code string) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// First, set all currencies to not be default
	_, err = tx.Exec("UPDATE currencies SET is_default = false")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then set the specified currency as default
	_, err = tx.Exec("UPDATE currencies SET is_default = true WHERE code = $1", code)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GetProductPrices retrieves all prices for a product in different currencies
func (r *CurrencyRepository) GetProductPrices(productID uint) ([]entity.ProductPrice, error) {
	query := `
		SELECT id, product_id, currency_code, price, created_at, updated_at
		FROM product_prices
		WHERE product_id = $1
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []entity.ProductPrice
	for rows.Next() {
		var price entity.ProductPrice

		err := rows.Scan(
			&price.ID,
			&price.ProductID,
			&price.CurrencyCode,
			&price.Price,
			&price.CreatedAt,
			&price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		prices = append(prices, price)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

// SetProductPrice sets or updates a price for a product in a specific currency
func (r *CurrencyRepository) SetProductPrice(price *entity.ProductPrice) error {
	query := `
		INSERT INTO product_prices (product_id, currency_code, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (product_id, currency_code) DO UPDATE SET
			price = EXCLUDED.price,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	now := time.Now()

	err := r.db.QueryRow(
		query,
		price.ProductID,
		price.CurrencyCode,
		price.Price,
		now,
		now,
	).Scan(&price.ID)

	return err
}

// DeleteProductPrice removes a price for a product in a specific currency
func (r *CurrencyRepository) DeleteProductPrice(productID uint, currencyCode string) error {
	query := "DELETE FROM product_prices WHERE product_id = $1 AND currency_code = $2"
	_, err := r.db.Exec(query, productID, currencyCode)
	return err
}

// GetProductVariantPrices retrieves all prices for a product variant in different currencies
func (r *CurrencyRepository) GetVariantPrices(variantID uint) ([]entity.ProductVariantPrice, error) {
	query := `
		SELECT id, variant_id, currency_code, price, created_at, updated_at
		FROM product_variant_prices
		WHERE variant_id = $1
	`

	rows, err := r.db.Query(query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []entity.ProductVariantPrice
	for rows.Next() {
		var price entity.ProductVariantPrice

		err := rows.Scan(
			&price.ID,
			&price.VariantID,
			&price.CurrencyCode,
			&price.Price,
			&price.CreatedAt,
			&price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		prices = append(prices, price)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

// SetProductVariantPrice sets or updates a price for a product variant in a specific currency
func (r *CurrencyRepository) SetVariantPrice(price *entity.ProductVariantPrice) error {
	query := `
		INSERT INTO product_variant_prices (variant_id, currency_code, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (variant_id, currency_code) DO UPDATE SET
			price = EXCLUDED.price,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	now := time.Now()

	err := r.db.QueryRow(
		query,
		price.VariantID,
		price.CurrencyCode,
		price.Price,
		now,
		now,
	).Scan(&price.ID)

	return err
}

// DeleteProductVariantPrice removes a price for a product variant in a specific currency
func (r *CurrencyRepository) DeleteVariantPrice(variantID uint, currencyCode string) error {
	query := "DELETE FROM product_variant_prices WHERE variant_id = $1 AND currency_code = $2"
	_, err := r.db.Exec(query, variantID, currencyCode)
	return err
}
