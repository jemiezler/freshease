package permissions

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

func (m *MockService) List(ctx context.Context) ([]*GetPermissionDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetPermissionDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetPermissionDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPermissionDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreatePermissionDTO) (*GetPermissionDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPermissionDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdatePermissionDTO) (*GetPermissionDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPermissionDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListPermissions(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns list of permissions",
			mockSetup: func(mockSvc *MockService) {
				permissions := []*GetPermissionDTO{
					{
						ID:          uuid.New(),
						Name:        "read_users",
						Description: "Permission to read user data",
					},
					{
						ID:          uuid.New(),
						Name:        "write_users",
						Description: "Permission to write user data",
					},
				}
				mockSvc.On("List", mock.Anything).Return(permissions, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetPermissionDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetPermissionDTO)(nil), errors.New("service error"))
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
			app.Get("/permissions", controller.ListPermissions)

			req := httptest.NewRequest(http.MethodGet, "/permissions", nil)
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

func TestController_GetPermission(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns permission",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				permission := &GetPermissionDTO{
					ID:          id,
					Name:        "read_users",
					Description: "Permission to read user data",
				}
				mockSvc.On("Get", mock.Anything, id).Return(permission, nil)
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
			name: "error - permission not found",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetPermissionDTO)(nil), errors.New("not found"))
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
			app.Get("/permissions/:id", controller.GetPermission)

			req := httptest.NewRequest(http.MethodGet, "/permissions/"+tt.id, nil)
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

func TestController_CreatePermission(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreatePermissionDTO
		mockSetup      func(*MockService, CreatePermissionDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - creates permission",
			requestBody: CreatePermissionDTO{
				ID:          uuid.New(),
				Name:        "read_users",
				Description: "Permission to read user data",
			},
			mockSetup: func(mockSvc *MockService, dto CreatePermissionDTO) {
				createdPermission := &GetPermissionDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Description: dto.Description,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreatePermissionDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description
				})).Return(createdPermission, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: CreatePermissionDTO{
				ID:          uuid.New(),
				Name:        "read_users",
				Description: "Permission to read user data",
			},
			mockSetup: func(mockSvc *MockService, dto CreatePermissionDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreatePermissionDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description
				})).Return((*GetPermissionDTO)(nil), errors.New("creation failed"))
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
			app.Post("/permissions", controller.CreatePermission)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(jsonBody))
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

func TestController_UpdatePermission(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    UpdatePermissionDTO
		mockSetup      func(*MockService, uuid.UUID, UpdatePermissionDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - updates permission",
			id:   uuid.New().String(),
			requestBody: UpdatePermissionDTO{
				Name:        stringPtr("updated_name"),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdatePermissionDTO) {
				updatedPermission := &GetPermissionDTO{
					ID:          id,
					Name:        "updated_name",
					Description: "Updated description",
				}
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdatePermissionDTO) bool {
					return actual.Name != nil && *actual.Name == "updated_name" &&
						actual.Description != nil && *actual.Description == "Updated description"
				})).Return(updatedPermission, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			requestBody: UpdatePermissionDTO{
				Name: stringPtr("updated_name"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdatePermissionDTO) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			requestBody: UpdatePermissionDTO{
				Name: stringPtr("updated_name"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdatePermissionDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdatePermissionDTO) bool {
					return actual.Name != nil && *actual.Name == "updated_name"
				})).Return((*GetPermissionDTO)(nil), errors.New("update failed"))
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
			app.Patch("/permissions/:id", controller.UpdatePermission)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/permissions/"+tt.id, bytes.NewBuffer(jsonBody))
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

func TestController_DeletePermission(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - deletes permission",
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
			app.Delete("/permissions/:id", controller.DeletePermission)

			req := httptest.NewRequest(http.MethodDelete, "/permissions/"+tt.id, nil)
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
