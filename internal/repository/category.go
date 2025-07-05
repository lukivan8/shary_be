package repository

import (
	"database/sql"
	"time"

	"shary_be/internal/models"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT * FROM categories ORDER BY name ASC`

	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create adds a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	query := `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id`

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	return r.db.QueryRow(
		query,
		category.Name,
	).Scan(&category.ID)
}

// Update updates an existing category
func (r *CategoryRepository) Update(category *models.Category) error {
	query := `
		UPDATE categories
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id`

	now := time.Now()
	category.UpdatedAt = now

	return r.db.QueryRow(
		query,
		category.Name,
		category.UpdatedAt,
		category.ID,
	).Scan(&category.ID)
}

// Delete deletes a category by ID
func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE id = $1`

	err := r.db.Get(&category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
