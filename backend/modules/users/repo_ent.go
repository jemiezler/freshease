package users

import (
	"context"
	"time"

	"freshease/backend/ent"
	"freshease/backend/ent/user"
	"freshease/backend/internal/common/errs"
	"freshease/backend/internal/common/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
			ID:          v.ID,
			Email:       v.Email,
			Name:        v.Name,
			Phone:       v.Phone,
			Bio:         helpers.PtrIfNotNil(v.Bio),
			Avatar:      helpers.PtrIfNotNil(v.Avatar),
			Cover:       helpers.PtrIfNotNil(v.Cover),
			DateOfBirth: helpers.TimeToISOString(v.DateOfBirth),
			Sex:         helpers.PtrIfNotNil(v.Sex),
			Goal:        helpers.PtrIfNotNil(v.Goal),
			HeightCm:    v.HeightCm,
			WeightKg:    v.WeightKg,
			Status:      helpers.PtrIfNotNil(v.Status),
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
		ID:          v.ID,
		Email:       v.Email,
		Name:        v.Name,
		Phone:       v.Phone,
		Bio:         helpers.PtrIfNotNil(v.Bio),
		Avatar:      helpers.PtrIfNotNil(v.Avatar),
		Cover:       helpers.PtrIfNotNil(v.Cover),
		DateOfBirth: helpers.TimeToISOString(v.DateOfBirth),
		Sex:         helpers.PtrIfNotNil(v.Sex),
		Goal:        helpers.PtrIfNotNil(v.Goal),
		HeightCm:    v.HeightCm,
		WeightKg:    v.WeightKg,
		Status:      helpers.PtrIfNotNil(v.Status),
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateUserDTO) (*GetUserDTO, error) {
	q := r.c.User.Create().
		SetEmail(dto.Email).
		SetName(dto.Name).
		SetNillablePhone(dto.Phone)

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
	if dto.DateOfBirth != nil {
		if parsedTime, err := time.Parse(time.RFC3339, *dto.DateOfBirth); err == nil {
			q.SetDateOfBirth(parsedTime)
		}
	}
	if dto.Sex != nil {
		q.SetNillableSex(dto.Sex)
	}
	if dto.Goal != nil {
		q.SetNillableGoal(dto.Goal)
	}
	if dto.HeightCm != nil {
		q.SetNillableHeightCm(dto.HeightCm)
	}
	if dto.WeightKg != nil {
		q.SetNillableWeightKg(dto.WeightKg)
	}
	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	q.SetPassword(string(hashed))

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetUserDTO{
		ID:          row.ID,
		Email:       row.Email,
		Name:        row.Name,
		Phone:       row.Phone,
		Bio:         helpers.PtrIfNotNil(row.Bio),
		Avatar:      helpers.PtrIfNotNil(row.Avatar),
		Cover:       helpers.PtrIfNotNil(row.Cover),
		DateOfBirth: helpers.TimeToISOString(row.DateOfBirth),
		Sex:         helpers.PtrIfNotNil(row.Sex),
		Goal:        helpers.PtrIfNotNil(row.Goal),
		HeightCm:    row.HeightCm,
		WeightKg:    row.WeightKg,
		Status:      helpers.PtrIfNotNil(row.Status),
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateUserDTO) (*GetUserDTO, error) {
	q := r.c.User.UpdateOneID(dto.ID)

	if dto.Email != nil {
		q.SetEmail(*dto.Email)
	}

	if dto.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		q.SetPassword(string(hashed))
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
	if dto.DateOfBirth != nil {
		if parsedTime, err := time.Parse(time.RFC3339, *dto.DateOfBirth); err == nil {
			q.SetDateOfBirth(parsedTime)
		}
	}
	if dto.Sex != nil {
		q.SetNillableSex(dto.Sex)
	}
	if dto.Goal != nil {
		q.SetNillableGoal(dto.Goal)
	}
	if dto.HeightCm != nil {
		q.SetNillableHeightCm(dto.HeightCm)
	}
	if dto.WeightKg != nil {
		q.SetNillableWeightKg(dto.WeightKg)
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
		ID:          row.ID,
		Email:       row.Email,
		Name:        row.Name,
		Phone:       row.Phone,
		Bio:         helpers.PtrIfNotNil(row.Bio),
		Avatar:      helpers.PtrIfNotNil(row.Avatar),
		Cover:       helpers.PtrIfNotNil(row.Cover),
		DateOfBirth: helpers.TimeToISOString(row.DateOfBirth),
		Sex:         helpers.PtrIfNotNil(row.Sex),
		Goal:        helpers.PtrIfNotNil(row.Goal),
		HeightCm:    row.HeightCm,
		WeightKg:    row.WeightKg,
		Status:      helpers.PtrIfNotNil(row.Status),
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.User.DeleteOneID(id).Exec(ctx)
}
