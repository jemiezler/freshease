package recipes

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

func (m *MockService) List(ctx context.Context) ([]*GetRecipeDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetRecipeDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetRecipeDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipeDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateRecipeDTO) (*GetRecipeDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipeDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateRecipeDTO) (*GetRecipeDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipeDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListRecipes(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns recipes list",
			mockSetup: func(mockSvc *MockService) {
				instructions := "Test instructions"
				expectedRecipes := []*GetRecipeDTO{
					{
						ID:           uuid.New(),
						Name:         "Recipe One",
						Instructions: &instructions,
						Kcal:         500,
					},
					{
						ID:           uuid.New(),
						Name:         "Recipe Two",
						Instructions: &instructions,
						Kcal:         600,
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedRecipes, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMessage: "Recipes Retrieved Successfully",
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetRecipeDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/recipes", controller.ListRecipes)

			req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
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

func TestController_GetRecipe(t *testing.T) {
	recipeID := uuid.New()

	tests := []struct {
		name           string
		recipeID       string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - returns recipe by ID",
			recipeID: recipeID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				instructions := "Test instructions"
				mockSvc.On("Get", mock.Anything, id).Return(&GetRecipeDTO{
					ID:           id,
					Name:         "Test Recipe",
					Instructions: &instructions,
					Kcal:         500,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMessage: "Recipe Retrieved Successfully",
		},
		{
			name:     "error - invalid UUID",
			recipeID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:     "error - recipe not found",
			recipeID: recipeID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedMessage: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.recipeID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.recipeID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/recipes/:id", controller.GetRecipe)

			req := httptest.NewRequest(http.MethodGet, "/recipes/"+tt.recipeID, nil)
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

func TestController_CreateRecipe(t *testing.T) {
	recipeID := uuid.New()

	tests := []struct {
		name           string
		requestBody    CreateRecipeDTO
		mockSetup      func(*MockService, CreateRecipeDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates recipe",
			requestBody: CreateRecipeDTO{
				ID:   recipeID,
				Name: "New Recipe",
				Kcal: 550,
			},
			mockSetup: func(mockSvc *MockService, dto CreateRecipeDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(&GetRecipeDTO{
					ID:   dto.ID,
					Name: dto.Name,
					Kcal: dto.Kcal,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Recipe Created Successfully",
		},
		{
			name: "error - service returns error",
			requestBody: CreateRecipeDTO{
				ID:   recipeID,
				Name: "New Recipe",
				Kcal: 550,
			},
			mockSetup: func(mockSvc *MockService, dto CreateRecipeDTO) {
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
			app.Post("/recipes", controller.CreateRecipe)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(body))
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

func TestController_UpdateRecipe(t *testing.T) {
	recipeID := uuid.New()
	newName := "Updated Recipe"

	tests := []struct {
		name           string
		recipeID       string
		requestBody    UpdateRecipeDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateRecipeDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - updates recipe",
			recipeID: recipeID.String(),
			requestBody: UpdateRecipeDTO{
				ID:   recipeID,
				Name: &newName,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRecipeDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return(&GetRecipeDTO{
					ID:   id,
					Name: newName,
					Kcal: 500,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Recipe Updated Successfully",
		},
		{
			name:     "error - invalid UUID",
			recipeID: "invalid-uuid",
			requestBody: UpdateRecipeDTO{
				Name: &newName,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateRecipeDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.recipeID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.recipeID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/recipes/:id", controller.UpdateRecipe)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/recipes/"+tt.recipeID, bytes.NewBuffer(body))
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

func TestController_DeleteRecipe(t *testing.T) {
	recipeID := uuid.New()

	tests := []struct {
		name           string
		recipeID       string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - deletes recipe",
			recipeID: recipeID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedMessage: "Recipe Deleted Successfully",
		},
		{
			name:     "error - invalid UUID",
			recipeID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:     "error - service returns error",
			recipeID: recipeID.String(),
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
			if tt.recipeID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.recipeID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/recipes/:id", controller.DeleteRecipe)

			req := httptest.NewRequest(http.MethodDelete, "/recipes/"+tt.recipeID, nil)
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

