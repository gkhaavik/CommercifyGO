package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// ShippingMethodRepository implements the shipping method repository interface using PostgreSQL
type ShippingMethodRepository struct {
	db *sql.DB
}

// NewShippingMethodRepository creates a new ShippingMethodRepository
func NewShippingMethodRepository(db *sql.DB) *ShippingMethodRepository {
	return &ShippingMethodRepository{db: db}
}

// Create creates a new shipping method
func (r *ShippingMethodRepository) Create(method *entity.ShippingMethod) error {
	query := `
		INSERT INTO shipping_methods (name, description, estimated_delivery_days, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		method.Name,
		method.Description,
		method.EstimatedDeliveryDays,
		method.Active,
		method.CreatedAt,
		method.UpdatedAt,
	).Scan(&method.ID)

	return err
}

// GetByID retrieves a shipping method by ID
func (r *ShippingMethodRepository) GetByID(id uint) (*entity.ShippingMethod, error) {
	query := `
		SELECT id, name, description, estimated_delivery_days, active, created_at, updated_at
		FROM shipping_methods
		WHERE id = $1
	`

	method := &entity.ShippingMethod{}
	err := r.db.QueryRow(query, id).Scan(
		&method.ID,
		&method.Name,
		&method.Description,
		&method.EstimatedDeliveryDays,
		&method.Active,
		&method.CreatedAt,
		&method.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("shipping method not found")
	}

	if err != nil {
		return nil, err
	}

	return method, nil
}

// List retrieves all shipping methods
func (r *ShippingMethodRepository) List(active bool) ([]*entity.ShippingMethod, error) {
	var query string
	var rows *sql.Rows
	var err error

	if active {
		query = `
			SELECT id, name, description, estimated_delivery_days, active, created_at, updated_at
			FROM shipping_methods
			WHERE active = true
			ORDER BY name
		`
		rows, err = r.db.Query(query)
	} else {
		query = `
			SELECT id, name, description, estimated_delivery_days, active, created_at, updated_at
			FROM shipping_methods
			ORDER BY name
		`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	methods := []*entity.ShippingMethod{}
	for rows.Next() {
		method := &entity.ShippingMethod{}
		err := rows.Scan(
			&method.ID,
			&method.Name,
			&method.Description,
			&method.EstimatedDeliveryDays,
			&method.Active,
			&method.CreatedAt,
			&method.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	return methods, nil
}

// Update updates a shipping method
func (r *ShippingMethodRepository) Update(method *entity.ShippingMethod) error {
	query := `
		UPDATE shipping_methods
		SET name = $1, description = $2, estimated_delivery_days = $3, active = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		method.Name,
		method.Description,
		method.EstimatedDeliveryDays,
		method.Active,
		time.Now(),
		method.ID,
	)

	return err
}

// Delete deletes a shipping method
func (r *ShippingMethodRepository) Delete(id uint) error {
	query := `DELETE FROM shipping_methods WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
