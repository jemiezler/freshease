package products

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetProductDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetProductDTO, error)
	Create(ctx context.Context, u *CreateProductDTO) (*GetProductDTO, error)
	Update(ctx context.Context, u *UpdateProductDTO) (*GetProductDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
