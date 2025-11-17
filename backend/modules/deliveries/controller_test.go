package deliveries

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

func (m *MockService) List(ctx context.Context) ([]*GetDeliveryDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetDeliveryDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetDeliveryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDeliveryDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateDeliveryDTO) (*GetDeliveryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDeliveryDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateDeliveryDTO) (*GetDeliveryDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDeliveryDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListDeliveries(t *testing.T) {
	tests := []struct {
		name            string
		mockSetup       func(*MockService)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "success - returns deliveries list",
			mockSetup: func(mockSvc *MockService) {
				orderID := uuid.New()
				trackingNo := "TRACK-001"
				expectedDeliveries := []*GetDeliveryDTO{
					{
						ID:         uuid.New(),
						Provider:   "fedex",
						TrackingNo: &trackingNo,
						Status:     "pending",
						OrderID:    orderID,
					},
					{
						ID:       uuid.New(),
						Provider: "ups",
						Status:   "in_transit",
						OrderID:  orderID,
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedDeliveries, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "Deliveries Retrieved Successfully",
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetDeliveryDTO(nil), errors.New("database error"))
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
			app.Get("/deliveries", controller.ListDeliveries)

			req := httptest.NewRequest(http.MethodGet, "/deliveries", nil)
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

func TestController_GetDelivery(t *testing.T) {
	deliveryID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name            string
		deliveryID      string
		mockSetup       func(*MockService, uuid.UUID)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:       "success - returns delivery by ID",
			deliveryID: deliveryID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				trackingNo := "TRACK-001"
				mockSvc.On("Get", mock.Anything, id).Return(&GetDeliveryDTO{
					ID:         id,
					Provider:   "fedex",
					TrackingNo: &trackingNo,
					Status:     "pending",
					OrderID:    orderID,
				}, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "Delivery Retrieved Successfully",
		},
		{
			name:       "error - invalid UUID",
			deliveryID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:       "error - delivery not found",
			deliveryID: deliveryID.String(),
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
			if tt.deliveryID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.deliveryID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/deliveries/:id", controller.GetDelivery)

			req := httptest.NewRequest(http.MethodGet, "/deliveries/"+tt.deliveryID, nil)
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

func TestController_CreateDelivery(t *testing.T) {
	deliveryID := uuid.New()
	orderID := uuid.New()
	trackingNo := "TRACK-002"
	eta := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name            string
		requestBody     CreateDeliveryDTO
		mockSetup       func(*MockService, CreateDeliveryDTO)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "success - creates delivery",
			requestBody: CreateDeliveryDTO{
				ID:         deliveryID,
				Provider:   "fedex",
				TrackingNo: &trackingNo,
				Status:     "pending",
				Eta:        &eta,
				OrderID:    orderID,
			},
			mockSetup: func(mockSvc *MockService, dto CreateDeliveryDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateDeliveryDTO) bool {
					return actual.ID == dto.ID && actual.Provider == dto.Provider
				})).Return(&GetDeliveryDTO{
					ID:         dto.ID,
					Provider:   dto.Provider,
					TrackingNo: dto.TrackingNo,
					Status:     dto.Status,
					Eta:        dto.Eta,
					OrderID:    dto.OrderID,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Delivery Created Successfully",
		},
		{
			name: "error - service returns error",
			requestBody: CreateDeliveryDTO{
				ID:       deliveryID,
				Provider: "fedex",
				Status:   "pending",
				OrderID:  orderID,
			},
			mockSetup: func(mockSvc *MockService, dto CreateDeliveryDTO) {
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
			app.Post("/deliveries", controller.CreateDelivery)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/deliveries", bytes.NewBuffer(body))
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

func TestController_UpdateDelivery(t *testing.T) {
	deliveryID := uuid.New()
	newStatus := "in_transit"

	tests := []struct {
		name            string
		deliveryID      string
		requestBody     UpdateDeliveryDTO
		mockSetup       func(*MockService, uuid.UUID, UpdateDeliveryDTO)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:       "success - updates delivery",
			deliveryID: deliveryID.String(),
			requestBody: UpdateDeliveryDTO{
				ID:     deliveryID,
				Status: &newStatus,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateDeliveryDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return(&GetDeliveryDTO{
					ID:       id,
					Status:   newStatus,
					Provider: "fedex",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Delivery Updated Successfully",
		},
		{
			name:       "error - invalid UUID",
			deliveryID: "invalid-uuid",
			requestBody: UpdateDeliveryDTO{
				Status: &newStatus,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateDeliveryDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.deliveryID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.deliveryID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/deliveries/:id", controller.UpdateDelivery)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/deliveries/"+tt.deliveryID, bytes.NewBuffer(body))
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

func TestController_DeleteDelivery(t *testing.T) {
	deliveryID := uuid.New()

	tests := []struct {
		name            string
		deliveryID      string
		mockSetup       func(*MockService, uuid.UUID)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:       "success - deletes delivery",
			deliveryID: deliveryID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus:  http.StatusAccepted,
			expectedMessage: "Delivery Deleted Successfully",
		},
		{
			name:       "error - invalid UUID",
			deliveryID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:       "error - service returns error",
			deliveryID: deliveryID.String(),
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
			if tt.deliveryID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.deliveryID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/deliveries/:id", controller.DeleteDelivery)

			req := httptest.NewRequest(http.MethodDelete, "/deliveries/"+tt.deliveryID, nil)
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

