package notifications

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

func (m *MockService) List(ctx context.Context) ([]*GetNotificationDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetNotificationDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetNotificationDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNotificationDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateNotificationDTO) (*GetNotificationDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNotificationDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateNotificationDTO) (*GetNotificationDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNotificationDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListNotifications(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns notifications list",
			mockSetup: func(mockSvc *MockService) {
				body := "Your order has been shipped"
				expectedItems := []*GetNotificationDTO{
					{
						ID:        uuid.New(),
						Title:     "Order Shipped",
						Body:      &body,
						Channel:   "email",
						Status:    "sent",
						UserID:    uuid.New(),
						CreatedAt: time.Now(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetNotificationDTO)(nil), errors.New("database error"))
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
			app.Get("/notifications", controller.ListNotifications)

			req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetNotification(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - returns notification by ID",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				body := "Your order has been shipped"
				item := &GetNotificationDTO{
					ID:        id,
					Title:     "Order Shipped",
					Body:      &body,
					Channel:   "email",
					Status:    "sent",
					UserID:    uuid.New(),
					CreatedAt: time.Now(),
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
			name:   "error - notification not found",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetNotificationDTO)(nil), errors.New("not found"))
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
			app.Get("/notifications/:id", controller.GetNotification)

			req := httptest.NewRequest(http.MethodGet, "/notifications/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateNotification(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateNotificationDTO
		mockSetup      func(*MockService, CreateNotificationDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new notification",
			requestBody: CreateNotificationDTO{
				ID:      uuid.New(),
				Title:   "Order Shipped",
				Body:    stringPtr("Your order has been shipped"),
				Channel: "email",
				Status:  "sent",
				UserID:  uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateNotificationDTO) {
				expectedItem := &GetNotificationDTO{
					ID:        dto.ID,
					Title:     dto.Title,
					Body:      dto.Body,
					Channel:   dto.Channel,
					Status:    dto.Status,
					UserID:    dto.UserID,
					CreatedAt: time.Now(),
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateNotificationDTO) bool {
					return actual.ID == dto.ID && actual.Title == dto.Title
				})).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateNotificationDTO{
				ID:      uuid.New(),
				Title:   "Order Shipped",
				Channel: "email",
				Status:  "sent",
				UserID:  uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateNotificationDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateNotificationDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetNotificationDTO)(nil), errors.New("creation failed"))
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
			app.Post("/notifications", controller.CreateNotification)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateNotification(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    UpdateNotificationDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateNotificationDTO)
		expectedStatus int
	}{
		{
			name:   "success - updates notification",
			itemID: uuid.New().String(),
			requestBody: UpdateNotificationDTO{
				ID:     uuid.New(), // This will be overwritten by service
				Status: stringPtr("read"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateNotificationDTO) {
				body := "Your order has been shipped"
				expectedItem := &GetNotificationDTO{
					ID:        id,
					Title:     "Order Shipped",
					Body:      &body,
					Channel:   "email",
					Status:    "read",
					UserID:    uuid.New(),
					CreatedAt: time.Now(),
				}
				// Service sets dto.ID = id, so use mock.Anything for DTO
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			requestBody:    UpdateNotificationDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateNotificationDTO) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			requestBody: UpdateNotificationDTO{
				ID:     uuid.New(), // This will be overwritten by service
				Status: stringPtr("read"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateNotificationDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return((*GetNotificationDTO)(nil), errors.New("update failed"))
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
			app.Patch("/notifications/:id", controller.UpdateNotification)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/notifications/"+tt.itemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteNotification(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - deletes notification",
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
			app.Delete("/notifications/:id", controller.DeleteNotification)

			req := httptest.NewRequest(http.MethodDelete, "/notifications/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

