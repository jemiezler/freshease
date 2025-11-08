package deliveries

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeliveryDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Provider    string     `json:"provider" validate:"required"`
	TrackingNo  *string    `json:"tracking_no,omitempty"`
	Status      string     `json:"status" validate:"required"`
	Eta         *time.Time `json:"eta,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	OrderID     uuid.UUID  `json:"order_id" validate:"required"`
}

type UpdateDeliveryDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Provider    *string    `json:"provider,omitempty"`
	TrackingNo  *string    `json:"tracking_no,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Eta         *time.Time `json:"eta,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}

type GetDeliveryDTO struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Provider    string     `json:"provider" validate:"required"`
	TrackingNo  *string    `json:"tracking_no,omitempty"`
	Status      string     `json:"status" validate:"required"`
	Eta         *time.Time `json:"eta,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	OrderID     uuid.UUID  `json:"order_id" validate:"required"`
}

