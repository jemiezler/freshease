package roles

type CreateRoleDTO struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name"  validate:"required,min=2,max=60"`
}

type UpdateRoleDTO struct {
	Email *string `json:"email" validate:"omitempty,email"`
	Name  *string `json:"name"  validate:"omitempty,min=2,max=60"`
}
