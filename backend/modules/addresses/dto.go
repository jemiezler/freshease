package addresses

import "github.com/google/uuid"

type CreateAddressDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Line1      string    `json:"line1" validate:"required"`
	Line2      *string   `json:"line2" validate:"omitempty"`
	City       string    `json:"city" validate:"required"`
	Province   string    `json:"province" validate:"required"`
	Country    string    `json:"country" validate:"required"`
	PostalCode string    `json:"postal_code" validate:"required"`
	IsDefault  bool      `json:"is_default"`
}

type UpdateAddressDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Line1      *string   `json:"line1" validate:"omitempty"`
	Line2      *string   `json:"line2" validate:"omitempty"`
	City       *string   `json:"city" validate:"omitempty"`
	Province   *string   `json:"province" validate:"omitempty"`
	Country    *string   `json:"country" validate:"omitempty"`
	PostalCode *string   `json:"postal_code" validate:"omitempty"`
	IsDefault  *bool     `json:"is_default" validate:"omitempty"`
}

type GetAddressDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Line1      string    `json:"line1" validate:"required"`
	Line2      string    `json:"line2" validate:"omitempty"`
	City       string    `json:"city" validate:"required"`
	Province   string    `json:"province" validate:"required"`
	Country    string    `json:"country" validate:"required"`
	PostalCode string    `json:"postal_code" validate:"required"`
	IsDefault  bool      `json:"is_default"`
}
