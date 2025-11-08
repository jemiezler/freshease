package vendors

import "github.com/google/uuid"

type CreateVendorDTO struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	Name    *string   `json:"name,omitempty"`
	Contact *string   `json:"contact,omitempty"`
}

type UpdateVendorDTO struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	Name    *string   `json:"name,omitempty"`
	Contact *string   `json:"contact,omitempty"`
}

type GetVendorDTO struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	Name    *string   `json:"name,omitempty"`
	Contact *string   `json:"contact,omitempty"`
}
