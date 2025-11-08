package orders

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetOrderDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetOrderDTO, error)
	Create(ctx context.Context, u *CreateOrderDTO) (*GetOrderDTO, error)
	Update(ctx context.Context, u *UpdateOrderDTO) (*GetOrderDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
