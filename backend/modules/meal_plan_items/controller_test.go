package meal_plan_items

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

func (m *MockService) List(ctx context.Context) ([]*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListMeal_plan_items(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns meal plan items list",
			mockSetup: func(mockSvc *MockService) {
				expectedItems := []*GetMeal_plan_itemDTO{
					{
						ID:         uuid.New(),
						Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						Slot:       "breakfast",
						MealPlanID: uuid.New(),
						RecipeID:   uuid.New(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetMeal_plan_itemDTO)(nil), errors.New("database error"))
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
			app.Get("/meal-plan-items", controller.ListMeal_plan_items)

			req := httptest.NewRequest(http.MethodGet, "/meal-plan-items", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetMeal_plan_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - returns meal plan item by ID",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				item := &GetMeal_plan_itemDTO{
					ID:         id,
					Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Slot:       "breakfast",
					MealPlanID: uuid.New(),
					RecipeID:   uuid.New(),
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
			name:   "error - meal plan item not found",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetMeal_plan_itemDTO)(nil), errors.New("not found"))
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
			app.Get("/meal-plan-items/:id", controller.GetMeal_plan_item)

			req := httptest.NewRequest(http.MethodGet, "/meal-plan-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateMeal_plan_item(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateMeal_plan_itemDTO
		mockSetup      func(*MockService, CreateMeal_plan_itemDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new meal plan item",
			requestBody: CreateMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "breakfast",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateMeal_plan_itemDTO) {
				expectedItem := &GetMeal_plan_itemDTO{
					ID:         dto.ID,
					Day:        dto.Day,
					Slot:       dto.Slot,
					MealPlanID: dto.MealPlanID,
					RecipeID:   dto.RecipeID,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateMeal_plan_itemDTO) bool {
					return actual.ID == dto.ID && actual.Slot == dto.Slot
				})).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "breakfast",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateMeal_plan_itemDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateMeal_plan_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetMeal_plan_itemDTO)(nil), errors.New("creation failed"))
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
			app.Post("/meal-plan-items", controller.CreateMeal_plan_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/meal-plan-items", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateMeal_plan_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    UpdateMeal_plan_itemDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateMeal_plan_itemDTO)
		expectedStatus int
	}{
		{
			name:   "success - updates meal plan item",
			itemID: uuid.New().String(),
			requestBody: UpdateMeal_plan_itemDTO{
				ID:   uuid.New(), // This will be overwritten by service
				Slot: stringPtr("lunch"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateMeal_plan_itemDTO) {
				expectedItem := &GetMeal_plan_itemDTO{
					ID:         id,
					Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Slot:       "lunch",
					MealPlanID: uuid.New(),
					RecipeID:   uuid.New(),
				}
				// Service sets dto.ID = id, so use mock.Anything for DTO
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			requestBody:    UpdateMeal_plan_itemDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateMeal_plan_itemDTO) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			requestBody: UpdateMeal_plan_itemDTO{
				ID:   uuid.New(), // This will be overwritten by service
				Slot: stringPtr("lunch"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateMeal_plan_itemDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return((*GetMeal_plan_itemDTO)(nil), errors.New("update failed"))
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
			app.Patch("/meal-plan-items/:id", controller.UpdateMeal_plan_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/meal-plan-items/"+tt.itemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteMeal_plan_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - deletes meal plan item",
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
			app.Delete("/meal-plan-items/:id", controller.DeleteMeal_plan_item)

			req := httptest.NewRequest(http.MethodDelete, "/meal-plan-items/"+tt.itemID, nil)
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

