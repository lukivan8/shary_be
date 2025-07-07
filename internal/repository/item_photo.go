package repository

import (
	"shary_be/internal/models"

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
