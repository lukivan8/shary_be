package repository

import (
	"database/sql"
	"fmt"
	"time"

	"shary_be/internal/models"

	"github.com/jmoiron/sqlx"
)

// ItemRepository handles database operations for items
type ItemRepository struct {
	db *sqlx.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create creates a new item in the database
func (r *ItemRepository) Create(item *models.Item) error {
	query := `
		INSERT INTO items (title, description, price, location, has_photos, author_id, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	return r.db.QueryRow(
		query,
		item.Title,
		item.Description,
		item.Price,
		item.Location,
		item.HasPhotos,
		item.AuthorID,
		item.CategoryID,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&item.ID)
}

// GetByID retrieves an item by ID
func (r *ItemRepository) GetByID(id int) (*models.Item, error) {
	var item models.Item
	query := `SELECT id, title, description, price, location, has_photos, author_id, category_id, created_at, updated_at FROM items WHERE id = $1`

	err := r.db.Get(&item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// GetAll retrieves all items with optional filtering
func (r *ItemRepository) GetAll(filter *models.ItemFilter) ([]models.Item, error) {
	var items []models.Item

	// Build dynamic query with filters
	query := `SELECT id, title, description, price, location FROM items WHERE 1=1`
	var args []interface{}
	argIndex := 1

	// Add filters
	if filter != nil {
		if filter.MinPrice != nil {
			query += fmt.Sprintf(" AND price >= $%d", argIndex)
			args = append(args, *filter.MinPrice)
			argIndex++
		}

		if filter.MaxPrice != nil {
			query += fmt.Sprintf(" AND price <= $%d", argIndex)
			args = append(args, *filter.MaxPrice)
			argIndex++
		}

		if filter.Location != nil && *filter.Location != "" {
			query += fmt.Sprintf(" AND LOWER(location) LIKE LOWER($%d)", argIndex)
			args = append(args, "%"+*filter.Location+"%")
			argIndex++
		}

		if filter.Search != nil && *filter.Search != "" {
			query += fmt.Sprintf(" AND (LOWER(title) LIKE LOWER($%d) OR LOWER(description) LIKE LOWER($%d))", argIndex, argIndex)
			args = append(args, "%"+*filter.Search+"%")
			argIndex++
		}
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Add pagination
	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argIndex)
			args = append(args, filter.Limit)
			argIndex++
		}

		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filter.Offset)
		}
	}

	err := r.db.Select(&items, query, args...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Update updates an item in the database
func (r *ItemRepository) Update(item *models.Item) error {
	query := `
		UPDATE items 
		SET title = $1, description = $2, price = $3, location = $4, has_photos = $5, author_id = $6, category_id = $7, updated_at = $8
		WHERE id = $9`

	item.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		item.Title,
		item.Description,
		item.Price,
		item.Location,
		item.HasPhotos,
		item.AuthorID,
		item.CategoryID,
		item.UpdatedAt,
		item.ID,
	)
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

// Delete deletes an item by ID
func (r *ItemRepository) Delete(id int) error {
	query := `DELETE FROM items WHERE id = $1`

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

// GetByLocation retrieves items by location
func (r *ItemRepository) GetByLocation(location string) ([]models.Item, error) {
	var items []models.Item
	query := `SELECT id, title, description, price, location, has_photos, author_id, category_id, created_at, updated_at FROM items WHERE LOWER(location) LIKE LOWER($1) ORDER BY created_at DESC`

	err := r.db.Select(&items, query, "%"+location+"%")
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetAvailableItems retrieves only available items
func (r *ItemRepository) GetAvailableItems() ([]models.Item, error) {
	var items []models.Item
	query := `SELECT id, title, description, price, location, has_photos, author_id, category_id, created_at, updated_at FROM items ORDER BY created_at DESC`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}
