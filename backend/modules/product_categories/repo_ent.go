package product_categories

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/product_category"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetProductCategoryDTO, error) {
	rows, err := r.c.Product_category.Query().Order(ent.Asc(product_category.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetProductCategoryDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetProductCategoryDTO{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Slug:        v.Slug,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error) {
	v, err := r.c.Product_category.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetProductCategoryDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Slug:        v.Slug,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	q := r.c.Product_category.
		Create().
		SetName(dto.Name).
		SetDescription(dto.Description).
		SetSlug(dto.Slug)

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetProductCategoryDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Slug:        row.Slug,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	q := r.c.Product_category.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
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

	return &GetProductCategoryDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Slug:        row.Slug,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Product_category.DeleteOneID(id).Exec(ctx)
}
