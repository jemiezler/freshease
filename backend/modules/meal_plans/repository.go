package meal_plans

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetMeal_planDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetMeal_planDTO, error)
	Create(ctx context.Context, u *CreateMeal_planDTO) (*GetMeal_planDTO, error)
	Update(ctx context.Context, u *UpdateMeal_planDTO) (*GetMeal_planDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
