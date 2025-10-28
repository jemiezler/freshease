package carts

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

func (m *MockService) List(ctx context.Context) ([]*GetCartDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetCartDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetCartDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCartDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateCartDTO) (*GetCartDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCartDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateCartDTO) (*GetCartDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCartDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListCarts(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns carts list",
			mockSetup: func(mockSvc *MockService) {
				carts := []*GetCartDTO{
					{
						ID:        uuid.New(),
						Status:    "pending",
						Total:     100.50,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Status:    "completed",
						Total:     250.75,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(carts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Carts Retrieved Successfully",
			},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetCartDTO)(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/carts", controller.ListCarts)

			req := httptest.NewRequest(http.MethodGet, "/carts", nil)
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

func TestController_GetCart(t *testing.T) {
	tests := []struct {
		name           string
		cartID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "success - returns cart by ID",
			cartID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				cart := &GetCartDTO{
					ID:        id,
					Status:    "pending",
					Total:     150.25,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(cart, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Cart Retrieved Successfully",
			},
		},
		{
			name:   "error - invalid UUID",
			cartID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid uuid",
			},
		},
		{
			name:   "error - cart not found",
			cartID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetCartDTO)(nil), errors.New("cart not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"message": "not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.cartID != "invalid-uuid" {
				cartID, err := uuid.Parse(tt.cartID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, cartID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/carts/:id", controller.GetCart)

			req := httptest.NewRequest(http.MethodGet, "/carts/"+tt.cartID, nil)
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

func TestController_CreateCart(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateCartDTO
		mockSetup      func(*MockService, CreateCartDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates new cart",
			requestBody: CreateCartDTO{
				Status: stringPtr("pending"),
				Total:  float64Ptr(99.99),
			},
			mockSetup: func(mockSvc *MockService, dto CreateCartDTO) {
				expectedCart := &GetCartDTO{
					ID:        uuid.New(),
					Status:    *dto.Status,
					Total:     *dto.Total,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateCartDTO) bool {
					return actual.Status != nil && *actual.Status == *dto.Status &&
						actual.Total != nil && *actual.Total == *dto.Total
				})).Return(expectedCart, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Cart Created Successfully",
			},
		},
		{
			name: "error - service returns error",
			requestBody: CreateCartDTO{
				Status: stringPtr("pending"),
				Total:  float64Ptr(50.00),
			},
			mockSetup: func(mockSvc *MockService, dto CreateCartDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateCartDTO) bool {
					return actual.Status != nil && *actual.Status == *dto.Status &&
						actual.Total != nil && *actual.Total == *dto.Total
				})).Return((*GetCartDTO)(nil), errors.New("creation failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "creation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/carts", controller.CreateCart)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/carts", bytes.NewBuffer(jsonBody))
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

func TestController_UpdateCart(t *testing.T) {
	tests := []struct {
		name           string
		cartID         string
		requestBody    UpdateCartDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateCartDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "success - updates cart",
			cartID: uuid.New().String(),
			requestBody: UpdateCartDTO{
				Status: stringPtr("completed"),
				Total:  float64Ptr(200.00),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCartDTO) {
				expectedCart := &GetCartDTO{
					ID:        id,
					Status:    *dto.Status,
					Total:     *dto.Total,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateCartDTO) bool {
					return actual.Status != nil && *actual.Status == *dto.Status &&
						actual.Total != nil && *actual.Total == *dto.Total
				})).Return(expectedCart, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Cart Updated Successfully",
			},
		},
		{
			name:   "error - invalid UUID",
			cartID: "invalid-uuid",
			requestBody: UpdateCartDTO{
				Status: stringPtr("completed"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCartDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid uuid",
			},
		},
		{
			name:   "error - service returns error",
			cartID: uuid.New().String(),
			requestBody: UpdateCartDTO{
				Status: stringPtr("cancelled"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCartDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateCartDTO) bool {
					return actual.Status != nil && *actual.Status == *dto.Status
				})).Return((*GetCartDTO)(nil), errors.New("update failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "update failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.cartID != "invalid-uuid" {
				cartID, err := uuid.Parse(tt.cartID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, cartID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/carts/:id", controller.UpdateCart)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/carts/"+tt.cartID, bytes.NewBuffer(jsonBody))
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

func TestController_DeleteCart(t *testing.T) {
	tests := []struct {
		name           string
		cartID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "success - deletes cart",
			cartID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody: map[string]interface{}{
				"message": "Cart Deleted Successfully",
			},
		},
		{
			name:   "error - invalid UUID",
			cartID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid uuid",
			},
		},
		{
			name:   "error - service returns error",
			cartID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "delete failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.cartID != "invalid-uuid" {
				cartID, err := uuid.Parse(tt.cartID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, cartID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/carts/:id", controller.DeleteCart)

			req := httptest.NewRequest(http.MethodDelete, "/carts/"+tt.cartID, nil)
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
