package products

import (
	"context"
	"mime/multipart"

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
	repo       Repository
	uploadsSvc uploads.Service
}

func NewService(r Repository, uploadsSvc uploads.Service) Service {
	return &service{
		repo:       r,
		uploadsSvc: uploadsSvc,
	}
}

func (s *service) List(ctx context.Context) ([]*GetProductDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetProductDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateProductDTO) (*GetProductDTO, error) {
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateProductDTO) (*GetProductDTO, error) {
	dto.ID = id
	return s.repo.Update(ctx, &dto)
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
