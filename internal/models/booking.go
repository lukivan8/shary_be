package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Rent struct {
	ID int `json:"id" db:"id"`
	ItemID int `json:"item_id" db:"item_id" validate:"required"`
	UserID int `json:"user_id" db:"user_id" validate:"required"`
	StartDate time.Time `json:"start_date" db:"start_date" validate:"required"`
	EndDate time.Time `json:"end_date" db:"end_date" validate:"required,gtfield=StartDate"`
	Price float64 `json:"price" db:"price"`
	statusID int `json:"status_id" db:"status_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BlockedPeriod struct {
	ID int `json:"id" db:"id"`
	ItemID int `json:"item_id" db:"item_id" validate:"required"`
	Reason string `json:"reason,omitempty" db:"reason"`
	StartDate time.Time `json:"start_date" db:"start_date" validate:"required"`
	EndDate time.Time `json:"end_date" db:"end_date" validate:"required,gtfield=StartDate"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type RentToUpdate struct {
	ID int `json:"id" db:"id"`
	ItemID int `json:"item_id" db:"item_id"`
	UserID int `json:"user_id" db:"user_id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate time.Time `json:"end_date" db:"end_date"`
	Price float64 `json:"price" db:"price"`
	statusID int `json:"status_id" db:"status_id"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateRentRequest struct {
	ItemID int `json:"item_id" validate:"required,min=1"`
	UserID int `json:"user_id" validate:"required,min=1"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateRentRequest struct {
	StartDate *time.Time `json:"start_date"`
	EndDate *time.Time `json:"end_date"`
}

type RentFilter struct {
	MinStartDate *time.Time `json:"min_start_date"`
	MaxStartDate *time.Time `json:"max_start_date"`
	MinEndDate *time.Time `json:"min_end_date"`
	MaxEndDate *time.Time `json:"max_end_date"`
	Price *float64 `json:"price"`
	StatusID *int `json:"status_id"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type Status struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (b *Rent) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}

func (c *CreateRentRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (u *UpdateRentRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}