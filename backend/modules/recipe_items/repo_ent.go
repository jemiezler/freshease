package recipe_items

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/recipe_item"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetRecipe_itemDTO, error) {
	rows, err := r.c.Recipe_item.Query().
		WithRecipe().
		WithProduct().
		Order(ent.Asc(recipe_item.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetRecipe_itemDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetRecipe_itemDTO{
			ID:     v.ID,
			Amount: v.Amount,
			Unit:   v.Unit,
		}
		if v.Edges.Recipe != nil {
			dto.RecipeID = v.Edges.Recipe.ID
		}
		if v.Edges.Product != nil {
			dto.ProductID = v.Edges.Product.ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetRecipe_itemDTO, error) {
	v, err := r.c.Recipe_item.Query().
		WithRecipe().
		WithProduct().
		Where(recipe_item.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetRecipe_itemDTO{
		ID:     v.ID,
		Amount: v.Amount,
		Unit:   v.Unit,
	}
	if v.Edges.Recipe != nil {
		dto.RecipeID = v.Edges.Recipe.ID
	}
	if v.Edges.Product != nil {
		dto.ProductID = v.Edges.Product.ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateRecipe_itemDTO) (*GetRecipe_itemDTO, error) {
	recipe, err := r.c.Recipe.Get(ctx, dto.RecipeID)
	if err != nil {
		return nil, err
	}
	product, err := r.c.Product.Get(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}

	row, err := r.c.Recipe_item.
		Create().
		SetID(dto.ID).
		SetAmount(dto.Amount).
		SetUnit(dto.Unit).
		SetRecipe(recipe).
		SetProduct(product).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetRecipe_itemDTO{
		ID:        row.ID,
		Amount:    row.Amount,
		Unit:      row.Unit,
		RecipeID:  dto.RecipeID,
		ProductID: dto.ProductID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateRecipe_itemDTO) (*GetRecipe_itemDTO, error) {
	q := r.c.Recipe_item.UpdateOneID(dto.ID)

	if dto.Amount != nil {
		q.SetAmount(*dto.Amount)
	}
	if dto.Unit != nil {
		q.SetUnit(*dto.Unit)
	}
	if dto.RecipeID != nil {
		recipe, err := r.c.Recipe.Get(ctx, *dto.RecipeID)
		if err != nil {
			return nil, err
		}
		q.SetRecipe(recipe)
	}
	if dto.ProductID != nil {
		product, err := r.c.Product.Get(ctx, *dto.ProductID)
		if err != nil {
			return nil, err
		}
		q.SetProduct(product)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Recipe_item.Query().
		WithRecipe().
		WithProduct().
		Where(recipe_item.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetRecipe_itemDTO{
		ID:     v.ID,
		Amount: v.Amount,
		Unit:   v.Unit,
	}
	if v.Edges.Recipe != nil {
		dtoOut.RecipeID = v.Edges.Recipe.ID
	}
	if v.Edges.Product != nil {
		dtoOut.ProductID = v.Edges.Product.ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Recipe_item.DeleteOneID(id).Exec(ctx)
}
