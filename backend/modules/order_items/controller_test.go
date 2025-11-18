package order_items

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

func (m *MockService) List(ctx context.Context) ([]*GetOrder_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetOrder_itemDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetOrder_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrder_itemDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateOrder_itemDTO) (*GetOrder_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrder_itemDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateOrder_itemDTO) (*GetOrder_itemDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrder_itemDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListOrder_items(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns order items list",
			mockSetup: func(mockSvc *MockService) {
				expectedItems := []*GetOrder_itemDTO{
					{
						ID:        uuid.New(),
						Qty:       2,
						UnitPrice: 10.99,
						LineTotal: 21.98,
						OrderID:   uuid.New(),
						ProductID: uuid.New(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetOrder_itemDTO)(nil), errors.New("database error"))
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
			app.Get("/order-items", controller.ListOrder_items)

			req := httptest.NewRequest(http.MethodGet, "/order-items", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetOrder_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - returns order item by ID",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				item := &GetOrder_itemDTO{
					ID:        id,
					Qty:       2,
					UnitPrice: 10.99,
					LineTotal: 21.98,
					OrderID:   uuid.New(),
					ProductID: uuid.New(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(item, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - order item not found",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetOrder_itemDTO)(nil), errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.itemID != "invalid-uuid" {
				itemID, err := uuid.Parse(tt.itemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, itemID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/order-items/:id", controller.GetOrder_item)

			req := httptest.NewRequest(http.MethodGet, "/order-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateOrder_item(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateOrder_itemDTO
		mockSetup      func(*MockService, CreateOrder_itemDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new order item",
			requestBody: CreateOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 10.99,
				LineTotal: 21.98,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateOrder_itemDTO) {
				expectedItem := &GetOrder_itemDTO{
					ID:        dto.ID,
					Qty:       dto.Qty,
					UnitPrice: dto.UnitPrice,
					LineTotal: dto.LineTotal,
					OrderID:   dto.OrderID,
					ProductID: dto.ProductID,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateOrder_itemDTO) bool {
					return actual.ID == dto.ID && actual.Qty == dto.Qty
				})).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 10.99,
				LineTotal: 21.98,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateOrder_itemDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateOrder_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetOrder_itemDTO)(nil), errors.New("creation failed"))
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
			app.Post("/order-items", controller.CreateOrder_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/order-items", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateOrder_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    UpdateOrder_itemDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateOrder_itemDTO)
		expectedStatus int
	}{
		{
			name:   "success - updates order item",
			itemID: uuid.New().String(),
			requestBody: UpdateOrder_itemDTO{
				ID:        uuid.New(), // This will be overwritten by service
				Qty:       intPtr(5),
				UnitPrice: float64Ptr(12.99),
				LineTotal: float64Ptr(64.95),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateOrder_itemDTO) {
				expectedItem := &GetOrder_itemDTO{
					ID:        id,
					Qty:       5,
					UnitPrice: 12.99,
					LineTotal: 64.95,
					OrderID:   uuid.New(),
					ProductID: uuid.New(),
				}
				// Service sets dto.ID = id, so use mock.Anything for DTO
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			requestBody:    UpdateOrder_itemDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateOrder_itemDTO) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			requestBody: UpdateOrder_itemDTO{
				ID:  uuid.New(), // This will be overwritten by service
				Qty: intPtr(5),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateOrder_itemDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return((*GetOrder_itemDTO)(nil), errors.New("update failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.itemID != "invalid-uuid" {
				itemID, err := uuid.Parse(tt.itemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, itemID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/order-items/:id", controller.UpdateOrder_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/order-items/"+tt.itemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteOrder_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - deletes order item",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.itemID != "invalid-uuid" {
				itemID, err := uuid.Parse(tt.itemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, itemID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/order-items/:id", controller.DeleteOrder_item)

			req := httptest.NewRequest(http.MethodDelete, "/order-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

