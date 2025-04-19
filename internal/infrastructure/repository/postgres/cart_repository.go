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
	query := `
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
