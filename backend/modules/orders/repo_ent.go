package orders

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/order"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetOrderDTO, error) {
	rows, err := r.c.Order.Query().
		WithUser().
		WithShippingAddress().
		WithBillingAddress().
		Order(ent.Asc(order.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetOrderDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetOrderDTO{
			ID:          v.ID,
			OrderNo:     v.OrderNo,
			Status:      v.Status,
			Subtotal:    v.Subtotal,
			ShippingFee: v.ShippingFee,
			Discount:    v.Discount,
			Total:       v.Total,
			PlacedAt:    v.PlacedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
			dto.UserID = v.Edges.User[0].ID
		}
		if len(v.Edges.ShippingAddress) > 0 && v.Edges.ShippingAddress[0] != nil {
			dto.ShippingAddressID = &v.Edges.ShippingAddress[0].ID
		}
		if len(v.Edges.BillingAddress) > 0 && v.Edges.BillingAddress[0] != nil {
			dto.BillingAddressID = &v.Edges.BillingAddress[0].ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetOrderDTO, error) {
	v, err := r.c.Order.Query().
		WithUser().
		WithShippingAddress().
		WithBillingAddress().
		Where(order.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetOrderDTO{
		ID:          v.ID,
		OrderNo:     v.OrderNo,
		Status:      v.Status,
		Subtotal:    v.Subtotal,
		ShippingFee: v.ShippingFee,
		Discount:    v.Discount,
		Total:       v.Total,
		PlacedAt:    v.PlacedAt,
		UpdatedAt:   v.UpdatedAt,
	}
	if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
		dto.UserID = v.Edges.User[0].ID
	}
	if len(v.Edges.ShippingAddress) > 0 && v.Edges.ShippingAddress[0] != nil {
		dto.ShippingAddressID = &v.Edges.ShippingAddress[0].ID
	}
	if len(v.Edges.BillingAddress) > 0 && v.Edges.BillingAddress[0] != nil {
		dto.BillingAddressID = &v.Edges.BillingAddress[0].ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateOrderDTO) (*GetOrderDTO, error) {
	user, err := r.c.User.Get(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}

	q := r.c.Order.
		Create().
		SetID(dto.ID).
		SetOrderNo(dto.OrderNo).
		SetStatus(dto.Status).
		SetSubtotal(dto.Subtotal).
		SetShippingFee(dto.ShippingFee).
		SetDiscount(dto.Discount).
		SetTotal(dto.Total).
		AddUser(user)

	if dto.PlacedAt != nil {
		q.SetPlacedAt(*dto.PlacedAt)
	}
	if dto.ShippingAddressID != nil {
		shippingAddr, err := r.c.Address.Get(ctx, *dto.ShippingAddressID)
		if err != nil {
			return nil, err
		}
		q.AddShippingAddress(shippingAddr)
	}
	if dto.BillingAddressID != nil {
		billingAddr, err := r.c.Address.Get(ctx, *dto.BillingAddressID)
		if err != nil {
			return nil, err
		}
		q.AddBillingAddress(billingAddr)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetOrderDTO{
		ID:                row.ID,
		OrderNo:           row.OrderNo,
		Status:            row.Status,
		Subtotal:          row.Subtotal,
		ShippingFee:       row.ShippingFee,
		Discount:          row.Discount,
		Total:             row.Total,
		PlacedAt:          row.PlacedAt,
		UpdatedAt:         row.UpdatedAt,
		UserID:            dto.UserID,
		ShippingAddressID: dto.ShippingAddressID,
		BillingAddressID:  dto.BillingAddressID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateOrderDTO) (*GetOrderDTO, error) {
	q := r.c.Order.UpdateOneID(dto.ID)

	if dto.OrderNo != nil {
		q.SetOrderNo(*dto.OrderNo)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Subtotal != nil {
		q.SetSubtotal(*dto.Subtotal)
	}
	if dto.ShippingFee != nil {
		q.SetShippingFee(*dto.ShippingFee)
	}
	if dto.Discount != nil {
		q.SetDiscount(*dto.Discount)
	}
	if dto.Total != nil {
		q.SetTotal(*dto.Total)
	}
	if dto.PlacedAt != nil {
		q.SetPlacedAt(*dto.PlacedAt)
	}
	if dto.ShippingAddressID != nil {
		shippingAddr, err := r.c.Address.Get(ctx, *dto.ShippingAddressID)
		if err != nil {
			return nil, err
		}
		q.AddShippingAddress(shippingAddr)
	}
	if dto.BillingAddressID != nil {
		billingAddr, err := r.c.Address.Get(ctx, *dto.BillingAddressID)
		if err != nil {
			return nil, err
		}
		q.AddBillingAddress(billingAddr)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Order.Query().
		WithUser().
		WithShippingAddress().
		WithBillingAddress().
		Where(order.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetOrderDTO{
		ID:          v.ID,
		OrderNo:     v.OrderNo,
		Status:      v.Status,
		Subtotal:    v.Subtotal,
		ShippingFee: v.ShippingFee,
		Discount:    v.Discount,
		Total:       v.Total,
		PlacedAt:    v.PlacedAt,
		UpdatedAt:   v.UpdatedAt,
	}
	if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
		dtoOut.UserID = v.Edges.User[0].ID
	}
	if len(v.Edges.ShippingAddress) > 0 && v.Edges.ShippingAddress[0] != nil {
		dtoOut.ShippingAddressID = &v.Edges.ShippingAddress[0].ID
	}
	if len(v.Edges.BillingAddress) > 0 && v.Edges.BillingAddress[0] != nil {
		dtoOut.BillingAddressID = &v.Edges.BillingAddress[0].ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Order.DeleteOneID(id).Exec(ctx)
}
