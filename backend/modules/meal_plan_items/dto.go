package meal_plan_items

import (
	"time"

	"github.com/google/uuid"
)

type CreateMeal_plan_itemDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Day        time.Time `json:"day" validate:"required"`
	Slot       string    `json:"slot" validate:"required"`
	MealPlanID uuid.UUID `json:"meal_plan_id" validate:"required"`
	RecipeID   uuid.UUID `json:"recipe_id" validate:"required"`
}

type UpdateMeal_plan_itemDTO struct {
	ID         uuid.UUID  `json:"id" validate:"required"`
	Day        *time.Time `json:"day,omitempty"`
	Slot       *string    `json:"slot,omitempty"`
	MealPlanID *uuid.UUID `json:"meal_plan_id,omitempty"`
	RecipeID   *uuid.UUID `json:"recipe_id,omitempty"`
}

type GetMeal_plan_itemDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Day        time.Time `json:"day" validate:"required"`
	Slot       string    `json:"slot" validate:"required"`
	MealPlanID uuid.UUID `json:"meal_plan_id" validate:"required"`
	RecipeID   uuid.UUID `json:"recipe_id" validate:"required"`
}
