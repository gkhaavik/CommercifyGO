package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// ShippingRateRepository implements the shipping rate repository interface using PostgreSQL
type ShippingRateRepository struct {
	db *sql.DB
}

// NewShippingRateRepository creates a new ShippingRateRepository
func NewShippingRateRepository(db *sql.DB) *ShippingRateRepository {
	return &ShippingRateRepository{db: db}
}

// Create creates a new shipping rate
func (r *ShippingRateRepository) Create(rate *entity.ShippingRate) error {
	query := `
		INSERT INTO shipping_rates (shipping_method_id, shipping_zone_id, base_rate, min_order_value, 
			free_shipping_threshold, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var freeShippingThresholdSQL sql.NullInt64
	if rate.FreeShippingThreshold != nil {
		freeShippingThresholdSQL.Int64 = *rate.FreeShippingThreshold
		freeShippingThresholdSQL.Valid = true
	}

	err := r.db.QueryRow(
		query,
		rate.ShippingMethodID,
		rate.ShippingZoneID,
		rate.BaseRate,
		rate.MinOrderValue,
		freeShippingThresholdSQL,
		rate.Active,
		rate.CreatedAt,
		rate.UpdatedAt,
	).Scan(&rate.ID)

	return err
}

// GetByID retrieves a shipping rate by ID
func (r *ShippingRateRepository) GetByID(id uint) (*entity.ShippingRate, error) {
	// First, get the basic shipping rate data
	query := `
		SELECT id, shipping_method_id, shipping_zone_id, base_rate, min_order_value, 
			free_shipping_threshold, active, created_at, updated_at
		FROM shipping_rates
		WHERE id = $1
	`

	var freeShippingThresholdSQL sql.NullInt64
	rate := &entity.ShippingRate{
		ShippingMethod: &entity.ShippingMethod{},
		ShippingZone:   &entity.ShippingZone{},
	}

	err := r.db.QueryRow(query, id).Scan(
		&rate.ID,
		&rate.ShippingMethodID,
		&rate.ShippingZoneID,
		&rate.BaseRate,
		&rate.MinOrderValue,
		&freeShippingThresholdSQL,
		&rate.Active,
		&rate.CreatedAt,
		&rate.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("shipping rate not found")
	}

	if err != nil {
		return nil, fmt.Errorf("database error fetching shipping rate: %w", err)
	}

	// Set free shipping threshold if available
	if freeShippingThresholdSQL.Valid {
		value := freeShippingThresholdSQL.Int64
		rate.FreeShippingThreshold = &value
	}

	// Now try to get the shipping method data (if it exists)
	methodQuery := `
		SELECT name, description, estimated_delivery_days, active
		FROM shipping_methods
		WHERE id = $1
	`

	err = r.db.QueryRow(methodQuery, rate.ShippingMethodID).Scan(
		&rate.ShippingMethod.Name,
		&rate.ShippingMethod.Description,
		&rate.ShippingMethod.EstimatedDeliveryDays,
		&rate.ShippingMethod.Active,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error fetching shipping method: %w", err)
	}

	// Set shipping method ID
	rate.ShippingMethod.ID = rate.ShippingMethodID

	// Try to get the shipping zone data (if it exists)
	zoneQuery := `
		SELECT name, description, countries, states, zip_codes, active
		FROM shipping_zones
		WHERE id = $1
	`

	var countriesJSON, statesJSON, zipCodesJSON []byte

	err = r.db.QueryRow(zoneQuery, rate.ShippingZoneID).Scan(
		&rate.ShippingZone.Name,
		&rate.ShippingZone.Description,
		&countriesJSON,
		&statesJSON,
		&zipCodesJSON,
		&rate.ShippingZone.Active,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error fetching shipping zone: %w", err)
	}

	// Set shipping zone ID
	rate.ShippingZone.ID = rate.ShippingZoneID

	// Only try to unmarshal zone JSON fields if we got them
	if err != sql.ErrNoRows {
		// Unmarshal shipping zone JSON fields
		if err := json.Unmarshal(countriesJSON, &rate.ShippingZone.Countries); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(statesJSON, &rate.ShippingZone.States); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(zipCodesJSON, &rate.ShippingZone.ZipCodes); err != nil {
			return nil, err
		}
	} else {
		// Initialize empty slices
		rate.ShippingZone.Countries = []string{}
		rate.ShippingZone.States = []string{}
		rate.ShippingZone.ZipCodes = []string{}
	}

	// Get weight-based rates
	weightRates, err := r.GetWeightBasedRates(rate.ID)
	if err != nil {
		return nil, err
	}
	rate.WeightBasedRates = weightRates

	// Get value-based rates
	valueRates, err := r.GetValueBasedRates(rate.ID)
	if err != nil {
		return nil, err
	}
	rate.ValueBasedRates = valueRates

	return rate, nil
}

// GetByMethodID retrieves shipping rates by method ID
func (r *ShippingRateRepository) GetByMethodID(methodID uint) ([]*entity.ShippingRate, error) {
	query := `
		SELECT id, shipping_method_id, shipping_zone_id, base_rate, min_order_value, 
			free_shipping_threshold, active, created_at, updated_at
		FROM shipping_rates
		WHERE shipping_method_id = $1
		ORDER BY base_rate
	`

	rows, err := r.db.Query(query, methodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := []*entity.ShippingRate{}
	for rows.Next() {
		var freeShippingThresholdSQL sql.NullInt64
		rate := &entity.ShippingRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.ShippingMethodID,
			&rate.ShippingZoneID,
			&rate.BaseRate,
			&rate.MinOrderValue,
			&freeShippingThresholdSQL,
			&rate.Active,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set free shipping threshold if available
		if freeShippingThresholdSQL.Valid {
			value := freeShippingThresholdSQL.Int64
			rate.FreeShippingThreshold = &value
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

// GetByZoneID retrieves shipping rates by zone ID
func (r *ShippingRateRepository) GetByZoneID(zoneID uint) ([]*entity.ShippingRate, error) {
	query := `
		SELECT id, shipping_method_id, shipping_zone_id, base_rate, min_order_value, 
			free_shipping_threshold, active, created_at, updated_at
		FROM shipping_rates
		WHERE shipping_zone_id = $1
		ORDER BY base_rate
	`

	rows, err := r.db.Query(query, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := []*entity.ShippingRate{}
	for rows.Next() {
		var freeShippingThresholdSQL sql.NullInt64
		rate := &entity.ShippingRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.ShippingMethodID,
			&rate.ShippingZoneID,
			&rate.BaseRate,
			&rate.MinOrderValue,
			&freeShippingThresholdSQL,
			&rate.Active,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set free shipping threshold if available
		if freeShippingThresholdSQL.Valid {
			value := freeShippingThresholdSQL.Int64
			rate.FreeShippingThreshold = &value
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

// GetAvailableRatesForAddress retrieves available shipping rates for a specific address
func (r *ShippingRateRepository) GetAvailableRatesForAddress(address entity.Address, orderValue int64) ([]*entity.ShippingRate, error) {
	// First, find applicable shipping zones for this address
	query := `
		SELECT id 
		FROM shipping_zones 
		WHERE active = true AND (
			(countries @> $1::jsonb) OR 
			(states @> $2::jsonb) OR
			($3 = ANY(SELECT jsonb_array_elements_text(zip_codes)))
		)
	`

	// Convert the address data into the format needed for matching
	countryArray := []string{address.Country}
	countryJSON, err := json.Marshal(countryArray)
	if err != nil {
		return nil, err
	}

	stateArray := []string{address.State}
	stateJSON, err := json.Marshal(stateArray)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, countryJSON, stateJSON, address.PostalCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the matching zone IDs
	var zoneIDs []interface{}
	var params []string
	i := 1
	for rows.Next() {
		var zoneID uint
		if err := rows.Scan(&zoneID); err != nil {
			return nil, err
		}
		zoneIDs = append(zoneIDs, zoneID)
		params = append(params, "$"+fmt.Sprint(i))
		i++
	}

	if len(zoneIDs) == 0 {
		return nil, errors.New("no shipping zones available for this address")
	}

	// Now get the shipping rates that match these zones and where the order value meets the minimum
	ratesQuery := `
		SELECT sr.id, sr.shipping_method_id, sr.shipping_zone_id, sr.base_rate, sr.min_order_value, 
			sr.free_shipping_threshold, sr.active, sr.created_at, sr.updated_at,
			sm.name, sm.description, sm.estimated_delivery_days, sm.active
		FROM shipping_rates sr
		JOIN shipping_methods sm ON sr.shipping_method_id = sm.id
		WHERE sr.shipping_zone_id IN (` + strings.Join(params, ",") + `)
		AND sr.active = true
		AND sm.active = true
		AND sr.min_order_value <= $` + fmt.Sprint(i) + `
		ORDER BY sr.base_rate
	`

	// Add order value to query params
	args := make([]interface{}, len(zoneIDs)+1)
	copy(args, zoneIDs)
	args[len(zoneIDs)] = orderValue

	rateRows, err := r.db.Query(ratesQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rateRows.Close()

	rates := []*entity.ShippingRate{}
	for rateRows.Next() {
		var freeShippingThresholdSQL sql.NullInt64
		rate := &entity.ShippingRate{
			ShippingMethod: &entity.ShippingMethod{},
		}
		err := rateRows.Scan(
			&rate.ID,
			&rate.ShippingMethodID,
			&rate.ShippingZoneID,
			&rate.BaseRate,
			&rate.MinOrderValue,
			&freeShippingThresholdSQL,
			&rate.Active,
			&rate.CreatedAt,
			&rate.UpdatedAt,
			&rate.ShippingMethod.Name,
			&rate.ShippingMethod.Description,
			&rate.ShippingMethod.EstimatedDeliveryDays,
			&rate.ShippingMethod.Active,
		)
		if err != nil {
			return nil, err
		}

		// Set shipping method ID
		rate.ShippingMethod.ID = rate.ShippingMethodID

		// Set free shipping threshold if available
		if freeShippingThresholdSQL.Valid {
			value := freeShippingThresholdSQL.Int64
			rate.FreeShippingThreshold = &value
		}

		// Check if free shipping applies based on order value
		if rate.FreeShippingThreshold != nil && orderValue >= *rate.FreeShippingThreshold {
			rate.BaseRate = 0
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

// CreateWeightBasedRate creates a new weight-based rate
func (r *ShippingRateRepository) CreateWeightBasedRate(weightRate *entity.WeightBasedRate) error {
	query := `
		INSERT INTO weight_based_rates (shipping_rate_id, min_weight, max_weight, rate)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		weightRate.ShippingRateID,
		weightRate.MinWeight,
		weightRate.MaxWeight,
		weightRate.Rate,
	).Scan(&weightRate.ID)

	// Set default timestamps
	weightRate.CreatedAt = time.Now()
	weightRate.UpdatedAt = time.Now()

	return err
}

