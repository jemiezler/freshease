package product_categories

import "github.com/google/uuid"

type CreateProductCategoryDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	ProductID  uuid.UUID `json:"product_id" validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

type UpdateProductCategoryDTO struct {
	ID         uuid.UUID  `json:"id" validate:"required"`
	ProductID  *uuid.UUID `json:"product_id,omitempty"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
}

type GetProductCategoryDTO struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	ProductID  uuid.UUID `json:"product_id" validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}
