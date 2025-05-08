package usecase

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductUseCase implements product-related use cases
type ProductUseCase struct {
	productRepo        repository.ProductRepository
	categoryRepo       repository.CategoryRepository
	productVariantRepo repository.ProductVariantRepository
	currencyRepo       repository.CurrencyRepository
}

// NewProductUseCase creates a new ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	productVariantRepo repository.ProductVariantRepository,
	currencyRepo repository.CurrencyRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:        productRepo,
		categoryRepo:       categoryRepo,
		productVariantRepo: productVariantRepo,
		currencyRepo:       currencyRepo,
	}
}

// CurrencyPriceInput represents a price in a specific currency
type CurrencyPriceInput struct {
	CurrencyCode string  `json:"currency_code"`
	Price        float64 `json:"price"`         // Price in dollars/currency unit
	ComparePrice float64 `json:"compare_price"` // Compare price in dollars/currency unit
}

// CreateProductInput contains the data needed to create a product (prices in dollars)
type CreateProductInput struct {
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Price          float64              `json:"price"`         // Price in dollars (default currency)
	ComparePrice   float64              `json:"compare_price"` // Compare price in dollars (default currency)
	Stock          int                  `json:"stock"`
	Weight         float64              `json:"weight"`
	CategoryID     uint                 `json:"category_id"`
	SellerID       uint                 `json:"seller_id"`
	Images         []string             `json:"images"`
	HasVariants    bool                 `json:"has_variants"`
	CurrencyPrices []CurrencyPriceInput `json:"currency_prices"` // Prices in other currencies
	Variants       []CreateVariantInput `json:"variants"`
}

// CreateVariantInput contains the data needed to create a product variant (prices in dollars)
type CreateVariantInput struct {
	SKU            string                    `json:"sku"`
	Price          float64                   `json:"price"`         // Price in dollars (default currency)
	ComparePrice   float64                   `json:"compare_price"` // Price in dollars (default currency)
	Stock          int                       `json:"stock"`
	Attributes     []entity.VariantAttribute `json:"attributes"`
	Images         []string                  `json:"images"`
	IsDefault      bool                      `json:"is_default"`
	CurrencyPrices []CurrencyPriceInput      `json:"currency_prices"` // Prices in other currencies
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(input CreateProductInput) (*entity.Product, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Convert price to cents
	priceCents := money.ToCents(input.Price)

	// Create product
	product, err := entity.NewProduct(
		input.Name,
		input.Description,
		priceCents, // Use cents
		input.Stock,
		input.Weight,
		input.CategoryID,
		input.SellerID,
		input.Images,
	)
	if err != nil {
		return nil, err
	}

	// Set has_variants flag
	product.HasVariants = input.HasVariants

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		product.Prices = make([]entity.ProductPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)
			comparePriceCents := money.ToCents(currPrice.ComparePrice)

			product.Prices = append(product.Prices, entity.ProductPrice{
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
				ComparePrice: comparePriceCents,
			})
		}
	}

	// Save product
	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	// If product has variants, create them
	if input.HasVariants && len(input.Variants) > 0 {
		variants := make([]*entity.ProductVariant, 0, len(input.Variants))

		for _, variantInput := range input.Variants {
			// Convert variant prices to cents
			variantPriceCents := money.ToCents(variantInput.Price)
			variantComparePriceCents := money.ToCents(variantInput.ComparePrice)

			variant, err := entity.NewProductVariant(
				product.ID,
				variantInput.SKU,
				variantPriceCents, // Use cents
				variantInput.Stock,
				variantInput.Attributes,
				variantInput.Images,
				variantInput.IsDefault,
			)
			if err != nil {
				return nil, err
			}

			if variantInput.ComparePrice > 0 {
				if err := variant.SetComparePrice(variantComparePriceCents); err != nil { // Use cents
					return nil, err
				}
			}

			// Process currency-specific prices for variant, if any
			if len(variantInput.CurrencyPrices) > 0 {
				variant.Prices = make([]entity.ProductVariantPrice, 0, len(variantInput.CurrencyPrices))

				for _, currPrice := range variantInput.CurrencyPrices {
					// Validate currency exists
					_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
					if err != nil {
						return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
					}

					// Convert price to cents
					priceCents := money.ToCents(currPrice.Price)
					comparePriceCents := money.ToCents(currPrice.ComparePrice)

					variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
						CurrencyCode: currPrice.CurrencyCode,
						Price:        priceCents,
						ComparePrice: comparePriceCents,
					})
				}
			}

			variants = append(variants, variant)
		}

		// Save each variant individually to process their currency prices too
		for _, variant := range variants {
			if err := uc.productVariantRepo.Create(variant); err != nil {
				return nil, err
			}
		}

		// Add variants to product
		product.Variants = variants
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (uc *ProductUseCase) GetProductByID(id uint) (*entity.Product, error) {
	// If product has variants, get them too
	return uc.productRepo.GetByIDWithVariants(id)
}

