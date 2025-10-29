package shop

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Get all active products with vendor and category info
	GetActiveProducts(ctx context.Context, filters ShopSearchFilters) ([]*ShopProductDTO, int, error)

	// Get product by ID with full details
	GetProductByID(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error)

	// Get all active categories
	GetActiveCategories(ctx context.Context) ([]*ShopCategoryDTO, error)

	// Get all active vendors
	GetActiveVendors(ctx context.Context) ([]*ShopVendorDTO, error)

	// Get category by ID
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error)

	// Get vendor by ID
	GetVendorByID(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error)
}
