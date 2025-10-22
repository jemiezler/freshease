package users

import "github.com/google/uuid"

type CreateUserDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	Name     string    `json:"name"  validate:"required,min=2,max=100"`
	Phone    string    `json:"phone" validate:"required,min=10,max=20"`
	Bio      *string   `json:"bio" validate:"omitempty,min=10,max=500"`
	Avatar   *string   `json:"avatar" validate:"omitempty,min=10,max=200"`
	Cover    *string   `json:"cover" validate:"omitempty,min=10,max=200"`
	Status   *string   `json:"status" validate:"omitempty"`
}

type UpdateUserDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Email    *string   `json:"email" validate:"omitempty,email"`
	Password *string   `json:"password" validate:"omitempty,min=8,max=100"`
	Name     *string   `json:"name"  validate:"omitempty,min=2,max=100"`
	Phone    *string   `json:"phone" validate:"omitempty,min=10,max=20"`
	Bio      *string   `json:"bio" validate:"omitempty,min=10,max=500"`
	Avatar   *string   `json:"avatar" validate:"omitempty,min=10,max=200"`
	Cover    *string   `json:"cover" validate:"omitempty,min=10,max=200"`
	Status   *string   `json:"status" validate:"omitempty"`
}

type GetUserDTO struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Phone  string    `json:"phone"`
	Bio    *string   `json:"bio"`
	Avatar *string   `json:"avatar"`
	Cover  *string   `json:"cover"`
	Status string    `json:"status"`
}
