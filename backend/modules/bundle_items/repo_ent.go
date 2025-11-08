package bundle_items

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/bundle_item"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetBundle_itemDTO, error) {
	rows, err := r.c.Bundle_item.Query().
		WithBundle().
		WithProduct().
		Order(ent.Asc(bundle_item.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetBundle_itemDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetBundle_itemDTO{
			ID:  v.ID,
			Qty: v.Qty,
		}
		if v.Edges.Bundle != nil {
			dto.BundleID = v.Edges.Bundle.ID
		}
		if v.Edges.Product != nil {
			dto.ProductID = v.Edges.Product.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetBundle_itemDTO, error) {
	v, err := r.c.Bundle_item.Query().
		WithBundle().
		WithProduct().
		Where(bundle_item.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetBundle_itemDTO{
		ID:  v.ID,
		Qty: v.Qty,
	}
	if v.Edges.Bundle != nil {
		dto.BundleID = v.Edges.Bundle.ID
	}
	if v.Edges.Product != nil {
		dto.ProductID = v.Edges.Product.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateBundle_itemDTO) (*GetBundle_itemDTO, error) {
	bundle, err := r.c.Bundle.Get(ctx, dto.BundleID)
	if err != nil {
		return nil, err
	}
	product, err := r.c.Product.Get(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}

	row, err := r.c.Bundle_item.
		Create().
		SetID(dto.ID).
		SetQty(dto.Qty).
		SetBundle(bundle).
		SetProduct(product).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetBundle_itemDTO{
		ID:        row.ID,
		Qty:       row.Qty,
		BundleID:  dto.BundleID,
		ProductID: dto.ProductID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateBundle_itemDTO) (*GetBundle_itemDTO, error) {
	q := r.c.Bundle_item.UpdateOneID(dto.ID)

	if dto.Qty != nil {
		q.SetQty(*dto.Qty)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Bundle_item.Query().
		WithBundle().
		WithProduct().
		Where(bundle_item.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetBundle_itemDTO{
		ID:  v.ID,
		Qty: v.Qty,
	}
	if v.Edges.Bundle != nil {
		dtoOut.BundleID = v.Edges.Bundle.ID
	}
	if v.Edges.Product != nil {
		dtoOut.ProductID = v.Edges.Product.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Bundle_item.DeleteOneID(id).Exec(ctx)
}
