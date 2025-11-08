package recipe_items

import "github.com/google/uuid"

type CreateRecipe_itemDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,min=0"`
	Unit      string    `json:"unit" validate:"required"`
	RecipeID  uuid.UUID `json:"recipe_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

type UpdateRecipe_itemDTO struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	Amount    *float64   `json:"amount,omitempty" validate:"omitempty,min=0"`
	Unit      *string    `json:"unit,omitempty"`
	RecipeID  *uuid.UUID `json:"recipe_id,omitempty"`
	ProductID *uuid.UUID `json:"product_id,omitempty"`
}

type GetRecipe_itemDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	Unit      string    `json:"unit" validate:"required"`
	RecipeID  uuid.UUID `json:"recipe_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}
