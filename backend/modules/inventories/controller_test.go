package inventories

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

func (m *MockService) List(ctx context.Context) ([]*GetInventoryDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetInventoryDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetInventoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetInventoryDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateInventoryDTO) (*GetInventoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetInventoryDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateInventoryDTO) (*GetInventoryDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetInventoryDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}

func TestController_ListInventories(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns list of inventories",
			mockSetup: func(mockSvc *MockService) {
				inventories := []*GetInventoryDTO{
					{
						ID:            uuid.New(),
						Quantity:      100,
						RestockAmount: 50,
						UpdatedAt:     time.Now(),
					},
					{
						ID:            uuid.New(),
						Quantity:      200,
						RestockAmount: 75,
						UpdatedAt:     time.Now(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(inventories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetInventoryDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetInventoryDTO)(nil), errors.New("service error"))
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
			app.Get("/inventories", controller.ListInventories)

			req := httptest.NewRequest(http.MethodGet, "/inventories", nil)
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

func TestController_GetInventory(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns inventory",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				inventory := &GetInventoryDTO{
					ID:            id,
					Quantity:      150,
					RestockAmount: 60,
					UpdatedAt:     time.Now(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(inventory, nil)
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
			name: "error - inventory not found",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetInventoryDTO)(nil), errors.New("not found"))
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
			app.Get("/inventories/:id", controller.GetInventory)

			req := httptest.NewRequest(http.MethodGet, "/inventories/"+tt.id, nil)
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

func TestController_CreateInventory(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateInventoryDTO
		mockSetup      func(*MockService, CreateInventoryDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - creates inventory",
			requestBody: CreateInventoryDTO{
				Quantity:      100,
				RestockAmount: 50,
				UpdatedAt:     time.Now(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateInventoryDTO) {
				createdInventory := &GetInventoryDTO{
					ID:            uuid.New(),
					Quantity:      dto.Quantity,
					RestockAmount: dto.RestockAmount,
					UpdatedAt:     time.Now(),
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateInventoryDTO) bool {
					return actual.Quantity == dto.Quantity &&
						actual.RestockAmount == dto.RestockAmount
				})).Return(createdInventory, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: CreateInventoryDTO{
				Quantity:      100,
				RestockAmount: 50,
				UpdatedAt:     time.Now(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateInventoryDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateInventoryDTO) bool {
					return actual.Quantity == dto.Quantity &&
						actual.RestockAmount == dto.RestockAmount
				})).Return((*GetInventoryDTO)(nil), errors.New("creation failed"))
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
			app.Post("/inventories", controller.CreateInventory)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/inventories", bytes.NewBuffer(jsonBody))
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

func TestController_UpdateInventory(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    UpdateInventoryDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateInventoryDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - updates inventory",
			id:   uuid.New().String(),
			requestBody: UpdateInventoryDTO{
				Quantity:      intPtr(200),
				RestockAmount: intPtr(75),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateInventoryDTO) {
				updatedInventory := &GetInventoryDTO{
					ID:            id,
					Quantity:      200,
					RestockAmount: 75,
					UpdatedAt:     time.Now(),
				}
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateInventoryDTO) bool {
					return actual.Quantity != nil && *actual.Quantity == 200 &&
						actual.RestockAmount != nil && *actual.RestockAmount == 75
				})).Return(updatedInventory, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			requestBody: UpdateInventoryDTO{
				Quantity: intPtr(200),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateInventoryDTO) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			requestBody: UpdateInventoryDTO{
				Quantity: intPtr(200),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateInventoryDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateInventoryDTO) bool {
					return actual.Quantity != nil && *actual.Quantity == 200
				})).Return((*GetInventoryDTO)(nil), errors.New("update failed"))
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
			app.Patch("/inventories/:id", controller.UpdateInventory)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/inventories/"+tt.id, bytes.NewBuffer(jsonBody))
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

func TestController_DeleteInventory(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - deletes inventory",
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
			app.Delete("/inventories/:id", controller.DeleteInventory)

			req := httptest.NewRequest(http.MethodDelete, "/inventories/"+tt.id, nil)
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
