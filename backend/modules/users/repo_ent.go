package users

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/user"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*User, error) {
	rows, err := r.c.User.Query().Order(ent.Asc(user.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*User, 0, len(rows))
	for _, v := range rows {
		out = append(out, &User{ID: v.ID, Email: v.Email, Name: v.Name})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	v, err := r.c.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{ID: v.ID, Email: v.Email, Name: v.Name}, nil
}

func (r *EntRepo) Create(ctx context.Context, u *User) error {
	newRow, err := r.c.User.
		Create().
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	if err != nil {
		return err
	}
	u.ID = newRow.ID
	return nil
}

func (r *EntRepo) Update(ctx context.Context, u *User) error {
	_, err := r.c.User.
		UpdateOneID(u.ID).
		SetEmail(u.Email).
		SetName(u.Name).
		Save(ctx)
	return err
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.User.DeleteOneID(id).Exec(ctx)
}
