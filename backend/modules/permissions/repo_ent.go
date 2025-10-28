package permissions

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/permission"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetPermissionDTO, error) {
	rows, err := r.c.Permission.Query().Order(ent.Asc(permission.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetPermissionDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetPermissionDTO{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetPermissionDTO, error) {
	v, err := r.c.Permission.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetPermissionDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreatePermissionDTO) (*GetPermissionDTO, error) {
	q := r.c.Permission.
		Create().
		SetName(dto.Name).
		SetDescription(dto.Description)

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetPermissionDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdatePermissionDTO) (*GetPermissionDTO, error) {
	q := r.c.Permission.UpdateOneID(dto.ID)

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

	return &GetPermissionDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Permission.DeleteOneID(id).Exec(ctx)
}
