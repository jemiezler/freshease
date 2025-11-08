package recipes

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetRecipeDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetRecipeDTO, error)
	Create(ctx context.Context, u *CreateRecipeDTO) (*GetRecipeDTO, error)
	Update(ctx context.Context, u *UpdateRecipeDTO) (*GetRecipeDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