// GetProductByCurrency retrieves a product by ID with prices in a specific currency
func (uc *ProductUseCase) GetProductByCurrency(id uint, currencyCode string) (*entity.Product, error) {
	// First get the product with all its data
	product, err := uc.productRepo.GetByIDWithVariants(id)
	if err != nil {
		return nil, err
	}

	// If no specific currency requested, return as is
	if currencyCode == "" {
		defaultCurr, err := uc.currencyRepo.GetDefault()
		if err != nil {
			return nil, err
		}

		currencyCode = defaultCurr.Code
	}

	// Validate currency exists
	currency, err := uc.currencyRepo.GetByCode(currencyCode)
	if err != nil {
		return nil, errors.New("invalid currency code: " + currencyCode)
	}

	currencyPrice, found := product.GetPriceInCurrency(currency.Code)
	if !found {
		return nil, errors.New("product not available in the requested currency")
	}

	// comparePrice, found := product.GetComparePriceInCurrency(currency.Code)
	// if found {
	// 	product.ComparePrice = comparePrice
	// }

	product.Price = currencyPrice

	return product, nil
}

// UpdateProductInput contains the data needed to update a product (prices in dollars)
type UpdateProductInput struct {
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Price          float64              `json:"price"`         // Price in dollars (default currency)
	ComparePrice   float64              `json:"compare_price"` // Compare price in dollars (default currency)
	Stock          int                  `json:"stock"`
	CategoryID     uint                 `json:"category_id"`
	Images         []string             `json:"images"`
	CurrencyPrices []CurrencyPriceInput `json:"currency_prices"` // Prices in other currencies
}

// UpdateProduct updates a product
func (uc *ProductUseCase) UpdateProduct(id uint, sellerID uint, input UpdateProductInput) (*entity.Product, error) {
	// Get product
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not the seller of this product")
	}

	// Validate category exists if changing
	if input.CategoryID != 0 && input.CategoryID != product.CategoryID {
		_, err := uc.categoryRepo.GetByID(input.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = input.CategoryID
	}

	// Update product fields
	if input.Name != "" {
		product.Name = input.Name
	}
	if input.Description != "" {
		product.Description = input.Description
	}
	if input.Price > 0 && !product.HasVariants {
		product.Price = money.ToCents(input.Price) // Convert to cents
	}
	if input.Stock >= 0 && !product.HasVariants {
		product.Stock = input.Stock
	}
	if len(input.Images) > 0 {
		product.Images = input.Images
	}

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		// Clear existing prices
		product.Prices = make([]entity.ProductPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)
			comparePriceCents := money.ToCents(currPrice.ComparePrice)

			product.Prices = append(product.Prices, entity.ProductPrice{
				ProductID:    product.ID,
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
				ComparePrice: comparePriceCents,
			})
		}
	}

	// Update product in repository
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	// If product has variants, get them
	if product.HasVariants {
		variants, err := uc.productVariantRepo.GetByProduct(product.ID)
		if err != nil {
			return nil, err
		}
		product.Variants = variants
	}

	return product, nil
}

