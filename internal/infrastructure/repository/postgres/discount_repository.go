package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// DiscountRepository implements the discount repository interface using PostgreSQL
type DiscountRepository struct {
	db *sql.DB
}

// NewDiscountRepository creates a new DiscountRepository
func NewDiscountRepository(db *sql.DB) repository.DiscountRepository {
	return &DiscountRepository{db: db}
}

// Create creates a new discount
func (r *DiscountRepository) Create(discount *entity.Discount) error {
	query := `
		INSERT INTO discounts (
			code, type, method, value, min_order_value, max_discount_value, 
			product_ids, category_ids, start_date, end_date, 
			usage_limit, current_usage, active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`

	productIDsJSON, err := json.Marshal(discount.ProductIDs)
	if err != nil {
		return err
	}

	categoryIDsJSON, err := json.Marshal(discount.CategoryIDs)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		discount.Code,
		discount.Type,
		discount.Method,
		discount.Value,
		discount.MinOrderValue,
		discount.MaxDiscountValue,
		productIDsJSON,
		categoryIDsJSON,
		discount.StartDate,
		discount.EndDate,
		discount.UsageLimit,
		discount.CurrentUsage,
		discount.Active,
		discount.CreatedAt,
		discount.UpdatedAt,
	).Scan(&discount.ID)

	return err
}

