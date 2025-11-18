package bundles

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) List(ctx context.Context) ([]*GetBundleDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetBundleDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetBundleDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundleDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateBundleDTO) (*GetBundleDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundleDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateBundleDTO) (*GetBundleDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundleDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListBundles(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns bundles list",
			mockSetup: func(mockSvc *MockService) {
				expectedBundles := []*GetBundleDTO{
					{
						ID:       uuid.New(),
						Name:     "Bundle One",
						Price:    99.99,
						IsActive: true,
					},
					{
						ID:       uuid.New(),
						Name:     "Bundle Two",
						Price:    149.99,
						IsActive: true,
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedBundles, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"message": "Bundles Retrieved Successfully"},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetBundleDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/bundles", controller.ListBundles)

			req := httptest.NewRequest(http.MethodGet, "/bundles", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetBundle(t *testing.T) {
	bundleID := uuid.New()

	tests := []struct {
		name           string
		bundleID       string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - returns bundle by ID",
			bundleID: bundleID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(&GetBundleDTO{
					ID:       id,
					Name:     "Test Bundle",
					Price:    99.99,
					IsActive: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"message": "Bundle Retrieved Successfully"},
		},
		{
			name:     "error - invalid UUID",
			bundleID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "invalid uuid"},
		},
		{
			name:     "error - bundle not found",
			bundleID: bundleID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"message": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.bundleID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.bundleID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/bundles/:id", controller.GetBundle)

			req := httptest.NewRequest(http.MethodGet, "/bundles/"+tt.bundleID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateBundle(t *testing.T) {
	bundleID := uuid.New()

	tests := []struct {
		name           string
		requestBody    CreateBundleDTO
		mockSetup      func(*MockService, CreateBundleDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates bundle",
			requestBody: CreateBundleDTO{
				ID:       bundleID,
				Name:     "New Bundle",
				Price:    199.99,
				IsActive: true,
			},
			mockSetup: func(mockSvc *MockService, dto CreateBundleDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(&GetBundleDTO{
					ID:       dto.ID,
					Name:     dto.Name,
					Price:    dto.Price,
					IsActive: dto.IsActive,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]interface{}{"message": "Bundle Created Successfully"},
		},
		{
			name: "error - service returns error",
			requestBody: CreateBundleDTO{
				ID:       bundleID,
				Name:     "New Bundle",
				Price:    199.99,
				IsActive: true,
			},
			mockSetup: func(mockSvc *MockService, dto CreateBundleDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(nil, errors.New("validation error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "validation error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/bundles", controller.CreateBundle)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/bundles", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateBundle(t *testing.T) {
	bundleID := uuid.New()
	newName := "Updated Bundle"

	tests := []struct {
		name           string
		bundleID       string
		requestBody    UpdateBundleDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateBundleDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - updates bundle",
			bundleID: bundleID.String(),
			requestBody: UpdateBundleDTO{
				ID:   bundleID,
				Name: &newName,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateBundleDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return(&GetBundleDTO{
					ID:       id,
					Name:     newName,
					Price:    99.99,
					IsActive: true,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]interface{}{"message": "Bundle Updated Successfully"},
		},
		{
			name:     "error - invalid UUID",
			bundleID: "invalid-uuid",
			requestBody: UpdateBundleDTO{
				Name: &newName,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateBundleDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "invalid uuid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.bundleID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.bundleID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/bundles/:id", controller.UpdateBundle)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/bundles/"+tt.bundleID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteBundle(t *testing.T) {
	bundleID := uuid.New()

	tests := []struct {
		name           string
		bundleID       string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - deletes bundle",
			bundleID: bundleID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   map[string]interface{}{"message": "Bundle Deleted Successfully"},
		},
		{
			name:     "error - invalid UUID",
			bundleID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "invalid uuid"},
		},
		{
			name:     "error - service returns error",
			bundleID: bundleID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.bundleID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.bundleID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/bundles/:id", controller.DeleteBundle)

			req := httptest.NewRequest(http.MethodDelete, "/bundles/"+tt.bundleID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			mockSvc.AssertExpectations(t)
		})
	}
}

