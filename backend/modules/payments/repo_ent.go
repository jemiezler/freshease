package payments

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/payment"
	"freshease/backend/internal/common/errs"
	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetPaymentDTO, error) {
	rows, err := r.c.Payment.Query().
		WithOrder().
		Order(ent.Asc(payment.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetPaymentDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetPaymentDTO{
			ID:          v.ID,
			Provider:    v.Provider,
			ProviderRef: v.ProviderRef,
			Status:      v.Status,
			Amount:      v.Amount,
			PaidAt:      v.PaidAt,
		}
		if len(v.Edges.Order) > 0 && v.Edges.Order[0] != nil {
			dto.OrderID = v.Edges.Order[0].ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetPaymentDTO, error) {
	v, err := r.c.Payment.Query().
		WithOrder().
		Where(payment.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetPaymentDTO{
		ID:          v.ID,
		Provider:    v.Provider,
		ProviderRef: v.ProviderRef,
		Status:      v.Status,
		Amount:      v.Amount,
		PaidAt:      v.PaidAt,
	}
	if len(v.Edges.Order) > 0 && v.Edges.Order[0] != nil {
		dto.OrderID = v.Edges.Order[0].ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreatePaymentDTO) (*GetPaymentDTO, error) {
	order, err := r.c.Order.Get(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	q := r.c.Payment.
		Create().
		SetID(dto.ID).
		SetProvider(dto.Provider).
		SetStatus(dto.Status).
		SetAmount(dto.Amount).
		AddOrder(order)

	if dto.ProviderRef != nil {
		q.SetProviderRef(*dto.ProviderRef)
	}
	if dto.PaidAt != nil {
		q.SetPaidAt(*dto.PaidAt)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetPaymentDTO{
		ID:          row.ID,
		Provider:    row.Provider,
		ProviderRef: row.ProviderRef,
		Status:      row.Status,
		Amount:      row.Amount,
		PaidAt:      row.PaidAt,
		OrderID:     dto.OrderID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdatePaymentDTO) (*GetPaymentDTO, error) {
	q := r.c.Payment.UpdateOneID(dto.ID)

	if dto.Provider != nil {
		q.SetProvider(*dto.Provider)
	}
	if dto.ProviderRef != nil {
		q.SetProviderRef(*dto.ProviderRef)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Amount != nil {
		q.SetAmount(*dto.Amount)
	}
	if dto.PaidAt != nil {
		q.SetPaidAt(*dto.PaidAt)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with order edge
	v, err := r.c.Payment.Query().
		WithOrder().
		Where(payment.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetPaymentDTO{
		ID:          v.ID,
		Provider:    v.Provider,
		ProviderRef: v.ProviderRef,
		Status:      v.Status,
		Amount:      v.Amount,
		PaidAt:      v.PaidAt,
	}
	if len(v.Edges.Order) > 0 && v.Edges.Order[0] != nil {
		dtoOut.OrderID = v.Edges.Order[0].ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Payment.DeleteOneID(id).Exec(ctx)
}
