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

func TestController_CreateOrder(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    CreateOrderDTO
		mockSetup      func(*MockService, CreateOrderDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new order",
			requestBody: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-001",
				Status:      "pending",
				Subtotal:    200.00,
				ShippingFee: 15.00,
				Discount:    10.00,
				Total:       205.00,
				UserID:      userID,
				PlacedAt:    &now,
			},
			mockSetup: func(mockSvc *MockService, dto CreateOrderDTO) {
				expectedOrder := &GetOrderDTO{
					ID:          dto.ID,
					OrderNo:     dto.OrderNo,
					Status:      dto.Status,
					Subtotal:    dto.Subtotal,
					ShippingFee: dto.ShippingFee,
					Discount:    dto.Discount,
					Total:       dto.Total,
					UserID:      dto.UserID,
					PlacedAt:   dto.PlacedAt,
					UpdatedAt:   time.Now(),
				}
				mockSvc.On("Create", mock.Anything, dto).Return(expectedOrder, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-002",
				Status:      "pending",
				Subtotal:    100.00,
				ShippingFee: 10.00,
				Discount:    0.00,
				Total:       110.00,
				UserID:      userID,
			},
			mockSetup: func(mockSvc *MockService, dto CreateOrderDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return((*GetOrderDTO)(nil), errors.New("creation failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/orders", controller.CreateOrder)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetOrder(t *testing.T) {
	tests := []struct {
		name           string
		orderID        string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:    "success - returns order by ID",
			orderID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				order := &GetOrderDTO{
					ID:          id,
					OrderNo:     "ORD-003",
					Status:      "pending",
					Subtotal:    150.00,
					ShippingFee: 12.00,
					Discount:    8.00,
					Total:       154.00,
					UserID:      uuid.New(),
					UpdatedAt:   time.Now(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(order, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid UUID",
			orderID:        "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "error - order not found",
			orderID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetOrderDTO)(nil), errors.New("order not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.orderID != "invalid-uuid" {
				orderID, err := uuid.Parse(tt.orderID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, orderID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/orders/:id", controller.GetOrder)

			req := httptest.NewRequest(http.MethodGet, "/orders/"+tt.orderID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_ListOrders(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns orders list",
			mockSetup: func(mockSvc *MockService) {
				orders := []*GetOrderDTO{
					{
						ID:          uuid.New(),
						OrderNo:     "ORD-001",
						Status:      "pending",
						Subtotal:    100.00,
						ShippingFee: 10.00,
						Discount:    5.00,
						Total:       105.00,
						UserID:      uuid.New(),
						UpdatedAt:   time.Now(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(orders, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetOrderDTO)(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
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

			mockSvc.AssertExpectations(t)
		})
	}
}
