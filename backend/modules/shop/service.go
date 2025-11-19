package shop

import (
	"context"

	"freshease/backend/modules/uploads"

	"github.com/google/uuid"
)

type Service interface {
	// Search products with filters
	SearchProducts(ctx context.Context, filters ShopSearchFilters) (*ShopSearchResponse, error)

	// Get product details by ID
	GetProduct(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error)

	// Get all categories
	GetCategories(ctx context.Context) ([]*ShopCategoryDTO, error)

	// Get all vendors
	GetVendors(ctx context.Context) ([]*ShopVendorDTO, error)

	// Get category by ID
	GetCategory(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error)

	// Get vendor by ID
	GetVendor(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error)
}

type service struct {
	repo       Repository
	uploadsSvc uploads.Service
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func NewServiceWithUploads(r Repository, uploadsSvc uploads.Service) Service {
	return &service{repo: r, uploadsSvc: uploadsSvc}
}

func (s *service) SearchProducts(ctx context.Context, filters ShopSearchFilters) (*ShopSearchResponse, error) {
	// Set default pagination values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Cap at 100 for performance
	}

	products, total, err := s.repo.GetActiveProducts(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Convert image object names to URLs if uploads service is available
	if s.uploadsSvc != nil {
		for _, product := range products {
			if product.ImageURL != "" {
				url, err := s.uploadsSvc.GetImageURL(ctx, product.ImageURL)
				if err == nil {
					product.ImageURL = url
				}
			}
		}
	}

	hasMore := filters.Offset+filters.Limit < total

	return &ShopSearchResponse{
		Products: products,
		Total:    total,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
		HasMore:  hasMore,
	}, nil
}

func (s *service) GetProduct(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert image object name to URL if uploads service is available
	if s.uploadsSvc != nil && product.ImageURL != "" {
		url, err := s.uploadsSvc.GetImageURL(ctx, product.ImageURL)
		if err == nil {
			product.ImageURL = url
		}
	}

	return product, nil
}

func (s *service) GetCategories(ctx context.Context) ([]*ShopCategoryDTO, error) {
	return s.repo.GetActiveCategories(ctx)
}

func (s *service) GetVendors(ctx context.Context) ([]*ShopVendorDTO, error) {
	return s.repo.GetActiveVendors(ctx)
}

func (s *service) GetCategory(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *service) GetVendor(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error) {
	return s.repo.GetVendorByID(ctx, id)
}
