package bundles

import "github.com/google/uuid"

type CreateBundleDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,min=0"`
	IsActive    bool      `json:"is_active"`
}

type UpdateBundleDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=2,max=60"`
	Description *string   `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty" validate:"omitempty,min=0"`
	IsActive    *bool     `json:"is_active,omitempty"`
}

type GetBundleDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required"`
	IsActive    bool      `json:"is_active" validate:"required"`
}
