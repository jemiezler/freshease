package roles

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

func (m *MockService) List(ctx context.Context) ([]*GetRoleDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetRoleDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetRoleDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRoleDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateRoleDTO) (*GetRoleDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRoleDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateRoleDTO) (*GetRoleDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRoleDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func TestController_ListRoles(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns list of roles",
			mockSetup: func(mockSvc *MockService) {
				roles := []*GetRoleDTO{
					{
						ID:          uuid.New(),
						Name:        "admin",
						Description: "Administrator role",
					},
					{
						ID:          uuid.New(),
						Name:        "user",
						Description: "Regular user role",
					},
				}
				mockSvc.On("List", mock.Anything).Return(roles, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetRoleDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetRoleDTO)(nil), errors.New("service error"))
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
			app.Get("/roles", controller.ListRoles)

			req := httptest.NewRequest(http.MethodGet, "/roles", nil)
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

func TestController_GetRole(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns role",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				role := &GetRoleDTO{
					ID:          id,
					Name:        "admin",
					Description: "Administrator role",
				}
				mockSvc.On("Get", mock.Anything, id).Return(role, nil)
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
			name: "error - role not found",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetRoleDTO)(nil), errors.New("not found"))
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
			app.Get("/roles/:id", controller.GetRole)

			req := httptest.NewRequest(http.MethodGet, "/roles/"+tt.id, nil)
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

func TestController_CreateRole(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateRoleDTO
		mockSetup      func(*MockService, CreateRoleDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - creates role",
			requestBody: CreateRoleDTO{
				ID:          uuid.New(),
				Name:        "admin",
				Description: "Administrator role",
			},
			mockSetup: func(mockSvc *MockService, dto CreateRoleDTO) {
				createdRole := &GetRoleDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Description: dto.Description,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateRoleDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description
				})).Return(createdRole, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: CreateRoleDTO{
				ID:          uuid.New(),
				Name:        "admin",
				Description: "Administrator role",
			},
			mockSetup: func(mockSvc *MockService, dto CreateRoleDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateRoleDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description
				})).Return((*GetRoleDTO)(nil), errors.New("creation failed"))
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
			app.Post("/roles", controller.CreateRole)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(jsonBody))
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

func TestController_UpdateRole(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    UpdateRoleDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateRoleDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - updates role",
			id:   uuid.New().String(),
			requestBody: UpdateRoleDTO{
				Name:        stringPtr("updated_admin"),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRoleDTO) {
				updatedRole := &GetRoleDTO{
					ID:          id,
					Name:        "updated_admin",
					Description: "Updated description",
				}
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateRoleDTO) bool {
					return actual.Name != nil && *actual.Name == "updated_admin" &&
						actual.Description != nil && *actual.Description == "Updated description"
				})).Return(updatedRole, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			requestBody: UpdateRoleDTO{
				Name: stringPtr("updated_admin"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRoleDTO) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			requestBody: UpdateRoleDTO{
				Name: stringPtr("updated_admin"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRoleDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateRoleDTO) bool {
					return actual.Name != nil && *actual.Name == "updated_admin"
				})).Return((*GetRoleDTO)(nil), errors.New("update failed"))
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
			app.Patch("/roles/:id", controller.UpdateRole)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/roles/"+tt.id, bytes.NewBuffer(jsonBody))
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

func TestController_DeleteRole(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - deletes role",
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
			app.Delete("/roles/:id", controller.DeleteRole)

			req := httptest.NewRequest(http.MethodDelete, "/roles/"+tt.id, nil)
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
