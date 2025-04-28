package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// ShippingZoneRepository implements the shipping zone repository interface using PostgreSQL
type ShippingZoneRepository struct {
	db *sql.DB
}

// NewShippingZoneRepository creates a new ShippingZoneRepository
func NewShippingZoneRepository(db *sql.DB) *ShippingZoneRepository {
	return &ShippingZoneRepository{db: db}
}

// Create creates a new shipping zone
func (r *ShippingZoneRepository) Create(zone *entity.ShippingZone) error {
	query := `
		INSERT INTO shipping_zones (name, description, countries, states, zip_codes, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	countriesJSON, err := json.Marshal(zone.Countries)
	if err != nil {
		return err
	}

	statesJSON, err := json.Marshal(zone.States)
	if err != nil {
		return err
	}

	zipCodesJSON, err := json.Marshal(zone.ZipCodes)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		zone.Name,
		zone.Description,
		countriesJSON,
		statesJSON,
		zipCodesJSON,
		zone.Active,
		zone.CreatedAt,
		zone.UpdatedAt,
	).Scan(&zone.ID)

	return err
}

// GetByID retrieves a shipping zone by ID
func (r *ShippingZoneRepository) GetByID(id uint) (*entity.ShippingZone, error) {
	query := `
		SELECT id, name, description, countries, states, zip_codes, active, created_at, updated_at
		FROM shipping_zones
		WHERE id = $1
	`

	var countriesJSON, statesJSON, zipCodesJSON []byte
	zone := &entity.ShippingZone{}
	err := r.db.QueryRow(query, id).Scan(
		&zone.ID,
		&zone.Name,
		&zone.Description,
		&countriesJSON,
		&statesJSON,
		&zipCodesJSON,
		&zone.Active,
		&zone.CreatedAt,
		&zone.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("shipping zone not found")
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON fields
	if err := json.Unmarshal(countriesJSON, &zone.Countries); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(statesJSON, &zone.States); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(zipCodesJSON, &zone.ZipCodes); err != nil {
		return nil, err
	}

	return zone, nil
}

// List retrieves all shipping zones
func (r *ShippingZoneRepository) List(active bool) ([]*entity.ShippingZone, error) {
	var query string
	var rows *sql.Rows
	var err error

	if active {
		query = `
			SELECT id, name, description, countries, states, zip_codes, active, created_at, updated_at
			FROM shipping_zones
			WHERE active = true
			ORDER BY name
		`
		rows, err = r.db.Query(query)
	} else {
		query = `
			SELECT id, name, description, countries, states, zip_codes, active, created_at, updated_at
			FROM shipping_zones
			ORDER BY name
		`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	zones := []*entity.ShippingZone{}
	for rows.Next() {
		var countriesJSON, statesJSON, zipCodesJSON []byte
		zone := &entity.ShippingZone{}
		err := rows.Scan(
			&zone.ID,
			&zone.Name,
			&zone.Description,
			&countriesJSON,
			&statesJSON,
			&zipCodesJSON,
			&zone.Active,
			&zone.CreatedAt,
			&zone.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal the JSON fields
		if err := json.Unmarshal(countriesJSON, &zone.Countries); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(statesJSON, &zone.States); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(zipCodesJSON, &zone.ZipCodes); err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

// Update updates a shipping zone
func (r *ShippingZoneRepository) Update(zone *entity.ShippingZone) error {
	query := `
		UPDATE shipping_zones
		SET name = $1, description = $2, countries = $3, states = $4, zip_codes = $5, 
			active = $6, updated_at = $7
		WHERE id = $8
	`

	countriesJSON, err := json.Marshal(zone.Countries)
	if err != nil {
		return err
	}

	statesJSON, err := json.Marshal(zone.States)
	if err != nil {
		return err
	}

	zipCodesJSON, err := json.Marshal(zone.ZipCodes)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		zone.Name,
		zone.Description,
		countriesJSON,
		statesJSON,
		zipCodesJSON,
		zone.Active,
		time.Now(),
		zone.ID,
	)

	return err
}

// Delete deletes a shipping zone
func (r *ShippingZoneRepository) Delete(id uint) error {
	query := `DELETE FROM shipping_zones WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
