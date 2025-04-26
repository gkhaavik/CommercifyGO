package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
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
		cartRepo.CreateWithID(cart)

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
		cartRepo.CreateWithID(cart)

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
		cartRepo.CreateWithID(cart)

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
		cartRepo.CreateWithID(cart)

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

func TestCartUseCase_AddToGuestCart(t *testing.T) {
	t.Run("Add item to guest cart successfully", func(t *testing.T) {
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

		// Create a test guest cart
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items:     []entity.CartItem{},
		}
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.AddToCartInput{
			ProductID: 1,
			Quantity:  2,
		}
		result, err := cartUseCase.AddToGuestCart(sessionID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, 2, result.Items[0].Quantity)
	})

	t.Run("Product not found for guest cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		sessionID := "test-session-123"
		input := usecase.AddToCartInput{
			ProductID: 999, // Non-existent product ID
			Quantity:  1,
		}
		result, err := cartUseCase.AddToGuestCart(sessionID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found")
	})
}

func TestCartUseCase_UpdateCartItem(t *testing.T) {
	t.Run("Update cart item successfully", func(t *testing.T) {
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

		// Create a test cart with item
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cart.AddItem(1, 2)
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.UpdateCartItemInput{
			ProductID: 1,
			Quantity:  5,
		}
		result, err := cartUseCase.UpdateCartItem(userID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, 5, result.Items[0].Quantity)
	})

	t.Run("Product not found when updating", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		userID := uint(1)
		input := usecase.UpdateCartItemInput{
			ProductID: 999, // Non-existent product ID
			Quantity:  3,
		}
		result, err := cartUseCase.UpdateCartItem(userID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found")
	})

	t.Run("Cart not found when updating", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product
		product := &entity.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 10.0,
			Stock: 10,
		}
		productRepo.Create(product)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent userID
		userID := uint(999)
		input := usecase.UpdateCartItemInput{
			ProductID: 1,
			Quantity:  3,
		}
		result, err := cartUseCase.UpdateCartItem(userID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart not found")
	})
}

func TestCartUseCase_UpdateGuestCartItem(t *testing.T) {
	t.Run("Update guest cart item successfully", func(t *testing.T) {
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

		// Create a test guest cart with item
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:        1,
					CartID:    1,
					ProductID: 1,
					Quantity:  2,
				},
			},
		}
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		input := usecase.UpdateCartItemInput{
			ProductID: 1,
			Quantity:  5,
		}
		result, err := cartUseCase.UpdateGuestCartItem(sessionID, input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, 5, result.Items[0].Quantity)
	})

	t.Run("Guest cart not found when updating", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test product
		product := &entity.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 10.0,
			Stock: 10,
		}
		productRepo.Create(product)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent sessionID
		sessionID := "non-existent-session"
		input := usecase.UpdateCartItemInput{
			ProductID: 1,
			Quantity:  3,
		}
		result, err := cartUseCase.UpdateGuestCartItem(sessionID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart not found")
	})
}

func TestCartUseCase_RemoveFromCart(t *testing.T) {
	t.Run("Remove item from cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test cart with multiple items
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cart.AddItem(1, 2)
		cart.AddItem(2, 1)
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		productID := uint(1)
		result, err := cartUseCase.RemoveFromCart(userID, productID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(2), result.Items[0].ProductID)
	})

	t.Run("Cart not found when removing", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent userID
		userID := uint(999)
		productID := uint(1)
		result, err := cartUseCase.RemoveFromCart(userID, productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart not found")
	})

	t.Run("Item not in cart", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test cart with an item
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cart.AddItem(1, 2)
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent product in cart
		productID := uint(2)
		result, err := cartUseCase.RemoveFromCart(userID, productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found in cart")
	})
}

func TestCartUseCase_RemoveFromGuestCart(t *testing.T) {
	t.Run("Remove item from guest cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test guest cart with multiple items
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:        1,
					CartID:    1,
					ProductID: 1,
					Quantity:  2,
				},
				{
					ID:        2,
					CartID:    1,
					ProductID: 2,
					Quantity:  1,
				},
			},
		}
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		productID := uint(1)
		result, err := cartUseCase.RemoveFromGuestCart(sessionID, productID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(2), result.Items[0].ProductID)
	})

	t.Run("Guest cart not found when removing", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent sessionID
		sessionID := "non-existent-session"
		productID := uint(1)
		result, err := cartUseCase.RemoveFromGuestCart(sessionID, productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart not found")
	})
}

func TestCartUseCase_ClearCart(t *testing.T) {
	t.Run("Clear cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test cart with items
		userID := uint(1)
		cart, _ := entity.NewCart(userID)
		cart.ID = 1
		cart.AddItem(1, 2)
		cart.AddItem(2, 3)
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		err := cartUseCase.ClearCart(userID)

		// Assert
		assert.NoError(t, err)

		// Verify cart is cleared
		updatedCart, _ := cartRepo.GetByUserID(userID)
		assert.Empty(t, updatedCart.Items)
	})

	t.Run("Cart not found when clearing", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent userID
		userID := uint(999)
		err := cartUseCase.ClearCart(userID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cart not found")
	})
}

func TestCartUseCase_ClearGuestCart(t *testing.T) {
	t.Run("Clear guest cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test guest cart with items
		sessionID := "test-session-123"
		cart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:        1,
					CartID:    1,
					ProductID: 1,
					Quantity:  2,
				},
				{
					ID:        2,
					CartID:    1,
					ProductID: 2,
					Quantity:  3,
				},
			},
		}
		cartRepo.CreateWithID(cart)

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute
		err := cartUseCase.ClearGuestCart(sessionID)

		// Assert
		assert.NoError(t, err)

		// Verify cart is cleared
		updatedCart, _ := cartRepo.GetBySessionID(sessionID)
		assert.Empty(t, updatedCart.Items)
	})

	t.Run("Guest cart not found when clearing", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent sessionID
		sessionID := "non-existent-session"
		err := cartUseCase.ClearGuestCart(sessionID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cart not found")
	})
}

func TestCartUseCase_ConvertGuestCartToUserCart(t *testing.T) {
	t.Run("Convert guest cart to user cart successfully", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create a test guest cart with items
		sessionID := "test-session-123"
		guestCart := &entity.Cart{
			ID:        1,
			SessionID: sessionID,
			Items: []entity.CartItem{
				{
					ID:        1,
					CartID:    1,
					ProductID: 1,
					Quantity:  2,
				},
			},
		}
		cartRepo.CreateWithID(guestCart)

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
		assert.Len(t, result.Items, 1)
		assert.Equal(t, uint(1), result.Items[0].ProductID)
		assert.Equal(t, 2, result.Items[0].Quantity)
	})

	t.Run("Guest cart not found when converting", func(t *testing.T) {
		// Setup mocks
		cartRepo := mock.NewMockCartRepository()
		productRepo := mock.NewMockProductRepository()

		// Create use case with mocks
		cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)

		// Execute with non-existent sessionID
		sessionID := "non-existent-session"
		userID := uint(1)
		result, err := cartUseCase.ConvertGuestCartToUserCart(sessionID, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart not found")
	})
}
