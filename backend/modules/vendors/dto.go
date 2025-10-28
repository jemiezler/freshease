package vendors

import (
	"time"

	"github.com/google/uuid"
)

type CreateVendorDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"omitempty,min=2,max=60"`
	Email       *string    `json:"email" validate:"omitempty,email"`
	Phone       *string    `json:"phone" validate:"omitempty"`
	Address     *string    `json:"address" validate:"omitempty"`
	City        *string    `json:"city" validate:"omitempty"`
	State       *string    `json:"state" validate:"omitempty"`
	Country     *string    `json:"country" validate:"omitempty"`
	PostalCode  *string    `json:"postal_code" validate:"omitempty"`
	Website     *string    `json:"website" validate:"omitempty,url"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	Description *string    `json:"description" validate:"omitempty"`
	IsActive    string     `json:"is_active" validate:"required"`
	CreatedAt   *time.Time `json:"created_at" validate:"omitempty"`
	UpdatedAt   *time.Time `json:"updated_at" validate:"omitempty"`
	DeletedAt   *time.Time `json:"deleted_at" validate:"omitempty"`
}

type UpdateVendorDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"omitempty,min=2,max=60"`
	Email       *string    `json:"email" validate:"omitempty,email"`
	Phone       *string    `json:"phone" validate:"omitempty"`
	Address     *string    `json:"address" validate:"omitempty"`
	City        *string    `json:"city" validate:"omitempty"`
	State       *string    `json:"state" validate:"omitempty"`
	Country     *string    `json:"country" validate:"omitempty"`
	PostalCode  *string    `json:"postal_code" validate:"omitempty"`
	Website     *string    `json:"website" validate:"omitempty,url"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	Description *string    `json:"description" validate:"omitempty"`
	IsActive    *string    `json:"is_active" validate:"omitempty,required"`
	DeletedAt   *time.Time `json:"deleted_at" validate:"omitempty"`
}

type GetVendorDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"omitempty,min=2,max=60"`
	Email       *string    `json:"email" validate:"omitempty,email"`
	Phone       *string    `json:"phone" validate:"omitempty"`
	Address     *string    `json:"address" validate:"omitempty"`
	City        *string    `json:"city" validate:"omitempty"`
	State       *string    `json:"state" validate:"omitempty"`
	Country     *string    `json:"country" validate:"omitempty"`
	PostalCode  *string    `json:"postal_code" validate:"omitempty"`
	Website     *string    `json:"website" validate:"omitempty,url"`
	LogoURL     *string    `json:"logo_url" validate:"omitempty,url"`
	Description *string    `json:"description" validate:"omitempty"`
	IsActive    string     `json:"is_active" validate:"required"`
	CreatedAt   time.Time  `json:"created_at" validate:"required"`
	UpdatedAt   time.Time  `json:"updated_at" validate:"required"`
	DeletedAt   *time.Time `json:"deleted_at" validate:"omitempty"`
}
