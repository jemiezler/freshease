package payments

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

func (m *MockService) List(ctx context.Context) ([]*GetPaymentDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetPaymentDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetPaymentDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPaymentDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreatePaymentDTO) (*GetPaymentDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPaymentDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdatePaymentDTO) (*GetPaymentDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPaymentDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListPayments(t *testing.T) {
	tests := []struct {
		name            string
		mockSetup       func(*MockService)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "success - returns payments list",
			mockSetup: func(mockSvc *MockService) {
				orderID := uuid.New()
				providerRef := "pay_001"
				expectedPayments := []*GetPaymentDTO{
					{
						ID:          uuid.New(),
						Provider:    "stripe",
						ProviderRef: &providerRef,
						Status:      "pending",
						Amount:      110.0,
						OrderID:     orderID,
					},
					{
						ID:       uuid.New(),
						Provider: "paypal",
						Status:   "completed",
						Amount:   110.0,
						OrderID:  orderID,
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedPayments, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "Payments Retrieved Successfully",
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetPaymentDTO(nil), errors.New("database error"))
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
			app.Get("/payments", controller.ListPayments)

			req := httptest.NewRequest(http.MethodGet, "/payments", nil)
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

func TestController_GetPayment(t *testing.T) {
	paymentID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name            string
		paymentID       string
		mockSetup       func(*MockService, uuid.UUID)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:      "success - returns payment by ID",
			paymentID: paymentID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				providerRef := "pay_001"
				mockSvc.On("Get", mock.Anything, id).Return(&GetPaymentDTO{
					ID:          id,
					Provider:    "stripe",
					ProviderRef: &providerRef,
					Status:      "pending",
					Amount:      110.0,
					OrderID:     orderID,
				}, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "Payment Retrieved Successfully",
		},
		{
			name:      "error - invalid UUID",
			paymentID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:      "error - payment not found",
			paymentID: paymentID.String(),
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
			if tt.paymentID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.paymentID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/payments/:id", controller.GetPayment)

			req := httptest.NewRequest(http.MethodGet, "/payments/"+tt.paymentID, nil)
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

func TestController_CreatePayment(t *testing.T) {
	paymentID := uuid.New()
	orderID := uuid.New()
	providerRef := "pay_002"
	paidAt := time.Now()

	tests := []struct {
		name           string
		requestBody    CreatePaymentDTO
		mockSetup      func(*MockService, CreatePaymentDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates payment",
			requestBody: CreatePaymentDTO{
				ID:          paymentID,
				Provider:    "stripe",
				ProviderRef: &providerRef,
				Status:      "completed",
				Amount:      110.0,
				PaidAt:      &paidAt,
				OrderID:     orderID,
			},
			mockSetup: func(mockSvc *MockService, dto CreatePaymentDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(&GetPaymentDTO{
					ID:          dto.ID,
					Provider:    dto.Provider,
					ProviderRef: dto.ProviderRef,
					Status:      dto.Status,
					Amount:      dto.Amount,
					PaidAt:      dto.PaidAt,
					OrderID:     dto.OrderID,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Payment Created Successfully",
		},
		{
			name: "error - service returns error",
			requestBody: CreatePaymentDTO{
				ID:       paymentID,
				Provider: "stripe",
				Status:   "pending",
				Amount:   110.0,
				OrderID:  orderID,
			},
			mockSetup: func(mockSvc *MockService, dto CreatePaymentDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(nil, errors.New("validation error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/payments", controller.CreatePayment)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
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

func TestController_UpdatePayment(t *testing.T) {
	paymentID := uuid.New()
	newStatus := "completed"

	tests := []struct {
		name           string
		paymentID      string
		requestBody    UpdatePaymentDTO
		mockSetup      func(*MockService, uuid.UUID, UpdatePaymentDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "success - updates payment",
			paymentID: paymentID.String(),
			requestBody: UpdatePaymentDTO{
				ID:     paymentID,
				Status: &newStatus,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdatePaymentDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return(&GetPaymentDTO{
					ID:       id,
					Status:   newStatus,
					Provider: "stripe",
					Amount:   110.0,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Payment Updated Successfully",
		},
		{
			name:      "error - invalid UUID",
			paymentID: "invalid-uuid",
			requestBody: UpdatePaymentDTO{
				Status: &newStatus,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdatePaymentDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.paymentID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.paymentID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/payments/:id", controller.UpdatePayment)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/payments/"+tt.paymentID, bytes.NewBuffer(body))
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

func TestController_DeletePayment(t *testing.T) {
	paymentID := uuid.New()

	tests := []struct {
		name           string
		paymentID      string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "success - deletes payment",
			paymentID: paymentID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedMessage: "Payment Deleted Successfully",
		},
		{
			name:      "error - invalid UUID",
			paymentID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:      "error - service returns error",
			paymentID: paymentID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.paymentID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.paymentID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/payments/:id", controller.DeletePayment)

			req := httptest.NewRequest(http.MethodDelete, "/payments/"+tt.paymentID, nil)
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

