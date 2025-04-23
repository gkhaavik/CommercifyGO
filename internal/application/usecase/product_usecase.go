package usecase

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductUseCase implements product-related use cases
type ProductUseCase struct {
	productRepo        repository.ProductRepository
	categoryRepo       repository.CategoryRepository
	productVariantRepo repository.ProductVariantRepository
}

// NewProductUseCase creates a new ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	productVariantRepo repository.ProductVariantRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:        productRepo,
		categoryRepo:       categoryRepo,
		productVariantRepo: productVariantRepo,
	}
}

// CreateProductInput contains the data needed to create a product
type CreateProductInput struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Stock       int                  `json:"stock"`
	CategoryID  uint                 `json:"category_id"`
	SellerID    uint                 `json:"seller_id"`
	Images      []string             `json:"images"`
	HasVariants bool                 `json:"has_variants"`
	Variants    []CreateVariantInput `json:"variants"`
}

// CreateVariantInput contains the data needed to create a product variant
type CreateVariantInput struct {
	SKU          string                    `json:"sku"`
	Price        float64                   `json:"price"`
	ComparePrice float64                   `json:"compare_price"`
	Stock        int                       `json:"stock"`
	Attributes   []entity.VariantAttribute `json:"attributes"`
	Images       []string                  `json:"images"`
	IsDefault    bool                      `json:"is_default"`
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(input CreateProductInput) (*entity.Product, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Create product
	product, err := entity.NewProduct(
		input.Name,
		input.Description,
		input.Price,
		input.Stock,
		input.CategoryID,
		input.SellerID,
		input.Images,
	)
	if err != nil {
		return nil, err
	}

	// Set has_variants flag
	product.HasVariants = input.HasVariants

	// Save product
	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	// If product has variants, create them
	if input.HasVariants && len(input.Variants) > 0 {
		variants := make([]*entity.ProductVariant, 0, len(input.Variants))

		for _, variantInput := range input.Variants {
			variant, err := entity.NewProductVariant(
				product.ID,
				variantInput.SKU,
				variantInput.Price,
				variantInput.Stock,
				variantInput.Attributes,
				variantInput.Images,
				variantInput.IsDefault,
			)
			if err != nil {
				return nil, err
			}

			if variantInput.ComparePrice > 0 {
				if err := variant.SetComparePrice(variantInput.ComparePrice); err != nil {
					return nil, err
				}
			}

			variants = append(variants, variant)
		}

		// Save variants in batch
		if err := uc.productVariantRepo.BatchCreate(variants); err != nil {
			return nil, err
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

// UpdateProductInput contains the data needed to update a product
type UpdateProductInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	CategoryID  uint     `json:"category_id"`
	Images      []string `json:"images"`
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
		product.Price = input.Price
	}
	if input.Stock >= 0 && !product.HasVariants {
		product.Stock = input.Stock
	}
	if len(input.Images) > 0 {
		product.Images = input.Images
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

// UpdateVariantInput contains the data needed to update a product variant
type UpdateVariantInput struct {
	SKU          string                    `json:"sku"`
	Price        float64                   `json:"price"`
	ComparePrice float64                   `json:"compare_price"`
	Stock        int                       `json:"stock"`
	Attributes   []entity.VariantAttribute `json:"attributes"`
	Images       []string                  `json:"images"`
	IsDefault    bool                      `json:"is_default"`
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
		variant.Price = input.Price
	}
	if input.ComparePrice > 0 {
		variant.ComparePrice = input.ComparePrice
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
		product.Price = variant.Price
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	return variant, nil
}

// AddVariantInput contains the data needed to add a variant to a product
type AddVariantInput struct {
	ProductID    uint                      `json:"product_id"`
	SKU          string                    `json:"sku"`
	Price        float64                   `json:"price"`
	ComparePrice float64                   `json:"compare_price"`
	Stock        int                       `json:"stock"`
	Attributes   []entity.VariantAttribute `json:"attributes"`
	Images       []string                  `json:"images"`
	IsDefault    bool                      `json:"is_default"`
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

	// Create variant
	variant, err := entity.NewProductVariant(
		input.ProductID,
		input.SKU,
		input.Price,
		input.Stock,
		input.Attributes,
		input.Images,
		input.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	if input.ComparePrice > 0 {
		if err := variant.SetComparePrice(input.ComparePrice); err != nil {
			return nil, err
		}
	}

	// If this is the first variant or it's set as default, update product
	isFirstVariant := !product.HasVariants
	if isFirstVariant || input.IsDefault {
		// If this is the first variant, set product to have variants
		if isFirstVariant {
			product.HasVariants = true
		}

		// If this is the default variant, update product price
		if input.IsDefault {
			product.Price = input.Price

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
	// Get product
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

// SearchProductsInput contains the data needed to search for products
type SearchProductsInput struct {
	Query      string  `json:"query"`
	CategoryID uint    `json:"category_id"`
	MinPrice   float64 `json:"min_price"`
	MaxPrice   float64 `json:"max_price"`
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
}

// SearchProducts searches for products based on criteria
func (uc *ProductUseCase) SearchProducts(input SearchProductsInput) ([]*entity.Product, error) {
	return uc.productRepo.Search(
		input.Query,
		input.CategoryID,
		input.MinPrice,
		input.MaxPrice,
		input.Offset,
		input.Limit,
	)
}

// ListProductsBySeller lists products by seller
func (uc *ProductUseCase) ListProductsBySeller(sellerID uint, offset, limit int) ([]*entity.Product, error) {
	return uc.productRepo.GetBySeller(sellerID, offset, limit)
}

func (uc *ProductUseCase) ListProducts(offset, limit int) ([]*entity.Product, error) {
	return uc.productRepo.List(offset, limit)
}

// ListProductsByCategory lists products by category
func (uc *ProductUseCase) ListCategories() ([]*entity.Category, error) {
	return uc.categoryRepo.List()
}

// ListProductsByCategoryAndSeller lists products by category and seller
