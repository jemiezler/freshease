package meal_plan_items

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/meal_plan_item"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetMeal_plan_itemDTO, error) {
	rows, err := r.c.Meal_plan_item.Query().
		WithMealPlan().
		WithRecipe().
		Order(ent.Asc(meal_plan_item.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetMeal_plan_itemDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetMeal_plan_itemDTO{
			ID:   v.ID,
			Day:  v.Day,
			Slot: v.Slot,
		}
		if v.Edges.MealPlan != nil {
			dto.MealPlanID = v.Edges.MealPlan.ID
		}
		if v.Edges.Recipe != nil {
			dto.RecipeID = v.Edges.Recipe.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetMeal_plan_itemDTO, error) {
	v, err := r.c.Meal_plan_item.Query().
		WithMealPlan().
		WithRecipe().
		Where(meal_plan_item.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetMeal_plan_itemDTO{
		ID:   v.ID,
		Day:  v.Day,
		Slot: v.Slot,
	}
	if v.Edges.MealPlan != nil {
		dto.MealPlanID = v.Edges.MealPlan.ID
	}
	if v.Edges.Recipe != nil {
		dto.RecipeID = v.Edges.Recipe.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error) {
	mealPlan, err := r.c.Meal_plan.Get(ctx, dto.MealPlanID)
	if err != nil {
		return nil, err
	}
	recipe, err := r.c.Recipe.Get(ctx, dto.RecipeID)
	if err != nil {
		return nil, err
	}

	row, err := r.c.Meal_plan_item.
		Create().
		SetID(dto.ID).
		SetDay(dto.Day).
		SetSlot(dto.Slot).
		SetMealPlan(mealPlan).
		SetRecipe(recipe).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetMeal_plan_itemDTO{
		ID:         row.ID,
		Day:        row.Day,
		Slot:       row.Slot,
		MealPlanID: dto.MealPlanID,
		RecipeID:   dto.RecipeID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error) {
	q := r.c.Meal_plan_item.UpdateOneID(dto.ID)

	if dto.Day != nil {
		q.SetDay(*dto.Day)
	}
	if dto.Slot != nil {
		q.SetSlot(*dto.Slot)
	}
	if dto.MealPlanID != nil {
		mealPlan, err := r.c.Meal_plan.Get(ctx, *dto.MealPlanID)
		if err != nil {
			return nil, err
		}
		q.SetMealPlan(mealPlan)
	}
	if dto.RecipeID != nil {
		recipe, err := r.c.Recipe.Get(ctx, *dto.RecipeID)
		if err != nil {
			return nil, err
		}
		q.SetRecipe(recipe)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Meal_plan_item.Query().
		WithMealPlan().
		WithRecipe().
		Where(meal_plan_item.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetMeal_plan_itemDTO{
		ID:   v.ID,
		Day:  v.Day,
		Slot: v.Slot,
	}
	if v.Edges.MealPlan != nil {
		dtoOut.MealPlanID = v.Edges.MealPlan.ID
	}
	if v.Edges.Recipe != nil {
		dtoOut.RecipeID = v.Edges.Recipe.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Meal_plan_item.DeleteOneID(id).Exec(ctx)
}
