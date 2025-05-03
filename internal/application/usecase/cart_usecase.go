package usecase

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CartUseCase implements cart-related use cases
type CartUseCase struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

// NewCartUseCase creates a new CartUseCase
func NewCartUseCase(cartRepo repository.CartRepository, productRepo repository.ProductRepository) *CartUseCase {
	return &CartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetOrCreateCart gets a user's cart or creates one if it doesn't exist
func (uc *CartUseCase) GetOrCreateCart(userID uint) (*entity.Cart, error) {
	cart, err := uc.cartRepo.GetByUserID(userID)
	if err == nil {
		return cart, nil
	}

	// Create new cart if not found
	newCart, err := entity.NewCart(userID)
	if err != nil {
		return nil, err
	}

	if err := uc.cartRepo.Create(newCart); err != nil {
		return nil, err
	}

	return newCart, nil
}

// AddToCartInput contains the data needed to add an item to a cart
type AddToCartInput struct {
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"` // Added variant ID
	Quantity  int  `json:"quantity"`
}

// AddToCart adds a product to a user's cart
func (uc *CartUseCase) AddToCart(userID uint, input AddToCartInput) (*entity.Cart, error) {
	// Check if product exists and has enough stock
	product, err := uc.productRepo.GetByIDWithVariants(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock based on whether it's a variant or a regular product
	isVariant := input.VariantID > 0
	if isVariant {
		variant := product.GetVariantByID(input.VariantID)

		if variant == nil {
			return nil, errors.New("product variant not found")
		}

		if !variant.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}

	} else {
		// Regular product
		if product.HasVariants {
			return nil, errors.New("please select a product variant")
		}

		if !product.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	}

	// Get cart
	cart, err := uc.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	if err := cart.AddItem(input.ProductID, input.VariantID, input.Quantity); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// UpdateCartItemInput contains the data needed to update a cart item
type UpdateCartItemInput struct {
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"` // Added variant ID
	Quantity  int  `json:"quantity"`
}

// UpdateCartItem updates the quantity of a product in a user's cart
func (uc *CartUseCase) UpdateCartItem(userID uint, input UpdateCartItemInput) (*entity.Cart, error) {
	// Check if product exists and has enough stock
	product, err := uc.productRepo.GetByIDWithVariants(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock based on whether it's a variant or a regular product
	isVariant := input.VariantID > 0
	if isVariant {
		// Find the variant and check its stock
		variant := product.GetVariantByID(input.VariantID)
		if variant == nil {
			return nil, errors.New("product variant not found")
		}

		if !variant.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	} else {
		// Regular product
		if product.HasVariants {
			return nil, errors.New("please select a product variant")
		}

		if !product.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	}

	// Get cart
	cart, err := uc.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Update item in cart
	if err := cart.UpdateItem(input.ProductID, input.VariantID, input.Quantity); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveFromCart removes a product from a user's cart
func (uc *CartUseCase) RemoveFromCart(userID uint, productID uint, variantID uint) (*entity.Cart, error) {
	// Get cart
	cart, err := uc.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Remove item from cart
	if err := cart.RemoveItem(productID, variantID); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// ClearCart removes all items from a user's cart
func (uc *CartUseCase) ClearCart(userID uint) error {
	// Get cart
	cart, err := uc.cartRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	// Clear cart
	cart.Clear()

	// Update cart in repository
	return uc.cartRepo.Update(cart)
}

// GetOrCreateGuestCart gets a guest cart or creates one if it doesn't exist
func (uc *CartUseCase) GetOrCreateGuestCart(sessionID string) (*entity.Cart, error) {
	cart, err := uc.cartRepo.GetBySessionID(sessionID)
	if err == nil {
		return cart, nil
	}

	// Create new cart if not found
	newCart, err := entity.NewGuestCart(sessionID)
	if err != nil {
		return nil, err
	}

	if err := uc.cartRepo.Create(newCart); err != nil {
		return nil, err
	}

	return newCart, nil
}

// AddToGuestCart adds a product to a guest's cart
func (uc *CartUseCase) AddToGuestCart(sessionID string, input AddToCartInput) (*entity.Cart, error) {
	// Check if product exists and has enough stock
	product, err := uc.productRepo.GetByIDWithVariants(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock based on whether it's a variant or a regular product
	isVariant := input.VariantID > 0
	if isVariant {
		variant := product.GetVariantByID(input.VariantID)
		if variant == nil {
			return nil, errors.New("product variant not found")
		}
		if !variant.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	} else {
		// Regular product
		if product.HasVariants {
			return nil, errors.New("please select a product variant")
		}

		if !product.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	}

	// Get cart
	cart, err := uc.GetOrCreateGuestCart(sessionID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	if err := cart.AddItem(input.ProductID, input.VariantID, input.Quantity); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// UpdateGuestCartItem updates the quantity of a product in a guest's cart
func (uc *CartUseCase) UpdateGuestCartItem(sessionID string, input UpdateCartItemInput) (*entity.Cart, error) {
	// Check if product exists and has enough stock
	product, err := uc.productRepo.GetByIDWithVariants(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock based on whether it's a variant or a regular product
	isVariant := input.VariantID > 0
	if isVariant {
		variant := product.GetVariantByID(input.VariantID)
		if variant == nil {
			return nil, errors.New("product variant not found")
		}
		if !variant.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	} else {
		// Regular product
		if product.HasVariants {
			return nil, errors.New("please select a product variant")
		}

		if !product.IsAvailable(input.Quantity) {
			return nil, errors.New("insufficient stock")
		}
	}

	// Get cart
	cart, err := uc.GetOrCreateGuestCart(sessionID)
	if err != nil {
		return nil, err
	}

	// Update item in cart
	if err := cart.UpdateItem(input.ProductID, input.VariantID, input.Quantity); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveFromGuestCart removes a product from a guest's cart
func (uc *CartUseCase) RemoveFromGuestCart(sessionID string, productID uint, variantID uint) (*entity.Cart, error) {
	// Get cart
	cart, err := uc.cartRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Remove item from cart
	if err := cart.RemoveItem(productID, variantID); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// ClearGuestCart removes all items from a guest's cart
func (uc *CartUseCase) ClearGuestCart(sessionID string) error {
	// Get cart
	cart, err := uc.cartRepo.GetBySessionID(sessionID)
	if err != nil {
		return errors.New("cart not found")
	}

	// Clear cart
	cart.Clear()

	// Update cart in repository
	return uc.cartRepo.Update(cart)
}

// ConvertGuestCartToUserCart converts a guest cart to a user cart
func (uc *CartUseCase) ConvertGuestCartToUserCart(sessionID string, userID uint) (*entity.Cart, error) {
	return uc.cartRepo.ConvertGuestCartToUserCart(sessionID, userID)
}
