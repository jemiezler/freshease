package notifications

import (
	"time"

	"github.com/google/uuid"
)

type CreateNotificationDTO struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	Title     string     `json:"title" validate:"required"`
	Body      *string    `json:"body,omitempty"`
	Channel   string     `json:"channel" validate:"required"`
	Status    string     `json:"status" validate:"required"`
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type UpdateNotificationDTO struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	Title   *string   `json:"title,omitempty"`
	Body    *string   `json:"body,omitempty"`
	Channel *string   `json:"channel,omitempty"`
	Status  *string   `json:"status,omitempty"`
}

type GetNotificationDTO struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	Title     string     `json:"title" validate:"required"`
	Body      *string    `json:"body,omitempty"`
	Channel   string     `json:"channel" validate:"required"`
	Status    string     `json:"status" validate:"required"`
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	CreatedAt time.Time  `json:"created_at" validate:"required"`
}
