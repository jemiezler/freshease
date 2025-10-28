package users

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

func (m *MockService) List(ctx context.Context) ([]*GetUserDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetUserDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetUserDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateUserDTO) (*GetUserDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*GetUserDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - returns users list",
			mockSetup: func(mockSvc *MockService) {
				expectedUsers := []*GetUserDTO{
					{
						ID:     uuid.New(),
						Email:  "user1@example.com",
						Name:   "User One",
						Status: "active",
					},
					{
						ID:     uuid.New(),
						Email:  "user2@example.com",
						Name:   "User Two",
						Status: "active",
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedUsers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*GetUserDTO{
				{
					ID:     uuid.New(),
					Email:  "user1@example.com",
					Name:   "User One",
					Status: "active",
				},
				{
					ID:     uuid.New(),
					Email:  "user2@example.com",
					Name:   "User Two",
					Status: "active",
				},
			},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetUserDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/users", controller.ListUsers)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.IsType(t, []interface{}{}, responseBody)
			} else {
				assert.IsType(t, map[string]interface{}{}, responseBody)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "success - returns user by ID",
			userID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				expectedUser := &GetUserDTO{
					ID:     id,
					Email:  "user@example.com",
					Name:   "Test User",
					Status: "active",
				}
				mockSvc.On("Get", mock.Anything, id).Return(expectedUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &GetUserDTO{
				ID:     uuid.New(),
				Email:  "user@example.com",
				Name:   "Test User",
				Status: "active",
			},
		},
		{
			name:           "error - invalid UUID",
			userID:         "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:   "error - user not found",
			userID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetUserDTO)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"message": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.userID != "invalid-uuid" {
				userID, err := uuid.Parse(tt.userID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, userID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/users/:id", controller.GetUser)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.IsType(t, map[string]interface{}{}, responseBody)
			} else {
				assert.IsType(t, map[string]interface{}{}, responseBody)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateUserDTO
		mockSetup      func(*MockService, CreateUserDTO)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - creates new user",
			requestBody: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			mockSetup: func(mockSvc *MockService, dto CreateUserDTO) {
				expectedUser := &GetUserDTO{
					ID:     dto.ID,
					Email:  dto.Email,
					Name:   dto.Name,
					Status: "active",
				}
				mockSvc.On("Create", mock.Anything, dto).Return(expectedUser, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &GetUserDTO{
				ID:     uuid.New(),
				Email:  "newuser@example.com",
				Name:   "New User",
				Status: "active",
			},
		},
		{
			name: "error - service returns error",
			requestBody: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			mockSetup: func(mockSvc *MockService, dto CreateUserDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return((*GetUserDTO)(nil), errors.New("email already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "email already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/users", controller.CreateUser)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.IsType(t, map[string]interface{}{}, responseBody)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    UpdateUserDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateUserDTO)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "success - updates user",
			userID: uuid.New().String(),
			requestBody: UpdateUserDTO{
				ID:    uuid.New(), // This will be overridden by the service
				Email: stringPtr("updated@example.com"),
				Name:  stringPtr("Updated User"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateUserDTO) {
				expectedUser := &GetUserDTO{
					ID:     id,
					Email:  *dto.Email,
					Name:   *dto.Name,
					Status: "active",
				}
				mockSvc.On("Update", mock.Anything, id, dto).Return(expectedUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &GetUserDTO{
				ID:     uuid.New(),
				Email:  "updated@example.com",
				Name:   "Updated User",
				Status: "active",
			},
		},
		{
			name:           "error - invalid UUID",
			userID:         "invalid-uuid",
			requestBody:    UpdateUserDTO{ID: uuid.New()},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateUserDTO) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:   "error - service returns error",
			userID: uuid.New().String(),
			requestBody: UpdateUserDTO{
				ID:    uuid.New(),
				Email: stringPtr("updated@example.com"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateUserDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return((*GetUserDTO)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "user not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.userID != "invalid-uuid" {
				userID, err := uuid.Parse(tt.userID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, userID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Put("/users/:id", controller.UpdateUser)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.IsType(t, map[string]interface{}{}, responseBody)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "success - deletes user",
			userID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "error - invalid UUID",
			userID:         "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:   "error - service returns error",
			userID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "user not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.userID != "invalid-uuid" {
				userID, err := uuid.Parse(tt.userID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, userID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/users/:id", controller.DeleteUser)

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusNoContent {
				assert.Equal(t, int64(0), resp.ContentLength)
			} else {
				var responseBody interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.IsType(t, map[string]interface{}{}, responseBody)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
