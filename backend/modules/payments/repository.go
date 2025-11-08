package payments

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetPaymentDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetPaymentDTO, error)
	Create(ctx context.Context, u *CreatePaymentDTO) (*GetPaymentDTO, error)
	Update(ctx context.Context, u *UpdatePaymentDTO) (*GetPaymentDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
