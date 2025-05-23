package usecase

import (
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CheckoutInput defines the input for creating/adding to a checkout
type CheckoutInput struct {
	ProductID   uint
	VariantID   uint
	Quantity    int
	Price       int64
	Weight      float64
	ProductName string
	VariantName string
	SKU         string
}

// UpdateCheckoutItemInput defines the input for updating a checkout item
type UpdateCheckoutItemInput struct {
	ProductID uint
	VariantID uint
	Quantity  int
}

// CheckoutUseCase implements checkout business logic
type CheckoutUseCase struct {
	checkoutRepo       repository.CheckoutRepository
	productRepo        repository.ProductRepository
	productVariantRepo repository.ProductVariantRepository
	shippingMethodRepo repository.ShippingMethodRepository
	shippingRateRepo   repository.ShippingRateRepository
	discountRepo       repository.DiscountRepository
	orderRepo          repository.OrderRepository
	currencyRepo       repository.CurrencyRepository
}

// NewCheckoutUseCase creates a new checkout use case
func NewCheckoutUseCase(
	checkoutRepo repository.CheckoutRepository,
	productRepo repository.ProductRepository,
	productVariantRepo repository.ProductVariantRepository,
	shippingMethodRepo repository.ShippingMethodRepository,
	shippingRateRepo repository.ShippingRateRepository,
	discountRepo repository.DiscountRepository,
	orderRepo repository.OrderRepository,
	currencyRepo repository.CurrencyRepository,
) *CheckoutUseCase {
	return &CheckoutUseCase{
		checkoutRepo:       checkoutRepo,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		shippingMethodRepo: shippingMethodRepo,
		shippingRateRepo:   shippingRateRepo,
		discountRepo:       discountRepo,
		orderRepo:          orderRepo,
		currencyRepo:       currencyRepo,
	}
}

// GetOrCreateCheckout retrieves or creates a checkout for a user
func (uc *CheckoutUseCase) GetOrCreateCheckout(userID uint) (*entity.Checkout, error) {
	// Try to get an existing active checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err == nil {
		// If found, return it
		return checkout, nil
	}

	// If not found, create a new one
	checkout, err = entity.NewCheckout(userID)
	if err != nil {
		return nil, err
	}

	// Set default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err == nil && defaultCurrency != nil {
		checkout.Currency = defaultCurrency.Code
	}

	// Save to repository
	err = uc.checkoutRepo.Create(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// GetOrCreateGuestCheckout retrieves or creates a checkout for a guest user
func (uc *CheckoutUseCase) GetOrCreateGuestCheckout(sessionID string) (*entity.Checkout, error) {
	// Try to get an existing active checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err == nil {
		// If found, return it
		return checkout, nil
	}

	// If not found, create a new one
	checkout, err = entity.NewGuestCheckout(sessionID)
	if err != nil {
		return nil, err
	}

	// Set default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err == nil && defaultCurrency != nil {
		checkout.Currency = defaultCurrency.Code
	}

	// Save to repository
	err = uc.checkoutRepo.Create(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// AddToCheckout adds a product to the user's checkout
func (uc *CheckoutUseCase) AddToCheckout(userID uint, input CheckoutInput) (*entity.Checkout, error) {
	// Get or create checkout
	checkout, err := uc.GetOrCreateCheckout(userID)
	if err != nil {
		return nil, err
	}

	// Get product details if not provided
	if input.ProductName == "" || input.Price == 0 {
		product, err := uc.productRepo.GetByID(input.ProductID)
		if err != nil {
			return nil, err
		}

		input.ProductName = product.Name
		input.Price = product.Price
		input.Weight = product.Weight

		// Check if product is active
		if !product.Active {
			return nil, errors.New("product is not available")
		}
	}

	// If variant is provided, get variant details
	if input.VariantID > 0 && (input.VariantName == "" || input.SKU == "") {
		variant, err := uc.productVariantRepo.GetByID(input.VariantID)
		if err != nil {
			return nil, err
		}

		// Extract variant name from attributes
		variantName := ""
		for _, attr := range variant.Attributes {
			if variantName == "" {
				variantName = attr.Value
			} else {
				variantName += " / " + attr.Value
			}
		}
		input.VariantName = variantName
		input.SKU = variant.SKU
		// Use variant price directly
		input.Price = variant.Price
		// Note: Weight is not available in the ProductVariant entity
	}

	// Add item to checkout
	err = checkout.AddItem(
		input.ProductID,
		input.VariantID,
		input.Quantity,
		input.Price,
		input.Weight,
		input.ProductName,
		input.VariantName,
		input.SKU,
	)
	if err != nil {
		return nil, err
	}

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// AddToGuestCheckout adds a product to a guest's checkout
func (uc *CheckoutUseCase) AddToGuestCheckout(sessionID string, input CheckoutInput) (*entity.Checkout, error) {
	// Get or create checkout
	checkout, err := uc.GetOrCreateGuestCheckout(sessionID)
	if err != nil {
		return nil, err
	}

	// Get product details if not provided
	if input.ProductName == "" || input.Price == 0 {
		product, err := uc.productRepo.GetByID(input.ProductID)
		if err != nil {
			return nil, err
		}

		input.ProductName = product.Name
		input.Price = product.Price
		input.Weight = product.Weight

		// Check if product is active
		if !product.Active {
			return nil, errors.New("product is not available")
		}
	}

	// If variant is provided, get variant details
	if input.VariantID > 0 && (input.VariantName == "" || input.SKU == "") {
		variant, err := uc.productVariantRepo.GetByID(input.VariantID)
		if err != nil {
			return nil, err
		}

		// Extract variant name from attributes
		variantName := ""
		for _, attr := range variant.Attributes {
			if variantName == "" {
				variantName = attr.Value
			} else {
				variantName += " / " + attr.Value
			}
		}
		input.VariantName = variantName
		input.SKU = variant.SKU
		// Use variant price directly
		input.Price = variant.Price
		// Note: Weight is not available in the ProductVariant entity
	}

	// Add item to checkout
	err = checkout.AddItem(
		input.ProductID,
		input.VariantID,
		input.Quantity,
		input.Price,
		input.Weight,
		input.ProductName,
		input.VariantName,
		input.SKU,
	)
	if err != nil {
		return nil, err
	}

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// UpdateCheckoutItem updates the quantity of a product in the user's checkout
func (uc *CheckoutUseCase) UpdateCheckoutItem(userID uint, input UpdateCheckoutItemInput) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update item in checkout
	err = checkout.UpdateItem(input.ProductID, input.VariantID, input.Quantity)
	if err != nil {
		return nil, err
	}

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// UpdateGuestCheckoutItem updates the quantity of a product in a guest's checkout
func (uc *CheckoutUseCase) UpdateGuestCheckoutItem(sessionID string, input UpdateCheckoutItemInput) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Update item in checkout
	err = checkout.UpdateItem(input.ProductID, input.VariantID, input.Quantity)
	if err != nil {
		return nil, err
	}

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// RemoveFromCheckout removes a product from the user's checkout
func (uc *CheckoutUseCase) RemoveFromCheckout(userID uint, productID uint, variantID uint) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Remove item from checkout
	err = checkout.RemoveItem(productID, variantID)
	if err != nil {
		return nil, err
	}

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// RemoveFromGuestCheckout removes a product from a guest's checkout
func (uc *CheckoutUseCase) RemoveFromGuestCheckout(sessionID string, productID uint, variantID uint) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Remove item from checkout
	err = checkout.RemoveItem(productID, variantID)
	if err != nil {
		return nil, err
	}

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ClearCheckout empties the user's checkout
func (uc *CheckoutUseCase) ClearCheckout(userID uint) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Clear checkout
	checkout.Clear()

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ClearGuestCheckout empties a guest's checkout
func (uc *CheckoutUseCase) ClearGuestCheckout(sessionID string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Clear checkout
	checkout.Clear()

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetShippingAddress sets the shipping address for the user's checkout
func (uc *CheckoutUseCase) SetShippingAddress(userID uint, address entity.Address) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set shipping address
	checkout.SetShippingAddress(address)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetGuestShippingAddress sets the shipping address for a guest's checkout
func (uc *CheckoutUseCase) SetGuestShippingAddress(sessionID string, address entity.Address) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Set shipping address
	checkout.SetShippingAddress(address)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetBillingAddress sets the billing address for the user's checkout
func (uc *CheckoutUseCase) SetBillingAddress(userID uint, address entity.Address) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set billing address
	checkout.SetBillingAddress(address)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetGuestBillingAddress sets the billing address for a guest's checkout
func (uc *CheckoutUseCase) SetGuestBillingAddress(sessionID string, address entity.Address) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Set billing address
	checkout.SetBillingAddress(address)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetCustomerDetails sets the customer details for the user's checkout
func (uc *CheckoutUseCase) SetCustomerDetails(userID uint, details entity.CustomerDetails) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set customer details
	checkout.SetCustomerDetails(details)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetGuestCustomerDetails sets the customer details for a guest's checkout
func (uc *CheckoutUseCase) SetGuestCustomerDetails(sessionID string, details entity.CustomerDetails) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Set customer details
	checkout.SetCustomerDetails(details)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetShippingMethod sets the shipping method for the user's checkout
func (uc *CheckoutUseCase) SetShippingMethod(userID uint, methodID uint) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get shipping method
	shippingMethod, err := uc.shippingMethodRepo.GetByID(methodID)
	if err != nil {
		return nil, err
	}

	// Calculate shipping cost
	var shippingCost int64
	if checkout.ShippingAddr.Street != "" && checkout.ShippingAddr.Country != "" {
		rates, err := uc.shippingRateRepo.GetAvailableRatesForAddress(checkout.ShippingAddr, checkout.TotalAmount)
		if err != nil {
			return nil, err
		}

		for _, rate := range rates {
			if rate.ShippingMethodID == methodID {
				shippingCost = rate.BaseRate

				// Check for weight-based rates
				weightRates, err := uc.shippingRateRepo.GetWeightBasedRates(rate.ID)
				if err == nil && len(weightRates) > 0 {
					for _, weightRate := range weightRates {
						if checkout.TotalWeight >= weightRate.MinWeight && (weightRate.MaxWeight == 0 || checkout.TotalWeight <= weightRate.MaxWeight) {
							shippingCost = weightRate.Rate
							break
						}
					}
				}

				// Check for value-based rates
				valueRates, err := uc.shippingRateRepo.GetValueBasedRates(rate.ID)
				if err == nil && len(valueRates) > 0 {
					for _, valueRate := range valueRates {
						if checkout.TotalAmount >= valueRate.MinOrderValue && (valueRate.MaxOrderValue == 0 || checkout.TotalAmount <= valueRate.MaxOrderValue) {
							shippingCost = valueRate.Rate
							break
						}
					}
				}

				break
			}
		}
	}

	// Set shipping method and cost
	checkout.SetShippingMethod(methodID, shippingCost)
	checkout.ShippingMethod = shippingMethod

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetGuestShippingMethod sets the shipping method for a guest's checkout
func (uc *CheckoutUseCase) SetGuestShippingMethod(sessionID string, methodID uint) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Get shipping method
	shippingMethod, err := uc.shippingMethodRepo.GetByID(methodID)
	if err != nil {
		return nil, err
	}

	// Calculate shipping cost
	var shippingCost int64
	if checkout.ShippingAddr.Street != "" && checkout.ShippingAddr.Country != "" {
		rates, err := uc.shippingRateRepo.GetAvailableRatesForAddress(checkout.ShippingAddr, checkout.TotalAmount)
		if err != nil {
			return nil, err
		}

		for _, rate := range rates {
			if rate.ShippingMethodID == methodID {
				shippingCost = rate.BaseRate

				// Check for weight-based rates
				weightRates, err := uc.shippingRateRepo.GetWeightBasedRates(rate.ID)
				if err == nil && len(weightRates) > 0 {
					for _, weightRate := range weightRates {
						if checkout.TotalWeight >= weightRate.MinWeight && (weightRate.MaxWeight == 0 || checkout.TotalWeight <= weightRate.MaxWeight) {
							shippingCost = weightRate.Rate
							break
						}
					}
				}

				// Check for value-based rates
				valueRates, err := uc.shippingRateRepo.GetValueBasedRates(rate.ID)
				if err == nil && len(valueRates) > 0 {
					for _, valueRate := range valueRates {
						if checkout.TotalAmount >= valueRate.MinOrderValue && (valueRate.MaxOrderValue == 0 || checkout.TotalAmount <= valueRate.MaxOrderValue) {
							shippingCost = valueRate.Rate
							break
						}
					}
				}

				break
			}
		}
	}

	// Set shipping method and cost
	checkout.SetShippingMethod(methodID, shippingCost)
	checkout.ShippingMethod = shippingMethod

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetPaymentProvider sets the payment provider for the user's checkout
func (uc *CheckoutUseCase) SetPaymentProvider(userID uint, provider string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set payment provider
	checkout.SetPaymentProvider(provider)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetGuestPaymentProvider sets the payment provider for a guest's checkout
func (uc *CheckoutUseCase) SetGuestPaymentProvider(sessionID string, provider string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Set payment provider
	checkout.SetPaymentProvider(provider)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ApplyDiscountCode applies a discount code to the user's checkout
func (uc *CheckoutUseCase) ApplyDiscountCode(userID uint, code string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get discount
	discount, err := uc.discountRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Check if discount is valid
	if !discount.IsValid() {
		return nil, errors.New("discount is not valid")
	}

	// Apply discount
	checkout.ApplyDiscount(discount)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ApplyGuestDiscountCode applies a discount code to a guest's checkout
func (uc *CheckoutUseCase) ApplyGuestDiscountCode(sessionID string, code string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Get discount
	discount, err := uc.discountRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Check if discount is valid
	if !discount.IsValid() {
		return nil, errors.New("discount is not valid")
	}

	// Apply discount
	checkout.ApplyDiscount(discount)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// RemoveDiscountCode removes a discount code from the user's checkout
func (uc *CheckoutUseCase) RemoveDiscountCode(userID uint) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Remove discount
	checkout.ApplyDiscount(nil)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// RemoveGuestDiscountCode removes a discount code from a guest's checkout
func (uc *CheckoutUseCase) RemoveGuestDiscountCode(sessionID string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Remove discount
	checkout.ApplyDiscount(nil)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ConvertGuestCheckoutToUserCheckout converts a guest checkout to a user checkout
func (uc *CheckoutUseCase) ConvertGuestCheckoutToUserCheckout(sessionID string, userID uint) (*entity.Checkout, error) {
	return uc.checkoutRepo.ConvertGuestCheckoutToUserCheckout(sessionID, userID)
}

// ExpireOldCheckouts marks expired checkouts as expired
func (uc *CheckoutUseCase) ExpireOldCheckouts() error {
	// Get expired checkouts
	expiredCheckouts, err := uc.checkoutRepo.GetExpiredCheckouts()
	if err != nil {
		return err
	}

	// Mark each as expired
	for _, checkout := range expiredCheckouts {
		checkout.MarkAsExpired()
		err = uc.checkoutRepo.Update(checkout)
		if err != nil {
			// Continue despite errors
			continue
		}
	}

	return nil
}

// CreateOrderFromCheckout creates an order from a checkout
func (uc *CheckoutUseCase) CreateOrderFromCheckout(checkoutID uint) (*entity.Order, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Validate checkout
	if len(checkout.Items) == 0 {
		return nil, errors.New("checkout has no items")
	}

	if checkout.ShippingAddr.Street == "" || checkout.ShippingAddr.Country == "" {
		return nil, errors.New("shipping address is required")
	}

	if checkout.BillingAddr.Street == "" || checkout.BillingAddr.Country == "" {
		return nil, errors.New("billing address is required")
	}

	if checkout.CustomerDetails.Email == "" || checkout.CustomerDetails.FullName == "" {
		return nil, errors.New("customer details are required")
	}

	// Convert checkout to order
	order := checkout.ToOrder()

	// Create order in repository
	err = uc.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	// Mark checkout as completed
	checkout.MarkAsCompleted(order.ID)
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		// Don't fail if we can't update checkout status
		// The order is already created
	}

	// Increment discount usage if a discount was applied
	if checkout.AppliedDiscount != nil {
		discount, err := uc.discountRepo.GetByID(checkout.AppliedDiscount.DiscountID)
		if err == nil {
			discount.IncrementUsage()
			uc.discountRepo.Update(discount)
		}
	}

	return order, nil
}

// CreateOrderFromUserCheckout creates an order from the active checkout belonging to a user
func (uc *CheckoutUseCase) CreateOrderFromUserCheckout(userID uint) (*entity.Order, error) {
	// Get checkout for the user
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Use the existing method to create the order from checkout
	return uc.CreateOrderFromCheckout(checkout.ID)
}

// CreateOrderFromGuestCheckout creates an order from the active checkout belonging to a guest session
func (uc *CheckoutUseCase) CreateOrderFromGuestCheckout(sessionID string) (*entity.Order, error) {
	// Get checkout for the session
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Use the existing method to create the order from checkout
	return uc.CreateOrderFromCheckout(checkout.ID)
}

// ExtendCheckoutExpiry extends the expiry time of a checkout
func (uc *CheckoutUseCase) ExtendCheckoutExpiry(checkoutID uint, duration time.Duration) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Extend expiry
	checkout.ExtendExpiry(duration)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// GetCheckoutByID retrieves a checkout by ID
func (uc *CheckoutUseCase) GetCheckoutByID(checkoutID uint) (*entity.Checkout, error) {
	return uc.checkoutRepo.GetByID(checkoutID)
}

// AbandonCheckout marks a checkout as abandoned
func (uc *CheckoutUseCase) AbandonCheckout(checkoutID uint) error {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return err
	}

	// Mark as abandoned
	checkout.MarkAsAbandoned()

	// Update checkout in repository
	return uc.checkoutRepo.Update(checkout)
}

// GetCheckoutsByStatus retrieves checkouts by status with pagination
func (uc *CheckoutUseCase) GetCheckoutsByStatus(status entity.CheckoutStatus, offset, limit int) ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetCheckoutsByStatus(status, offset, limit)
}

// GetAllCheckouts retrieves all checkouts with pagination
func (uc *CheckoutUseCase) GetAllCheckouts(offset, limit int) ([]*entity.Checkout, error) {
	// If no specific status is requested, get checkouts regardless of status
	return uc.checkoutRepo.GetCheckoutsByStatus("", offset, limit)
}

// DeleteCheckout deletes a checkout by ID
func (uc *CheckoutUseCase) DeleteCheckout(checkoutID uint) error {
	return uc.checkoutRepo.Delete(checkoutID)
}

// GetExpiredCheckouts retrieves all expired checkouts
func (uc *CheckoutUseCase) GetExpiredCheckouts() ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetExpiredCheckouts()
}

// GetAbandonedCheckouts retrieves all abandoned checkouts
func (uc *CheckoutUseCase) GetAbandonedCheckouts(offset, limit int) ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetCheckoutsByStatus(entity.CheckoutStatusAbandoned, offset, limit)
}

// GetCheckoutsByUserID retrieves all checkouts for a user with pagination
func (uc *CheckoutUseCase) GetCheckoutsByUserID(userID uint, offset, limit int) ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetCompletedCheckoutsByUserID(userID, offset, limit)
}
