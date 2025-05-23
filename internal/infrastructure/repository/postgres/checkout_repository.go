package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CheckoutRepository implements the checkout repository interface using PostgreSQL
type CheckoutRepository struct {
	db *sql.DB
}

// NewCheckoutRepository creates a new CheckoutRepository
func NewCheckoutRepository(db *sql.DB) repository.CheckoutRepository {
	return &CheckoutRepository{db: db}
}

// Create creates a new checkout
func (r *CheckoutRepository) Create(checkout *entity.Checkout) error {
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

	// Marshal JSON fields
	shippingAddrJSON, err := json.Marshal(checkout.ShippingAddr)
	if err != nil {
		return err
	}

	billingAddrJSON, err := json.Marshal(checkout.BillingAddr)
	if err != nil {
		return err
	}

	customerDetailsJSON, err := json.Marshal(checkout.CustomerDetails)
	if err != nil {
		return err
	}

	var appliedDiscountJSON []byte
	if checkout.AppliedDiscount != nil {
		appliedDiscountJSON, err = json.Marshal(checkout.AppliedDiscount)
		if err != nil {
			return err
		}
	}

	// Insert checkout
	var query string
	if checkout.SessionID != "" {
		// Guest checkout
		query = `
			INSERT INTO checkouts (
				session_id, status, shipping_address, billing_address, shipping_method_id, 
				payment_provider, total_amount, shipping_cost, total_weight, customer_details, 
				currency, discount_code, discount_amount, final_amount, applied_discount, 
				created_at, updated_at, last_activity_at, expires_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
			) RETURNING id
		`
		err = tx.QueryRow(
			query,
			checkout.SessionID,
			checkout.Status,
			shippingAddrJSON,
			billingAddrJSON,
			nullableUint(checkout.ShippingMethodID),
			nullableString(checkout.PaymentProvider),
			checkout.TotalAmount,
			checkout.ShippingCost,
			checkout.TotalWeight,
			customerDetailsJSON,
			checkout.Currency,
			nullableString(checkout.DiscountCode),
			checkout.DiscountAmount,
			checkout.FinalAmount,
			nullableBytes(appliedDiscountJSON),
			checkout.CreatedAt,
			checkout.UpdatedAt,
			checkout.LastActivityAt,
			checkout.ExpiresAt,
		).Scan(&checkout.ID)
	} else {
		// User checkout
		query = `
			INSERT INTO checkouts (
				user_id, status, shipping_address, billing_address, shipping_method_id, 
				payment_provider, total_amount, shipping_cost, total_weight, customer_details, 
				currency, discount_code, discount_amount, final_amount, applied_discount, 
				created_at, updated_at, last_activity_at, expires_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
			) RETURNING id
		`
		err = tx.QueryRow(
			query,
			checkout.UserID,
			checkout.Status,
			shippingAddrJSON,
			billingAddrJSON,
			nullableUint(checkout.ShippingMethodID),
			nullableString(checkout.PaymentProvider),
			checkout.TotalAmount,
			checkout.ShippingCost,
			checkout.TotalWeight,
			customerDetailsJSON,
			checkout.Currency,
			nullableString(checkout.DiscountCode),
			checkout.DiscountAmount,
			checkout.FinalAmount,
			nullableBytes(appliedDiscountJSON),
			checkout.CreatedAt,
			checkout.UpdatedAt,
			checkout.LastActivityAt,
			checkout.ExpiresAt,
		).Scan(&checkout.ID)
	}

	if err != nil {
		return err
	}

	// Insert checkout items
	if len(checkout.Items) > 0 {
		for i := range checkout.Items {
			checkout.Items[i].CheckoutID = checkout.ID

			if checkout.Items[i].ProductVariantID == 0 {
				// If variant ID is 0, use NULL for the database
				query := `
					INSERT INTO checkout_items (
						checkout_id, product_id, product_variant_id, quantity, price, 
						weight, product_name, variant_name, sku, created_at, updated_at
					) VALUES (
						$1, $2, NULL, $3, $4, $5, $6, $7, $8, $9, $10
					) RETURNING id
				`
				err = tx.QueryRow(
					query,
					checkout.Items[i].CheckoutID,
					checkout.Items[i].ProductID,
					checkout.Items[i].Quantity,
					checkout.Items[i].Price,
					checkout.Items[i].Weight,
					checkout.Items[i].ProductName,
					nullableString(checkout.Items[i].VariantName),
					nullableString(checkout.Items[i].SKU),
					checkout.Items[i].CreatedAt,
					checkout.Items[i].UpdatedAt,
				).Scan(&checkout.Items[i].ID)
			} else {
				// If variant ID is not 0, use it in the query
				query := `
					INSERT INTO checkout_items (
						checkout_id, product_id, product_variant_id, quantity, price, 
						weight, product_name, variant_name, sku, created_at, updated_at
					) VALUES (
						$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
					) RETURNING id
				`
				err = tx.QueryRow(
					query,
					checkout.Items[i].CheckoutID,
					checkout.Items[i].ProductID,
					checkout.Items[i].ProductVariantID,
					checkout.Items[i].Quantity,
					checkout.Items[i].Price,
					checkout.Items[i].Weight,
					checkout.Items[i].ProductName,
					nullableString(checkout.Items[i].VariantName),
					nullableString(checkout.Items[i].SKU),
					checkout.Items[i].CreatedAt,
					checkout.Items[i].UpdatedAt,
				).Scan(&checkout.Items[i].ID)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetByID retrieves a checkout by ID
func (r *CheckoutRepository) GetByID(checkoutID uint) (*entity.Checkout, error) {
	// Get checkout
	query := `
		SELECT 
			id, COALESCE(user_id, 0), COALESCE(session_id, ''), status, 
			shipping_address, billing_address, COALESCE(shipping_method_id, 0), COALESCE(payment_provider, ''), 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, COALESCE(discount_code, ''), discount_amount, final_amount, applied_discount, 
			created_at, updated_at, last_activity_at, expires_at, 
			completed_at, COALESCE(converted_order_id, 0)
		FROM checkouts
		WHERE id = $1
	`

	var (
		checkout            entity.Checkout
		userID              sql.NullInt64
		sessionID           sql.NullString
		shippingMethodID    sql.NullInt64
		paymentProvider     sql.NullString
		discountCode        sql.NullString
		appliedDiscountJSON []byte
		shippingAddrJSON    []byte
		billingAddrJSON     []byte
		customerDetailsJSON []byte
		completedAt         sql.NullTime
		convertedOrderID    sql.NullInt64
	)

	err := r.db.QueryRow(query, checkoutID).Scan(
		&checkout.ID,
		&userID,
		&sessionID,
		&checkout.Status,
		&shippingAddrJSON,
		&billingAddrJSON,
		&shippingMethodID,
		&paymentProvider,
		&checkout.TotalAmount,
		&checkout.ShippingCost,
		&checkout.TotalWeight,
		&customerDetailsJSON,
		&checkout.Currency,
		&discountCode,
		&checkout.DiscountAmount,
		&checkout.FinalAmount,
		&appliedDiscountJSON,
		&checkout.CreatedAt,
		&checkout.UpdatedAt,
		&checkout.LastActivityAt,
		&checkout.ExpiresAt,
		&completedAt,
		&convertedOrderID,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("checkout not found")
	}

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if userID.Valid {
		checkout.UserID = uint(userID.Int64)
	}

	if sessionID.Valid {
		checkout.SessionID = sessionID.String
	}

	if shippingMethodID.Valid {
		checkout.ShippingMethodID = uint(shippingMethodID.Int64)
	}

	if paymentProvider.Valid {
		checkout.PaymentProvider = paymentProvider.String
	}

	if discountCode.Valid {
		checkout.DiscountCode = discountCode.String
	}

	if completedAt.Valid {
		checkout.CompletedAt = &completedAt.Time
	}

	if convertedOrderID.Valid {
		checkout.ConvertedOrderID = uint(convertedOrderID.Int64)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(shippingAddrJSON, &checkout.ShippingAddr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(billingAddrJSON, &checkout.BillingAddr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(customerDetailsJSON, &checkout.CustomerDetails); err != nil {
		return nil, err
	}

	if len(appliedDiscountJSON) > 0 {
		checkout.AppliedDiscount = &entity.AppliedDiscount{}
		if err := json.Unmarshal(appliedDiscountJSON, checkout.AppliedDiscount); err != nil {
			return nil, err
		}
	}

	// Get checkout items
	itemsQuery := `
		SELECT 
			id, checkout_id, product_id, COALESCE(product_variant_id, 0), 
			quantity, price, weight, product_name, 
			COALESCE(variant_name, ''), COALESCE(sku, ''), 
			created_at, updated_at
		FROM checkout_items
		WHERE checkout_id = $1
	`

	rows, err := r.db.Query(itemsQuery, checkout.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkout.Items = []entity.CheckoutItem{}
	for rows.Next() {
		var item entity.CheckoutItem
		var variantID sql.NullInt64
		var variantName sql.NullString
		var sku sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.CheckoutID,
			&item.ProductID,
			&variantID,
			&item.Quantity,
			&item.Price,
			&item.Weight,
			&item.ProductName,
			&variantName,
			&sku,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if variantID.Valid {
			item.ProductVariantID = uint(variantID.Int64)
		}

		if variantName.Valid {
			item.VariantName = variantName.String
		}

		if sku.Valid {
			item.SKU = sku.String
		}

		checkout.Items = append(checkout.Items, item)
	}

	return &checkout, nil
}

// GetByUserID retrieves an active checkout by user ID
func (r *CheckoutRepository) GetByUserID(userID uint) (*entity.Checkout, error) {
	query := `
		SELECT id
		FROM checkouts
		WHERE user_id = $1 AND status = 'active'
		ORDER BY last_activity_at DESC
		LIMIT 1
	`

	var checkoutID uint
	err := r.db.QueryRow(query, userID).Scan(&checkoutID)
	if err == sql.ErrNoRows {
		return nil, errors.New("checkout not found")
	}

	if err != nil {
		return nil, err
	}

	return r.GetByID(checkoutID)
}

// GetBySessionID retrieves an active checkout by session ID
func (r *CheckoutRepository) GetBySessionID(sessionID string) (*entity.Checkout, error) {
	query := `
		SELECT id
		FROM checkouts
		WHERE session_id = $1 AND status = 'active'
		ORDER BY last_activity_at DESC
		LIMIT 1
	`

	var checkoutID uint
	err := r.db.QueryRow(query, sessionID).Scan(&checkoutID)
	if err == sql.ErrNoRows {
		return nil, errors.New("checkout not found")
	}

	if err != nil {
		return nil, err
	}

	return r.GetByID(checkoutID)
}

// Update updates a checkout
func (r *CheckoutRepository) Update(checkout *entity.Checkout) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Marshal JSON fields
	shippingAddrJSON, err := json.Marshal(checkout.ShippingAddr)
	if err != nil {
		return err
	}

	billingAddrJSON, err := json.Marshal(checkout.BillingAddr)
	if err != nil {
		return err
	}

	customerDetailsJSON, err := json.Marshal(checkout.CustomerDetails)
	if err != nil {
		return err
	}

	var appliedDiscountJSON []byte
	if checkout.AppliedDiscount != nil {
		appliedDiscountJSON, err = json.Marshal(checkout.AppliedDiscount)
		if err != nil {
			return err
		}
	}

	// Update checkout
	query := `
		UPDATE checkouts 
		SET 
			status = $1, 
			shipping_address = $2, 
			billing_address = $3, 
			shipping_method_id = $4, 
			payment_provider = $5, 
			total_amount = $6, 
			shipping_cost = $7, 
			total_weight = $8, 
			customer_details = $9, 
			currency = $10, 
			discount_code = $11, 
			discount_amount = $12, 
			final_amount = $13, 
			applied_discount = $14, 
			updated_at = $15, 
			last_activity_at = $16, 
			expires_at = $17,
			completed_at = $18,
			converted_order_id = $19
		WHERE id = $20
	`

	_, err = tx.Exec(
		query,
		checkout.Status,
		shippingAddrJSON,
		billingAddrJSON,
		nullableUint(checkout.ShippingMethodID),
		nullableString(checkout.PaymentProvider),
		checkout.TotalAmount,
		checkout.ShippingCost,
		checkout.TotalWeight,
		customerDetailsJSON,
		checkout.Currency,
		nullableString(checkout.DiscountCode),
		checkout.DiscountAmount,
		checkout.FinalAmount,
		nullableBytes(appliedDiscountJSON),
		checkout.UpdatedAt,
		checkout.LastActivityAt,
		checkout.ExpiresAt,
		nullableTime(checkout.CompletedAt),
		nullableUint(checkout.ConvertedOrderID),
		checkout.ID,
	)
	if err != nil {
		return err
	}

	// Delete all existing checkout items
	_, err = tx.Exec("DELETE FROM checkout_items WHERE checkout_id = $1", checkout.ID)
	if err != nil {
		return err
	}

	// Insert checkout items
	for _, item := range checkout.Items {
		var itemQuery string
		var err error

		if item.ProductVariantID == 0 {
			// If variant ID is 0, use NULL for the database
			itemQuery = `
				INSERT INTO checkout_items (
					checkout_id, product_id, product_variant_id, quantity, price, 
					weight, product_name, variant_name, sku, created_at, updated_at
				) VALUES (
					$1, $2, NULL, $3, $4, $5, $6, $7, $8, $9, $10
				)
			`
			_, err = tx.Exec(
				itemQuery,
				checkout.ID,
				item.ProductID,
				item.Quantity,
				item.Price,
				item.Weight,
				item.ProductName,
				nullableString(item.VariantName),
				nullableString(item.SKU),
				item.CreatedAt,
				item.UpdatedAt,
			)
		} else {
			// If variant ID is not 0, use it in the query
			itemQuery = `
				INSERT INTO checkout_items (
					checkout_id, product_id, product_variant_id, quantity, price, 
					weight, product_name, variant_name, sku, created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
				)
			`
			_, err = tx.Exec(
				itemQuery,
				checkout.ID,
				item.ProductID,
				item.ProductVariantID,
				item.Quantity,
				item.Price,
				item.Weight,
				item.ProductName,
				nullableString(item.VariantName),
				nullableString(item.SKU),
				item.CreatedAt,
				item.UpdatedAt,
			)
		}
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

// Delete deletes a checkout
func (r *CheckoutRepository) Delete(checkoutID uint) error {
	query := `DELETE FROM checkouts WHERE id = $1`
	_, err := r.db.Exec(query, checkoutID)
	return err
}

// ConvertGuestCheckoutToUserCheckout converts a guest checkout to a user checkout
func (r *CheckoutRepository) ConvertGuestCheckoutToUserCheckout(sessionID string, userID uint) (*entity.Checkout, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Check if user already has a checkout
	var existingCheckoutID uint
	err = tx.QueryRow("SELECT id FROM checkouts WHERE user_id = $1 AND status = 'active'", userID).Scan(&existingCheckoutID)

	// Get the guest checkout
	guestCheckout, err := r.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if err == nil && existingCheckoutID > 0 {
		// If user already has a checkout, merge the guest checkout into the user's checkout
		userCheckout, err := r.GetByID(existingCheckoutID)
		if err != nil {
			return nil, err
		}

		// Add items from guest checkout to user checkout
		for _, item := range guestCheckout.Items {
			found := false
			// Check if the item already exists in user's checkout
			for i, userItem := range userCheckout.Items {
				if userItem.ProductID == item.ProductID && userItem.ProductVariantID == item.ProductVariantID {
					// Update quantity if product and variant already exist
					userCheckout.Items[i].Quantity += item.Quantity
					found = true
					break
				}
			}
			if !found {
				// Add new item if product and variant don't exist in user checkout
				userCheckout.Items = append(userCheckout.Items, item)
			}
		}

		// If guest checkout has shipping/billing addresses and user's doesn't, copy them
		if guestCheckout.ShippingAddr.Street != "" && userCheckout.ShippingAddr.Street == "" {
			userCheckout.ShippingAddr = guestCheckout.ShippingAddr
		}
		if guestCheckout.BillingAddr.Street != "" && userCheckout.BillingAddr.Street == "" {
			userCheckout.BillingAddr = guestCheckout.BillingAddr
		}

		// If guest checkout has customer details and user's doesn't, copy them
		if guestCheckout.CustomerDetails.Email != "" && userCheckout.CustomerDetails.Email == "" {
			userCheckout.CustomerDetails = guestCheckout.CustomerDetails
		}

		// If guest checkout has a shipping method and user's doesn't, copy it
		if guestCheckout.ShippingMethodID > 0 && userCheckout.ShippingMethodID == 0 {
			userCheckout.ShippingMethodID = guestCheckout.ShippingMethodID
			userCheckout.ShippingCost = guestCheckout.ShippingCost
		}

		// If guest checkout has a discount and user's doesn't, copy it
		if guestCheckout.AppliedDiscount != nil && userCheckout.AppliedDiscount == nil {
			userCheckout.DiscountCode = guestCheckout.DiscountCode
			userCheckout.DiscountAmount = guestCheckout.DiscountAmount
			userCheckout.AppliedDiscount = guestCheckout.AppliedDiscount
		}

		// Recalculate totals
		userCheckout.CalculateTotals()
		userCheckout.UpdatedAt = time.Now()
		userCheckout.LastActivityAt = time.Now()

		// Update the user checkout
		err = r.Update(userCheckout)
		if err != nil {
			return nil, err
		}

		// Delete the guest checkout
		err = r.Delete(guestCheckout.ID)
		if err != nil {
			return nil, err
		}

		return userCheckout, nil
	} else {
		// If user doesn't have a checkout, convert the guest checkout to a user checkout
		query := `
			UPDATE checkouts
			SET user_id = $1, session_id = NULL
			WHERE id = $2
		`
		_, err = tx.Exec(query, userID, guestCheckout.ID)
		if err != nil {
			return nil, err
		}

		guestCheckout.UserID = userID
		guestCheckout.SessionID = ""
		return guestCheckout, nil
	}
}

// GetExpiredCheckouts retrieves all checkouts that have expired
func (r *CheckoutRepository) GetExpiredCheckouts() ([]*entity.Checkout, error) {
	query := `
		SELECT id
		FROM checkouts
		WHERE status = 'active' AND expires_at < NOW()
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkoutIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		checkoutIDs = append(checkoutIDs, id)
	}

	if len(checkoutIDs) == 0 {
		return []*entity.Checkout{}, nil
	}

	checkouts := make([]*entity.Checkout, 0, len(checkoutIDs))
	for _, id := range checkoutIDs {
		checkout, err := r.GetByID(id)
		if err != nil {
			continue
		}
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetCheckoutsByStatus retrieves checkouts by status
func (r *CheckoutRepository) GetCheckoutsByStatus(status entity.CheckoutStatus, offset, limit int) ([]*entity.Checkout, error) {
	query := `
		SELECT id
		FROM checkouts
		WHERE status = $1
		ORDER BY last_activity_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkoutIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		checkoutIDs = append(checkoutIDs, id)
	}

	if len(checkoutIDs) == 0 {
		return []*entity.Checkout{}, nil
	}

	checkouts := make([]*entity.Checkout, 0, len(checkoutIDs))
	for _, id := range checkoutIDs {
		checkout, err := r.GetByID(id)
		if err != nil {
			continue
		}
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetActiveCheckoutsByUserID retrieves all active checkouts for a user
func (r *CheckoutRepository) GetActiveCheckoutsByUserID(userID uint) ([]*entity.Checkout, error) {
	query := `
		SELECT id
		FROM checkouts
		WHERE user_id = $1 AND status = 'active'
		ORDER BY last_activity_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkoutIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		checkoutIDs = append(checkoutIDs, id)
	}

	if len(checkoutIDs) == 0 {
		return []*entity.Checkout{}, nil
	}

	checkouts := make([]*entity.Checkout, 0, len(checkoutIDs))
	for _, id := range checkoutIDs {
		checkout, err := r.GetByID(id)
		if err != nil {
			continue
		}
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetCompletedCheckoutsByUserID retrieves all completed checkouts for a user
func (r *CheckoutRepository) GetCompletedCheckoutsByUserID(userID uint, offset, limit int) ([]*entity.Checkout, error) {
	query := `
		SELECT id
		FROM checkouts
		WHERE user_id = $1 AND status = 'completed'
		ORDER BY completed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkoutIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		checkoutIDs = append(checkoutIDs, id)
	}

	if len(checkoutIDs) == 0 {
		return []*entity.Checkout{}, nil
	}

	checkouts := make([]*entity.Checkout, 0, len(checkoutIDs))
	for _, id := range checkoutIDs {
		checkout, err := r.GetByID(id)
		if err != nil {
			continue
		}
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// Helper functions for handling NULL values in database

// nullableString returns a sql.NullString from a string
func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// nullableUint returns a sql.NullInt64 from a uint
func nullableUint(u uint) sql.NullInt64 {
	if u == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(u), Valid: true}
}

// nullableTime returns a sql.NullTime from a *time.Time
func nullableTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// nullableBytes returns a byte slice or nil
func nullableBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	return b
}
