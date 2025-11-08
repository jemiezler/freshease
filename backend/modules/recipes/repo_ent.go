package recipes

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/recipe"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetRecipeDTO, error) {
	rows, err := r.c.Recipe.Query().Order(ent.Asc(recipe.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetRecipeDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetRecipeDTO{
			ID:           v.ID,
			Name:         v.Name,
			Instructions: v.Instructions,
			Kcal:         v.Kcal,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetRecipeDTO, error) {
	v, err := r.c.Recipe.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetRecipeDTO{
		ID:           v.ID,
		Name:         v.Name,
		Instructions: v.Instructions,
		Kcal:         v.Kcal,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateRecipeDTO) (*GetRecipeDTO, error) {
	q := r.c.Recipe.
		Create().
		SetID(dto.ID).
		SetName(dto.Name).
		SetKcal(dto.Kcal)

	if dto.Instructions != nil {
		q.SetInstructions(*dto.Instructions)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetRecipeDTO{
		ID:           row.ID,
		Name:         row.Name,
		Instructions: row.Instructions,
		Kcal:         row.Kcal,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateRecipeDTO) (*GetRecipeDTO, error) {
	q := r.c.Recipe.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Instructions != nil {
		q.SetInstructions(*dto.Instructions)
	}
	if dto.Kcal != nil {
		q.SetKcal(*dto.Kcal)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetRecipeDTO{
		ID:           row.ID,
		Name:         row.Name,
		Instructions: row.Instructions,
		Kcal:         row.Kcal,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Recipe.DeleteOneID(id).Exec(ctx)
}
