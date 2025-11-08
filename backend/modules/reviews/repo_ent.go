package reviews

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/review"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetReviewDTO, error) {
	rows, err := r.c.Review.Query().
		WithUser().
		WithProduct().
		Order(ent.Asc(review.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetReviewDTO, 0, len(rows))
	for _, v := range rows {
		dto := &GetReviewDTO{
			ID:        v.ID,
			Rating:    v.Rating,
			Comment:   v.Comment,
			CreatedAt: v.CreatedAt,
		}
		if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
			dto.UserID = v.Edges.User[0].ID
		}
		if len(v.Edges.Product) > 0 && v.Edges.Product[0] != nil {
			dto.ProductID = v.Edges.Product[0].ID
		}
		out = append(out, dto)
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetReviewDTO, error) {
	v, err := r.c.Review.Query().
		WithUser().
		WithProduct().
		Where(review.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	dto := &GetReviewDTO{
		ID:        v.ID,
		Rating:    v.Rating,
		Comment:   v.Comment,
		CreatedAt: v.CreatedAt,
	}
	if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
		dto.UserID = v.Edges.User[0].ID
	}
	if len(v.Edges.Product) > 0 && v.Edges.Product[0] != nil {
		dto.ProductID = v.Edges.Product[0].ID
	}
	return dto, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateReviewDTO) (*GetReviewDTO, error) {
	user, err := r.c.User.Get(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}
	product, err := r.c.Product.Get(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}

	q := r.c.Review.
		Create().
		SetID(dto.ID).
		SetRating(dto.Rating).
		AddUser(user).
		AddProduct(product)

	if dto.Comment != nil {
		q.SetComment(*dto.Comment)
	}
	if dto.CreatedAt != nil {
		q.SetCreatedAt(*dto.CreatedAt)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetReviewDTO{
		ID:        row.ID,
		Rating:    row.Rating,
		Comment:   row.Comment,
		CreatedAt: row.CreatedAt,
		UserID:    dto.UserID,
		ProductID: dto.ProductID,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateReviewDTO) (*GetReviewDTO, error) {
	q := r.c.Review.UpdateOneID(dto.ID)

	if dto.Rating != nil {
		q.SetRating(*dto.Rating)
	}
	if dto.Comment != nil {
		q.SetComment(*dto.Comment)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	v, err := r.c.Review.Query().
		WithUser().
		WithProduct().
		Where(review.ID(row.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dtoOut := &GetReviewDTO{
		ID:        v.ID,
		Rating:    v.Rating,
		Comment:   v.Comment,
		CreatedAt: v.CreatedAt,
	}
	if len(v.Edges.User) > 0 && v.Edges.User[0] != nil {
		dtoOut.UserID = v.Edges.User[0].ID
	}
	if len(v.Edges.Product) > 0 && v.Edges.Product[0] != nil {
		dtoOut.ProductID = v.Edges.Product[0].ID
	}
	return dtoOut, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Review.DeleteOneID(id).Exec(ctx)
}