// CreateValueBasedRate creates a new value-based rate
func (r *ShippingRateRepository) CreateValueBasedRate(valueRate *entity.ValueBasedRate) error {
	query := `
		INSERT INTO value_based_rates (shipping_rate_id, min_order_value, max_order_value, rate)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		valueRate.ShippingRateID,
		valueRate.MinOrderValue,
		valueRate.MaxOrderValue,
		valueRate.Rate,
	).Scan(&valueRate.ID)

	// Set default timestamps
	valueRate.CreatedAt = time.Now()
	valueRate.UpdatedAt = time.Now()

	return err
}

// GetWeightBasedRates retrieves weight-based rates for a shipping rate
func (r *ShippingRateRepository) GetWeightBasedRates(rateID uint) ([]entity.WeightBasedRate, error) {
	query := `
		SELECT id, shipping_rate_id, min_weight, max_weight, rate
		FROM weight_based_rates
		WHERE shipping_rate_id = $1
		ORDER BY min_weight
	`

	rows, err := r.db.Query(query, rateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := []entity.WeightBasedRate{}
	for rows.Next() {
		rate := entity.WeightBasedRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.ShippingRateID,
			&rate.MinWeight,
			&rate.MaxWeight,
			&rate.Rate,
		)
		if err != nil {
			return nil, err
		}

		// Set default timestamps since they're not in the DB
		rate.CreatedAt = time.Now()
		rate.UpdatedAt = time.Now()

		rates = append(rates, rate)
	}

	return rates, nil
}

// GetValueBasedRates retrieves value-based rates for a shipping rate
func (r *ShippingRateRepository) GetValueBasedRates(rateID uint) ([]entity.ValueBasedRate, error) {
	query := `
		SELECT id, shipping_rate_id, min_order_value, max_order_value, rate
		FROM value_based_rates
		WHERE shipping_rate_id = $1
		ORDER BY min_order_value
	`

	rows, err := r.db.Query(query, rateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := []entity.ValueBasedRate{}
	for rows.Next() {
		rate := entity.ValueBasedRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.ShippingRateID,
			&rate.MinOrderValue,
			&rate.MaxOrderValue,
			&rate.Rate,
		)
		if err != nil {
			return nil, err
		}

		// Set default timestamps since they're not in the DB
		rate.CreatedAt = time.Now()
		rate.UpdatedAt = time.Now()

		rates = append(rates, rate)
	}

	return rates, nil
}

// Update updates a shipping rate
func (r *ShippingRateRepository) Update(rate *entity.ShippingRate) error {
	query := `
		UPDATE shipping_rates
		SET shipping_method_id = $1, shipping_zone_id = $2, base_rate = $3, min_order_value = $4,
			free_shipping_threshold = $5, active = $6, updated_at = $7
		WHERE id = $8
	`

	var freeShippingThresholdSQL sql.NullInt64
	if rate.FreeShippingThreshold != nil {
		freeShippingThresholdSQL.Int64 = *rate.FreeShippingThreshold
		freeShippingThresholdSQL.Valid = true
	}

	_, err := r.db.Exec(
		query,
		rate.ShippingMethodID,
		rate.ShippingZoneID,
		rate.BaseRate,
		rate.MinOrderValue,
		freeShippingThresholdSQL,
		rate.Active,
		time.Now(),
		rate.ID,
	)

	return err
}

// Delete deletes a shipping rate
func (r *ShippingRateRepository) Delete(id uint) error {
	// Start a transaction to delete related records as well
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete weight-based rates first
	_, err = tx.Exec("DELETE FROM weight_based_rates WHERE shipping_rate_id = $1", id)
	if err != nil {
		return err
	}

	// Delete value-based rates
	_, err = tx.Exec("DELETE FROM value_based_rates WHERE shipping_rate_id = $1", id)
	if err != nil {
		return err
	}

	// Delete the shipping rate itself
	_, err = tx.Exec("DELETE FROM shipping_rates WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
