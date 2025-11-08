package categories

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetCategoryDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetCategoryDTO, error)
	Create(ctx context.Context, u *CreateCategoryDTO) (*GetCategoryDTO, error)
	Update(ctx context.Context, u *UpdateCategoryDTO) (*GetCategoryDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

