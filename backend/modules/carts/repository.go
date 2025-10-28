package carts

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetCartDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetCartDTO, error)
	Create(ctx context.Context, u *CreateCartDTO) (*GetCartDTO, error)
	Update(ctx context.Context, u *UpdateCartDTO) (*GetCartDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
