package users

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string     `json:"email"`
	Name  string     `json:"name"`
}

type Repository interface {
	List(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
