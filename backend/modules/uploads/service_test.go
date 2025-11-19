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
	"github.com/stretchr/testify/require"
)

// MockMinIOClient is a mock implementation of MinIOClient interface
type MockMinIOClient struct {
	mock.Mock
}

func (m *MockMinIOClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockMinIOClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*minio.Object), args.Error(1)
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
			svc, err := NewService(tt.cfg)
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, svc)
			}
		})
	}
}

func TestService_UploadImage(t *testing.T) {
	tests := []struct {
		name          string
		file          *multipart.FileHeader
		folder        string
		mockSetup     func(*MockMinIOClient, string, string)
		expectedError bool
		errorContains string
	}{
		{
			name:   "success - upload jpg image",
			file:   createTestFileHeader("test.jpg", 100, []byte("fake image content")),
			folder: "images",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				mockClient.On("PutObject", mock.Anything, bucket, mock.MatchedBy(func(name string) bool {
					return strings.Contains(name, "images/") && strings.HasSuffix(name, ".jpg")
				}), mock.Anything, int64(100), mock.Anything).Return(minio.UploadInfo{}, nil)
			},
			expectedError: false,
		},
		{
			name:   "error - invalid file type",
			file:   createTestFileHeader("test.txt", 100, []byte("text content")),
			folder: "images",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				// No mock setup - should fail before service call
			},
			expectedError: true,
			errorContains: "invalid file type",
		},
		{
			name:   "error - file too large",
			file:   createTestFileHeader("test.jpg", 11*1024*1024, make([]byte, 11*1024*1024)),
			folder: "images",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				// No mock setup - should fail before service call
			},
			expectedError: true,
			errorContains: "file size exceeds",
		},
		{
			name:   "error - MinIO upload fails",
			file:   createTestFileHeader("test.jpg", 100, []byte("fake image content")),
			folder: "images",
			mockSetup: func(mockClient *MockMinIOClient, bucket, filename string) {
				mockClient.On("PutObject", mock.Anything, bucket, mock.MatchedBy(func(name string) bool {
					return strings.Contains(name, "images/") && strings.HasSuffix(name, ".jpg")
				}), mock.Anything, int64(100), mock.Anything).Return(minio.UploadInfo{}, errors.New("upload failed"))
			},
			expectedError: true,
			errorContains: "failed to upload file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinIOClient)
			svc := NewServiceWithClient(mockClient, "test-bucket")
			tt.mockSetup(mockClient, "test-bucket", "")

			objectName, err := svc.UploadImage(context.Background(), tt.file, tt.folder)
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, objectName)
				assert.Contains(t, objectName, tt.folder)
			}
		})
	}
}

func TestService_DeleteImage(t *testing.T) {
	tests := []struct {
		name          string
		objectName    string
		mockSetup     func(*MockMinIOClient, string, string)
		expectedError bool
	}{
		{
			name:       "success - delete existing file",
			objectName: "images/test.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("RemoveObject", mock.Anything, bucket, objectName, mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:       "error - MinIO delete fails",
			objectName: "images/test.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("RemoveObject", mock.Anything, bucket, objectName, mock.Anything).Return(errors.New("delete failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinIOClient)
			svc := NewServiceWithClient(mockClient, "test-bucket")
			tt.mockSetup(mockClient, "test-bucket", tt.objectName)

			err := svc.DeleteImage(context.Background(), tt.objectName)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetImageURL(t *testing.T) {
	tests := []struct {
		name          string
		publicBaseURL string
		objectName    string
		mockSetup     func(*MockMinIOClient, string, string)
		expectedError bool
	}{
		{
			name:          "success - with public base URL",
			publicBaseURL: "https://example.com/storage",
			objectName:    "images/test.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				// No MinIO call needed when publicBaseURL is set
			},
			expectedError: false,
		},
		{
			name:          "success - without public base URL (presigned)",
			publicBaseURL: "",
			objectName:    "images/test.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				testURL, _ := url.Parse("http://localhost:9000/test-bucket/" + objectName)
				mockClient.On("PresignedGetObject", mock.Anything, bucket, objectName, mock.Anything, mock.Anything).Return(testURL, nil)
			},
			expectedError: false,
		},
		{
			name:          "error - MinIO URL generation fails",
			publicBaseURL: "",
			objectName:    "images/test.jpg",
			mockSetup: func(mockClient *MockMinIOClient, bucket, objectName string) {
				mockClient.On("PresignedGetObject", mock.Anything, bucket, objectName, mock.Anything, mock.Anything).Return(nil, errors.New("URL generation failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinIOClient)
			svc := NewServiceWithClient(mockClient, "test-bucket")
			// Set publicBaseURL via reflection or create a new service with it
			// For now, we'll test the presigned URL path
			if tt.publicBaseURL == "" {
				tt.mockSetup(mockClient, "test-bucket", tt.objectName)
			}

			url, err := svc.GetImageURL(context.Background(), tt.objectName)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, url)
			}
		})
	}
}
