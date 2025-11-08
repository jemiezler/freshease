package bundle_items

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetBundle_itemDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetBundle_itemDTO, error)
	Create(ctx context.Context, u *CreateBundle_itemDTO) (*GetBundle_itemDTO, error)
	Update(ctx context.Context, u *UpdateBundle_itemDTO) (*GetBundle_itemDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
