package repository

import (
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