package cart_items

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetCart_itemDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetCart_itemDTO, error)
	Create(ctx context.Context, u *CreateCart_itemDTO) (*GetCart_itemDTO, error)
	Update(ctx context.Context, u *UpdateCart_itemDTO) (*GetCart_itemDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
