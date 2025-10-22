package permissions

import (
	"context"

	"github.com/google/uuid"
	"freshease/backend/ent"
	"freshease/backend/ent/permission"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*Permission, error) {
	rows, err := r.c.Permission.Query().Order(ent.Asc(permission.FieldID)).All(ctx)
	if err != nil { return nil, err }
	out := make([]*Permission, 0, len(rows))
	for _, v := range rows {
		out = append(out, &Permission{ID: v.ID, Email: v.Email, Name: v.Name})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	v, err := r.c.Permission.Get(ctx, id)
	if err != nil { return nil, err }
	return &Permission{ID: v.ID, Email: v.Email, Name: v.Name}, nil
}

func (r *EntRepo) Create(ctx context.Context, u *Permission) error {
	newRow, err := r.c.Permission.
		Create().
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	if err != nil { return err }
	u.ID = newRow.ID
	return nil
}

func (r *EntRepo) Update(ctx context.Context, u *Permission) error {
	_, err := r.c.Permission.
		UpdateOneID(u.ID).
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	return err
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Permission.DeleteOneID(id).Exec(ctx)
}
