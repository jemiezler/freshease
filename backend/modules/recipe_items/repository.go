package recipe_items

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetRecipe_itemDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetRecipe_itemDTO, error)
	Create(ctx context.Context, u *CreateRecipe_itemDTO) (*GetRecipe_itemDTO, error)
	Update(ctx context.Context, u *UpdateRecipe_itemDTO) (*GetRecipe_itemDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
