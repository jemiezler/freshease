package addresses

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/address"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetAddressDTO, error) {
	rows, err := r.c.Address.Query().Order(ent.Asc(address.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetAddressDTO, 0, len(rows))
	for _, v := range rows {
		line2 := ""
		if v.Line2 != nil {
			line2 = *v.Line2
		}
		out = append(out, &GetAddressDTO{
			ID:        v.ID,
			Line1:     v.Line1,
			Line2:     line2,
			City:      v.City,
			Province:  v.Province,
			Country:   v.Country,
			Zip:       v.Zip,
			IsDefault: v.IsDefault,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetAddressDTO, error) {
	v, err := r.c.Address.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	line2 := ""
	if v.Line2 != nil {
		line2 = *v.Line2
	}
	return &GetAddressDTO{
		ID:        v.ID,
		Line1:     v.Line1,
		Line2:     line2,
		City:      v.City,
		Province:  v.Province,
		Country:   v.Country,
		Zip:       v.Zip,
		IsDefault: v.IsDefault,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateAddressDTO) (*GetAddressDTO, error) {
	q := r.c.Address.
		Create().
		SetID(dto.ID).
		SetLine1(dto.Line1).
		SetCity(dto.City).
		SetProvince(dto.Province).
		SetCountry(dto.Country).
		SetZip(dto.Zip).
		SetIsDefault(dto.IsDefault)

	if dto.Line2 != nil {
		q.SetLine2(*dto.Line2)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	line2 := ""
	if row.Line2 != nil {
		line2 = *row.Line2
	}

	return &GetAddressDTO{
		ID:        row.ID,
		Line1:     row.Line1,
		Line2:     line2,
		City:      row.City,
		Province:  row.Province,
		Country:   row.Country,
		Zip:       row.Zip,
		IsDefault: row.IsDefault,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateAddressDTO) (*GetAddressDTO, error) {
	q := r.c.Address.UpdateOneID(dto.ID)

	if dto.Line1 != nil {
		q.SetLine1(*dto.Line1)
	}
	if dto.Line2 != nil {
		q.SetLine2(*dto.Line2)
	}
	if dto.City != nil {
		q.SetCity(*dto.City)
	}
	if dto.Province != nil {
		q.SetProvince(*dto.Province)
	}
	if dto.Country != nil {
		q.SetCountry(*dto.Country)
	}
	if dto.Zip != nil {
		q.SetZip(*dto.Zip)
	}
	if dto.IsDefault != nil {
		q.SetIsDefault(*dto.IsDefault)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	line2 := ""
	if row.Line2 != nil {
		line2 = *row.Line2
	}

	return &GetAddressDTO{
		ID:        row.ID,
		Line1:     row.Line1,
		Line2:     line2,
		City:      row.City,
		Province:  row.Province,
		Country:   row.Country,
		Zip:       row.Zip,
		IsDefault: row.IsDefault,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Address.DeleteOneID(id).Exec(ctx)
}
