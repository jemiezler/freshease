package permissions

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetPermissionDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetPermissionDTO, error)
	Create(ctx context.Context, u *CreatePermissionDTO) (*GetPermissionDTO, error)
	Update(ctx context.Context, u *UpdatePermissionDTO) (*GetPermissionDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
