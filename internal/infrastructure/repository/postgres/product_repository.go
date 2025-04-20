package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// ProductRepository implements the product repository interface using PostgreSQL
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(product *entity.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, category_id, seller_id, images, has_variants, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.SellerID,
		imagesJSON,
		product.HasVariants,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)
	if err != nil {
		return err
	}

	// Generate and set the product number
	product.SetProductNumber(product.ID)

	// Update the product with the generated product number
	_, err = r.db.Exec(
		"UPDATE products SET product_number = $1 WHERE id = $2",
		product.ProductNumber,
		product.ID,
	)

	return err
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(id uint) (*entity.Product, error) {
	query := `
		SELECT id, product_number, name, description, price, stock, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var imagesJSON []byte
	product := &entity.Product{}
	var productNumber sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&productNumber,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.SellerID,
		&imagesJSON,
		&product.HasVariants,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}

	if err != nil {
		return nil, err
	}

	// Set product number if valid
	if productNumber.Valid {
		product.ProductNumber = productNumber.String
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
		return nil, err
	}

	return product, nil
}

// GetByIDWithVariants retrieves a product by ID including its variants
func (r *ProductRepository) GetByIDWithVariants(id uint) (*entity.Product, error) {
	// First get the product
	product, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// If product has variants, fetch them
	if product.HasVariants {
		query := `
			SELECT id, product_id, sku, price, compare_price, stock, attributes, images, is_default, created_at, updated_at
			FROM product_variants
			WHERE product_id = $1
			ORDER BY is_default DESC, id ASC
		`

		rows, err := r.db.Query(query, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		product.Variants = []*entity.ProductVariant{}
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

			product.Variants = append(product.Variants, variant)
		}
	}

	return product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *entity.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, category_id = $5, images = $6, has_variants = $7, updated_at = $8, product_number = $9
		WHERE id = $10
	`

	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
		imagesJSON,
		product.HasVariants,
		time.Now(),
		product.ProductNumber,
		product.ID,
	)

	return err
}

// Delete deletes a product
func (r *ProductRepository) Delete(id uint) error {
	// Start a transaction to delete product and its variants
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

	// Delete variants first (if any)
	_, err = tx.Exec("DELETE FROM product_variants WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Delete the product
	_, err = tx.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves a list of products with pagination
func (r *ProductRepository) List(offset, limit int) ([]*entity.Product, error) {
	query := `
		SELECT id, product_number, name, description, price, stock, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*entity.Product{}
	for rows.Next() {
		var imagesJSON []byte
		product := &entity.Product{}
		var productNumber sql.NullString

		err := rows.Scan(
			&product.ID,
			&productNumber,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.SellerID,
			&imagesJSON,
			&product.HasVariants,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set product number if valid
		if productNumber.Valid {
			product.ProductNumber = productNumber.String
		}

		// Unmarshal images JSON
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// Search searches for products based on criteria
func (r *ProductRepository) Search(query string, categoryID uint, minPrice, maxPrice float64, offset, limit int) ([]*entity.Product, error) {
	// Build the SQL query dynamically based on search criteria
	sqlQuery := `
		SELECT id, product_number, name, description, price, stock, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Add search conditions
	if query != "" {
		sqlQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+query+"%")
		argCount++
	}

	if categoryID > 0 {
		sqlQuery += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, categoryID)
		argCount++
	}

	if minPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, minPrice)
		argCount++
	}

	if maxPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, maxPrice)
		argCount++
	}

	// Add pagination
	sqlQuery += " ORDER BY id LIMIT $" + strconv.Itoa(argCount) + " OFFSET $" + strconv.Itoa(argCount+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*entity.Product{}
	for rows.Next() {
		var imagesJSON []byte
		product := &entity.Product{}
		var productNumber sql.NullString

		err := rows.Scan(
			&product.ID,
			&productNumber,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.SellerID,
			&imagesJSON,
			&product.HasVariants,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set product number if valid
		if productNumber.Valid {
			product.ProductNumber = productNumber.String
		}

		// Unmarshal images JSON
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// GetBySeller retrieves products by seller ID
func (r *ProductRepository) GetBySeller(sellerID uint, offset, limit int) ([]*entity.Product, error) {
	query := `
		SELECT id, product_number, name, description, price, stock, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		WHERE seller_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, sellerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*entity.Product{}
	for rows.Next() {
		var imagesJSON []byte
		product := &entity.Product{}
		var productNumber sql.NullString

		err := rows.Scan(
			&product.ID,
			&productNumber,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.SellerID,
			&imagesJSON,
			&product.HasVariants,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set product number if valid
		if productNumber.Valid {
			product.ProductNumber = productNumber.String
		}

		// Unmarshal images JSON
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
