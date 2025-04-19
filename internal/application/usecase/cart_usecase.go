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
	ProductID uint
	Quantity  int
}

// AddToCart adds a product to a user's cart
func (uc *CartUseCase) AddToCart(userID uint, input AddToCartInput) (*entity.Cart, error) {
	// Check if product exists and has enough stock
	product, err := uc.productRepo.GetByID(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if !product.IsAvailable(input.Quantity) {
		return nil, errors.New("insufficient stock")
	}

	// Get or create cart
	cart, err := uc.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	if err := cart.AddItem(input.ProductID, input.Quantity); err != nil {
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
	ProductID uint
	Quantity  int
}

// UpdateCartItem updates the quantity of a product in a user's cart
func (uc *CartUseCase) UpdateCartItem(userID uint, input UpdateCartItemInput) (*entity.Cart, error) {
	// Check if product exists and has enough stock
	product, err := uc.productRepo.GetByID(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if !product.IsAvailable(input.Quantity) {
		return nil, errors.New("insufficient stock")
	}

	// Get cart
	cart, err := uc.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Update item in cart
	if err := cart.UpdateItem(input.ProductID, input.Quantity); err != nil {
		return nil, err
	}

	// Update cart in repository
	if err := uc.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveFromCart removes a product from a user's cart
func (uc *CartUseCase) RemoveFromCart(userID uint, productID uint) (*entity.Cart, error) {
	// Get cart
	cart, err := uc.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Remove item from cart
	if err := cart.RemoveItem(productID); err != nil {
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
