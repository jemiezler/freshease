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
			ID:          v.ID,
			Name:        v.Name,
			Email:       v.Email,
			Phone:       v.Phone,
			Address:     v.Address,
			City:        v.City,
			State:       v.State,
			Country:     v.Country,
			PostalCode:  v.PostalCode,
			Website:     v.Website,
			LogoURL:     v.LogoURL,
			Description: v.Description,
			IsActive:    v.IsActive,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			DeletedAt:   v.DeletedAt,
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
		ID:          v.ID,
		Name:        v.Name,
		Email:       v.Email,
		Phone:       v.Phone,
		Address:     v.Address,
		City:        v.City,
		State:       v.State,
		Country:     v.Country,
		PostalCode:  v.PostalCode,
		Website:     v.Website,
		LogoURL:     v.LogoURL,
		Description: v.Description,
		IsActive:    v.IsActive,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		DeletedAt:   v.DeletedAt,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateVendorDTO) (*GetVendorDTO, error) {
	q := r.c.Vendor.Create()

	// Set ID if provided
	if dto.ID != uuid.Nil {
		q.SetID(dto.ID)
	}

	// Nillable fields
	q.SetNillableName(dto.Name).
		SetNillableEmail(dto.Email).
		SetNillablePhone(dto.Phone).
		SetNillableAddress(dto.Address).
		SetNillableCity(dto.City).
		SetNillableState(dto.State).
		SetNillableCountry(dto.Country).
		SetNillablePostalCode(dto.PostalCode).
		SetNillableWebsite(dto.Website).
		SetNillableLogoURL(dto.LogoURL).
		SetNillableDescription(dto.Description).
		SetNillableDeletedAt(dto.DeletedAt)

	// Non-nillable required field
	q.SetIsActive(dto.IsActive)

	// Optional override for created/updated timestamps if provided
	if dto.CreatedAt != nil {
		q.SetCreatedAt(*dto.CreatedAt)
	}
	if dto.UpdatedAt != nil {
		q.SetUpdatedAt(*dto.UpdatedAt)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetVendorDTO{
		ID:          row.ID,
		Name:        row.Name,
		Email:       row.Email,
		Phone:       row.Phone,
		Address:     row.Address,
		City:        row.City,
		State:       row.State,
		Country:     row.Country,
		PostalCode:  row.PostalCode,
		Website:     row.Website,
		LogoURL:     row.LogoURL,
		Description: row.Description,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		DeletedAt:   row.DeletedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateVendorDTO) (*GetVendorDTO, error) {
	q := r.c.Vendor.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetNillableName(dto.Name)
	}
	if dto.Email != nil {
		q.SetNillableEmail(dto.Email)
	}
	if dto.Phone != nil {
		q.SetNillablePhone(dto.Phone)
	}
	if dto.Address != nil {
		q.SetNillableAddress(dto.Address)
	}
	if dto.City != nil {
		q.SetNillableCity(dto.City)
	}
	if dto.State != nil {
		q.SetNillableState(dto.State)
	}
	if dto.Country != nil {
		q.SetNillableCountry(dto.Country)
	}
	if dto.PostalCode != nil {
		q.SetNillablePostalCode(dto.PostalCode)
	}
	if dto.Website != nil {
		q.SetNillableWebsite(dto.Website)
	}
	if dto.LogoURL != nil {
		q.SetNillableLogoURL(dto.LogoURL)
	}
	if dto.Description != nil {
		q.SetNillableDescription(dto.Description)
	}
	if dto.IsActive != nil {
		q.SetIsActive(*dto.IsActive)
	}
	if dto.DeletedAt != nil {
		q.SetNillableDeletedAt(dto.DeletedAt)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetVendorDTO{
		ID:          row.ID,
		Name:        row.Name,
		Email:       row.Email,
		Phone:       row.Phone,
		Address:     row.Address,
		City:        row.City,
		State:       row.State,
		Country:     row.Country,
		PostalCode:  row.PostalCode,
		Website:     row.Website,
		LogoURL:     row.LogoURL,
		Description: row.Description,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		DeletedAt:   row.DeletedAt,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Vendor.DeleteOneID(id).Exec(ctx)
}
