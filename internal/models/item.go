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
	ImageURL    string    `json:"image_url" db:"image_url" validate:"required,url"`
	PricePerDay float64   `json:"price_per_day" db:"price_per_day" validate:"required,min=0.01"`
	Location    string    `json:"location" db:"location" validate:"required,min=1,max=500"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemRequest represents the request to create a new item
type CreateItemRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=200"`
	Description string  `json:"description" validate:"required,min=10,max=2000"`
	ImageURL    string  `json:"image_url" validate:"required,url"`
	PricePerDay float64 `json:"price_per_day" validate:"required,min=0.01"`
	Location    string  `json:"location" validate:"required,min=1,max=500"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=10,max=2000"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"omitempty,url"`
	PricePerDay *float64 `json:"price_per_day,omitempty" validate:"omitempty,min=0.01"`
	Location    *string  `json:"location,omitempty" validate:"omitempty,min=1,max=500"`
	IsAvailable *bool    `json:"is_available,omitempty"`
}

// ItemFilter represents filters for listing items
type ItemFilter struct {
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	Location    *string  `json:"location,omitempty"`
	IsAvailable *bool    `json:"is_available,omitempty"`
	Search      *string  `json:"search,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	Offset      int      `json:"offset,omitempty"`
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
