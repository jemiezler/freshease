package product_categories

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetProductCategoryDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error)
	Create(ctx context.Context, u *CreateProductCategoryDTO) (*GetProductCategoryDTO, error)
	Update(ctx context.Context, u *UpdateProductCategoryDTO) (*GetProductCategoryDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
