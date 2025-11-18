package uploads

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"strings"
	"testing"
	"time"

	"freshease/backend/internal/common/config"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMinIOClient is a mock implementation of MinIOClient interface
type MockMinIOClient struct {
	mock.Mock
}

func (m *MockMinIOClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockMinIOClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Error(0)
}

func (m *MockMinIOClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	args := m.Called(ctx, bucketName, objectName, expiry, reqParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*url.URL), args.Error(1)
}

func (m *MockMinIOClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	args := m.Called(ctx, bucketName)
	return args.Bool(0), args.Error(1)
}

func (m *MockMinIOClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	args := m.Called(ctx, bucketName, opts)
	return args.Error(0)
}

// createTestFileHeader creates a multipart.FileHeader for testing
func createTestFileHeader(filename string, size int64, content []byte) *multipart.FileHeader {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil
	}
	fileWriter.Write(content)
	writer.Close()

	reader := multipart.NewReader(&body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	if err != nil {
		return nil
	}
	defer form.RemoveAll()

	if len(form.File["file"]) == 0 {
		return nil
	}

	return form.File["file"][0]
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name          string
		cfg           config.MinIOConfig
		expectedError bool
		errorContains string
	}{
		{
			name: "error - invalid endpoint",
			cfg: config.MinIOConfig{
				Endpoint:        "",
				AccessKeyID:     "test",
				SecretAccessKey: "test",
				Bucket:          "test-bucket",
				UseSSL:          false,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(tt.cfg)
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UploadImage(t *testing.T) {
	tests := []struct {
		name          string
		fileHeader    *multipart.FileHeader
		folder        string
		mockSetup     func(*MockMinIOClient, string, string)
		expectedError bool
		errorContains string
	}{
		{
			name: "error - invalid file type",
			fileHeader: &multipart.FileHeader{
				Filename: "test.pdf",
				Size:     1024,
			},
			folder:        "images",
			mockSetup:     func(*MockMinIOClient, string, string) {},
			expectedError: true,
			errorContains: "invalid file type",
		},
		{
			name: "error - file size exceeds limit",
			fileHeader: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     11 * 1024 * 1024, // 11MB
			},
			folder:        "images",
			mockSetup:     func(*MockMinIOClient, string, string) {},
			expectedError: true,
			errorContains: "file size exceeds 10MB",
		},
		{
			name: "error - failed to open file",
			fileHeader: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
			},
			folder:        "images",
			mockSetup:     func(*MockMinIOClient, string, string) {},
			expectedError: true,
			errorContains: "failed to open file",
		},
		{
			name: "error - MinIO upload fails",
			fileHeader: createTestFileHeader("test.jpg", 1024, []byte("fake image content")),
			folder:     "images",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				mockClient.On("PutObject", mock.Anything, bucket, mock.MatchedBy(func(name string) bool {
					return strings.HasPrefix(name, "images/") && strings.HasSuffix(name, ".jpg")
				}), mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, errors.New("upload failed"))
			},
			expectedError: true,
			errorContains: "failed to upload file",
		},
		{
			name: "success - uploads image",
			fileHeader: createTestFileHeader("test.jpg", 1024, []byte("fake image content")),
			folder:     "images",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				mockClient.On("PutObject", mock.Anything, bucket, mock.MatchedBy(func(name string) bool {
					return strings.HasPrefix(name, "images/") && strings.HasSuffix(name, ".jpg")
				}), mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, nil)
			},
			expectedError: false,
		},
		{
			name: "success - uploads image with empty folder",
			fileHeader: createTestFileHeader("test.png", 2048, []byte("fake image content")),
			folder:     "",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				mockClient.On("PutObject", mock.Anything, bucket, mock.MatchedBy(func(name string) bool {
					return strings.HasPrefix(name, "/") && strings.HasSuffix(name, ".png")
				}), mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinIOClient)
			tt.mockSetup(mockClient, "test-bucket", "")

			svc := NewServiceWithClient(mockClient, "test-bucket")
			ctx := context.Background()

			result, err := svc.UploadImage(ctx, tt.fileHeader, tt.folder)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.True(t, strings.HasPrefix(result, tt.folder) || (tt.folder == "" && strings.HasPrefix(result, "/")))
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestService_DeleteImage(t *testing.T) {
	tests := []struct {
		name          string
		objectName    string
		mockSetup     func(*MockMinIOClient, string, string)
		expectedError bool
		errorContains string
	}{
		{
			name:       "success - deletes image",
			objectName: "images/test-uuid.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("RemoveObject", mock.Anything, bucket, objectName, mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:       "error - MinIO delete fails",
			objectName: "images/test-uuid.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("RemoveObject", mock.Anything, bucket, objectName, mock.Anything).Return(errors.New("delete failed"))
			},
			expectedError: true,
			errorContains: "failed to delete file",
		},
		{
			name:       "success - deletes nested path",
			objectName: "users/avatars/test-uuid.png",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("RemoveObject", mock.Anything, bucket, objectName, mock.Anything).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinIOClient)
			tt.mockSetup(mockClient, "test-bucket", tt.objectName)

			svc := NewServiceWithClient(mockClient, "test-bucket")
			ctx := context.Background()

			err := svc.DeleteImage(ctx, tt.objectName)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestService_GetImageURL(t *testing.T) {
	tests := []struct {
		name          string
		objectName    string
		mockSetup     func(*MockMinIOClient, string, string)
		expectedError bool
		errorContains string
		expectedURL   string
	}{
		{
			name:       "success - generates presigned URL",
			objectName: "images/test-uuid.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				testURL, _ := url.Parse("https://example.com/images/test-uuid.jpg?X-Amz-Algorithm=...")
				mockClient.On("PresignedGetObject", mock.Anything, bucket, objectName, 7*24*time.Hour, mock.Anything).Return(testURL, nil)
			},
			expectedError: false,
			expectedURL:  "https://example.com/images/test-uuid.jpg?X-Amz-Algorithm=...",
		},
		{
			name:       "error - MinIO URL generation fails",
			objectName: "images/test-uuid.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("PresignedGetObject", mock.Anything, bucket, objectName, 7*24*time.Hour, mock.Anything).Return(nil, errors.New("URL generation failed"))
			},
			expectedError: true,
			errorContains: "failed to generate presigned URL",
		},
		{
			name:       "success - generates URL for nested path",
			objectName: "users/avatars/test-uuid.png",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				testURL, _ := url.Parse("https://example.com/users/avatars/test-uuid.png?X-Amz-Algorithm=...")
				mockClient.On("PresignedGetObject", mock.Anything, bucket, objectName, 7*24*time.Hour, mock.Anything).Return(testURL, nil)
			},
			expectedError: false,
			expectedURL:  "https://example.com/users/avatars/test-uuid.png?X-Amz-Algorithm=...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinIOClient)
			tt.mockSetup(mockClient, "test-bucket", tt.objectName)

			svc := NewServiceWithClient(mockClient, "test-bucket")
			ctx := context.Background()

			result, err := svc.GetImageURL(ctx, tt.objectName)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				if tt.expectedURL != "" {
					assert.Contains(t, result, tt.objectName)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// Helper function to extract file extension (similar to service logic)
func getFileExtension(filename string) string {
	if len(filename) == 0 {
		return ""
	}
	ext := ""
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i:]
			break
		}
	}
	return ext
}
