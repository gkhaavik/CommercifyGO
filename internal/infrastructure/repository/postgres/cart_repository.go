package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// CartRepository implements the cart repository interface using PostgreSQL
type CartRepository struct {
	db *sql.DB
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// Create creates a new cart
func (r *CartRepository) Create(cart *entity.Cart) error {
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

	// Insert cart
	var query string
	if cart.SessionID != "" {
		// Guest cart
		query = `
			INSERT INTO carts (session_id, created_at, updated_at)
			VALUES ($1, $2, $3)
			RETURNING id
		`
		err = tx.QueryRow(
			query,
			cart.SessionID,
			cart.CreatedAt,
			cart.UpdatedAt,
		).Scan(&cart.ID)
	} else {
		// User cart
		query = `
			INSERT INTO carts (user_id, created_at, updated_at)
			VALUES ($1, $2, $3)
			RETURNING id
		`
		err = tx.QueryRow(
			query,
			cart.UserID,
			cart.CreatedAt,
			cart.UpdatedAt,
		).Scan(&cart.ID)
	}

	if err != nil {
		return err
	}

	// Insert cart items if any
	if len(cart.Items) > 0 {
		for i := range cart.Items {
			cart.Items[i].CartID = cart.ID
			query := `
				INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`
			err = tx.QueryRow(
				query,
				cart.Items[i].CartID,
				cart.Items[i].ProductID,
				cart.Items[i].Quantity,
				cart.Items[i].CreatedAt,
				cart.Items[i].UpdatedAt,
			).Scan(&cart.Items[i].ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetByUserID retrieves a cart by user ID
func (r *CartRepository) GetByUserID(userID uint) (*entity.Cart, error) {
	// Get cart
	query := `
		SELECT id, user_id, created_at, updated_at
		FROM carts
		WHERE user_id = $1
	`

	cart := &entity.Cart{}
	err := r.db.QueryRow(query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("cart not found")
	}

	if err != nil {
		return nil, err
	}

	// Get cart items
	query = `
		SELECT id, cart_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
	`

	rows, err := r.db.Query(query, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart.Items = []entity.CartItem{}
	for rows.Next() {
		item := entity.CartItem{}
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cart.Items = append(cart.Items, item)
	}

	return cart, nil
}

// GetBySessionID retrieves a cart by session ID
func (r *CartRepository) GetBySessionID(sessionID string) (*entity.Cart, error) {
	// Get cart
	query := `
		SELECT id, session_id, created_at, updated_at
		FROM carts
		WHERE session_id = $1
	`

	cart := &entity.Cart{}
	err := r.db.QueryRow(query, sessionID).Scan(
		&cart.ID,
		&cart.SessionID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("cart not found")
	}

	if err != nil {
		return nil, err
	}

	// Get cart items
	query = `
		SELECT id, cart_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
	`

	rows, err := r.db.Query(query, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart.Items = []entity.CartItem{}
	for rows.Next() {
		item := entity.CartItem{}
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cart.Items = append(cart.Items, item)
	}

	return cart, nil
}

// Update updates a cart
func (r *CartRepository) Update(cart *entity.Cart) error {
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

	// Update cart
	query := `
		UPDATE carts
		SET updated_at = $1
		WHERE id = $2
	`

	_, err = tx.Exec(
		query,
		time.Now(),
		cart.ID,
	)
	if err != nil {
		return err
	}

	// Delete all cart items
	query = `DELETE FROM cart_items WHERE cart_id = $1`
	_, err = tx.Exec(query, cart.ID)
	if err != nil {
		return err
	}

	// Insert new cart items
	for i := range cart.Items {
		cart.Items[i].CartID = cart.ID
		now := time.Now()
		query := `
			INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`
		err = tx.QueryRow(
			query,
			cart.Items[i].CartID,
			cart.Items[i].ProductID,
			cart.Items[i].Quantity,
			now,
			now,
		).Scan(&cart.Items[i].ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a cart
func (r *CartRepository) Delete(id uint) error {
	query := `DELETE FROM carts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ConvertGuestCartToUserCart converts a guest cart to a user cart
func (r *CartRepository) ConvertGuestCartToUserCart(sessionID string, userID uint) (*entity.Cart, error) {
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

	// Check if user already has a cart
	var existingCartID uint
	err = tx.QueryRow("SELECT id FROM carts WHERE user_id = $1", userID).Scan(&existingCartID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get the guest cart
	guestCart, err := r.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if err == nil && existingCartID > 0 {
		// If user already has a cart, merge the guest cart into the user's cart
		// First, get the user's cart
		userCart, err := r.GetByUserID(userID)
		if err != nil {
			return nil, err
		}

		// Add items from guest cart to user cart
		for _, item := range guestCart.Items {
			found := false
			for i, userItem := range userCart.Items {
				if userItem.ProductID == item.ProductID {
					// Update quantity if product already exists
					userCart.Items[i].Quantity += item.Quantity
					found = true
					break
				}
			}
			if !found {
				// Add new item if product doesn't exist in user cart
				userCart.Items = append(userCart.Items, entity.CartItem{
					CartID:    userCart.ID,
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			}
		}

		// Update the user cart
		err = r.Update(userCart)
		if err != nil {
			return nil, err
		}

		// Delete the guest cart
		err = r.Delete(guestCart.ID)
		if err != nil {
			return nil, err
		}

		return userCart, nil
	} else {
		// If user doesn't have a cart, convert the guest cart to a user cart
		query := `
			UPDATE carts
			SET user_id = $1, session_id = NULL
			WHERE id = $2
		`
		_, err = tx.Exec(query, userID, guestCart.ID)
		if err != nil {
			return nil, err
		}

		guestCart.UserID = userID
		guestCart.SessionID = ""
		return guestCart, nil
	}
}
