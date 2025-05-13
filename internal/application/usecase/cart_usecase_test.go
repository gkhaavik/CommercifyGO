package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/testutil/mock"
)

func TestCartUseCase_GetOrCreateCart(t *testing.T) {
	t.Run("Get existing cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		result, err := cartUseCase.GetOrCreateCart(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, cart.ID, result.ID)
		assert.Equal(t, cart.UserID, result.UserID)
	})

	t.Run("Create new cart when not found", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		userID := uint(2)
		result, err := cartUseCase.GetOrCreateCart(userID)

		// Assert
		assert.NoError(t, err)
		assert.NotZero(t, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Empty(t, result.Items)
	})
}

func TestCartUseCase_GetOrCreateGuestCart(t *testing.T) {
	t.Run("Get existing guest cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test guest cart
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items:     []entity.CartItem{},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		result, err := cartUseCase.GetOrCreateGuestCart(sessionID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, cart.ID, result.ID)
		assert.Equal(t, sessionID, result.SessionID)
	})

	t.Run("Create new guest cart when not found", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		sessionID := "new-session-456"
		result, err := cartUseCase.GetOrCreateGuestCart(sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotZero(t, result.ID)
		assert.Equal(t, sessionID, result.SessionID)
		assert.Empty(t, result.Items)
	})
}

func TestCartUseCase_AddToCart(t *testing.T) {
	t.Run("Add item to cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with sufficient stock
		product := &entity.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 10.0,
			Stock: 10,
		}
		productRepo.Create(product)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			Quantity:  2,
		}
		result, err := cartUseCase.AddToCart(userID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, 2, result.Items[0].Quantity)
	})

	t.Run("Product not found", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		userID := uint(1)
		input := usecase.AddToCartInput{
			ProductID: 999, // Non-existent product ID
			Quantity:  1,
		}
		result, err := cartUseCase.AddToCart(userID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found")
	})

	t.Run("Insufficient stock", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with limited stock
		product := &entity.Product{
			ID:    1,
			Name:  "Limited Stock Product",
			Price: 10.0,
			Stock: 3,
		}
		productRepo.Create(product)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			Quantity:  5, // More than available stock
		}
		result, err := cartUseCase.AddToCart(userID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient stock")
	})
}

func TestCartUseCase_AddToCartWithVariant(t *testing.T) {
	t.Run("Add item with variant to cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		variant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "VAR-001",
			Price:     money.ToCents(12.5),
			Stock:     8,
		}

		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       10.0,
			Stock:       10,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant},
		}
		productRepo.Create(product)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 1,
			Quantity:  2,
		}
		result, err := cartUseCase.AddToCart(userID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, uint(1), result.Items[0].ProductVariantID)
		assert.Equal(t, 2, result.Items[0].Quantity)
	})

	t.Run("Variant not found", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		product := &entity.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 10.0,
			Stock: 10,
		}
		productRepo.Create(product)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 999, // Non-existent variant ID
			Quantity:  2,
		}
		result, err := cartUseCase.AddToCart(userID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product variant not found")
	})

	t.Run("Insufficient variant stock", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		variant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "VAR-001",
			Price:     money.ToCents(12.5),
			Stock:     3,
		}

		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       money.ToCents(10.0),
			Stock:       10,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant},
		}
		productRepo.Create(product)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 1,
			Quantity:  5, // More than available stock
		}
		result, err := cartUseCase.AddToCart(userID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient stock")
	})
}

