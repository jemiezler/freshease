package bundles

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetBundleDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetBundleDTO, error)
	Create(ctx context.Context, u *CreateBundleDTO) (*GetBundleDTO, error)
	Update(ctx context.Context, u *UpdateBundleDTO) (*GetBundleDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
