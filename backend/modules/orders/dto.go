package orders

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderDTO struct {
	ID                uuid.UUID  `json:"id" validate:"required"`
	OrderNo           string     `json:"order_no" validate:"required"`
	Status            string     `json:"status" validate:"required"`
	Subtotal          float64    `json:"subtotal" validate:"required,min=0"`
	ShippingFee       float64    `json:"shipping_fee" validate:"required,min=0"`
	Discount          float64    `json:"discount" validate:"required,min=0"`
	Total             float64    `json:"total" validate:"required,min=0"`
	PlacedAt          *time.Time `json:"placed_at,omitempty"`
	UserID            uuid.UUID  `json:"user_id" validate:"required"`
	ShippingAddressID *uuid.UUID `json:"shipping_address_id,omitempty"`
	BillingAddressID  *uuid.UUID `json:"billing_address_id,omitempty"`
}

type UpdateOrderDTO struct {
	ID                uuid.UUID  `json:"id" validate:"required"`
	OrderNo           *string    `json:"order_no,omitempty"`
	Status            *string    `json:"status,omitempty"`
	Subtotal          *float64   `json:"subtotal,omitempty" validate:"omitempty,min=0"`
	ShippingFee       *float64   `json:"shipping_fee,omitempty" validate:"omitempty,min=0"`
	Discount          *float64   `json:"discount,omitempty" validate:"omitempty,min=0"`
	Total             *float64   `json:"total,omitempty" validate:"omitempty,min=0"`
	PlacedAt          *time.Time `json:"placed_at,omitempty"`
	ShippingAddressID *uuid.UUID `json:"shipping_address_id,omitempty"`
	BillingAddressID  *uuid.UUID `json:"billing_address_id,omitempty"`
}

type GetOrderDTO struct {
	ID                uuid.UUID  `json:"id" validate:"required"`
	OrderNo           string     `json:"order_no" validate:"required"`
	Status            string     `json:"status" validate:"required"`
	Subtotal          float64    `json:"subtotal" validate:"required"`
	ShippingFee       float64    `json:"shipping_fee" validate:"required"`
	Discount          float64    `json:"discount" validate:"required"`
	Total             float64    `json:"total" validate:"required"`
	PlacedAt          *time.Time `json:"placed_at,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at" validate:"required"`
	UserID            uuid.UUID  `json:"user_id" validate:"required"`
	ShippingAddressID *uuid.UUID `json:"shipping_address_id,omitempty"`
	BillingAddressID  *uuid.UUID `json:"billing_address_id,omitempty"`
}
