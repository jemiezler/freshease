package vendors

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetVendorDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetVendorDTO, error)
	Create(ctx context.Context, u *CreateVendorDTO) (*GetVendorDTO, error)
	Update(ctx context.Context, u *UpdateVendorDTO) (*GetVendorDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
