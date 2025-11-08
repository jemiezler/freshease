package products

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/product"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetProductDTO, error) {
	rows, err := r.c.Product.Query().Order(ent.Asc(product.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetProductDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetProductDTO{
			ID:          v.ID,
			Name:        v.Name,
			SKU:         v.Sku,
			Price:       v.Price,
			Description: v.Description,
			UnitLabel:   v.UnitLabel,
			IsActive:    v.IsActive,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetProductDTO, error) {
	v, err := r.c.Product.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetProductDTO{
		ID:          v.ID,
		Name:        v.Name,
		SKU:         v.Sku,
		Price:       v.Price,
		Description: v.Description,
		UnitLabel:   v.UnitLabel,
		IsActive:    v.IsActive,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateProductDTO) (*GetProductDTO, error) {
	// Create product first
	q := r.c.Product.
		Create().
		SetID(dto.ID).
		SetName(dto.Name).
		SetSku(dto.SKU).
		SetPrice(dto.Price).
		SetUnitLabel(dto.UnitLabel).
		SetIsActive(dto.IsActive).
		SetCreatedAt(dto.CreatedAt).
		SetUpdatedAt(dto.UpdatedAt)

	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Note: Inventory creation requires product and vendor IDs
	// This should be handled separately or passed in the DTO

	return &GetProductDTO{
		ID:          row.ID,
		Name:        row.Name,
		SKU:         row.Sku,
		Price:       row.Price,
		Description: row.Description,
		UnitLabel:   row.UnitLabel,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateProductDTO) (*GetProductDTO, error) {
	q := r.c.Product.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.SKU != nil {
		q.SetSku(*dto.SKU)
	}
	if dto.Price != nil {
		q.SetPrice(*dto.Price)
	}
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}
	if dto.UnitLabel != nil {
		q.SetUnitLabel(*dto.UnitLabel)
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

	return &GetProductDTO{
		ID:          row.ID,
		Name:        row.Name,
		SKU:         row.Sku,
		Price:       row.Price,
		Description: row.Description,
		UnitLabel:   row.UnitLabel,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Product.DeleteOneID(id).Exec(ctx)
}
