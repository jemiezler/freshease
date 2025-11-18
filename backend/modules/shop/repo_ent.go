package shop

import (
	"context"
	"strings"

	"freshease/backend/ent"
	"freshease/backend/ent/category"
	"freshease/backend/ent/product"
	"freshease/backend/ent/product_category"
	"freshease/backend/ent/vendor"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) GetActiveProducts(ctx context.Context, filters ShopSearchFilters) ([]*ShopProductDTO, int, error) {
	query := r.c.Product.Query().
		Where(product.IsActive(true)).
		WithVendor().
		WithProductCategories(func(q *ent.ProductCategoryQuery) {
			q.WithCategory()
		}).
		WithInventories()

	// Apply filters
	if filters.CategoryID != nil {
		query = query.Where(product.HasProductCategoriesWith(product_category.HasCategoryWith(category.ID(*filters.CategoryID))))
	}
	if filters.VendorID != nil {
		query = query.Where(product.HasVendorWith(vendor.ID(*filters.VendorID)))
	}
	if filters.MinPrice != nil {
		query = query.Where(product.PriceGTE(*filters.MinPrice))
	}
	if filters.MaxPrice != nil {
		query = query.Where(product.PriceLTE(*filters.MaxPrice))
	}
	if filters.SearchTerm != nil && *filters.SearchTerm != "" {
		searchTerm := strings.ToLower(*filters.SearchTerm)
		query = query.Where(
			product.Or(
				product.NameContains(searchTerm),
				product.DescriptionContains(searchTerm),
			),
		)
	}
	if filters.InStock != nil && *filters.InStock {
		query = query.Where(product.HasInventoriesWith())
	}

	// Get total count for pagination
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Execute query
	products, err := query.Order(ent.Asc(product.FieldName)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Convert to DTOs
	result := make([]*ShopProductDTO, 0, len(products))
	for _, p := range products {
		dto := &ShopProductDTO{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			UnitLabel: p.UnitLabel,
			IsActive:  boolToString(p.IsActive),
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
		if p.Description != nil {
			dto.Description = *p.Description
		} else {
			dto.Description = ""
		}

		// Add vendor info
		if p.Edges.Vendor != nil {
			dto.VendorID = p.Edges.Vendor.ID
			dto.VendorName = getStringValue(p.Edges.Vendor.Name)
		}

		// Add category info (product_categories is a many-to-many relationship)
		if len(p.Edges.ProductCategories) > 0 {
			// Get the first category
			pc := p.Edges.ProductCategories[0]
			if pc.Edges.Category != nil {
				dto.CategoryID = pc.Edges.Category.ID
				dto.CategoryName = pc.Edges.Category.Name
			}
		}

		// Add inventory info
		if len(p.Edges.Inventories) > 0 && p.Edges.Inventories[0] != nil {
			dto.StockQuantity = p.Edges.Inventories[0].Quantity
			dto.IsInStock = p.Edges.Inventories[0].Quantity > 0
		}

		// Add image object name (path, not URL)
		// Clients should use /api/uploads/{object_name} to get presigned URLs
		if p.ImageURL != nil {
			dto.ImageURL = *p.ImageURL
		} else {
			dto.ImageURL = ""
		}

		result = append(result, dto)
	}

	return result, total, nil
}

func (r *EntRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error) {
	p, err := r.c.Product.Query().
		Where(product.ID(id), product.IsActive(true)).
		WithVendor().
		WithProductCategories(func(q *ent.ProductCategoryQuery) {
			q.WithCategory()
		}).
		WithInventories().
		First(ctx)
	if err != nil {
		return nil, err
	}

	dto := &ShopProductDTO{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		UnitLabel: p.UnitLabel,
		IsActive:  boolToString(p.IsActive),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if p.Description != nil {
		dto.Description = *p.Description
	} else {
		dto.Description = ""
	}

	// Add image object name (path, not URL)
	// Clients should use /api/uploads/{object_name} to get presigned URLs
	if p.ImageURL != nil {
		dto.ImageURL = *p.ImageURL
	} else {
		dto.ImageURL = ""
	}

	// Add vendor info
	if p.Edges.Vendor != nil {
		dto.VendorID = p.Edges.Vendor.ID
		dto.VendorName = getStringValue(p.Edges.Vendor.Name)
	}

	// Add category info (product_categories is a many-to-many relationship)
	if len(p.Edges.ProductCategories) > 0 {
		// Get the first category
		pc := p.Edges.ProductCategories[0]
		if pc.Edges.Category != nil {
			dto.CategoryID = pc.Edges.Category.ID
			dto.CategoryName = pc.Edges.Category.Name
		}
	}

	// Add inventory info
	if len(p.Edges.Inventories) > 0 && p.Edges.Inventories[0] != nil {
		dto.StockQuantity = p.Edges.Inventories[0].Quantity
		dto.IsInStock = p.Edges.Inventories[0].Quantity > 0
	}

	return dto, nil
}

func (r *EntRepo) GetActiveCategories(ctx context.Context) ([]*ShopCategoryDTO, error) {
	categories, err := r.c.Category.Query().
		Order(ent.Asc(category.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ShopCategoryDTO, 0, len(categories))
	for _, c := range categories {
		result = append(result, &ShopCategoryDTO{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Slug, // Using slug as description since DTO doesn't have slug field
		})
	}

	return result, nil
}

func (r *EntRepo) GetActiveVendors(ctx context.Context) ([]*ShopVendorDTO, error) {
	vendors, err := r.c.Vendor.Query().
		Order(ent.Asc(vendor.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ShopVendorDTO, 0, len(vendors))
	for _, v := range vendors {
		result = append(result, &ShopVendorDTO{
			ID:   v.ID,
			Name: getStringValue(v.Name),
		})
	}

	return result, nil
}

func (r *EntRepo) GetCategoryByID(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error) {
	c, err := r.c.Category.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ShopCategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Slug, // Using slug as description since DTO doesn't have slug field
	}, nil
}

func (r *EntRepo) GetVendorByID(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error) {
	v, err := r.c.Vendor.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ShopVendorDTO{
		ID:   v.ID,
		Name: getStringValue(v.Name),
	}, nil
}

// Helper function to safely dereference string pointers
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to convert bool to string
func boolToString(b bool) string {
	if b {
		return "active"
	}
	return "inactive"
}
