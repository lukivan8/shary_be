package repository

import (
	"database/sql"
	"fmt"
	"strings"
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
func (r *ItemRepository) Create(item *models.Item, photos []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	itemQuery := `
		INSERT INTO items (title, description, price, location, has_photos, author_id, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	item.HasPhotos = len(photos) > 0

	err = tx.QueryRow(
		itemQuery,
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

	if err != nil {
		return err
	}

	if item.HasPhotos {
		photoQuery := `
			INSERT INTO item_photos (item_id, url, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`

		for _, photoURL := range photos {
			_, err := tx.Exec(
				photoQuery,
				item.ID,
				photoURL,
				now,
				now,
			)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// GetByID retrieves an item by ID
func (r *ItemRepository) GetByID(id int) (*models.ItemResponse, error) {
	var item models.ItemResponse
	query := `
		SELECT 
			i.id, i.title, i.description, i.price, i.location, i.has_photos,
			i.author_id, 
			c.id AS "category.id", c.name AS "category.name",
			i.created_at, i.updated_at,
			COALESCE(array_agg(p.url) FILTER (WHERE p.url IS NOT NULL), '{}') AS photos
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.id
		LEFT JOIN item_photos p ON i.id = p.item_id
		WHERE i.id = $1
		GROUP BY i.id, c.id, c.name
	`
	err := r.db.Get(&item, query, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetAll retrieves all items with optional filtering
func (r *ItemRepository) GetAll(filter *models.ItemFilter) ([]models.ItemResponse, error) {
	var items []models.ItemResponse

	var queryBuilder strings.Builder

	// Build dynamic query with filters
	queryBuilder.WriteString(`
        SELECT
            i.id, i.title, i.description, i.price, i.location, i.has_photos, i.author_id, i.created_at, i.updated_at,
            c.id AS "category.id",
            c.name AS "category.name"
        FROM items i
        LEFT JOIN categories c ON i.category_id = c.id
        WHERE 1=1
    `)

	var args []interface{}

	// Add filters
	if filter != nil {
		if filter.MinPrice != nil {
			queryBuilder.WriteString(" AND i.price >= ?")
			args = append(args, *filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			queryBuilder.WriteString(" AND i.price <= ?")
			args = append(args, *filter.MaxPrice)
		}
		if filter.Location != nil && *filter.Location != "" {
			queryBuilder.WriteString(" AND LOWER(i.location) LIKE LOWER(?)")
			args = append(args, "%"+*filter.Location+"%")
		}
		if filter.Search != nil && *filter.Search != "" {
			queryBuilder.WriteString(" AND (LOWER(i.title) LIKE LOWER(?) OR LOWER(i.description) LIKE LOWER(?))")
			searchTerm := "%" + *filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
		}
		if filter.CategoryID != nil {
			queryBuilder.WriteString(" AND i.category_id = ?")
			args = append(args, *filter.CategoryID)
		}
	}

	// Add ordering
	queryBuilder.WriteString(" ORDER BY i.created_at DESC")

	// Add pagination
	if filter != nil {
		if filter.Limit > 0 {
			queryBuilder.WriteString(" LIMIT ?")
			args = append(args, filter.Limit)
		}
		if filter.Offset > 0 {
			queryBuilder.WriteString(" OFFSET ?")
			args = append(args, filter.Offset)
		}
	}

	query := r.db.Rebind(queryBuilder.String())

	err := r.db.Select(&items, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get all items with filter: %w", err)
	}

	return items, nil
}

// Update updates an item in the database
func (r *ItemRepository) Update(tx *sql.Tx, item *models.ItemToUpdate) error {
	query := `
		UPDATE items 
		SET title = $1, description = $2, price = $3, location = $4, has_photos = $5, category_id = $6, updated_at = $7
		WHERE id = $8`

	item.UpdatedAt = time.Now()

	result, err := tx.Exec(query,
		item.Title,
		item.Description,
		item.Price,
		item.Location,
		item.HasPhotos,
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
func (r *ItemRepository) GetByLocation(location string) ([]models.ItemResponse, error) {
	var items []models.ItemResponse
	query := `SELECT i.id, i.title, i.description, i.price, i.location, i.has_photos, i.author_id, i.category_id, i.created_at, i.updated_at, c.id AS "category.id", c.name AS "category.name" FROM items i LEFT JOIN categories c ON i.category_id = c.id WHERE LOWER(i.location) LIKE LOWER($1) ORDER BY i.created_at DESC`

	err := r.db.Select(&items, query, "%"+location+"%")
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetAvailableItems retrieves only available items
func (r *ItemRepository) GetAvailableItems() ([]models.ItemResponse, error) {
	var items []models.ItemResponse
	query := `SELECT i.id, i.title, i.description, i.price, i.location, i.has_photos, i.author_id, i.category_id, i.created_at, i.updated_at, c.id AS "category.id", c.name AS "category.name" FROM items i LEFT JOIN categories c ON i.category_id = c.id ORDER BY i.created_at DESC`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetByCategory gets items by category with category info
func (r *ItemRepository) GetByCategory(categoryID int) ([]models.ItemResponse, error) {
	var items []models.ItemResponse

	const query = `
        SELECT
            i.id,
            i.title,
            i.description,
            i.price,
            i.location,
            i.has_photos,
            i.author_id,
            i.created_at,
            i.updated_at,
            c.id AS "category.id",
            c.name AS "category.name"
        FROM
            items i
        INNER JOIN
            categories c ON i.category_id = c.id
        WHERE
            i.category_id = $1
        ORDER BY
            i.created_at DESC`

	err := r.db.Select(&items, query, categoryID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetPhotosByItemID retrieves all photos for an item
func (r *ItemRepository) GetPhotosByItemID(itemID int) ([]models.ItemPhoto, error) {
	var photos []models.ItemPhoto
	query := `SELECT * FROM item_photos WHERE item_id = $1`

	err := r.db.Select(&photos, query, itemID)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

// AddPhotos inserts new photos for an item within a transaction.
func (r *ItemRepository) AddPhotos(tx *sql.Tx, itemID int, photoURLs []string) error {
	if len(photoURLs) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(`INSERT INTO item_photos (item_id, url) VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, url := range photoURLs {
		if _, err := stmt.Exec(itemID, url); err != nil {
			return err
		}
	}
	return nil
}

// DeletePhotos removes photos by their IDs within a transaction.
func (r *ItemRepository) DeletePhotos(tx *sql.Tx, photoIDs []int, itemID int) error {
	if len(photoIDs) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(`DELETE FROM item_photos WHERE id = $1 AND item_id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, photoID := range photoIDs {
		if _, err := stmt.Exec(photoID, itemID); err != nil {
			return err
		}
	}

	return nil
}

// CountPhotosByItemID counts the total number of photos for an item.
func (r *ItemRepository) CountPhotosByItemID(tx *sql.Tx, itemID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM item_photos WHERE item_id = $1"
	err := tx.QueryRow(query, itemID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
