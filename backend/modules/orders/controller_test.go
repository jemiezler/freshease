package orders

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

func (m *MockService) List(ctx context.Context) ([]*GetOrderDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetOrderDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetOrderDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateOrderDTO) (*GetOrderDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateOrderDTO) (*GetOrderDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListOrders(t *testing.T) {
	tests := []struct {
		name            string
		mockSetup       func(*MockService)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "success - returns orders list",
			mockSetup: func(mockSvc *MockService) {
				userID := uuid.New()
				expectedOrders := []*GetOrderDTO{
					{
						ID:      uuid.New(),
						OrderNo: "ORD-001",
						Status:  "pending",
						Total:   110.0,
						UserID:  userID,
					},
					{
						ID:      uuid.New(),
						OrderNo: "ORD-002",
						Status:  "completed",
						Total:   210.0,
						UserID:  userID,
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedOrders, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "Orders Retrieved Successfully",
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetOrderDTO(nil), errors.New("database error"))
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/orders", controller.ListOrders)

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetOrder(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name            string
		orderID         string
		mockSetup       func(*MockService, uuid.UUID)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:    "success - returns order by ID",
			orderID: orderID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(&GetOrderDTO{
					ID:      id,
					OrderNo: "ORD-001",
					Status:  "pending",
					Total:   110.0,
					UserID:  userID,
				}, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "Order Retrieved Successfully",
		},
		{
			name:    "error - invalid UUID",
			orderID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:    "error - order not found",
			orderID: orderID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(nil, errors.New("not found"))
			},
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.orderID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.orderID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/orders/:id", controller.GetOrder)

			req := httptest.NewRequest(http.MethodGet, "/orders/"+tt.orderID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateOrder(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name            string
		requestBody     CreateOrderDTO
		mockSetup       func(*MockService, CreateOrderDTO)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "success - creates order",
			requestBody: CreateOrderDTO{
				ID:          orderID,
				OrderNo:     "ORD-003",
				Status:      "pending",
				Subtotal:    150.0,
				ShippingFee: 15.0,
				Discount:    5.0,
				Total:       160.0,
				PlacedAt:    &now,
				UserID:      userID,
			},
			mockSetup: func(mockSvc *MockService, dto CreateOrderDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateOrderDTO) bool {
					return actual.ID == dto.ID && actual.OrderNo == dto.OrderNo
				})).Return(&GetOrderDTO{
					ID:      dto.ID,
					OrderNo: dto.OrderNo,
					Status:  dto.Status,
					Total:   dto.Total,
					UserID:  dto.UserID,
				}, nil)
			},
			expectedStatus:  http.StatusCreated,
			expectedMessage: "Order Created Successfully",
		},
		{
			name: "error - service returns error",
			requestBody: CreateOrderDTO{
				ID:          orderID,
				OrderNo:     "ORD-004",
				Status:      "pending",
				Subtotal:    150.0,
				ShippingFee: 15.0,
				Discount:    5.0,
				Total:       160.0,
				UserID:      userID,
			},
			mockSetup: func(mockSvc *MockService, dto CreateOrderDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateOrderDTO) bool {
					return actual.ID == dto.ID && actual.OrderNo == dto.OrderNo
				})).Return(nil, errors.New("validation error"))
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/orders", controller.CreateOrder)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateOrder(t *testing.T) {
	orderID := uuid.New()
	newStatus := "completed"

	tests := []struct {
		name            string
		orderID         string
		requestBody     UpdateOrderDTO
		mockSetup       func(*MockService, uuid.UUID, UpdateOrderDTO)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:    "success - updates order",
			orderID: orderID.String(),
			requestBody: UpdateOrderDTO{
				ID:     orderID,
				Status: &newStatus,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateOrderDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return(&GetOrderDTO{
					ID:     id,
					Status: newStatus,
					Total:  110.0,
				}, nil)
			},
			expectedStatus:  http.StatusCreated,
			expectedMessage: "Order Updated Successfully",
		},
		{
			name:    "error - invalid UUID",
			orderID: "invalid-uuid",
			requestBody: UpdateOrderDTO{
				Status: &newStatus,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateOrderDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.orderID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.orderID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/orders/:id", controller.UpdateOrder)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/orders/"+tt.orderID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteOrder(t *testing.T) {
	orderID := uuid.New()

	tests := []struct {
		name            string
		orderID         string
		mockSetup       func(*MockService, uuid.UUID)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:    "success - deletes order",
			orderID: orderID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus:  http.StatusAccepted,
			expectedMessage: "Order Deleted Successfully",
		},
		{
			name:    "error - invalid UUID",
			orderID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:    "error - service returns error",
			orderID: orderID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("database error"))
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.orderID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.orderID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/orders/:id", controller.DeleteOrder)

			req := httptest.NewRequest(http.MethodDelete, "/orders/"+tt.orderID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			mockSvc.AssertExpectations(t)
		})
	}
}