func TestCartUseCase_AddMultipleVariantsToCart(t *testing.T) {
	t.Run("Add multiple variants of same product to cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		variant1 := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "TSHIRT-S-RED",
			Price:     20.0,
			Stock:     10,
		}
		variant2 := &entity.ProductVariant{
			ID:        2,
			ProductID: 1,
			SKU:       "TSHIRT-M-RED",
			Price:     20.0,
			Stock:     10,
		}

		variant3 := &entity.ProductVariant{
			ID:        3,
			ProductID: 1,
			SKU:       "TSHIRT-L-BLUE",
			Price:     22.0,
			Stock:     5,
		}
		product := &entity.Product{
			ID:          1,
			Name:        "T-Shirt",
			Price:       20.0,
			Stock:       50,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant1, variant2, variant3},
		}
		productRepo.Create(product)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Add first variant
		input1 := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 1, // Small Red
			Quantity:  2,
		}
		result1, err1 := cartUseCase.AddToCart(userID, input1)
		assert.NoError(t, err1)
		assert.Len(t, result1.Items, 1)

		// Add second variant of same product
		input2 := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 2, // Medium Red
			Quantity:  1,
		}
		result2, err2 := cartUseCase.AddToCart(userID, input2)
		assert.NoError(t, err2)
		assert.Len(t, result2.Items, 2)

		// Add third variant of same product
		input3 := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 3, // Large Blue
			Quantity:  3,
		}
		result3, err3 := cartUseCase.AddToCart(userID, input3)
		assert.NoError(t, err3)
		assert.Len(t, result3.Items, 3)

		// Verify each variant is correctly added
		foundVariants := make(map[uint]bool)
		for _, item := range result3.Items {
			assert.Equal(t, uint(1), item.ProductID)
			foundVariants[item.ProductVariantID] = true
		}

		assert.True(t, foundVariants[1], "Small variant should be in cart")
		assert.True(t, foundVariants[2], "Medium variant should be in cart")
		assert.True(t, foundVariants[3], "Large variant should be in cart")
	})
}

func TestCartUseCase_MixedCartWithVariantsAndRegularProducts(t *testing.T) {
	t.Run("Create mixed cart with variants and regular products", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a product with variants
		variant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "TSHIRT-M-BLUE",
			Price:     22.0,
			Stock:     15,
		}
		variantProduct := &entity.Product{
			ID:          1,
			Name:        "T-Shirt",
			Price:       20.0,
			Stock:       50,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant},
		}
		productRepo.Create(variantProduct)

		// Create a regular product (no variants)
		regularProduct := &entity.Product{
			ID:    2,
			Name:  "Water Bottle",
			Price: 15.0,
			Stock: 20,
		}
		productRepo.Create(regularProduct)

		// Create a test cart
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Add regular product
		input1 := usecase.AddToCartInput{
			ProductID: 2, // Water Bottle
			Quantity:  2,
		}
		result1, err1 := cartUseCase.AddToCart(userID, input1)
		assert.NoError(t, err1)
		assert.Len(t, result1.Items, 1)
		assert.Equal(t, uint(2), result1.Items[0].ProductID)
		assert.Equal(t, uint(0), result1.Items[0].ProductVariantID)

		// Add product with variant
		input2 := usecase.AddToCartInput{
			ProductID: 1, // T-Shirt
			VariantID: 1, // Medium Blue
			Quantity:  1,
		}
		result2, err2 := cartUseCase.AddToCart(userID, input2)
		assert.NoError(t, err2)
		assert.Len(t, result2.Items, 2)

		// Check that both items are in the cart
		hasRegularProduct := false
		hasVariantProduct := false

		for _, item := range result2.Items {
			if item.ProductID == 2 && item.ProductVariantID == 0 {
				hasRegularProduct = true
				assert.Equal(t, 2, item.Quantity)
			}
			if item.ProductID == 1 && item.ProductVariantID == 1 {
				hasVariantProduct = true
				assert.Equal(t, 1, item.Quantity)
			}
		}

		assert.True(t, hasRegularProduct, "Regular product should be in cart")
		assert.True(t, hasVariantProduct, "Product variant should be in cart")
	})
}

