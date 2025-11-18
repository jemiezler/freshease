package products

import (
	"context"
	"mime/multipart"

	"freshease/backend/modules/product_categories"
	"freshease/backend/modules/uploads"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*GetProductDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetProductDTO, error)
	Create(ctx context.Context, dto CreateProductDTO) (*GetProductDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateProductDTO) (*GetProductDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UploadProductImage(ctx context.Context, file *multipart.FileHeader) (string, error)
	GetProductImageURL(ctx context.Context, objectName string) (string, error)
}

type service struct {
	repo              Repository
	uploadsSvc        uploads.Service
	productCategorySvc product_categories.Service
}

func NewService(r Repository, uploadsSvc uploads.Service) Service {
	return &service{
		repo:       r,
		uploadsSvc: uploadsSvc,
	}
}

// NewServiceWithProductCategories creates a service with product categories support
func NewServiceWithProductCategories(r Repository, uploadsSvc uploads.Service, productCategorySvc product_categories.Service) Service {
	return &service{
		repo:              r,
		uploadsSvc:        uploadsSvc,
		productCategorySvc: productCategorySvc,
	}
}

func (s *service) List(ctx context.Context) ([]*GetProductDTO, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert image object names to URLs
	for _, product := range products {
		if product.ImageURL != nil && *product.ImageURL != "" {
			url, err := s.uploadsSvc.GetImageURL(ctx, *product.ImageURL)
			if err == nil {
				product.ImageURL = &url
			}
		}
	}
	
	return products, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetProductDTO, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Convert image object name to URL
	if product.ImageURL != nil && *product.ImageURL != "" {
		url, err := s.uploadsSvc.GetImageURL(ctx, *product.ImageURL)
		if err == nil {
			product.ImageURL = &url
		}
	}
	
	return product, nil
}

func (s *service) Create(ctx context.Context, dto CreateProductDTO) (*GetProductDTO, error) {
	// Create the product first
	product, err := s.repo.Create(ctx, &dto)
	if err != nil {
		return nil, err
	}

	// If product categories service is available and category IDs are provided, create associations
	if s.productCategorySvc != nil && len(dto.CategoryIDs) > 0 {
		for _, categoryID := range dto.CategoryIDs {
			_, err := s.productCategorySvc.Create(ctx, product_categories.CreateProductCategoryDTO{
				ID:         uuid.New(),
				ProductID:  product.ID,
				CategoryID: categoryID,
			})
			if err != nil {
				// Return error if category creation fails - this ensures data consistency
				return nil, err
			}
		}
	}

	// Convert image object name to URL before returning
	if product.ImageURL != nil && *product.ImageURL != "" {
		url, err := s.uploadsSvc.GetImageURL(ctx, *product.ImageURL)
		if err == nil {
			product.ImageURL = &url
		}
	}

	return product, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateProductDTO) (*GetProductDTO, error) {
	dto.ID = id
	product, err := s.repo.Update(ctx, &dto)
	if err != nil {
		return nil, err
	}
	
	// Convert image object name to URL
	if product.ImageURL != nil && *product.ImageURL != "" {
		url, err := s.uploadsSvc.GetImageURL(ctx, *product.ImageURL)
		if err == nil {
			product.ImageURL = &url
		}
	}
	
	return product, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// UploadProductImage uploads a product image to MinIO
func (s *service) UploadProductImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return s.uploadsSvc.UploadImage(ctx, file, "products")
}

// GetProductImageURL generates a presigned URL for a product image
func (s *service) GetProductImageURL(ctx context.Context, objectName string) (string, error) {
	return s.uploadsSvc.GetImageURL(ctx, objectName)
}
