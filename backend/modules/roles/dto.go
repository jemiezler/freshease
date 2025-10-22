package roles

import "github.com/google/uuid"

type CreateRoleDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
}

type UpdateRoleDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        *string   `json:"name" validate:"omitempty,min=2,max=60"`
	Description *string   `json:"description" validate:"omitempty"`
}

type GetRoleDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
}
