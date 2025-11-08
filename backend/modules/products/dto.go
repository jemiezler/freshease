package products

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductDTO struct {
	ID            uuid.UUID  `json:"id" validate:"required"`
	Name          string     `json:"name" validate:"required,min=2,max=60"`
	Price         float64    `json:"price" validate:"required,gt=0"`
	Description   string     `json:"description" validate:"required"`
	ImageURL      string     `json:"image_url" validate:"required,url"`
	UnitLabel     string     `json:"unit_label" validate:"required"`
	IsActive      string     `json:"is_active" validate:"required"`
	CreatedAt     time.Time  `json:"created_at" validate:"required"`
	UpdatedAt     time.Time  `json:"updated_at" validate:"required"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" validate:"omitempty"`
	Quantity      int        `json:"quantity" validate:"required,gt=0"`
	RestockAmount int        `json:"restock_amount" validate:"required,gt=0"`
}

type UpdateProductDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"omitempty,min=2,max=60"`
	Price       *float64   `json:"price" validate:"omitempty,gt=0"`
	Description *string    `json:"description" validate:"omitempty"`
	ImageURL    *string    `json:"image_url" validate:"omitempty,url"`
	UnitLabel   *string    `json:"unit_label" validate:"omitempty"`
	IsActive    *string    `json:"is_active" validate:"omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" validate:"omitempty"`
}

type GetProductDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        string     `json:"name" validate:"required,min=2,max=60"`
	Price       float64    `json:"price" validate:"required,gt=0"`
	Description string     `json:"description" validate:"required"`
	ImageURL    string     `json:"image_url" validate:"required,url"`
	UnitLabel   string     `json:"unit_label" validate:"required"`
	IsActive    string     `json:"is_active" validate:"required"`
	CreatedAt   time.Time  `json:"created_at" validate:"required"`
	UpdatedAt   time.Time  `json:"updated_at" validate:"required"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" validate:"omitempty"`
}
