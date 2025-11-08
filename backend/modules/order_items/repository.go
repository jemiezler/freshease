package order_items

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetOrder_itemDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetOrder_itemDTO, error)
	Create(ctx context.Context, u *CreateOrder_itemDTO) (*GetOrder_itemDTO, error)
	Update(ctx context.Context, u *UpdateOrder_itemDTO) (*GetOrder_itemDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
