package payments

import (
	"time"

	"github.com/google/uuid"
)

type CreatePaymentDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Provider    string     `json:"provider" validate:"required"`
	ProviderRef *string    `json:"provider_ref,omitempty"`
	Status      string     `json:"status" validate:"required"`
	Amount      float64    `json:"amount" validate:"required,min=0"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	OrderID     uuid.UUID  `json:"order_id" validate:"required"`
}

type UpdatePaymentDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Provider    *string    `json:"provider,omitempty"`
	ProviderRef *string    `json:"provider_ref,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Amount      *float64   `json:"amount,omitempty" validate:"omitempty,min=0"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

type GetPaymentDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Provider    string     `json:"provider" validate:"required"`
	ProviderRef *string    `json:"provider_ref,omitempty"`
	Status      string     `json:"status" validate:"required"`
	Amount      float64    `json:"amount" validate:"required"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	OrderID     uuid.UUID  `json:"order_id" validate:"required"`
}