// UpdateVariantInput contains the data needed to update a product variant (prices in dollars)
type UpdateVariantInput struct {
	SKU            string                    `json:"sku"`
	Price          float64                   `json:"price"`         // Price in dollars
	ComparePrice   float64                   `json:"compare_price"` // Price in dollars
	Stock          int                       `json:"stock"`
	Attributes     []entity.VariantAttribute `json:"attributes"`
	Images         []string                  `json:"images"`
	IsDefault      bool                      `json:"is_default"`
	CurrencyPrices []CurrencyPriceInput      `json:"currency_prices"` // Prices in other currencies
}

// UpdateVariant updates a product variant
func (uc *ProductUseCase) UpdateVariant(productID uint, variantID uint, sellerID uint, input UpdateVariantInput) (*entity.ProductVariant, error) {
	// Get product to check ownership
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not the seller of this product")
	}

	// Get variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return nil, err
	}

	// Check if variant belongs to the product
	if variant.ProductID != productID {
		return nil, errors.New("variant does not belong to this product")
	}

	// Update variant fields
	if input.SKU != "" {
		variant.SKU = input.SKU
	}
	if input.Price > 0 {
		variant.Price = money.ToCents(input.Price) // Convert to cents
	}
	if input.ComparePrice > 0 {
		variant.ComparePrice = money.ToCents(input.ComparePrice) // Convert to cents
	}
	if input.Stock >= 0 {
		variant.Stock = input.Stock
	}
	if len(input.Attributes) > 0 {
		variant.Attributes = input.Attributes
	}
	if len(input.Images) > 0 {
		variant.Images = input.Images
	}

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		// Clear existing prices
		variant.Prices = make([]entity.ProductVariantPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)
			comparePriceCents := money.ToCents(currPrice.ComparePrice)

			variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
				VariantID:    variant.ID,
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
				ComparePrice: comparePriceCents,
			})
		}
	}

	// Handle default status
	if input.IsDefault != variant.IsDefault {
		// If setting this variant as default, unset any other default variants
		if input.IsDefault {
			variants, err := uc.productVariantRepo.GetByProduct(productID)
			if err != nil {
				return nil, err
			}

			for _, v := range variants {
				if v.ID != variantID && v.IsDefault {
					v.IsDefault = false
					if err := uc.productVariantRepo.Update(v); err != nil {
						return nil, err
					}
				}
			}
		}

		variant.IsDefault = input.IsDefault
	}

	// Update variant in repository
	if err := uc.productVariantRepo.Update(variant); err != nil {
		return nil, err
	}

	// If this is the default variant, update the product price
	if variant.IsDefault {
		product.Price = variant.Price // Already in cents
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	return variant, nil
}

// AddVariantInput contains the data needed to add a variant to a product (prices in dollars)
type AddVariantInput struct {
	ProductID      uint                      `json:"product_id"`
	SKU            string                    `json:"sku"`
	Price          float64                   `json:"price"`         // Price in dollars
	ComparePrice   float64                   `json:"compare_price"` // Price in dollars
	Stock          int                       `json:"stock"`
	Attributes     []entity.VariantAttribute `json:"attributes"`
	Images         []string                  `json:"images"`
	IsDefault      bool                      `json:"is_default"`
	CurrencyPrices []CurrencyPriceInput      `json:"currency_prices"` // Prices in other currencies
}

