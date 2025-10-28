package roles

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/role"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetRoleDTO, error) {
	rows, err := r.c.Role.Query().Order(ent.Asc(role.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetRoleDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetRoleDTO{ID: v.ID, Name: v.Name, Description: v.Description})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetRoleDTO, error) {
	v, err := r.c.Role.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetRoleDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateRoleDTO) (*GetRoleDTO, error) {
	q := r.c.Role.Create().
		SetName(dto.Name).
		SetDescription(dto.Description)

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetRoleDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateRoleDTO) (*GetRoleDTO, error) {
	q := r.c.Role.UpdateOneID(dto.ID)

	// Set only provided fields
	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetRoleDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Role.DeleteOneID(id).Exec(ctx)
}
