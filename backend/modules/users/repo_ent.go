package users

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/user"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetUserDTO, error) {
	rows, err := r.c.User.
		Query().
		Order(ent.Asc(user.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetUserDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetUserDTO{
			ID:     v.ID,
			Email:  v.Email,
			Name:   v.Name,
			Phone:  v.Phone,
			Bio:    ptrIfNotNil(v.Bio),
			Avatar: ptrIfNotNil(v.Avatar),
			Cover:  ptrIfNotNil(v.Cover),
			Status: v.Status,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetUserDTO, error) {
	v, err := r.c.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetUserDTO{
		ID:     v.ID,
		Email:  v.Email,
		Name:   v.Name,
		Phone:  v.Phone,
		Bio:    ptrIfNotNil(v.Bio),
		Avatar: ptrIfNotNil(v.Avatar),
		Cover:  ptrIfNotNil(v.Cover),
		Status: v.Status,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateUserDTO) (*GetUserDTO, error) {
	q := r.c.User.Create().
		SetEmail(dto.Email).
		SetPassword(dto.Password). // TODO: hash before storing!
		SetName(dto.Name).
		SetPhone(dto.Phone)

	// Optional/nillable fields
	if dto.Bio != nil {
		q.SetNillableBio(dto.Bio)
	}
	if dto.Avatar != nil {
		q.SetNillableAvatar(dto.Avatar)
	}
	if dto.Cover != nil {
		q.SetNillableCover(dto.Cover)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetUserDTO{
		ID:     row.ID,
		Email:  row.Email,
		Name:   row.Name,
		Phone:  row.Phone,
		Bio:    ptrIfNotNil(row.Bio),
		Avatar: ptrIfNotNil(row.Avatar),
		Cover:  ptrIfNotNil(row.Cover),
		Status: row.Status,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateUserDTO) (*GetUserDTO, error) {
	q := r.c.User.UpdateOneID(dto.ID)

	// Set only provided fields
	if dto.Email != nil {
		q.SetEmail(*dto.Email)
	}
	if dto.Password != nil {
		// TODO: hash before storing
		q.SetPassword(*dto.Password)
	}
	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Phone != nil {
		q.SetPhone(*dto.Phone)
	}
	if dto.Bio != nil {
		q.SetNillableBio(dto.Bio)
	}
	if dto.Avatar != nil {
		q.SetNillableAvatar(dto.Avatar)
	}
	if dto.Cover != nil {
		q.SetNillableCover(dto.Cover)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetUserDTO{
		ID:     row.ID,
		Email:  row.Email,
		Name:   row.Name,
		Phone:  row.Phone,
		Bio:    ptrIfNotNil(row.Bio),
		Avatar: ptrIfNotNil(row.Avatar),
		Cover:  ptrIfNotNil(row.Cover),
		Status: row.Status,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.User.DeleteOneID(id).Exec(ctx)
}

// helper: turn sql-nullable string into *string for DTO
func ptrIfNotNil(s *string) *string { return s }
