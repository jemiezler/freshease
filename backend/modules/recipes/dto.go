package recipes

import "github.com/google/uuid"

type CreateRecipeDTO struct {
	ID           uuid.UUID `json:"id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Instructions *string   `json:"instructions,omitempty"`
	Kcal         int       `json:"kcal" validate:"min=0"`
}

type UpdateRecipeDTO struct {
	ID           uuid.UUID `json:"id" validate:"required"`
	Name         *string    `json:"name,omitempty"`
	Instructions *string   `json:"instructions,omitempty"`
	Kcal         *int       `json:"kcal,omitempty" validate:"omitempty,min=0"`
}

type GetRecipeDTO struct {
	ID           uuid.UUID `json:"id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Instructions *string   `json:"instructions,omitempty"`
	Kcal         int       `json:"kcal" validate:"required"`
}
