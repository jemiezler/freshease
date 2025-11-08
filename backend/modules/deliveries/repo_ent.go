package deliveries

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/delivery"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetDeliveryDTO, error) {
	rows, err := r.c.Delivery.Query().WithOrder().Order(ent.Asc(delivery.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetDeliveryDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetDeliveryDTO{
			ID:         v.ID,
			Provider:   v.Provider,
			TrackingNo: v.TrackingNo,
			Status:     v.Status,
			Eta:        v.Eta,
			DeliveredAt: v.DeliveredAt,
		}
		if len(v.Edges.Order) > 0 && v.Edges.Order[0] != nil {
			dto.OrderID = v.Edges.Order[0].ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetDeliveryDTO, error) {
	v, err := r.c.Delivery.Query().WithOrder().Where(delivery.ID(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetDeliveryDTO{
		ID:         v.ID,
		Provider:   v.Provider,
		TrackingNo: v.TrackingNo,
		Status:     v.Status,
		Eta:        v.Eta,
		DeliveredAt: v.DeliveredAt,
	}
	if len(v.Edges.Order) > 0 && v.Edges.Order[0] != nil {
		dto.OrderID = v.Edges.Order[0].ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateDeliveryDTO) (*GetDeliveryDTO, error) {
	order, err := r.c.Order.Get(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	q := r.c.Delivery.
		Create().
		SetID(dto.ID).
		SetProvider(dto.Provider).
		SetStatus(dto.Status).
		AddOrder(order)

	if dto.TrackingNo != nil {
		q.SetTrackingNo(*dto.TrackingNo)
	}
	if dto.Eta != nil {
		q.SetEta(*dto.Eta)
	}
	if dto.DeliveredAt != nil {
		q.SetDeliveredAt(*dto.DeliveredAt)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetDeliveryDTO{
		ID:          row.ID,
		Provider:    row.Provider,
		TrackingNo:  row.TrackingNo,
		Status:      row.Status,
		Eta:         row.Eta,
		DeliveredAt: row.DeliveredAt,
		OrderID:     dto.OrderID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateDeliveryDTO) (*GetDeliveryDTO, error) {
	q := r.c.Delivery.UpdateOneID(dto.ID)

	if dto.Provider != nil {
		q.SetProvider(*dto.Provider)
	}
	if dto.TrackingNo != nil {
		q.SetTrackingNo(*dto.TrackingNo)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Eta != nil {
		q.SetEta(*dto.Eta)
	}
	if dto.DeliveredAt != nil {
		q.SetDeliveredAt(*dto.DeliveredAt)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetDeliveryDTO{
		ID:          row.ID,
		Provider:    row.Provider,
		TrackingNo:  row.TrackingNo,
		Status:      row.Status,
		Eta:         row.Eta,
		DeliveredAt: row.DeliveredAt,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Delivery.DeleteOneID(id).Exec(ctx)
}

