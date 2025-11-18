package meal_plans

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

func (m *MockService) List(ctx context.Context) ([]*GetMeal_planDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetMeal_planDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetMeal_planDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_planDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateMeal_planDTO) (*GetMeal_planDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_planDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateMeal_planDTO) (*GetMeal_planDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_planDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListMeal_plans(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns meal plans list",
			mockSetup: func(mockSvc *MockService) {
				goal := "Weight loss"
				expectedItems := []*GetMeal_planDTO{
					{
						ID:        uuid.New(),
						WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						Goal:      &goal,
						UserID:    uuid.New(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetMeal_planDTO)(nil), errors.New("database error"))
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
			app.Get("/meal-plans", controller.ListMeal_plans)

			req := httptest.NewRequest(http.MethodGet, "/meal-plans", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetMeal_plan(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - returns meal plan by ID",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				goal := "Weight loss"
				item := &GetMeal_planDTO{
					ID:        id,
					WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Goal:      &goal,
					UserID:    uuid.New(),
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
			name:   "error - meal plan not found",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetMeal_planDTO)(nil), errors.New("not found"))
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
			app.Get("/meal-plans/:id", controller.GetMeal_plan)

			req := httptest.NewRequest(http.MethodGet, "/meal-plans/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateMeal_plan(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateMeal_planDTO
		mockSetup      func(*MockService, CreateMeal_planDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new meal plan",
			requestBody: CreateMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Goal:      stringPtr("Weight loss"),
				UserID:    uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateMeal_planDTO) {
				expectedItem := &GetMeal_planDTO{
					ID:        dto.ID,
					WeekStart: dto.WeekStart,
					Goal:      dto.Goal,
					UserID:    dto.UserID,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateMeal_planDTO) bool {
					return actual.ID == dto.ID && actual.UserID == dto.UserID
				})).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UserID:    uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateMeal_planDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateMeal_planDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetMeal_planDTO)(nil), errors.New("creation failed"))
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
			app.Post("/meal-plans", controller.CreateMeal_plan)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/meal-plans", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateMeal_plan(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    UpdateMeal_planDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateMeal_planDTO)
		expectedStatus int
	}{
		{
			name:   "success - updates meal plan",
			itemID: uuid.New().String(),
			requestBody: UpdateMeal_planDTO{
				ID:   uuid.New(), // This will be overwritten by service
				Goal: stringPtr("Muscle gain"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateMeal_planDTO) {
				expectedItem := &GetMeal_planDTO{
					ID:        id,
					WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Goal:      dto.Goal,
					UserID:    uuid.New(),
				}
				// Service sets dto.ID = id, so use mock.Anything for DTO
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			requestBody:    UpdateMeal_planDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateMeal_planDTO) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			requestBody: UpdateMeal_planDTO{
				ID:   uuid.New(), // This will be overwritten by service
				Goal: stringPtr("Muscle gain"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateMeal_planDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return((*GetMeal_planDTO)(nil), errors.New("update failed"))
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
			app.Patch("/meal-plans/:id", controller.UpdateMeal_plan)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/meal-plans/"+tt.itemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteMeal_plan(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - deletes meal plan",
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
			app.Delete("/meal-plans/:id", controller.DeleteMeal_plan)

			req := httptest.NewRequest(http.MethodDelete, "/meal-plans/"+tt.itemID, nil)
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

