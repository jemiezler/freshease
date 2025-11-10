package password

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string      `json:"accessToken"`
	User        UserResponse `json:"user"`
	Message     string      `json:"message"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role,omitempty"`
}

type InitAdminRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type InitAdminResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

