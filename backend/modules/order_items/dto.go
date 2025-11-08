package order_items

import "github.com/google/uuid"

type CreateOrder_itemDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Qty       int       `json:"qty" validate:"required,min=1"`
	UnitPrice float64   `json:"unit_price" validate:"required,min=0"`
	LineTotal float64   `json:"line_total" validate:"required,min=0"`
	OrderID   uuid.UUID `json:"order_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

type UpdateOrder_itemDTO struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	Qty       *int       `json:"qty,omitempty" validate:"omitempty,min=1"`
	UnitPrice *float64   `json:"unit_price,omitempty" validate:"omitempty,min=0"`
	LineTotal *float64   `json:"line_total,omitempty" validate:"omitempty,min=0"`
	OrderID   *uuid.UUID `json:"order_id,omitempty"`
	ProductID *uuid.UUID `json:"product_id,omitempty"`
}

type GetOrder_itemDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Qty       int       `json:"qty" validate:"required"`
	UnitPrice float64   `json:"unit_price" validate:"required"`
	LineTotal float64   `json:"line_total" validate:"required"`
	OrderID   uuid.UUID `json:"order_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}
