package carts

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/cart"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetCartDTO, error) {
	rows, err := r.c.Cart.Query().Order(ent.Asc(cart.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetCartDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetCartDTO{
			ID:        v.ID,
			Status:    v.Status,
			Total:     v.Total,
			Subtotal: v.Subtotal,
		Discount: v.Discount,
			UpdatedAt: v.UpdatedAt,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetCartDTO, error) {
	v, err := r.c.Cart.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetCartDTO{
		ID:        v.ID,
		Status:    v.Status,
		Total:     v.Total,
		Subtotal: v.Subtotal,
		Discount: v.Discount,
		UpdatedAt: v.UpdatedAt,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateCartDTO) (*GetCartDTO, error) {
	q := r.c.Cart.Create()

	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Total != nil {
		q.SetTotal(*dto.Total)
	}
	if dto.UserID != nil {
		user, err := r.c.User.Get(ctx, *dto.UserID)
		if err != nil {
			return nil, err
		}
		q.AddUser(user)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCartDTO{
		ID:        row.ID,
		Status:    row.Status,
		Subtotal:  row.Subtotal,
		Discount:  row.Discount,
		Total:     row.Total,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateCartDTO) (*GetCartDTO, error) {
	q := r.c.Cart.UpdateOneID(dto.ID)

	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Total != nil {
		q.SetTotal(*dto.Total)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCartDTO{
		ID:        row.ID,
		Status:    row.Status,
		Subtotal:  row.Subtotal,
		Discount:  row.Discount,
		Total:     row.Total,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Cart.DeleteOneID(id).Exec(ctx)
}