// AddVariant adds a new variant to a product
func (uc *ProductUseCase) AddVariant(sellerID uint, input AddVariantInput) (*entity.ProductVariant, error) {
	// Get product to check ownership
	product, err := uc.productRepo.GetByID(input.ProductID)
	if err != nil {
		return nil, err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not the seller of this product")
	}

	// Convert prices to cents
	priceCents := money.ToCents(input.Price)
	comparePriceCents := money.ToCents(input.ComparePrice)

	// Create variant
	variant, err := entity.NewProductVariant(
		input.ProductID,
		input.SKU,
		priceCents, // Use cents
		input.Stock,
		input.Attributes,
		input.Images,
		input.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	if input.ComparePrice > 0 {
		if err := variant.SetComparePrice(comparePriceCents); err != nil { // Use cents
			return nil, err
		}
	}

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		variant.Prices = make([]entity.ProductVariantPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)
			comparePriceCents := money.ToCents(currPrice.ComparePrice)

			variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
				ComparePrice: comparePriceCents,
			})
		}
	}

	// If this is the first variant or it's set as default, update product
	isFirstVariant := !product.HasVariants
	if isFirstVariant || input.IsDefault {
		// If this is the first variant, set product to have variants
		if isFirstVariant {
			product.HasVariants = true
		}

		// If this is the default variant, update product price and weight
		if input.IsDefault {
			product.Price = priceCents // Update product price with cents

			// If there are other variants, unset their default status
			if !isFirstVariant {
				variants, err := uc.productVariantRepo.GetByProduct(input.ProductID)
				if err != nil {
					return nil, err
				}

				for _, v := range variants {
					if v.IsDefault {
						v.IsDefault = false
						if err := uc.productVariantRepo.Update(v); err != nil {
							return nil, err
						}
					}
				}
			}
		}

		// Update product
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	// Save variant
	if err := uc.productVariantRepo.Create(variant); err != nil {
		return nil, err
	}

	return variant, nil
}

// DeleteVariant deletes a product variant
func (uc *ProductUseCase) DeleteVariant(productID uint, variantID uint, sellerID uint) error {
	// Get product to check ownership
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return errors.New("unauthorized: not the seller of this product")
	}

	// Get variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return err
	}

	// Check if variant belongs to the product
	if variant.ProductID != productID {
		return errors.New("variant does not belong to this product")
	}

	// Get all variants for the product
	variants, err := uc.productVariantRepo.GetByProduct(productID)
	if err != nil {
		return err
	}

	// Check if this is the only variant
	if len(variants) == 1 {
		return errors.New("cannot delete the only variant of a product")
	}

	// If this is the default variant, set another variant as default
	if variant.IsDefault {
		for _, v := range variants {
			if v.ID != variantID {
				v.IsDefault = true
				product.Price = v.Price

				// Update the new default variant
				if err := uc.productVariantRepo.Update(v); err != nil {
					return err
				}

				// Update product price
				if err := uc.productRepo.Update(product); err != nil {
					return err
				}

				break
			}
		}
	}

	// Delete variant
	return uc.productVariantRepo.Delete(variantID)
}

// DeleteProduct deletes a product
func (uc *ProductUseCase) DeleteProduct(id uint, sellerID uint) error {
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return errors.New("unauthorized: not the seller of this product")
	}

	// Delete product (this will cascade delete variants)
	return uc.productRepo.Delete(id)
}

// SearchProductsInput contains the data needed to search for products (prices in dollars)
type SearchProductsInput struct {
	Query        string  `json:"query"`
	CategoryID   uint    `json:"category_id"`
	MinPrice     float64 `json:"min_price"`     // Price in dollars
	MaxPrice     float64 `json:"max_price"`     // Price in dollars
	CurrencyCode string  `json:"currency_code"` // Optional currency code for prices
	Offset       int     `json:"offset"`
	Limit        int     `json:"limit"`
}