func TestCartUseCase_RemoveFromCartWithVariant(t *testing.T) {
	t.Run("Remove item with variant from cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test cart with multiple items including variants
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cart.AddItem(1, 1, 2) // Product 1, Variant 1
		cart.AddItem(1, 2, 1) // Product 1, Variant 2
		cart.AddItem(2, 0, 3) // Product 2, No variant
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute - remove specific variant
		productID := uint(1)
		variantID := uint(1)
		result, err := cartUseCase.RemoveFromCart(userID, productID, variantID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 2)

		// Verify the correct variant was removed
		hasVariant1 := false
		for _, item := range result.Items {
			if item.ProductID == 1 && item.ProductVariantID == 1 {
				hasVariant1 = true
			}
		}
		assert.False(t, hasVariant1, "Variant 1 should be removed")

		// Verify other items remain
		hasVariant2 := false
		hasProduct2 := false
		for _, item := range result.Items {
			if item.ProductID == 1 && item.ProductVariantID == 2 {
				hasVariant2 = true
			}
			if item.ProductID == 2 {
				hasProduct2 = true
			}
		}
		assert.True(t, hasVariant2, "Variant 2 should remain")
		assert.True(t, hasProduct2, "Product 2 should remain")
	})

	t.Run("Variant not in cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test cart with an item
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cart.AddItem(1, 1, 2) // Product 1, Variant 1
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute - try to remove non-existent variant
		productID := uint(1)
		variantID := uint(2) // Non-existent variant in cart
		result, err := cartUseCase.RemoveFromCart(userID, productID, variantID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found in cart")
	})
}

func TestCartUseCase_AddToGuestCartWithVariant(t *testing.T) {
	t.Run("Add item with variant to guest cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		variant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "VAR-001",
			Price:     money.ToCents(12.5),
			Stock:     8,
		}
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       10.0,
			Stock:       10,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant},
		}
		productRepo.Create(product)

		// Create a test guest cart
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items:     []entity.CartItem{},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 1,
			Quantity:  2,
		}
		result, err := cartUseCase.AddToGuestCart(sessionID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, uint(1), result.Items[0].ProductVariantID)
		assert.Equal(t, 2, result.Items[0].Quantity)
	})

	t.Run("Insufficient variant stock in guest cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		variant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "VAR-001",
			Price:     money.ToCents(12.5),
			Stock:     3,
		}
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       money.ToCents(10.0),
			Stock:       10,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant},
		}
		productRepo.Create(product)

		// Create a test guest cart
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items:     []entity.CartItem{},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			VariantID: 1,
			Quantity:  5, // More than available stock
		}
		result, err := cartUseCase.AddToGuestCart(sessionID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient stock")
	})
}

func TestCartUseCase_UpdateGuestCartItemWithVariant(t *testing.T) {
	t.Run("Update guest cart item with variant successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product with variants
		variant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "VAR-001",
			Price:     money.ToCents(12.5),
			Stock:     8,
		}
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       money.ToCents(10.0),
			Stock:       10,
			HasVariants: true,
			Variants:    []*entity.ProductVariant{variant},
		}
		productRepo.Create(product)

		// Create a test guest cart with variant item
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:               1,
					CartID:           1,
					ProductID:        1,
					ProductVariantID: 1,
					Quantity:         2,
				},
			},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.UpdateCartItemInput{
			ProductID: 1,
			VariantID: 1,
			Quantity:  5,
		}
		result, err := cartUseCase.UpdateGuestCartItem(sessionID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, uint(1), result.Items[0].ProductVariantID)
		assert.Equal(t, 5, result.Items[0].Quantity)
	})
}

func TestCartUseCase_RemoveFromGuestCartWithVariant(t *testing.T) {
	t.Run("Remove variant from guest cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test guest cart with multiple items including variants
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:               1,
					CartID:           1,
					ProductID:        1,
					ProductVariantID: 1,
					Quantity:         2,
				},
				{
					ID:               2,
					CartID:           1,
					ProductID:        1,
					ProductVariantID: 2,
					Quantity:         1,
				},
				{
					ID:        3,
					CartID:    1,
					ProductID: 2,
					Quantity:  3,
				},
			},
		}
		cartRepo.Create(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute - remove specific variant
		productID := uint(1)
		variantID := uint(1)
		result, err := cartUseCase.RemoveFromGuestCart(sessionID, productID, variantID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 2)

		// Verify the correct variant was removed
		hasVariant1 := false
		for _, item := range result.Items {
			if item.ProductID == 1 && item.ProductVariantID == 1 {
				hasVariant1 = true
			}
		}
		assert.False(t, hasVariant1, "Variant 1 should be removed")

		// Verify other items remain
		hasVariant2 := false
		hasProduct2 := false
		for _, item := range result.Items {
			if item.ProductID == 1 && item.ProductVariantID == 2 {
				hasVariant2 = true
			}
			if item.ProductID == 2 {
				hasProduct2 = true
			}
		}
		assert.True(t, hasVariant2, "Variant 2 should remain")
		assert.True(t, hasProduct2, "Product 2 should remain")
	})
}

