package inventories

import (
	"context"
	"time"

	"freshease/backend/ent"
	"freshease/backend/ent/inventory"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetInventoryDTO, error) {
	rows, err := r.c.Inventory.Query().Order(ent.Asc(inventory.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetInventoryDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetInventoryDTO{
			ID:            v.ID,
			Quantity:      v.Quantity,
			ReorderLevel: v.ReorderLevel,
			UpdatedAt:     v.UpdatedAt,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetInventoryDTO, error) {
	v, err := r.c.Inventory.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetInventoryDTO{
		ID:            v.ID,
		Quantity:      v.Quantity,
		ReorderLevel: v.ReorderLevel,
		UpdatedAt:     v.UpdatedAt,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateInventoryDTO) (*GetInventoryDTO, error) {
	q := r.c.Inventory.
		Create().
		SetQuantity(dto.Quantity).
		SetReorderLevel(dto.ReorderLevel)

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetInventoryDTO{
		ID:            row.ID,
		Quantity:      row.Quantity,
		ReorderLevel: row.ReorderLevel,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateInventoryDTO) (*GetInventoryDTO, error) {
	q := r.c.Inventory.UpdateOneID(dto.ID)

	if dto.Quantity != nil {
		q.SetQuantity(*dto.Quantity)
	}
	if dto.ReorderLevel != nil {
		q.SetReorderLevel(*dto.ReorderLevel)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure UpdatedAt is populated; if ent didn't set it in-memory, fetch current time.
	updatedAt := row.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	return &GetInventoryDTO{
		ID:            row.ID,
		Quantity:      row.Quantity,
		ReorderLevel: row.ReorderLevel,
		UpdatedAt:     updatedAt,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Inventory.DeleteOneID(id).Exec(ctx)
}
