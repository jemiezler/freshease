package categories

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Slug      string    `json:"slug" validate:"required,min=2,max=100"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

type UpdateCategoryDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      *string   `json:"name" validate:"omitempty,min=2,max=100"`
	Slug      *string   `json:"slug" validate:"omitempty,min=2,max=100"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

type GetCategoryDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Slug      string    `json:"slug" validate:"required"`
}

