package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CategoryRepository implements the category repository interface using PostgreSQL
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *entity.Category) error {
	query := `
		INSERT INTO categories (name, description, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		category.Name,
		category.Description,
		category.ParentID,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID)

	return err
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id uint) (*entity.Category, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	category := &entity.Category{}
	var parentID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&parentID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}

	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		parentIDUint := uint(parentID.Int64)
		category.ParentID = &parentIDUint
	}

	return category, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *entity.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2, parent_id = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		query,
		category.Name,
		category.Description,
		category.ParentID,
		time.Now(),
		category.ID,
	)

	return err
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves all categories
func (r *CategoryRepository) List() ([]*entity.Category, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM categories
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*entity.Category{}
	for rows.Next() {
		category := &entity.Category{}
		var parentID sql.NullInt64

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if parentID.Valid {
			parentIDUint := uint(parentID.Int64)
			category.ParentID = &parentIDUint
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetChildren retrieves child categories for a parent category
func (r *CategoryRepository) GetChildren(parentID uint) ([]*entity.Category, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM categories
		WHERE parent_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*entity.Category{}
	for rows.Next() {
		category := &entity.Category{}
		var parentIDNull sql.NullInt64

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&parentIDNull,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if parentIDNull.Valid {
			parentIDUint := uint(parentIDNull.Int64)
			category.ParentID = &parentIDUint
		}

		categories = append(categories, category)
	}

	return categories, nil
}
