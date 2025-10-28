package product_categories

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

func (m *MockService) List(ctx context.Context) ([]*GetProductCategoryDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetProductCategoryDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductCategoryDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductCategoryDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductCategoryDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func TestController_ListProduct_categories(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns list of product categories",
			mockSetup: func(mockSvc *MockService) {
				categories := []*GetProductCategoryDTO{
					{
						ID:          uuid.New(),
						Name:        "Fruits",
						Description: "Fresh fruits and vegetables",
						Slug:        "fruits",
					},
					{
						ID:          uuid.New(),
						Name:        "Vegetables",
						Description: "Fresh vegetables",
						Slug:        "vegetables",
					},
				}
				mockSvc.On("List", mock.Anything).Return(categories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetProductCategoryDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetProductCategoryDTO)(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/product_categories", controller.ListProduct_categories)

			req := httptest.NewRequest(http.MethodGet, "/product_categories", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetProduct_category(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - returns product category",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				category := &GetProductCategoryDTO{
					ID:          id,
					Name:        "Fruits",
					Description: "Fresh fruits and vegetables",
					Slug:        "fruits",
				}
				mockSvc.On("Get", mock.Anything, id).Return(category, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - product category not found",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetProductCategoryDTO)(nil), errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.id != "invalid-uuid" {
				id, _ := uuid.Parse(tt.id)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/product_categories/:id", controller.GetProduct_category)

			req := httptest.NewRequest(http.MethodGet, "/product_categories/"+tt.id, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateProduct_category(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateProductCategoryDTO
		mockSetup      func(*MockService, CreateProductCategoryDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - creates product category",
			requestBody: CreateProductCategoryDTO{
				ID:          uuid.New(),
				Name:        "Fruits",
				Description: "Fresh fruits and vegetables",
				Slug:        "fruits",
			},
			mockSetup: func(mockSvc *MockService, dto CreateProductCategoryDTO) {
				createdCategory := &GetProductCategoryDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Description: dto.Description,
					Slug:        dto.Slug,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateProductCategoryDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description &&
						actual.Slug == dto.Slug
				})).Return(createdCategory, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: CreateProductCategoryDTO{
				ID:          uuid.New(),
				Name:        "Fruits",
				Description: "Fresh fruits and vegetables",
				Slug:        "fruits",
			},
			mockSetup: func(mockSvc *MockService, dto CreateProductCategoryDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateProductCategoryDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description &&
						actual.Slug == dto.Slug
				})).Return((*GetProductCategoryDTO)(nil), errors.New("creation failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/product_categories", controller.CreateProduct_category)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/product_categories", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateProduct_category(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    UpdateProductCategoryDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateProductCategoryDTO)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - updates product category",
			id:   uuid.New().String(),
			requestBody: UpdateProductCategoryDTO{
				Name:        stringPtr("Updated Fruits"),
				Description: stringPtr("Updated description"),
				Slug:        stringPtr("updated-fruits"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateProductCategoryDTO) {
				updatedCategory := &GetProductCategoryDTO{
					ID:          id,
					Name:        "Updated Fruits",
					Description: "Updated description",
					Slug:        "updated-fruits",
				}
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateProductCategoryDTO) bool {
					return actual.Name != nil && *actual.Name == "Updated Fruits" &&
						actual.Description != nil && *actual.Description == "Updated description" &&
						actual.Slug != nil && *actual.Slug == "updated-fruits"
				})).Return(updatedCategory, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			requestBody: UpdateProductCategoryDTO{
				Name: stringPtr("Updated Fruits"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateProductCategoryDTO) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			requestBody: UpdateProductCategoryDTO{
				Name: stringPtr("Updated Fruits"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateProductCategoryDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.MatchedBy(func(actual UpdateProductCategoryDTO) bool {
					return actual.Name != nil && *actual.Name == "Updated Fruits"
				})).Return((*GetProductCategoryDTO)(nil), errors.New("update failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.id != "invalid-uuid" {
				id, _ := uuid.Parse(tt.id)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/product_categories/:id", controller.UpdateProduct_category)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/product_categories/"+tt.id, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "data")
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteProduct_category(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - deletes product category",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedError:  false,
		},
		{
			name: "error - invalid UUID",
			id:   "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed - should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "error - service returns error",
			id:   uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.id != "invalid-uuid" {
				id, _ := uuid.Parse(tt.id)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/product_categories/:id", controller.DeleteProduct_category)

			req := httptest.NewRequest(http.MethodDelete, "/product_categories/"+tt.id, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "message")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
