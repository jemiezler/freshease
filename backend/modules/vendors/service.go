package vendors

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"freshease/backend/modules/uploads"
)

type Service interface {
	List(ctx context.Context) ([]*GetVendorDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetVendorDTO, error)
	Create(ctx context.Context, dto CreateVendorDTO) (*GetVendorDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateVendorDTO) (*GetVendorDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UploadVendorLogo(ctx context.Context, file *multipart.FileHeader) (string, error)
	GetVendorImageURL(ctx context.Context, objectName string) (string, error)
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

func (s *service) List(ctx context.Context) ([]*GetVendorDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetVendorDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateVendorDTO) (*GetVendorDTO, error) {
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateVendorDTO) (*GetVendorDTO, error) {
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// UploadVendorLogo uploads a vendor logo image to MinIO
func (s *service) UploadVendorLogo(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return s.uploadsSvc.UploadImage(ctx, file, "vendors/logos")
}

// GetVendorImageURL generates a presigned URL for a vendor image
func (s *service) GetVendorImageURL(ctx context.Context, objectName string) (string, error) {
	return s.uploadsSvc.GetImageURL(ctx, objectName)
}
