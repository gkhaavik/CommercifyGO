package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductRepository is the PostgreSQL implementation of the ProductRepository interface
type ProductRepository struct {
	db                *sql.DB
	variantRepository repository.ProductVariantRepository
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *sql.DB, variantRepository repository.ProductVariantRepository) repository.ProductRepository {
	return &ProductRepository{
		db:                db,
		variantRepository: variantRepository,
	}
}

// Create creates a new product
func (r *ProductRepository) Create(product *entity.Product) error {
	query := `
	INSERT INTO products (name, description, price, stock, weight, category_id, seller_id, images, has_variants, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		product.Weight,
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

	// Update the product number in the database
	updateQuery := "UPDATE products SET product_number = $1 WHERE id = $2"
	_, err = r.db.Exec(updateQuery, product.ProductNumber, product.ID)
	if err != nil {
		return err
	}

	// If the product has currency-specific prices, save them
	if len(product.Prices) > 0 {
		for i := range product.Prices {
			product.Prices[i].ProductID = product.ID
			if err = r.createProductPrice(&product.Prices[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// createProductPrice creates a product price entry for a specific currency
func (r *ProductRepository) createProductPrice(price *entity.ProductPrice) error {
	query := `
		INSERT INTO product_prices (product_id, currency_code, price, compare_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (product_id, currency_code) DO UPDATE SET
		price = EXCLUDED.price,
		compare_price = EXCLUDED.compare_price,
		updated_at = EXCLUDED.updated_at
		RETURNING id
		`

	var comparePrice sql.NullInt64
	if price.ComparePrice > 0 {
		comparePrice.Int64 = price.ComparePrice
		comparePrice.Valid = true
	}

	now := time.Now()

	return r.db.QueryRow(
		query,
		price.ProductID,
		price.CurrencyCode,
		price.Price,
		comparePrice,
		now,
		now,
	).Scan(&price.ID)
}

// GetByID gets a product by ID
func (r *ProductRepository) GetByID(id uint) (*entity.Product, error) {
	query := `
			SELECT id, product_number, name, description, price, stock, weight, category_id, seller_id, images, has_variants, created_at, updated_at
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
		&product.Weight,
		&product.CategoryID,
		&product.SellerID,
		&imagesJSON,
		&product.HasVariants,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
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

	// Load currency-specific prices
	prices, err := r.getProductPrices(product.ID)
	if err != nil {
		return nil, err
	}
	product.Prices = prices

	return product, nil
}

// getProductPrices retrieves all prices for a product in different currencies
func (r *ProductRepository) getProductPrices(productID uint) ([]entity.ProductPrice, error) {
	query := `
			SELECT id, product_id, currency_code, price, compare_price, created_at, updated_at
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
		var comparePrice sql.NullInt64

		err := rows.Scan(
			&price.ID,
			&price.ProductID,
			&price.CurrencyCode,
			&price.Price,
			&comparePrice,
			&price.CreatedAt,
			&price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if comparePrice.Valid {
			price.ComparePrice = comparePrice.Int64
		}

		prices = append(prices, price)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

// GetByIDWithVariants gets a product by ID with variants
func (r *ProductRepository) GetByIDWithVariants(productId uint) (*entity.Product, error) {
	// Get the base product
	product, err := r.GetByID(productId)
	if err != nil {
		return nil, err
	}

	// If product has variants, get them
	if product.HasVariants {
		variants, err := r.variantRepository.GetByProduct(productId)
		if err != nil {
			return nil, err
		}

		product.Variants = variants
	}

	return product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *entity.Product) error {
	query := `
			UPDATE products
			SET name = $1, description = $2, price = $3, stock = $4, weight = $5, category_id = $6, 
		    images = $7, has_variants = $8, updated_at = $9
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
		product.Weight,
		product.CategoryID,
		imagesJSON,
		product.HasVariants,
		time.Now(),
		product.ID,
	)
	if err != nil {
		return err
	}

	// Update currency-specific prices if they exist
	if len(product.Prices) > 0 {
		// Use an upsert query to update or insert prices
		query := `
			INSERT INTO product_prices (product_id, currency_code, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (product_id, currency_code)
			DO UPDATE SET price = EXCLUDED.price
		`
		for _, price := range product.Prices {
			_, err := r.db.Exec(query, product.ID, price.CurrencyCode, price.Price)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete deletes a product
func (r *ProductRepository) Delete(id uint) error {
	// Start a transaction to delete variants as well
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete variants first
	_, err = tx.Exec("DELETE FROM product_variants WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Delete the product
	_, err = tx.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// List lists products with pagination
func (r *ProductRepository) List(offset, limit int) ([]*entity.Product, error) {
	query := `
		SELECT id, product_number, name, description, price, stock, weight, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
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
			&product.Weight,
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

		// Load currency-specific prices
		prices, err := r.getProductPrices(product.ID)
		if err != nil {
			return nil, err
		}
		product.Prices = prices

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Search searches for products based on criteria (prices in cents)
func (r *ProductRepository) Search(query string, categoryID uint, minPriceCents, maxPriceCents int64, offset, limit int) ([]*entity.Product, error) {
	// Build dynamic query parts
	searchQuery := `
		SELECT id, product_number, name, description, price, stock, weight, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	queryParams := []interface{}{}
	paramCounter := 1

	if query != "" {
		searchQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", paramCounter, paramCounter)
		queryParams = append(queryParams, "%"+query+"%")
		paramCounter++
	}

	if categoryID > 0 {
		searchQuery += fmt.Sprintf(" AND category_id = $%d", paramCounter)
		queryParams = append(queryParams, categoryID)
		paramCounter++
	}

	if minPriceCents > 0 {
		searchQuery += fmt.Sprintf(" AND price >= $%d", paramCounter)
		queryParams = append(queryParams, minPriceCents) // Use cents
		paramCounter++
	}

	if maxPriceCents > 0 {
		searchQuery += fmt.Sprintf(" AND price <= $%d", paramCounter)
		queryParams = append(queryParams, maxPriceCents) // Use cents
		paramCounter++
	}

	// Add pagination
	searchQuery += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(paramCounter) + " OFFSET $" + strconv.Itoa(paramCounter+1)
	queryParams = append(queryParams, limit, offset)

	// Execute query
	rows, err := r.db.Query(searchQuery, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse results
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
			&product.Price, // Reads int64 directly
			&product.Stock,
			&product.Weight,
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

		// Load currency-specific prices
		prices, err := r.getProductPrices(product.ID)
		if err != nil {
			return nil, err
		}
		product.Prices = prices

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) Count() (int, error) {
	query := `
		SELECT COUNT(*) FROM products
	`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductRepository) CountBySeller(sellerID uint) (int, error) {
	query := `
		SELECT COUNT(*) FROM products
		WHERE seller_id = $1
	`

	var count int
	err := r.db.QueryRow(query, sellerID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductRepository) CountSearch(searchQuery string, categoryID uint, minPriceCents, maxPriceCents int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM products
		WHERE 1=1
	`

	queryParams := []any{}
	paramCounter := 1

	if searchQuery != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", paramCounter, paramCounter)
		queryParams = append(queryParams, "%"+searchQuery+"%")
		paramCounter++
	}

	if categoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", paramCounter)
		queryParams = append(queryParams, categoryID)
		paramCounter++
	}

	if minPriceCents > 0 {
		query += fmt.Sprintf(" AND price >= $%d", paramCounter)
		queryParams = append(queryParams, minPriceCents)
		paramCounter++
	}

	if maxPriceCents > 0 {
		query += fmt.Sprintf(" AND price <= $%d", paramCounter)
		queryParams = append(queryParams, maxPriceCents)
		paramCounter++
	}

	var count int
	err := r.db.QueryRow(query, queryParams...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetBySeller gets products by seller ID with pagination
func (r *ProductRepository) GetBySeller(sellerID uint, offset, limit int) ([]*entity.Product, error) {
	query := `
		SELECT id, product_number, name, description, price, stock, weight, category_id, seller_id, images, has_variants, created_at, updated_at
		FROM products
		WHERE seller_id = $1
		ORDER BY created_at DESC
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
			&product.Weight,
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

		// Load currency-specific prices
		prices, err := r.getProductPrices(product.ID)
		if err != nil {
			return nil, err
		}
		product.Prices = prices

		products = append(products, product)
	}

	return products, nil
}
