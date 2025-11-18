package categories

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

func (m *MockService) List(ctx context.Context) ([]*GetCategoryDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetCategoryDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetCategoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCategoryDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateCategoryDTO) (*GetCategoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCategoryDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateCategoryDTO) (*GetCategoryDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCategoryDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListCategories(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns categories list",
			mockSetup: func(mockSvc *MockService) {
				expectedCategories := []*GetCategoryDTO{
					{
						ID:   uuid.New(),
						Name: "Category One",
						Slug: "category-one",
					},
					{
						ID:   uuid.New(),
						Name: "Category Two",
						Slug: "category-two",
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedCategories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"message": "Categories Retrieved Successfully"},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetCategoryDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/categories", controller.ListCategories)

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetCategory(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name           string
		categoryID     string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:       "success - returns category by ID",
			categoryID: categoryID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(&GetCategoryDTO{
					ID:   id,
					Name: "Test Category",
					Slug: "test-category",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"message": "Category Retrieved Successfully"},
		},
		{
			name:       "error - invalid UUID",
			categoryID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "invalid uuid"},
		},
		{
			name:       "error - category not found",
			categoryID: categoryID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"message": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.categoryID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.categoryID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/categories/:id", controller.GetCategory)

			req := httptest.NewRequest(http.MethodGet, "/categories/"+tt.categoryID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateCategory(t *testing.T) {
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    CreateCategoryDTO
		mockSetup      func(*MockService, CreateCategoryDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates category",
			requestBody: CreateCategoryDTO{
				ID:        categoryID,
				Name:      "New Category",
				Slug:      "new-category",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mockSvc *MockService, dto CreateCategoryDTO) {
				mockSvc.On("Create", mock.Anything, mock.Anything).Return(&GetCategoryDTO{
					ID:   dto.ID,
					Name: dto.Name,
					Slug: dto.Slug,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]interface{}{"message": "Category Created Successfully"},
		},
		{
			name: "error - service returns error",
			requestBody: CreateCategoryDTO{
				ID:        categoryID,
				Name:      "New Category",
				Slug:      "new-category",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mockSvc *MockService, dto CreateCategoryDTO) {
				mockSvc.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("validation error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "validation error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/categories", controller.CreateCategory)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateCategory(t *testing.T) {
	categoryID := uuid.New()
	newName := "Updated Category"

	tests := []struct {
		name           string
		categoryID     string
		requestBody    UpdateCategoryDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateCategoryDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:       "success - updates category",
			categoryID: categoryID.String(),
			requestBody: UpdateCategoryDTO{
				ID:        categoryID,
				Name:      &newName,
				UpdatedAt: time.Now(),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCategoryDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return(&GetCategoryDTO{
					ID:   id,
					Name: newName,
					Slug: "test-category",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]interface{}{"message": "Category Updated Successfully"},
		},
		{
			name:       "error - invalid UUID",
			categoryID: "invalid-uuid",
			requestBody: UpdateCategoryDTO{
				Name:      &newName,
				UpdatedAt: time.Now(),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCategoryDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "invalid uuid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.categoryID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.categoryID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/categories/:id", controller.UpdateCategory)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/categories/"+tt.categoryID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteCategory(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name           string
		categoryID     string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:       "success - deletes category",
			categoryID: categoryID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   map[string]interface{}{"message": "Category Deleted Successfully"},
		},
		{
			name:       "error - invalid UUID",
			categoryID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "invalid uuid"},
		},
		{
			name:       "error - service returns error",
			categoryID: categoryID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.categoryID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.categoryID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/categories/:id", controller.DeleteCategory)

			req := httptest.NewRequest(http.MethodDelete, "/categories/"+tt.categoryID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])

			mockSvc.AssertExpectations(t)
		})
	}
}

