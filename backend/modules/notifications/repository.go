package notifications

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*GetNotificationDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*GetNotificationDTO, error)
	Create(ctx context.Context, u *CreateNotificationDTO) (*GetNotificationDTO, error)
	Update(ctx context.Context, u *UpdateNotificationDTO) (*GetNotificationDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
