package products

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductDTO struct {
	ID            uuid.UUID  `json:"id" validate:"required"`
	Name          string     `json:"name" validate:"required,min=2,max=60"`
	SKU           string     `json:"sku" validate:"required"`
	Price         float64    `json:"price" validate:"required,gt=0"`
	Description   *string    `json:"description,omitempty"`
	UnitLabel     string     `json:"unit_label" validate:"required"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at" validate:"required"`
	UpdatedAt     time.Time  `json:"updated_at" validate:"required"`
	Quantity      int        `json:"quantity" validate:"required,gt=0"`
	ReorderLevel  int        `json:"reorder_level" validate:"required,gt=0"`
}

type UpdateProductDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"omitempty,min=2,max=60"`
	SKU         *string    `json:"sku" validate:"omitempty"`
	Price       *float64   `json:"price" validate:"omitempty,gt=0"`
	Description *string    `json:"description,omitempty"`
	UnitLabel   *string    `json:"unit_label" validate:"omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

type GetProductDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        string     `json:"name" validate:"required"`
	SKU         string     `json:"sku" validate:"required"`
	Price       float64    `json:"price" validate:"required"`
	Description *string    `json:"description,omitempty"`
	UnitLabel   string     `json:"unit_label" validate:"required"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at" validate:"required"`
	UpdatedAt   time.Time  `json:"updated_at" validate:"required"`
}
