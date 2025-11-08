package meal_plans

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/meal_plan"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetMeal_planDTO, error) {
	rows, err := r.c.Meal_plan.Query().
		WithUser().
		Order(ent.Asc(meal_plan.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetMeal_planDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetMeal_planDTO{
			ID:        v.ID,
			WeekStart: v.WeekStart,
			Goal:      v.Goal,
		}
		if v.Edges.User != nil {
			dto.UserID = v.Edges.User.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetMeal_planDTO, error) {
	v, err := r.c.Meal_plan.Query().
		WithUser().
		Where(meal_plan.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetMeal_planDTO{
		ID:        v.ID,
		WeekStart: v.WeekStart,
		Goal:      v.Goal,
	}
	if v.Edges.User != nil {
		dto.UserID = v.Edges.User.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateMeal_planDTO) (*GetMeal_planDTO, error) {
	user, err := r.c.User.Get(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}

	q := r.c.Meal_plan.
		Create().
		SetID(dto.ID).
		SetWeekStart(dto.WeekStart).
		SetUser(user)
	if dto.Goal != nil {
		q.SetGoal(*dto.Goal)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetMeal_planDTO{
		ID:        row.ID,
		WeekStart: row.WeekStart,
		Goal:      row.Goal,
		UserID:    dto.UserID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateMeal_planDTO) (*GetMeal_planDTO, error) {
	q := r.c.Meal_plan.UpdateOneID(dto.ID)

	if dto.WeekStart != nil {
		q.SetWeekStart(*dto.WeekStart)
	}
	if dto.Goal != nil {
		q.SetGoal(*dto.Goal)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with user edge
	v, err := r.c.Meal_plan.Query().
		WithUser().
		Where(meal_plan.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetMeal_planDTO{
		ID:        v.ID,
		WeekStart: v.WeekStart,
		Goal:      v.Goal,
	}
	if v.Edges.User != nil {
		dtoOut.UserID = v.Edges.User.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Meal_plan.DeleteOneID(id).Exec(ctx)
}
