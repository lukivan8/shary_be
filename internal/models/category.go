package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Category represents a category of rent items
type Category struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=50"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateCategoryRequest represents the request to create a new category
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
}

// Validate validates the struct using go-playground/validator
func (c *Category) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// Validate validates the CreateCategoryRequest
func (cc *CreateCategoryRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(cc)
}

// Validate validates the UpdateCategoryRequest
func (uc *UpdateCategoryRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(uc)
}
