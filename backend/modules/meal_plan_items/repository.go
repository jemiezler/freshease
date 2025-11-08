package meal_plan_items

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetMeal_plan_itemDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetMeal_plan_itemDTO, error)
	Create(ctx context.Context, u *CreateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error)
	Update(ctx context.Context, u *UpdateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
