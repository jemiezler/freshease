package vendors

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/vendor"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetVendorDTO, error) {
	rows, err := r.c.Vendor.Query().Order(ent.Asc(vendor.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetVendorDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetVendorDTO{
			ID:      v.ID,
			Name:    v.Name,
			Contact: v.Contact,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetVendorDTO, error) {
	v, err := r.c.Vendor.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetVendorDTO{
		ID:      v.ID,
		Name:    v.Name,
		Contact: v.Contact,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateVendorDTO) (*GetVendorDTO, error) {
	q := r.c.Vendor.Create()

	// Set ID if provided
	if dto.ID != uuid.Nil {
		q.SetID(dto.ID)
	}

	// Nillable fields
	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Contact != nil {
		q.SetContact(*dto.Contact)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetVendorDTO{
		ID:      row.ID,
		Name:    row.Name,
		Contact: row.Contact,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateVendorDTO) (*GetVendorDTO, error) {
	q := r.c.Vendor.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Contact != nil {
		q.SetContact(*dto.Contact)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetVendorDTO{
		ID:      row.ID,
		Name:    row.Name,
		Contact: row.Contact,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Vendor.DeleteOneID(id).Exec(ctx)
}
