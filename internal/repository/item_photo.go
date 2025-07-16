package repository

import (
	"database/sql"
	"shary_be/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ItemPhotoRepository struct {
	db *sqlx.DB
}

func NewItemPhotoRepository(db *sqlx.DB) *ItemPhotoRepository {
	return &ItemPhotoRepository{db: db}
}

// GetPhotosByItemID retrieves all photos for an item
func (r *ItemPhotoRepository) GetPhotosByItemID(itemID int) ([]models.ItemPhoto, error) {
	var photos []models.ItemPhoto
	query := `SELECT * FROM item_photos WHERE item_id = $1`

	err := r.db.Select(&photos, query, itemID)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

// Add adds a new photos for an item
func (r *ItemPhotoRepository) Add(tx *sqlx.Tx, itemID int, urls []string) error {
	query := `
		INSERT INTO item_photos (item_id, url)
		SELECT $1, unnest($2::text[])`

	result, err := tx.Exec(query, itemID, pq.Array(urls))
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

// Delete deletes a photos by ID
func (r *ItemPhotoRepository) Delete(tx *sqlx.Tx, ids []int) error {
	query := `DELETE FROM item_photos WHERE id = ANY($1)`

	result, err := tx.Exec(query, pq.Array(ids))
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

// CountByItemID counts photos by item ID, using the provided querier (e.g., tx or db).
func (r *ItemPhotoRepository) CountByItemID(querier sqlx.Ext, itemID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM item_photos WHERE item_id = $1`

	err := sqlx.Get(querier, &count, query, itemID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
