package bundle_items

import "github.com/google/uuid"

type CreateBundle_itemDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Qty      int       `json:"qty" validate:"required,min=1"`
	BundleID uuid.UUID `json:"bundle_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

type UpdateBundle_itemDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Qty      *int      `json:"qty" validate:"omitempty,min=1"`
}

type GetBundle_itemDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Qty      int       `json:"qty" validate:"required"`
	BundleID uuid.UUID `json:"bundle_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}
