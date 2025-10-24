package addresses

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetAddressDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetAddressDTO, error)
	Create(ctx context.Context, u *CreateAddressDTO) (*GetAddressDTO, error)
	Update(ctx context.Context, u *UpdateAddressDTO) (*GetAddressDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