// SearchProducts searches for products based on criteria
func (uc *ProductUseCase) SearchProducts(input SearchProductsInput) ([]*entity.Product, error) {
	// If currency is specified and not the default, convert price ranges
	var minPriceCents, maxPriceCents int64

	// TODO: Default currency should be in memory
	defaultCurr, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return nil, err
	}

	if input.CurrencyCode != "" && input.CurrencyCode != defaultCurr.Code {
		// Get the currency
		currency, err := uc.currencyRepo.GetByCode(input.CurrencyCode)
		if err != nil {
			return nil, errors.New("invalid currency code: " + input.CurrencyCode)
		}

		// Convert min/max prices to default currency using exchange rate
		defaultPrice := input.MinPrice / currency.ExchangeRate
		minPriceCents = money.ToCents(defaultPrice)

		defaultPrice = input.MaxPrice / currency.ExchangeRate
		maxPriceCents = money.ToCents(defaultPrice)
	} else {
		// Convert min/max prices to cents for repository search
		minPriceCents = money.ToCents(input.MinPrice)
		maxPriceCents = money.ToCents(input.MaxPrice)
	}

	return uc.productRepo.Search(
		input.Query,
		input.CategoryID,
		minPriceCents, // Pass cents
		maxPriceCents, // Pass cents
		input.Offset,
		input.Limit,
	)
}

// ListProductsBySeller lists products by seller
func (uc *ProductUseCase) ListProductsBySeller(sellerID uint, offset, limit int) ([]*entity.Product, error) {
	return uc.productRepo.GetBySeller(sellerID, offset, limit)
}

// ListProducts lists all products with pagination
func (uc *ProductUseCase) ListProducts(offset, limit int) ([]*entity.Product, error) {
	return uc.productRepo.List(offset, limit)
}

// ListCategories lists all product categories
func (uc *ProductUseCase) ListCategories() ([]*entity.Category, error) {
	return uc.categoryRepo.List()
}

// SetProductCurrencyPrices sets currency-specific prices for a product
func (uc *ProductUseCase) SetProductCurrencyPrices(productID uint, sellerID uint, currencyPrices []CurrencyPriceInput) error {
	// Get product to check ownership
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return errors.New("unauthorized: not the seller of this product")
	}

	// Clear existing currency prices
	product.Prices = make([]entity.ProductPrice, 0, len(currencyPrices))

	// Add new currency prices
	for _, currPrice := range currencyPrices {
		// Validate currency exists
		_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
		if err != nil {
			return errors.New("invalid currency code: " + currPrice.CurrencyCode)
		}

		// Convert prices to cents
		priceCents := money.ToCents(currPrice.Price)
		comparePriceCents := money.ToCents(currPrice.ComparePrice)

		product.Prices = append(product.Prices, entity.ProductPrice{
			ProductID:    productID,
			CurrencyCode: currPrice.CurrencyCode,
			Price:        priceCents,
			ComparePrice: comparePriceCents,
		})
	}

	// Update product in repository
	return uc.productRepo.Update(product)
}

// SetVariantCurrencyPrices sets currency-specific prices for a product variant
func (uc *ProductUseCase) SetVariantCurrencyPrices(productID uint, variantID uint, sellerID uint, currencyPrices []CurrencyPriceInput) error {
	// Get product to check ownership
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	// Check if user is the seller of the product
	if product.SellerID != sellerID {
		return errors.New("unauthorized: not the seller of this product")
	}

	// Get variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return err
	}

	// Check if variant belongs to the product
	if variant.ProductID != productID {
		return errors.New("variant does not belong to this product")
	}

	// Clear existing currency prices
	variant.Prices = make([]entity.ProductVariantPrice, 0, len(currencyPrices))

	// Add new currency prices
	for _, currPrice := range currencyPrices {
		// Validate currency exists
		_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
		if err != nil {
			return errors.New("invalid currency code: " + currPrice.CurrencyCode)
		}

		// Convert prices to cents
		priceCents := money.ToCents(currPrice.Price)
		comparePriceCents := money.ToCents(currPrice.ComparePrice)

		variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
			VariantID:    variantID,
			CurrencyCode: currPrice.CurrencyCode,
			Price:        priceCents,
			ComparePrice: comparePriceCents,
		})
	}

	// Update variant in repository
	return uc.productVariantRepo.Update(variant)
}
