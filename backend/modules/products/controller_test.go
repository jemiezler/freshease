package products

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

func (m *MockService) List(ctx context.Context) ([]*GetProductDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetProductDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetProductDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateProductDTO) (*GetProductDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateProductDTO) (*GetProductDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListProducts(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns products list",
			mockSetup: func(mockSvc *MockService) {
				expectedProducts := []*GetProductDTO{
					{
						ID:          uuid.New(),
						Name:        "Product One",
						Price:       99.99,
						Description: "First product",
						ImageURL:    "https://example.com/image1.jpg",
						UnitLabel:   "kg",
						IsActive:    "true",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						Name:        "Product Two",
						Price:       149.99,
						Description: "Second product",
						ImageURL:    "https://example.com/image2.jpg",
						UnitLabel:   "piece",
						IsActive:    "true",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedProducts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Products Retrieved Successfully",
			},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetProductDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/products", controller.ListProducts)

			req := httptest.NewRequest(http.MethodGet, "/products", nil)
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

func TestController_GetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "success - returns product by ID",
			productID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				expectedProduct := &GetProductDTO{
					ID:          id,
					Name:        "Test Product",
					Price:       99.99,
					Description: "Test product description",
					ImageURL:    "https://example.com/image.jpg",
					UnitLabel:   "kg",
					IsActive:    "true",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(expectedProduct, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Product Retrieved Successfully",
			},
		},
		{
			name:           "error - invalid UUID",
			productID:      "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid uuid",
			},
		},
		{
			name:      "error - product not found",
			productID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetProductDTO)(nil), errors.New("product not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"message": "not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.productID != "invalid-uuid" {
				productID, err := uuid.Parse(tt.productID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, productID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/products/:id", controller.GetProduct)

			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
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

func TestController_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateProductDTO
		mockSetup      func(*MockService, CreateProductDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates new product",
			requestBody: CreateProductDTO{
				ID:          uuid.New(),
				Name:        "New Product",
				Price:       199.99,
				Description: "New product description",
				ImageURL:    "https://example.com/new-image.jpg",
				UnitLabel:   "kg",
				IsActive:    "true",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateProductDTO) {
				expectedProduct := &GetProductDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Price:       dto.Price,
					Description: dto.Description,
					ImageURL:    dto.ImageURL,
					UnitLabel:   dto.UnitLabel,
					IsActive:    dto.IsActive,
					CreatedAt:   dto.CreatedAt,
					UpdatedAt:   dto.UpdatedAt,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateProductDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Price == dto.Price &&
						actual.Description == dto.Description &&
						actual.ImageURL == dto.ImageURL &&
						actual.UnitLabel == dto.UnitLabel &&
						actual.IsActive == dto.IsActive
				})).Return(expectedProduct, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Product Created Successfully",
			},
		},
		{
			name: "error - service returns error",
			requestBody: CreateProductDTO{
				ID:          uuid.New(),
				Name:        "New Product",
				Price:       199.99,
				Description: "New product description",
				ImageURL:    "https://example.com/new-image.jpg",
				UnitLabel:   "kg",
				IsActive:    "true",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateProductDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateProductDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Price == dto.Price &&
						actual.Description == dto.Description &&
						actual.ImageURL == dto.ImageURL &&
						actual.UnitLabel == dto.UnitLabel &&
						actual.IsActive == dto.IsActive
				})).Return((*GetProductDTO)(nil), errors.New("name already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "name already exists",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/products", controller.CreateProduct)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
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

func TestController_UpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		requestBody    UpdateProductDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateProductDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "success - updates product",
			productID: uuid.New().String(),
			requestBody: UpdateProductDTO{
				Name:        stringPtr("Updated Product"),
				Price:       float64Ptr(299.99),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateProductDTO) {
				expectedProduct := &GetProductDTO{
					ID:          id,
					Name:        *dto.Name,
					Price:       *dto.Price,
					Description: *dto.Description,
					ImageURL:    "https://example.com/image.jpg",
					UnitLabel:   "kg",
					IsActive:    "true",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockSvc.On("Update", mock.Anything, id, dto).Return(expectedProduct, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Product Updated Successfully",
			},
		},
		{
			name:           "error - invalid UUID",
			productID:      "invalid-uuid",
			requestBody:    UpdateProductDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateProductDTO) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid uuid",
			},
		},
		{
			name:      "error - service returns error",
			productID: uuid.New().String(),
			requestBody: UpdateProductDTO{
				Name: stringPtr("Updated Product"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateProductDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return((*GetProductDTO)(nil), errors.New("product not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "product not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.productID != "invalid-uuid" {
				productID, err := uuid.Parse(tt.productID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, productID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/products/:id", controller.UpdateProduct)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/products/"+tt.productID, bytes.NewBuffer(jsonBody))
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

func TestController_DeleteProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "success - deletes product",
			productID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody: map[string]interface{}{
				"message": "Product Deleted Successfully",
			},
		},
		{
			name:           "error - invalid UUID",
			productID:      "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid uuid",
			},
		},
		{
			name:      "error - service returns error",
			productID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("product not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "product not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.productID != "invalid-uuid" {
				productID, err := uuid.Parse(tt.productID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, productID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/products/:id", controller.DeleteProduct)

			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.productID, nil)
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

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
