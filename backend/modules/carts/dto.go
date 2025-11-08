package carts

import (
	"time"

	"github.com/google/uuid"
)

type CreateCartDTO struct {
	Status *string  `json:"status,omitempty" validate:"omitempty"`
	Total  *float64 `json:"total,omitempty" validate:"omitempty"`
}

type UpdateCartDTO struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	Status *string   `json:"status,omitempty" validate:"omitempty"`
	Total  *float64  `json:"total,omitempty" validate:"omitempty"`
}

type GetCartDTO struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid"`
	Status    string    `json:"status" validate:"required"`
	Subtotal  float64   `json:"subtotal" validate:"required"`
	Discount  float64   `json:"discount" validate:"required"`
	Total     float64   `json:"total" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
