package uploads

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, folder)
	return args.String(0), args.Error(1)
}

func (m *MockService) DeleteImage(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

func (m *MockService) GetImageURL(ctx context.Context, objectName string) (string, error) {
	args := m.Called(ctx, objectName)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetImage(ctx context.Context, objectName string) (io.ReadCloser, *minio.ObjectInfo, error) {
	args := m.Called(ctx, objectName)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(io.ReadCloser), nil, args.Error(2)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(*minio.ObjectInfo), args.Error(2)
}

func TestController_UploadImage(t *testing.T) {
	tests := []struct {
		name           string
		folder         string
		mockSetup      func(*MockService, string)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "error - missing file",
			folder: "",
			mockSetup: func(mockSvc *MockService, folder string) {
				// No mock setup - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - upload service returns error",
			folder: "images",
			mockSetup: func(mockSvc *MockService, folder string) {
				mockSvc.On("UploadImage", mock.Anything, mock.Anything, folder).Return("", errors.New("upload failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "error - URL generation fails",
			folder: "images",
			mockSetup: func(mockSvc *MockService, folder string) {
				objectName := "images/test-uuid.jpg"
				mockSvc.On("UploadImage", mock.Anything, mock.Anything, folder).Return(objectName, nil)
				mockSvc.On("GetImageURL", mock.Anything, objectName).Return("", errors.New("URL generation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			expectedFolder := tt.folder
			if expectedFolder == "" {
				expectedFolder = "images"
			}
			tt.mockSetup(mockSvc, expectedFolder)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/uploads/images", controller.UploadImage)

			// Create multipart form data only for tests that need it
			var req *http.Request
			if tt.name == "error - missing file" {
				// Request without file
				req = httptest.NewRequest(http.MethodPost, "/uploads/images", nil)
			} else {
				// Create multipart form data with file
				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				fileWriter, err := writer.CreateFormFile("file", "test.jpg")
				require.NoError(t, err)
				fileWriter.Write([]byte("fake image content"))
				if tt.folder != "" {
					writer.WriteField("folder", tt.folder)
				}
				writer.Close()
				req = httptest.NewRequest(http.MethodPost, "/uploads/images", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && tt.name != "error - missing file" {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UploadImageToFolder(t *testing.T) {
	tests := []struct {
		name           string
		folder         string
		mockSetup      func(*MockService, *multipart.FileHeader, string)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success - uploads image to folder",
			folder: "products",
			mockSetup: func(mockSvc *MockService, file *multipart.FileHeader, folder string) {
				objectName := "products/test-uuid.jpg"
				url := "https://example.com/products/test-uuid.jpg"
				mockSvc.On("UploadImage", mock.Anything, mock.Anything, folder).Return(objectName, nil)
				mockSvc.On("GetImageURL", mock.Anything, objectName).Return(url, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "error - missing folder parameter",
			folder: "",
			mockSetup: func(mockSvc *MockService, file *multipart.FileHeader, folder string) {
				// No mock setup needed - route won't match
			},
			expectedStatus: http.StatusNotFound, // Route doesn't match without folder param
			expectedError:  true,
		},
		{
			name:   "error - upload service returns error",
			folder: "products",
			mockSetup: func(mockSvc *MockService, file *multipart.FileHeader, folder string) {
				mockSvc.On("UploadImage", mock.Anything, mock.Anything, folder).Return("", errors.New("upload failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/uploads/images/:folder", controller.UploadImageToFolder)

			// Create multipart form data
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			
			// Add file field
			fileWriter, err := writer.CreateFormFile("file", "test.jpg")
			require.NoError(t, err)
			fileWriter.Write([]byte("fake image content"))
			
			writer.Close()

			url := "/uploads/images/" + tt.folder
			if tt.folder == "" {
				url = "/uploads/images/"
			}

			req := httptest.NewRequest(http.MethodPost, url, &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create a mock file header for the mock setup
			if tt.folder != "" {
				mockFile := &multipart.FileHeader{
					Filename: "test.jpg",
					Size:     1024,
				}
				tt.mockSetup(mockSvc, mockFile, tt.folder)
			}

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && tt.folder != "" {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockService, string)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns image URL",
			path: "products/test-uuid.jpg",
			mockSetup: func(mockSvc *MockService, objectName string) {
				url := "https://example.com/products/test-uuid.jpg"
				// The controller extracts path from wildcard, so objectName will be "products/test-uuid.jpg"
				mockSvc.On("GetImageURL", mock.Anything, "products/test-uuid.jpg").Return(url, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - handles URL-encoded path",
			path: "users%2Favatars%2Ftest-uuid.jpg",
			mockSetup: func(mockSvc *MockService, objectName string) {
				url := "https://example.com/users/avatars/test-uuid.jpg"
				mockSvc.On("GetImageURL", mock.Anything, "users/avatars/test-uuid.jpg").Return(url, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "error - missing path parameter",
			path:           "",
			mockSetup:      func(mockSvc *MockService, objectName string) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - path is 'images' (reserved)",
			path: "images",
			mockSetup: func(mockSvc *MockService, objectName string) {
				// No mock setup - should return 404 before service call
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			path: "products/test-uuid.jpg",
			mockSetup: func(mockSvc *MockService, objectName string) {
				mockSvc.On("GetImageURL", mock.Anything, objectName).Return("", errors.New("failed to generate URL"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.path != "" && tt.path != "images" {
				// Decode the path for the mock
				decodedPath := strings.ReplaceAll(tt.path, "%2F", "/")
				tt.mockSetup(mockSvc, decodedPath)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			// Register route under /uploads to match actual app structure
			app.Get("/uploads/*", controller.GetImage)

			url := "/uploads/" + tt.path
			if tt.path == "" {
				url = "/uploads/"
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && tt.path != "" && tt.path != "images" {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "image_url")
				assert.Contains(t, responseBody, "object_name")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockService, string)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - deletes image",
			path: "products%2Ftest-uuid.jpg", // URL-encode the slash for Fiber route
			mockSetup: func(mockSvc *MockService, objectName string) {
				// The controller will decode the path, so use decoded path
				mockSvc.On("DeleteImage", mock.Anything, "products/test-uuid.jpg").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - handles URL-encoded path",
			path: "users%2Favatars%2Ftest-uuid.jpg",
			mockSetup: func(mockSvc *MockService, objectName string) {
				mockSvc.On("DeleteImage", mock.Anything, "users/avatars/test-uuid.jpg").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "error - missing path parameter",
			path:           "",
			mockSetup:      func(mockSvc *MockService, objectName string) {},
			expectedStatus: http.StatusNotFound, // Fiber returns 404 when route doesn't match
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			path: "products%2Ftest-uuid.jpg", // URL-encode the slash for Fiber route
			mockSetup: func(mockSvc *MockService, objectName string) {
				// The controller will decode the path, so use decoded path
				mockSvc.On("DeleteImage", mock.Anything, "products/test-uuid.jpg").Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.path != "" {
				// Decode the path for the mock - this is what the controller will receive
				decodedPath := strings.ReplaceAll(tt.path, "%2F", "/")
				tt.mockSetup(mockSvc, decodedPath)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			// Use the same route pattern as the controller Register method
			app.Delete("/images/:path", controller.DeleteImage)

			url := "/images/" + tt.path
			if tt.path == "" {
				// For empty path, route won't match - Fiber requires a path param
				url = "/images/"
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && tt.path != "" {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "message")
			}

			// Only assert expectations if path is not empty (service won't be called for empty path)
			// Also skip assertion if we expect an error that prevents service call
			if tt.path != "" && tt.expectedStatus != http.StatusNotFound {
				mockSvc.AssertExpectations(t)
			}
		})
	}
}

