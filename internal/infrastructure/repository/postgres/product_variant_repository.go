package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// ProductVariantRepository implements the product variant repository interface using PostgreSQL
type ProductVariantRepository struct {
	db *sql.DB
}

// NewProductVariantRepository creates a new ProductVariantRepository
func NewProductVariantRepository(db *sql.DB) *ProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

// Create creates a new product variant
func (r *ProductVariantRepository) Create(variant *entity.ProductVariant) error {
	query := `
		INSERT INTO product_variants (
			product_id, sku, price, compare_price, stock, attributes, images, is_default, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return err
	}

	imagesJSON, err := json.Marshal(variant.Images)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		variant.ProductID,
		variant.SKU,
		variant.Price,
		variant.ComparePrice,
		variant.Stock,
		attributesJSON,
		imagesJSON,
		variant.IsDefault,
		variant.CreatedAt,
		variant.UpdatedAt,
	).Scan(&variant.ID)

	return err
}

// GetByID retrieves a product variant by ID
func (r *ProductVariantRepository) GetByID(id uint) (*entity.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, price, compare_price, stock, attributes, images, is_default, created_at, updated_at
		FROM product_variants
		WHERE id = $1
	`

	var attributesJSON, imagesJSON []byte
	variant := &entity.ProductVariant{}
	err := r.db.QueryRow(query, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.Price,
		&variant.ComparePrice,
		&variant.Stock,
		&attributesJSON,
		&imagesJSON,
		&variant.IsDefault,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product variant not found")
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal attributes JSON
	if err := json.Unmarshal(attributesJSON, &variant.Attributes); err != nil {
		return nil, err
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &variant.Images); err != nil {
		return nil, err
	}

	return variant, nil
}

// GetBySKU retrieves a product variant by SKU
func (r *ProductVariantRepository) GetBySKU(sku string) (*entity.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, price, compare_price, stock, attributes, images, is_default, created_at, updated_at
		FROM product_variants
		WHERE sku = $1
	`

	var attributesJSON, imagesJSON []byte
	variant := &entity.ProductVariant{}
	err := r.db.QueryRow(query, sku).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.Price,
		&variant.ComparePrice,
		&variant.Stock,
		&attributesJSON,
		&imagesJSON,
		&variant.IsDefault,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product variant not found")
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal attributes JSON
	if err := json.Unmarshal(attributesJSON, &variant.Attributes); err != nil {
		return nil, err
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &variant.Images); err != nil {
		return nil, err
	}

	return variant, nil
}

// GetByProduct retrieves all variants for a product
func (r *ProductVariantRepository) GetByProduct(productID uint) ([]*entity.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, price, compare_price, stock, attributes, images, is_default, created_at, updated_at
		FROM product_variants
		WHERE product_id = $1
		ORDER BY is_default DESC, id ASC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []*entity.ProductVariant{}
	for rows.Next() {
		var attributesJSON, imagesJSON []byte
		variant := &entity.ProductVariant{}
		err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.SKU,
			&variant.Price,
			&variant.ComparePrice,
			&variant.Stock,
			&attributesJSON,
			&imagesJSON,
			&variant.IsDefault,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal attributes JSON
		if err := json.Unmarshal(attributesJSON, &variant.Attributes); err != nil {
			return nil, err
		}

		// Unmarshal images JSON
		if err := json.Unmarshal(imagesJSON, &variant.Images); err != nil {
			return nil, err
		}

		variants = append(variants, variant)
	}

	return variants, nil
}

// Update updates a product variant
func (r *ProductVariantRepository) Update(variant *entity.ProductVariant) error {
	query := `
		UPDATE product_variants
		SET sku = $1, price = $2, compare_price = $3, stock = $4, attributes = $5, images = $6, is_default = $7, updated_at = $8
		WHERE id = $9
	`

	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return err
	}

	imagesJSON, err := json.Marshal(variant.Images)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		variant.SKU,
		variant.Price,
		variant.ComparePrice,
		variant.Stock,
		attributesJSON,
		imagesJSON,
		variant.IsDefault,
		time.Now(),
		variant.ID,
	)

	return err
}

// Delete deletes a product variant
func (r *ProductVariantRepository) Delete(id uint) error {
	query := `DELETE FROM product_variants WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// BatchCreate creates multiple product variants in a single transaction
func (r *ProductVariantRepository) BatchCreate(variants []*entity.ProductVariant) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		INSERT INTO product_variants (
			product_id, sku, price, compare_price, stock, attributes, images, is_default, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, variant := range variants {
		attributesJSON, err := json.Marshal(variant.Attributes)
		if err != nil {
			return err
		}

		imagesJSON, err := json.Marshal(variant.Images)
		if err != nil {
			return err
		}

		err = stmt.QueryRow(
			variant.ProductID,
			variant.SKU,
			variant.Price,
			variant.ComparePrice,
			variant.Stock,
			attributesJSON,
			imagesJSON,
			variant.IsDefault,
			variant.CreatedAt,
			variant.UpdatedAt,
		).Scan(&variant.ID)

		if err != nil {
			return err
		}
	}

	return nil
}
