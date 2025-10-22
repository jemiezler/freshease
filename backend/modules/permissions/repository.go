package permissions

import (
	"context"

	"github.com/google/uuid"
)

type Permission struct {
	ID    uuid.UUID `json:"id"`
	Email string     `json:"email"`
	Name  string     `json:"name"`
}

type Repository interface {
	List(ctx context.Context) ([]*Permission, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	Create(ctx context.Context, u *Permission) error
	Update(ctx context.Context, u *Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}
