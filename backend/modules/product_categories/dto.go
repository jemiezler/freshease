package product_categories

import "github.com/google/uuid"

type CreateProductCategoryDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
	Slug        string    `json:"slug" validate:"required"`
}

type UpdateProductCategoryDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        *string   `json:"name" validate:"omitempty,min=2,max=60"`
	Description *string   `json:"description" validate:"omitempty"`
	Slug        *string   `json:"slug" validate:"omitempty"`
}

type GetProductCategoryDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=60"`
	Description string    `json:"description" validate:"required"`
	Slug        string    `json:"slug" validate:"required"`
}
