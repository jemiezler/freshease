package roles

import (
	"context"

	"github.com/google/uuid"
	"freshease/backend/ent"
	"freshease/backend/ent/role"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*Role, error) {
	rows, err := r.c.Role.Query().Order(ent.Asc(role.FieldID)).All(ctx)
	if err != nil { return nil, err }
	out := make([]*Role, 0, len(rows))
	for _, v := range rows {
		out = append(out, &Role{ID: v.ID, Email: v.Email, Name: v.Name})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	v, err := r.c.Role.Get(ctx, id)
	if err != nil { return nil, err }
	return &Role{ID: v.ID, Email: v.Email, Name: v.Name}, nil
}

func (r *EntRepo) Create(ctx context.Context, u *Role) error {
	newRow, err := r.c.Role.
		Create().
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	if err != nil { return err }
	u.ID = newRow.ID
	return nil
}

func (r *EntRepo) Update(ctx context.Context, u *Role) error {
	_, err := r.c.Role.
		UpdateOneID(u.ID).
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	return err
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Role.DeleteOneID(id).Exec(ctx)
}
