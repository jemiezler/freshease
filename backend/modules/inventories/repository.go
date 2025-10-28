package inventories

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetInventoryDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetInventoryDTO, error)
	Create(ctx context.Context, u *CreateInventoryDTO) (*GetInventoryDTO, error)
	Update(ctx context.Context, u *UpdateInventoryDTO) (*GetInventoryDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
