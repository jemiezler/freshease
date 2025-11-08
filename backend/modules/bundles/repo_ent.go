package bundles

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/bundle"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetBundleDTO, error) {
	rows, err := r.c.Bundle.Query().Order(ent.Asc(bundle.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetBundleDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetBundleDTO{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			IsActive:    v.IsActive,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetBundleDTO, error) {
	v, err := r.c.Bundle.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetBundleDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Price:       v.Price,
		IsActive:    v.IsActive,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateBundleDTO) (*GetBundleDTO, error) {
	q := r.c.Bundle.
		Create().
		SetID(dto.ID).
		SetName(dto.Name).
		SetPrice(dto.Price).
		SetIsActive(dto.IsActive)
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetBundleDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Price:       row.Price,
		IsActive:    row.IsActive,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateBundleDTO) (*GetBundleDTO, error) {
	q := r.c.Bundle.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}
	if dto.Price != nil {
		q.SetPrice(*dto.Price)
	}
	if dto.IsActive != nil {
		q.SetIsActive(*dto.IsActive)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetBundleDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Price:       row.Price,
		IsActive:    row.IsActive,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Bundle.DeleteOneID(id).Exec(ctx)
}
