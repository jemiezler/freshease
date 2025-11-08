package cart_items

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/cart_item"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetCart_itemDTO, error) {
	rows, err := r.c.Cart_item.Query().
		WithCart().
		WithProduct().
		Order(ent.Asc(cart_item.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetCart_itemDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetCart_itemDTO{
			ID:        v.ID,
			Qty:       v.Qty,
			UnitPrice: v.UnitPrice,
			LineTotal: v.LineTotal,
		}
		if v.Edges.Cart != nil {
			dto.CartID = v.Edges.Cart.ID
		}
		if v.Edges.Product != nil {
			dto.ProductID = v.Edges.Product.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetCart_itemDTO, error) {
	v, err := r.c.Cart_item.Query().
		WithCart().
		WithProduct().
		Where(cart_item.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetCart_itemDTO{
		ID:        v.ID,
		Qty:       v.Qty,
		UnitPrice: v.UnitPrice,
		LineTotal: v.LineTotal,
	}
	if v.Edges.Cart != nil {
		dto.CartID = v.Edges.Cart.ID
	}
	if v.Edges.Product != nil {
		dto.ProductID = v.Edges.Product.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateCart_itemDTO) (*GetCart_itemDTO, error) {
	cart, err := r.c.Cart.Get(ctx, dto.CartID)
	if err != nil {
		return nil, err
	}
	product, err := r.c.Product.Get(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}

	row, err := r.c.Cart_item.
		Create().
		SetID(dto.ID).
		SetQty(dto.Qty).
		SetUnitPrice(dto.UnitPrice).
		SetLineTotal(dto.LineTotal).
		SetCart(cart).
		SetProduct(product).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// reload with edges populated
	v, err := r.c.Cart_item.Query().
		WithCart().
		WithProduct().
		Where(cart_item.IDEQ(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetCart_itemDTO{
		ID:        v.ID,
		Qty:       v.Qty,
		UnitPrice: v.UnitPrice,
		LineTotal: v.LineTotal,
	}
	if v.Edges.Cart != nil {
		dtoOut.CartID = v.Edges.Cart.ID
	}
	if v.Edges.Product != nil {
		dtoOut.ProductID = v.Edges.Product.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateCart_itemDTO) (*GetCart_itemDTO, error) {
	q := r.c.Cart_item.UpdateOneID(dto.ID)

	if dto.Qty != nil {
		q.SetQty(*dto.Qty)
	}
	if dto.UnitPrice != nil {
		q.SetUnitPrice(*dto.UnitPrice)
	}
	if dto.LineTotal != nil {
		q.SetLineTotal(*dto.LineTotal)
	}
	if dto.CartID != nil {
		cart, err := r.c.Cart.Get(ctx, *dto.CartID)
		if err != nil {
			return nil, err
		}
		q.SetCart(cart)
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

	// reload with edges populated
	v, err := r.c.Cart_item.Query().
		WithCart().
		WithProduct().
		Where(cart_item.IDEQ(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetCart_itemDTO{
		ID:        v.ID,
		Qty:       v.Qty,
		UnitPrice: v.UnitPrice,
		LineTotal: v.LineTotal,
	}
	if v.Edges.Cart != nil {
		dtoOut.CartID = v.Edges.Cart.ID
	}
	if v.Edges.Product != nil {
		dtoOut.ProductID = v.Edges.Product.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Cart_item.DeleteOneID(id).Exec(ctx)
}
