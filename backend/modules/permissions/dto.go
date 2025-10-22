package permissions

import "github.com/google/uuid"

type CreatePermissionDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
}

type UpdatePermissionDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        *string   `json:"name" validate:"omitempty,min=2,max=60"`
	Description *string   `json:"description" validate:"omitempty"`
}

type GetPermissionDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
}
