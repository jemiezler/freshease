package order_items

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/order_item"
	"freshease/backend/internal/common/errs"
	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetOrder_itemDTO, error) {
	rows, err := r.c.Order_item.Query().
		WithOrder().
		WithProduct().
		Order(ent.Asc(order_item.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetOrder_itemDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetOrder_itemDTO{
			ID:        v.ID,
			Qty:       v.Qty,
			UnitPrice: v.UnitPrice,
			LineTotal: v.LineTotal,
		}
		if v.Edges.Order != nil {
			dto.OrderID = v.Edges.Order.ID
		}
		if v.Edges.Product != nil {
			dto.ProductID = v.Edges.Product.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetOrder_itemDTO, error) {
	v, err := r.c.Order_item.Query().
		WithOrder().
		WithProduct().
		Where(order_item.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetOrder_itemDTO{
		ID:        v.ID,
		Qty:       v.Qty,
		UnitPrice: v.UnitPrice,
		LineTotal: v.LineTotal,
	}
	if v.Edges.Order != nil {
		dto.OrderID = v.Edges.Order.ID
	}
	if v.Edges.Product != nil {
		dto.ProductID = v.Edges.Product.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateOrder_itemDTO) (*GetOrder_itemDTO, error) {
	order, err := r.c.Order.Get(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}
	product, err := r.c.Product.Get(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}

	row, err := r.c.Order_item.
		Create().
		SetID(dto.ID).
		SetQty(dto.Qty).
		SetUnitPrice(dto.UnitPrice).
		SetLineTotal(dto.LineTotal).
		SetOrder(order).
		SetProduct(product).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetOrder_itemDTO{
		ID:        row.ID,
		Qty:       row.Qty,
		UnitPrice: row.UnitPrice,
		LineTotal: row.LineTotal,
		OrderID:   dto.OrderID,
		ProductID: dto.ProductID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateOrder_itemDTO) (*GetOrder_itemDTO, error) {
	q := r.c.Order_item.UpdateOneID(dto.ID)

	if dto.Qty != nil {
		q.SetQty(*dto.Qty)
	}
	if dto.UnitPrice != nil {
		q.SetUnitPrice(*dto.UnitPrice)
	}
	if dto.LineTotal != nil {
		q.SetLineTotal(*dto.LineTotal)
	}
	if dto.OrderID != nil {
		order, err := r.c.Order.Get(ctx, *dto.OrderID)
		if err != nil {
			return nil, err
		}
		q.SetOrder(order)
	}
	if dto.ProductID != nil {
		product, err := r.c.Product.Get(ctx, *dto.ProductID)
		if err != nil {
			return nil, err
		}
		q.SetProduct(product)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Order_item.Query().
		WithOrder().
		WithProduct().
		Where(order_item.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetOrder_itemDTO{
		ID:        v.ID,
		Qty:       v.Qty,
		UnitPrice: v.UnitPrice,
		LineTotal: v.LineTotal,
	}
	if v.Edges.Order != nil {
		dtoOut.OrderID = v.Edges.Order.ID
	}
	if v.Edges.Product != nil {
		dtoOut.ProductID = v.Edges.Product.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Order_item.DeleteOneID(id).Exec(ctx)
}
