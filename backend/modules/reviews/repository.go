package reviews

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetReviewDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetReviewDTO, error)
	Create(ctx context.Context, u *CreateReviewDTO) (*GetReviewDTO, error)
	Update(ctx context.Context, u *UpdateReviewDTO) (*GetReviewDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
