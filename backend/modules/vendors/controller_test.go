package vendors

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (m *MockService) List(ctx context.Context) ([]*GetVendorDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetVendorDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetVendorDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVendorDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateVendorDTO) (*GetVendorDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVendorDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateVendorDTO) (*GetVendorDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVendorDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestController_ListVendors(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns list of vendors",
			mockSetup: func(mockSvc *MockService) {
				vendors := []*GetVendorDTO{
					{
						ID:        uuid.New(),
						Name:      stringPtr("Vendor 1"),
						Email:     stringPtr("vendor1@example.com"),
						IsActive:  "true",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Name:      stringPtr("Vendor 2"),
						Email:     stringPtr("vendor2@example.com"),
						IsActive:  "true",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(vendors, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetVendorDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetVendorDTO)(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/vendors", controller.ListVendors)

			req := httptest.NewRequest(http.MethodGet, "/vendors", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetVendor(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns vendor",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				vendor := &GetVendorDTO{
					ID:        id,
					Name:      stringPtr("Test Vendor"),
					Email:     stringPtr("test@example.com"),
					IsActive:  "true",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(vendor, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - vendor not found",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetVendorDTO)(nil), errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.id != "invalid-uuid" {
				id, _ := uuid.Parse(tt.id)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/vendors/:id", controller.GetVendor)

			req := httptest.NewRequest(http.MethodGet, "/vendors/"+tt.id, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateVendor(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateVendorDTO
		mockSetup      func(*MockService, CreateVendorDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - creates vendor",
			requestBody: CreateVendorDTO{
				ID:       uuid.New(),
				IsActive: "true",
			},
			mockSetup: func(mockSvc *MockService, dto CreateVendorDTO) {
				createdVendor := &GetVendorDTO{
					ID:        dto.ID,
					IsActive:  dto.IsActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.IsActive == dto.IsActive
				})).Return(createdVendor, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: CreateVendorDTO{
				ID:       uuid.New(),
				IsActive: "true",
			},
			mockSetup: func(mockSvc *MockService, dto CreateVendorDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.IsActive == dto.IsActive
				})).Return((*GetVendorDTO)(nil), errors.New("creation failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/vendors", controller.CreateVendor)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/vendors", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateVendor(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    UpdateVendorDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateVendorDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - updates vendor",
			id:   uuid.New().String(),
			requestBody: UpdateVendorDTO{
				Name:     stringPtr("Updated Vendor"),
				IsActive: stringPtr("false"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateVendorDTO) {
				updatedVendor := &GetVendorDTO{
					ID:        id,
					Name:      stringPtr("Updated Vendor"),
					IsActive:  "false",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateVendorDTO) bool {
					return actual.Name != nil && *actual.Name == "Updated Vendor" &&
						actual.IsActive != nil && *actual.IsActive == "false"
				})).Return(updatedVendor, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			requestBody: UpdateVendorDTO{
				Name: stringPtr("Updated Vendor"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateVendorDTO) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			requestBody: UpdateVendorDTO{
				Name: stringPtr("Updated Vendor"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateVendorDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateVendorDTO) bool {
					return actual.Name != nil && *actual.Name == "Updated Vendor"
				})).Return((*GetVendorDTO)(nil), errors.New("update failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.id != "invalid-uuid" {
				id, _ := uuid.Parse(tt.id)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/vendors/:id", controller.UpdateVendor)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/vendors/"+tt.id, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteVendor(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - deletes vendor",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.id != "invalid-uuid" {
				id, _ := uuid.Parse(tt.id)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/vendors/:id", controller.DeleteVendor)

			req := httptest.NewRequest(http.MethodDelete, "/vendors/"+tt.id, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
