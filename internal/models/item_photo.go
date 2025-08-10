package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ItemPhoto struct {
	ID        int       `json:"id" db:"id" validate:"required,min=1"`
	ItemID    int       `json:"item_id" db:"item_id" validate:"required,min=1"`
	URL       string    `json:"url" db:"url" validate:"required,url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateItemPhotoRequest struct {
	ItemID int      `json:"item_id" validate:"required,min=1"`
	Photos []string `json:"photos" validate:"required,min=1"`
}

type DeleteItemPhotosRequest struct {
	PhotoIDs []int `json:"photo_ids" validate:"required,min=1"`
}

func (i *ItemPhoto) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}

func (c *CreateItemPhotoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
