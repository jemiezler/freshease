package carts

import (
	"time"

	"github.com/google/uuid"
)

type CreateCartDTO struct {
	Status *string     `json:"status,omitempty" validate:"omitempty"`
	Total  *float64    `json:"total,omitempty" validate:"omitempty"`
	UserID *uuid.UUID  `json:"user_id,omitempty" validate:"omitempty,uuid"`
}

type UpdateCartDTO struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	Status *string   `json:"status,omitempty" validate:"omitempty"`
	Total  *float64  `json:"total,omitempty" validate:"omitempty"`
}

type CartItemDTO struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductImage *string   `json:"product_image,omitempty"`
	ProductPrice float64   `json:"product_price"`
	Quantity    int       `json:"quantity"`
	LineTotal   float64   `json:"line_total"`
}

type GetCartDTO struct {
	ID            uuid.UUID     `json:"id" validate:"required,uuid"`
	Status        string        `json:"status" validate:"required"`
	Subtotal      float64       `json:"subtotal" validate:"required"`
	Discount      float64       `json:"discount" validate:"required"`
	Total         float64       `json:"total" validate:"required"`
	Shipping      float64       `json:"shipping"`
	Tax           float64       `json:"tax"`
	Items         []CartItemDTO  `json:"items"`
	PromoCode     *string       `json:"promo_code,omitempty"`
	PromoDiscount float64       `json:"promo_discount"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" validate:"required"`
}

// Request DTOs for cart operations
type AddToCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateCartItemRequest struct {
	CartItemID string `json:"cart_item_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
}

type ApplyPromoRequest struct {
	PromoCode string `json:"promo_code" validate:"required"`
}
