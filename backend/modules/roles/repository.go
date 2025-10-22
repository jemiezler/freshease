package roles

import (
	"context"

	"github.com/google/uuid"
)

type Role struct {
	ID    uuid.UUID `json:"id"`
	Email string     `json:"email"`
	Name  string     `json:"name"`
}

type Repository interface {
	List(ctx context.Context) ([]*Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Role, error)
	Create(ctx context.Context, u *Role) error
	Update(ctx context.Context, u *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}
