package categories

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/category"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetCategoryDTO, error) {
	rows, err := r.c.Category.Query().Order(ent.Asc(category.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetCategoryDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetCategoryDTO{
			ID:   v.ID,
			Name: v.Name,
			Slug: v.Slug,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetCategoryDTO, error) {
	v, err := r.c.Category.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetCategoryDTO{
		ID:   v.ID,
		Name: v.Name,
		Slug: v.Slug,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateCategoryDTO) (*GetCategoryDTO, error) {
	row, err := r.c.Category.
		Create().
		SetID(dto.ID).
		SetName(dto.Name).
		SetSlug(dto.Slug).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCategoryDTO{
		ID:   row.ID,
		Name: row.Name,
		Slug: row.Slug,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateCategoryDTO) (*GetCategoryDTO, error) {
	q := r.c.Category.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Slug != nil {
		q.SetSlug(*dto.Slug)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCategoryDTO{
		ID:   row.ID,
		Name: row.Name,
		Slug: row.Slug,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Category.DeleteOneID(id).Exec(ctx)
}

