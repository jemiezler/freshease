package notifications

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/notification"
	"freshease/backend/internal/common/errs"
	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetNotificationDTO, error) {
	rows, err := r.c.Notification.Query().
		WithUser().
		Order(ent.Asc(notification.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetNotificationDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetNotificationDTO{
			ID:        v.ID,
			Title:     v.Title,
			Body:      v.Body,
			Channel:   v.Channel,
			Status:    v.Status,
			CreatedAt: v.CreatedAt,
		}
		if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
			dto.UserID = v.Edges.User[0].ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetNotificationDTO, error) {
	v, err := r.c.Notification.Query().
		WithUser().
		Where(notification.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetNotificationDTO{
		ID:        v.ID,
		Title:     v.Title,
		Body:      v.Body,
		Channel:   v.Channel,
		Status:    v.Status,
		CreatedAt: v.CreatedAt,
	}
	if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
		dto.UserID = v.Edges.User[0].ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateNotificationDTO) (*GetNotificationDTO, error) {
	user, err := r.c.User.Get(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}

	q := r.c.Notification.
		Create().
		SetID(dto.ID).
		SetTitle(dto.Title).
		SetChannel(dto.Channel).
		SetStatus(dto.Status).
		AddUser(user)

	if dto.Body != nil {
		q.SetBody(*dto.Body)
	}
	if dto.CreatedAt != nil {
		q.SetCreatedAt(*dto.CreatedAt)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetNotificationDTO{
		ID:        row.ID,
		Title:     row.Title,
		Body:      row.Body,
		Channel:   row.Channel,
		Status:    row.Status,
		UserID:    dto.UserID,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateNotificationDTO) (*GetNotificationDTO, error) {
	q := r.c.Notification.UpdateOneID(dto.ID)

	if dto.Title != nil {
		q.SetTitle(*dto.Title)
	}
	if dto.Body != nil {
		q.SetBody(*dto.Body)
	}
	if dto.Channel != nil {
		q.SetChannel(*dto.Channel)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with user edge
	v, err := r.c.Notification.Query().
		WithUser().
		Where(notification.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetNotificationDTO{
		ID:        v.ID,
		Title:     v.Title,
		Body:      v.Body,
		Channel:   v.Channel,
		Status:    v.Status,
		CreatedAt: v.CreatedAt,
	}
	if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
		dtoOut.UserID = v.Edges.User[0].ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Notification.DeleteOneID(id).Exec(ctx)
}
