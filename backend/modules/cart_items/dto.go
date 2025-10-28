package cart_items

import "github.com/google/uuid" // match schema

type CreateCart_itemDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
	Cart        uuid.UUID `json:"cart" validate:"required"`
}

type UpdateCart_itemDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"omitempty,min=2,max=60"`
	Description *string    `json:"description" validate:"omitempty"`
	Cart        *uuid.UUID `json:"cart" validate:"omitempty"`
}

type GetCart_itemDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
	Cart        uuid.UUID `json:"cart" validate:"required"`
}
