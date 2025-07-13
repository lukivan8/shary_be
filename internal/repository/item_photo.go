package repository

import (
	"database/sql"
	"shary_be/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
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

// Add creates a new photo for an item
func (r *ItemPhotoRepository) Add(itemID int, url string) error {
	query := `
		INSERT INTO item_photos (item_id, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`

	now := time.Now()

	_, err := r.db.Exec(query, itemID, url, now, now)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a photo by ID
func (r *ItemPhotoRepository) Delete(id int) error {
	query := `DELETE FROM item_photos WHERE id = $1`

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

// CounByItemID counts photos by item ID
func (r *ItemPhotoRepository) CountByItemID(itemID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM item_photos WHERE item_id = $1`

	err := r.db.Get(&count, query, itemID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
