package inventories

import (
	"time"

	"github.com/google/uuid"
)

type CreateInventoryDTO struct {
	Quantity      int       `json:"quantity" validate:"required,gt=0"`
	ReorderLevel int       `json:"reorder_level" validate:"required,gt=0"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type UpdateInventoryDTO struct {
	ID            uuid.UUID  `json:"id" validate:"required"`
	Quantity      *int       `json:"quantity" validate:"omitempty,gt=0"`
	ReorderLevel *int       `json:"reorder_level" validate:"omitempty,gt=0"`
	UpdatedAt     *time.Time `json:"updated_at" validate:"omitempty"`
}

type GetInventoryDTO struct {
	ID            uuid.UUID `json:"id" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required,gt=0"`
	ReorderLevel int       `json:"reorder_level" validate:"required,gt=0"`
	UpdatedAt     time.Time `json:"updated_at" validate:"required"`
}