func TestCartUseCase_ConvertGuestCartToUserCartWithVariants(t *testing.T) {
	t.Run("Convert guest cart with variants to user cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test guest cart with items including variants
		sessionID := "test-session-123"
		guestCart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:               1,
					CartID:           1,
					ProductID:        1,
					ProductVariantID: 1,
					Quantity:         2,
				},
				{
					ID:        2,
					CartID:    1,
					ProductID: 2,
					Quantity:  3,
				},
			},
		}
		cartRepo.Create(guestCart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		userID := uint(1)
		result, err := cartUseCase.ConvertGuestCartToUserCart(sessionID, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Empty(t, result.SessionID) // SessionID should be cleared
		assert.Len(t, result.Items, 2)

		// Check that both items were transferred correctly
		hasVariantProduct := false
		hasRegularProduct := false
		for _, item := range result.Items {
			if item.ProductID == 1 && item.ProductVariantID == 1 {
				hasVariantProduct = true
				assert.Equal(t, 2, item.Quantity)
			}
			if item.ProductID == 2 && item.ProductVariantID == 0 {
				hasRegularProduct = true
				assert.Equal(t, 3, item.Quantity)
			}
		}

		assert.True(t, hasVariantProduct, "Product with variant should be in user cart")
		assert.True(t, hasRegularProduct, "Regular product should be in user cart")
	})
}

func TestCartUseCase_MergingCartWithVariants(t *testing.T) {
	t.Run("Merge guest cart with variants into existing user cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a user cart with some items
		userID := uint(1)
		userCart, _ := entity.NewCart(userID)
		userCart.ID = 1
		userCart.AddItem(1, 2, 1) // Product 1, Variant 2
		userCart.AddItem(3, 0, 2) // Product 3, No variant
		cartRepo.Create(userCart)

		// Create a guest cart with some items
		sessionID := "test-session-123"
		guestCart := &entity.Cart{
			ID:        2,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:               3,
					CartID:           2,
					ProductID:        1,
					ProductVariantID: 1,
					Quantity:         3,
				},
				{
					ID:               4,
					CartID:           2,
					ProductID:        1,
					ProductVariantID: 2,
					Quantity:         2, // Same variant as in user cart
				},
				{
					ID:        5,
					CartID:    2,
					ProductID: 2,
					Quantity:  1, // New product
				},
			},
		}
		cartRepo.Create(guestCart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		result, err := cartUseCase.ConvertGuestCartToUserCart(sessionID, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Len(t, result.Items, 4)

		// Check merged items
		itemQuantities := make(map[string]int)
		for _, item := range result.Items {
			key := ""
			if item.ProductVariantID > 0 {
				key = fmt.Sprintf("%d-%d", item.ProductID, item.ProductVariantID)
			} else {
				key = fmt.Sprintf("%d", item.ProductID)
			}
			itemQuantities[key] = item.Quantity
		}

		// Variant 2 of Product 1 should have combined quantity
		assert.Equal(t, 3, itemQuantities["1-2"], "Quantities should be combined for existing variant")

		// New variant should be added
		assert.Equal(t, 3, itemQuantities["1-1"], "New variant should be added")

		// Original product in user cart should remain
		assert.Equal(t, 2, itemQuantities["3"], "Original product in user cart should remain")

		// New product from guest cart should be added
		assert.Equal(t, 1, itemQuantities["2"], "New product from guest cart should be added")
	})
}
