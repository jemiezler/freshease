package permissions

type CreatePermissionDTO struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name"  validate:"required,min=2,max=60"`
}

type UpdatePermissionDTO struct {
	Email *string `json:"email" validate:"omitempty,email"`
	Name  *string `json:"name"  validate:"omitempty,min=2,max=60"`
}
