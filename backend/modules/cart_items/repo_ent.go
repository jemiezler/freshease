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
	rows, err := r.c.Cart_item.Query().Order(ent.Asc(cart_item.FieldID)).WithCart().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetCart_itemDTO, 0, len(rows))
	for _, v := range rows {
		var cartID uuid.UUID
		if v.Edges.Cart != nil {
			cartID = v.Edges.Cart.ID
		}
		out = append(out, &GetCart_itemDTO{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Cart:        cartID,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetCart_itemDTO, error) {
	v, err := r.c.Cart_item.Query().Where(cart_item.IDEQ(id)).WithCart().Only(ctx)
	if err != nil {
		return nil, err
	}
	var cartID uuid.UUID
	if v.Edges.Cart != nil {
		cartID = v.Edges.Cart.ID
	}
	return &GetCart_itemDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Cart:        cartID,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateCart_itemDTO) (*GetCart_itemDTO, error) {
	q := r.c.Cart_item.
		Create().
		SetID(dto.ID).
		SetName(dto.Name).
		SetDescription(dto.Description).
		SetCartID(dto.Cart)

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// reload with cart edge populated
	v, err := r.c.Cart_item.Query().Where(cart_item.IDEQ(row.ID)).WithCart().Only(ctx)
	if err != nil {
		return nil, err
	}

	var cartID uuid.UUID
	if v.Edges.Cart != nil {
		cartID = v.Edges.Cart.ID
	}

	return &GetCart_itemDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Cart:        cartID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateCart_itemDTO) (*GetCart_itemDTO, error) {
	q := r.c.Cart_item.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}
	if dto.Cart != nil {
		q.SetCartID(*dto.Cart)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// reload with cart edge populated
	v, err := r.c.Cart_item.Query().Where(cart_item.IDEQ(row.ID)).WithCart().Only(ctx)
	if err != nil {
		return nil, err
	}

	var cartID uuid.UUID
	if v.Edges.Cart != nil {
		cartID = v.Edges.Cart.ID
	}

	return &GetCart_itemDTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Cart:        cartID,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Cart_item.DeleteOneID(id).Exec(ctx)
}
