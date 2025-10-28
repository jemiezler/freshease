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
	Total     float64   `json:"total" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
