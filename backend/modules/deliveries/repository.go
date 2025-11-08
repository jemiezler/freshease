package deliveries

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetDeliveryDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetDeliveryDTO, error)
	Create(ctx context.Context, u *CreateDeliveryDTO) (*GetDeliveryDTO, error)
	Update(ctx context.Context, u *UpdateDeliveryDTO) (*GetDeliveryDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

