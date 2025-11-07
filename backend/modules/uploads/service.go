package uploads

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"freshease/backend/internal/common/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Service interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	DeleteImage(ctx context.Context, objectName string) error
	GetImageURL(ctx context.Context, objectName string) (string, error)
}

type service struct {
	minioClient *minio.Client
	bucket      string
	endpoint    string
	useSSL      bool
}

func NewService(cfg config.MinIOConfig) (Service, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &service{
		minioClient: minioClient,
		bucket:      cfg.Bucket,
		endpoint:    cfg.Endpoint,
		useSSL:      cfg.UseSSL,
	}, nil
}

func (s *service) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return "", fmt.Errorf("invalid file type. Allowed types: %v", allowedExts)
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return "", fmt.Errorf("file size exceeds 10MB limit")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload to MinIO
	contentType := "image/" + strings.TrimPrefix(ext, ".")
	if contentType == "image/jpg" {
		contentType = "image/jpeg"
	}

	_, err = s.minioClient.PutObject(ctx, s.bucket, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the object name (path)
	return filename, nil
}

func (s *service) DeleteImage(ctx context.Context, objectName string) error {
	err := s.minioClient.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *service) GetImageURL(ctx context.Context, objectName string) (string, error) {
	// Generate presigned URL valid for 7 days
	expiry := 7 * 24 * time.Hour
	reqParams := make(url.Values)

	presignedURL, err := s.minioClient.PresignedGetObject(ctx, s.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}
