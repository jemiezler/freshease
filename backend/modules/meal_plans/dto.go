package meal_plans

import (
	"time"

	"github.com/google/uuid"
)

type CreateMeal_planDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	WeekStart time.Time `json:"week_start" validate:"required"`
	Goal      *string   `json:"goal,omitempty"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
}

type UpdateMeal_planDTO struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	WeekStart *time.Time `json:"week_start,omitempty"`
	Goal      *string    `json:"goal,omitempty"`
}

type GetMeal_planDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	WeekStart time.Time `json:"week_start" validate:"required"`
	Goal      *string   `json:"goal,omitempty"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
}
