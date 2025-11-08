package reviews

import (
	"time"

	"github.com/google/uuid"
)

type CreateReviewDTO struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	Rating    int        `json:"rating" validate:"required,min=1,max=5"`
	Comment   *string    `json:"comment,omitempty"`
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	ProductID uuid.UUID  `json:"product_id" validate:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type UpdateReviewDTO struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	Rating  *int      `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment *string   `json:"comment,omitempty"`
}

type GetReviewDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Rating    int       `json:"rating" validate:"required"`
	Comment   *string   `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}
