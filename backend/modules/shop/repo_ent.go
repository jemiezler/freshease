package shop

import (
	"context"
	"strings"

	"freshease/backend/ent"
	"freshease/backend/ent/product"
	"freshease/backend/ent/product_category"
	"freshease/backend/ent/vendor"

	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) GetActiveProducts(ctx context.Context, filters ShopSearchFilters) ([]*ShopProductDTO, int, error) {
	query := r.c.Product.Query().
		Where(product.IsActive("active")).
		WithVendor().
		WithCatagory().
		WithInventory()

	// Apply filters
	if filters.CategoryID != nil {
		query = query.Where(product.HasCatagoryWith(product_category.ID(*filters.CategoryID)))
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
		query = query.Where(product.HasInventoryWith())
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
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			ImageURL:    p.ImageURL,
			UnitLabel:   p.UnitLabel,
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}

		// Add vendor info
		if len(p.Edges.Vendor) > 0 {
			dto.VendorID = p.Edges.Vendor[0].ID
			dto.VendorName = getStringValue(p.Edges.Vendor[0].Name)
		}

		// Add category info
		if len(p.Edges.Catagory) > 0 {
			dto.CategoryID = p.Edges.Catagory[0].ID
			dto.CategoryName = p.Edges.Catagory[0].Name
		}

		// Add inventory info
		if p.Edges.Inventory != nil {
			dto.StockQuantity = p.Edges.Inventory.Quantity
			dto.IsInStock = p.Edges.Inventory.Quantity > 0
		}

		result = append(result, dto)
	}

	return result, total, nil
}

func (r *EntRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error) {
	p, err := r.c.Product.Query().
		Where(product.ID(id), product.IsActive("active")).
		WithVendor().
		WithCatagory().
		WithInventory().
		First(ctx)
	if err != nil {
		return nil, err
	}

	dto := &ShopProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		UnitLabel:   p.UnitLabel,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	// Add vendor info
	if len(p.Edges.Vendor) > 0 {
		dto.VendorID = p.Edges.Vendor[0].ID
		dto.VendorName = getStringValue(p.Edges.Vendor[0].Name)
	}

	// Add category info
	if len(p.Edges.Catagory) > 0 {
		dto.CategoryID = p.Edges.Catagory[0].ID
		dto.CategoryName = p.Edges.Catagory[0].Name
	}

	// Add inventory info
	if p.Edges.Inventory != nil {
		dto.StockQuantity = p.Edges.Inventory.Quantity
		dto.IsInStock = p.Edges.Inventory.Quantity > 0
	}

	return dto, nil
}

func (r *EntRepo) GetActiveCategories(ctx context.Context) ([]*ShopCategoryDTO, error) {
	categories, err := r.c.Product_category.Query().
		Order(ent.Asc(product_category.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ShopCategoryDTO, 0, len(categories))
	for _, c := range categories {
		result = append(result, &ShopCategoryDTO{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return result, nil
}

func (r *EntRepo) GetActiveVendors(ctx context.Context) ([]*ShopVendorDTO, error) {
	vendors, err := r.c.Vendor.Query().
		Where(vendor.IsActive("active")).
		Order(ent.Asc(vendor.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ShopVendorDTO, 0, len(vendors))
	for _, v := range vendors {
		result = append(result, &ShopVendorDTO{
			ID:          v.ID,
			Name:        getStringValue(v.Name),
			Email:       getStringValue(v.Email),
			Phone:       getStringValue(v.Phone),
			Address:     getStringValue(v.Address),
			City:        getStringValue(v.City),
			State:       getStringValue(v.State),
			Country:     getStringValue(v.Country),
			PostalCode:  getStringValue(v.PostalCode),
			Website:     getStringValue(v.Website),
			LogoURL:     getStringValue(v.LogoURL),
			Description: getStringValue(v.Description),
			IsActive:    v.IsActive,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return result, nil
}

func (r *EntRepo) GetCategoryByID(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error) {
	c, err := r.c.Product_category.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ShopCategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}, nil
}

func (r *EntRepo) GetVendorByID(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error) {
	v, err := r.c.Vendor.Query().
		Where(vendor.ID(id), vendor.IsActive("active")).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &ShopVendorDTO{
		ID:          v.ID,
		IsActive:    v.IsActive,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		Name:        getStringValue(v.Name),
		Email:       getStringValue(v.Email),
		Phone:       getStringValue(v.Phone),
		Address:     getStringValue(v.Address),
		City:        getStringValue(v.City),
		State:       getStringValue(v.State),
		Country:     getStringValue(v.Country),
		PostalCode:  getStringValue(v.PostalCode),
		Website:     getStringValue(v.Website),
		LogoURL:     getStringValue(v.LogoURL),
		Description: getStringValue(v.Description),
	}, nil
}

// Helper function to safely dereference string pointers
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
