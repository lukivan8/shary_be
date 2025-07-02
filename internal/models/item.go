package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Item represents an item available for rent
type Item struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required,min=1,max=200"`
	Description string    `json:"description" db:"description" validate:"required,min=10,max=2000"`
	Price       float64   `json:"price" db:"price" validate:"required,min=0.01"`
	Location    string    `json:"location" db:"location" validate:"required,min=1,max=500"`
	HasPhotos   bool      `json:"has_photos" db:"has_photos"`
	AuthorID    int       `json:"author_id" db:"author_id"`
	CategoryID  int       `json:"category_id" db:"category_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemRequest represents the request to create a new item
type CreateItemRequest struct {
	Title       string   `json:"title" validate:"required,min=1,max=200"`
	Description string   `json:"description" validate:"required,min=10,max=2000"`
	Price       float64  `json:"price" validate:"required,min=0.01"`
	Location    string   `json:"location" validate:"required,min=1,max=500"`
	Photos      []string `json:"photos" validate:"omitempty,min=1,max=10"`
	CategoryID  int      `json:"category_id" validate:"required,min=1"`
	AuthorID    int      `json:"author_id" validate:"required,min=1"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Title       *string   `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=10,max=2000"`
	Price       *float64  `json:"price,omitempty" validate:"omitempty,min=0.01"`
	Location    *string   `json:"location,omitempty" validate:"omitempty,min=1,max=500"`
	CategoryID  *int      `json:"category_id,omitempty" validate:"omitempty,min=1"`
	Photos      *[]string `json:"photos,omitempty" validate:"omitempty,min=1,max=10"`
}

// ItemFilter represents filters for listing items
type ItemFilter struct {
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	Location   *string  `json:"location,omitempty"`
	Search     *string  `json:"search,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
	CategoryID *int     `json:"category_id,omitempty" validate:"omitempty,min=1"`
}

// Validate validates the struct using go-playground/validator
func (i *Item) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}

// Validate validates the CreateItemRequest
func (c *CreateItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// Validate validates the UpdateItemRequest
func (u *UpdateItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
