package permissions

import "github.com/google/uuid"

type CreatePermissionDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Code        string    `json:"code" validate:"required,min=2,max=60"`
	Description *string   `json:"description,omitempty"`
}

type UpdatePermissionDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Code        *string   `json:"code" validate:"omitempty,min=2,max=60"`
	Description *string   `json:"description,omitempty"`
}

type GetPermissionDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	Description *string   `json:"description,omitempty"`
}
