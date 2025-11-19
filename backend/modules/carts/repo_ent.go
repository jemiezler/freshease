package carts

import (
	"context"

	"freshease/backend/ent"
	"freshease/backend/ent/cart"
	"freshease/backend/ent/user"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*GetCartDTO, error) {
	rows, err := r.c.Cart.Query().Order(ent.Asc(cart.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*GetCartDTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &GetCartDTO{
			ID:        v.ID,
			Status:    v.Status,
			Total:     v.Total,
			Subtotal: v.Subtotal,
		Discount: v.Discount,
			UpdatedAt: v.UpdatedAt,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*GetCartDTO, error) {
	v, err := r.c.Cart.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetCartDTO{
		ID:        v.ID,
		Status:    v.Status,
		Total:     v.Total,
		Subtotal: v.Subtotal,
		Discount: v.Discount,
		UpdatedAt: v.UpdatedAt,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *CreateCartDTO) (*GetCartDTO, error) {
	q := r.c.Cart.Create()

	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Total != nil {
		q.SetTotal(*dto.Total)
	}
	if dto.UserID != nil {
		user, err := r.c.User.Get(ctx, *dto.UserID)
		if err != nil {
			return nil, err
		}
		q.AddUser(user)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCartDTO{
		ID:        row.ID,
		Status:    row.Status,
		Subtotal:  row.Subtotal,
		Discount:  row.Discount,
		Total:     row.Total,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *UpdateCartDTO) (*GetCartDTO, error) {
	q := r.c.Cart.UpdateOneID(dto.ID)

	if dto.Status != nil {
		q.SetStatus(*dto.Status)
	}
	if dto.Total != nil {
		q.SetTotal(*dto.Total)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCartDTO{
		ID:        row.ID,
		Status:    row.Status,
		Subtotal:  row.Subtotal,
		Discount:  row.Discount,
		Total:     row.Total,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.Cart.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error) {
	v, err := r.c.Cart.Query().
		Where(cart.HasUserWith(user.ID(userID))).
		WithItems(func(q *ent.CartItemQuery) {
			q.WithProduct()
		}).
		Order(ent.Desc(cart.FieldUpdatedAt)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return r.cartToDTO(v), nil
}

func (r *EntRepo) GetOrCreateCartForUser(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error) {
	// Try to find existing cart
	cartEntity, err := r.c.Cart.Query().
		Where(cart.HasUserWith(user.ID(userID))).
		WithItems(func(q *ent.CartItemQuery) {
			q.WithProduct()
		}).
		Order(ent.Desc(cart.FieldUpdatedAt)).
		First(ctx)
	
	if err == nil {
		return r.cartToDTO(cartEntity), nil
	}
	
	// If not found, create a new cart
	userEntity, err := r.c.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	newCart, err := r.c.Cart.Create().
		SetStatus("pending").
		SetSubtotal(0.0).
		SetDiscount(0.0).
		SetTotal(0.0).
		AddUser(userEntity).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	
	// Reload with items (empty)
	v, err := r.c.Cart.Query().
		Where(cart.ID(newCart.ID)).
		WithItems(func(q *ent.CartItemQuery) {
			q.WithProduct()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	
	return r.cartToDTO(v), nil
}

// Helper function to convert ent.Cart to GetCartDTO with items
func (r *EntRepo) cartToDTO(c *ent.Cart) *GetCartDTO {
	dto := &GetCartDTO{
		ID:            c.ID,
		Status:        c.Status,
		Subtotal:      c.Subtotal,
		Discount:      c.Discount,
		Total:         c.Total,
		Shipping:      0.0, // Will be calculated in service
		Tax:           0.0, // Will be calculated in service
		Items:         []CartItemDTO{},
		PromoCode:     nil,
		PromoDiscount: 0.0,
		CreatedAt:     c.UpdatedAt, // Using UpdatedAt as fallback
		UpdatedAt:     c.UpdatedAt,
	}
	
		// Convert cart items
		if c.Edges.Items != nil {
			dto.Items = make([]CartItemDTO, 0, len(c.Edges.Items))
			for _, item := range c.Edges.Items {
				itemDTO := CartItemDTO{
					ID:          item.ID,
					Quantity:    item.Qty,
					ProductPrice: item.UnitPrice,
					LineTotal:   item.LineTotal,
					ProductName: "", // Default empty
					ProductImage: nil, // Default nil
				}
				
				if item.Edges.Product != nil {
					itemDTO.ProductID = item.Edges.Product.ID
					itemDTO.ProductName = item.Edges.Product.Name
					if item.Edges.Product.ImageURL != nil {
						imageURL := *item.Edges.Product.ImageURL
						itemDTO.ProductImage = &imageURL
					} else {
						emptyStr := ""
						itemDTO.ProductImage = &emptyStr
					}
				} else {
					// If product is nil, set empty values
					emptyStr := ""
					itemDTO.ProductImage = &emptyStr
					// ProductID will be zero UUID which is fine
				}
				
				dto.Items = append(dto.Items, itemDTO)
			}
		}
	
	return dto
}
