package users

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"freshease/backend/modules/uploads"
)

type Service interface {
	List(ctx context.Context) ([]*GetUserDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetUserDTO, error)
	Create(ctx context.Context, dto CreateUserDTO) (*GetUserDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*GetUserDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UploadUserAvatar(ctx context.Context, file *multipart.FileHeader) (string, error)
	UploadUserCover(ctx context.Context, file *multipart.FileHeader) (string, error)
	GetUserImageURL(ctx context.Context, objectName string) (string, error)
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

func (s *service) List(ctx context.Context) ([]*GetUserDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetUserDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateUserDTO) (*GetUserDTO, error) {
	// pass the inbound DTO straight to the repo
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*GetUserDTO, error) {
	// ensure the DTO has the ID (path param is source of truth)
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// UploadUserAvatar uploads a user avatar image to MinIO
func (s *service) UploadUserAvatar(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return s.uploadsSvc.UploadImage(ctx, file, "users/avatars")
}

// UploadUserCover uploads a user cover image to MinIO
func (s *service) UploadUserCover(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return s.uploadsSvc.UploadImage(ctx, file, "users/covers")
}

// GetUserImageURL generates a presigned URL for a user image
func (s *service) GetUserImageURL(ctx context.Context, objectName string) (string, error) {
	return s.uploadsSvc.GetImageURL(ctx, objectName)
}
