package roles

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetRoleDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetRoleDTO, error)
	Create(ctx context.Context, u *CreateRoleDTO) (*GetRoleDTO, error)
	Update(ctx context.Context, u *UpdateRoleDTO) (*GetRoleDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
