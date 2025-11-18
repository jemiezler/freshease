package recipe_items

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

func (m *MockService) List(ctx context.Context) ([]*GetRecipe_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetRecipe_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateRecipe_itemDTO) (*GetRecipe_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateRecipe_itemDTO) (*GetRecipe_itemDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListRecipe_items(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns recipe items list",
			mockSetup: func(mockSvc *MockService) {
				expectedItems := []*GetRecipe_itemDTO{
					{
						ID:        uuid.New(),
						Amount:    2.5,
						Unit:      "cups",
						RecipeID:  uuid.New(),
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
				mockSvc.On("List", mock.Anything).Return(([]*GetRecipe_itemDTO)(nil), errors.New("database error"))
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
			app.Get("/recipe-items", controller.ListRecipe_items)

			req := httptest.NewRequest(http.MethodGet, "/recipe-items", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetRecipe_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - returns recipe item by ID",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				item := &GetRecipe_itemDTO{
					ID:        id,
					Amount:    2.5,
					Unit:      "cups",
					RecipeID:  uuid.New(),
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
			name:   "error - recipe item not found",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetRecipe_itemDTO)(nil), errors.New("not found"))
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
			app.Get("/recipe-items/:id", controller.GetRecipe_item)

			req := httptest.NewRequest(http.MethodGet, "/recipe-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateRecipe_item(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateRecipe_itemDTO
		mockSetup      func(*MockService, CreateRecipe_itemDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new recipe item",
			requestBody: CreateRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    2.5,
				Unit:      "cups",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateRecipe_itemDTO) {
				expectedItem := &GetRecipe_itemDTO{
					ID:        dto.ID,
					Amount:    dto.Amount,
					Unit:      dto.Unit,
					RecipeID:  dto.RecipeID,
					ProductID: dto.ProductID,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateRecipe_itemDTO) bool {
					return actual.ID == dto.ID && actual.Amount == dto.Amount
				})).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    2.5,
				Unit:      "cups",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateRecipe_itemDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateRecipe_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetRecipe_itemDTO)(nil), errors.New("creation failed"))
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
			app.Post("/recipe-items", controller.CreateRecipe_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/recipe-items", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateRecipe_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    UpdateRecipe_itemDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateRecipe_itemDTO)
		expectedStatus int
	}{
		{
			name:   "success - updates recipe item",
			itemID: uuid.New().String(),
			requestBody: UpdateRecipe_itemDTO{
				ID:     uuid.New(), // This will be overwritten by service
				Amount: float64Ptr(3.0),
				Unit:   stringPtr("tbsp"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRecipe_itemDTO) {
				expectedItem := &GetRecipe_itemDTO{
					ID:        id,
					Amount:    3.0,
					Unit:      "tbsp",
					RecipeID:  uuid.New(),
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
			requestBody:    UpdateRecipe_itemDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateRecipe_itemDTO) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			requestBody: UpdateRecipe_itemDTO{
				ID:     uuid.New(), // This will be overwritten by service
				Amount: float64Ptr(3.0),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRecipe_itemDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return((*GetRecipe_itemDTO)(nil), errors.New("update failed"))
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
			app.Patch("/recipe-items/:id", controller.UpdateRecipe_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/recipe-items/"+tt.itemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteRecipe_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - deletes recipe item",
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
			app.Delete("/recipe-items/:id", controller.DeleteRecipe_item)

			req := httptest.NewRequest(http.MethodDelete, "/recipe-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

