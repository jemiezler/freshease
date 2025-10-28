package users

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetUserDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetUserDTO, error)
	Create(ctx context.Context, u *CreateUserDTO) (*GetUserDTO, error)
	Update(ctx context.Context, u *UpdateUserDTO) (*GetUserDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