// GetByID retrieves a discount by ID
func (r *DiscountRepository) GetByID(discountID uint) (*entity.Discount, error) {
	query := `
		SELECT id, code, type, method, value, min_order_value, max_discount_value, 
			product_ids, category_ids, start_date, end_date, 
			usage_limit, current_usage, active, created_at, updated_at
		FROM discounts
		WHERE id = $1
	`

	var productIDsJSON, categoryIDsJSON []byte
	discount := &entity.Discount{}

	err := r.db.QueryRow(query, discountID).Scan(
		&discount.ID,
		&discount.Code,
		&discount.Type,
		&discount.Method,
		&discount.Value,
		&discount.MinOrderValue,
		&discount.MaxDiscountValue,
		&productIDsJSON,
		&categoryIDsJSON,
		&discount.StartDate,
		&discount.EndDate,
		&discount.UsageLimit,
		&discount.CurrentUsage,
		&discount.Active,
		&discount.CreatedAt,
		&discount.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("discount not found")
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal product IDs
	if err := json.Unmarshal(productIDsJSON, &discount.ProductIDs); err != nil {
		return nil, err
	}

	// Unmarshal category IDs
	if err := json.Unmarshal(categoryIDsJSON, &discount.CategoryIDs); err != nil {
		return nil, err
	}

	return discount, nil
}

// GetByCode retrieves a discount by code
func (r *DiscountRepository) GetByCode(code string) (*entity.Discount, error) {
	query := `
		SELECT id, code, type, method, value, min_order_value, max_discount_value, 
			product_ids, category_ids, start_date, end_date, 
			usage_limit, current_usage, active, created_at, updated_at
		FROM discounts
		WHERE code = $1
	`

	var productIDsJSON, categoryIDsJSON []byte
	discount := &entity.Discount{}

	err := r.db.QueryRow(query, code).Scan(
		&discount.ID,
		&discount.Code,
		&discount.Type,
		&discount.Method,
		&discount.Value,
		&discount.MinOrderValue,
		&discount.MaxDiscountValue,
		&productIDsJSON,
		&categoryIDsJSON,
		&discount.StartDate,
		&discount.EndDate,
		&discount.UsageLimit,
		&discount.CurrentUsage,
		&discount.Active,
		&discount.CreatedAt,
		&discount.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("discount not found")
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal product IDs
	if err := json.Unmarshal(productIDsJSON, &discount.ProductIDs); err != nil {
		return nil, err
	}

	// Unmarshal category IDs
	if err := json.Unmarshal(categoryIDsJSON, &discount.CategoryIDs); err != nil {
		return nil, err
	}

	return discount, nil
}

// Update updates a discount
func (r *DiscountRepository) Update(discount *entity.Discount) error {
	query := `
		UPDATE discounts
		SET code = $1, type = $2, method = $3, value = $4, min_order_value = $5, 
			max_discount_value = $6, product_ids = $7, category_ids = $8, 
			start_date = $9, end_date = $10, usage_limit = $11, 
			current_usage = $12, active = $13, updated_at = $14
		WHERE id = $15
	`

	productIDsJSON, err := json.Marshal(discount.ProductIDs)
	if err != nil {
		return err
	}

	categoryIDsJSON, err := json.Marshal(discount.CategoryIDs)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		discount.Code,
		discount.Type,
		discount.Method,
		discount.Value,
		discount.MinOrderValue,
		discount.MaxDiscountValue,
		productIDsJSON,
		categoryIDsJSON,
		discount.StartDate,
		discount.EndDate,
		discount.UsageLimit,
		discount.CurrentUsage,
		discount.Active,
		time.Now(),
		discount.ID,
	)

	return err
}

// Delete deletes a discount
func (r *DiscountRepository) Delete(discountID uint) error {
	query := `DELETE FROM discounts WHERE id = $1`
	_, err := r.db.Exec(query, discountID)
	return err
}

// List retrieves a list of discounts with pagination
func (r *DiscountRepository) List(offset, limit int) ([]*entity.Discount, error) {
	query := `
		SELECT id, code, type, method, value, min_order_value, max_discount_value, 
			product_ids, category_ids, start_date, end_date, 
			usage_limit, current_usage, active, created_at, updated_at
		FROM discounts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	discounts := []*entity.Discount{}
	for rows.Next() {
		var productIDsJSON, categoryIDsJSON []byte
		discount := &entity.Discount{}

		err := rows.Scan(
			&discount.ID,
			&discount.Code,
			&discount.Type,
			&discount.Method,
			&discount.Value,
			&discount.MinOrderValue,
			&discount.MaxDiscountValue,
			&productIDsJSON,
			&categoryIDsJSON,
			&discount.StartDate,
			&discount.EndDate,
			&discount.UsageLimit,
			&discount.CurrentUsage,
			&discount.Active,
			&discount.CreatedAt,
			&discount.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal product IDs
		if err := json.Unmarshal(productIDsJSON, &discount.ProductIDs); err != nil {
			return nil, err
		}

		// Unmarshal category IDs
		if err := json.Unmarshal(categoryIDsJSON, &discount.CategoryIDs); err != nil {
			return nil, err
		}

		discounts = append(discounts, discount)
	}

	return discounts, nil
}

// ListActive retrieves a list of active discounts with pagination
func (r *DiscountRepository) ListActive(offset, limit int) ([]*entity.Discount, error) {
	query := `
		SELECT id, code, type, method, value, min_order_value, max_discount_value, 
			product_ids, category_ids, start_date, end_date, 
			usage_limit, current_usage, active, created_at, updated_at
		FROM discounts
		WHERE active = true 
		AND start_date <= NOW() 
		AND end_date >= NOW()
		AND (usage_limit = 0 OR current_usage < usage_limit)
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	discounts := []*entity.Discount{}
	for rows.Next() {
		var productIDsJSON, categoryIDsJSON []byte
		discount := &entity.Discount{}

		err := rows.Scan(
			&discount.ID,
			&discount.Code,
			&discount.Type,
			&discount.Method,
			&discount.Value,
			&discount.MinOrderValue,
			&discount.MaxDiscountValue,
			&productIDsJSON,
			&categoryIDsJSON,
			&discount.StartDate,
			&discount.EndDate,
			&discount.UsageLimit,
			&discount.CurrentUsage,
			&discount.Active,
			&discount.CreatedAt,
			&discount.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal product IDs
		if err := json.Unmarshal(productIDsJSON, &discount.ProductIDs); err != nil {
			return nil, err
		}

		// Unmarshal category IDs
		if err := json.Unmarshal(categoryIDsJSON, &discount.CategoryIDs); err != nil {
			return nil, err
		}

		discounts = append(discounts, discount)
	}

	return discounts, nil
}

// IncrementUsage increments the usage count of a discount
func (r *DiscountRepository) IncrementUsage(discountID uint) error {
	query := `
		UPDATE discounts
		SET current_usage = current_usage + 1, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, time.Now(), discountID)
	return err
}
