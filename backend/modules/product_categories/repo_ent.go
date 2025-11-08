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
	rows, err := r.c.Product_category.Query().
		WithProduct().
		WithCategory().
		Order(ent.Asc(product_category.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetProductCategoryDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetProductCategoryDTO{
			ID: v.ID,
		}
		if v.Edges.Product != nil {
			dto.ProductID = v.Edges.Product.ID
		}
		if v.Edges.Category != nil {
			dto.CategoryID = v.Edges.Category.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error) {
	v, err := r.c.Product_category.Query().
		WithProduct().
		WithCategory().
		Where(product_category.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetProductCategoryDTO{
		ID: v.ID,
	}
	if v.Edges.Product != nil {
		dto.ProductID = v.Edges.Product.ID
	}
	if v.Edges.Category != nil {
		dto.CategoryID = v.Edges.Category.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	product, err := r.c.Product.Get(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}
	category, err := r.c.Category.Get(ctx, dto.CategoryID)
	if err != nil {
		return nil, err
	}

	row, err := r.c.Product_category.
		Create().
		SetID(dto.ID).
		SetProduct(product).
		SetCategory(category).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetProductCategoryDTO{
		ID:         row.ID,
		ProductID:  dto.ProductID,
		CategoryID: dto.CategoryID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	// Product_category is a join table with no fields to update
	// Only edges can be updated
	q := r.c.Product_category.UpdateOneID(dto.ID)
	hasChanges := false

	if dto.ProductID != nil {
		product, err := r.c.Product.Get(ctx, *dto.ProductID)
		if err != nil {
			return nil, err
		}
		q.SetProduct(product)
		hasChanges = true
	}
	if dto.CategoryID != nil {
		category, err := r.c.Category.Get(ctx, *dto.CategoryID)
		if err != nil {
			return nil, err
		}
		q.SetCategory(category)
		hasChanges = true
	}

	if !hasChanges {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Product_category.Query().
		WithProduct().
		WithCategory().
		Where(product_category.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetProductCategoryDTO{
		ID: v.ID,
	}
	if v.Edges.Product != nil {
		dtoOut.ProductID = v.Edges.Product.ID
	}
	if v.Edges.Category != nil {
		dtoOut.CategoryID = v.Edges.Category.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Product_category.DeleteOneID(id).Exec(ctx)
}
