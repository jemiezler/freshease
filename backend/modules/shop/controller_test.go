package shop

import (
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
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) SearchProducts(ctx context.Context, filters ShopSearchFilters) (*ShopSearchResponse, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopSearchResponse), args.Error(1)
}

func (m *MockService) GetProduct(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopProductDTO), args.Error(1)
}

func (m *MockService) GetCategories(ctx context.Context) ([]*ShopCategoryDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*ShopCategoryDTO), args.Error(1)
}

func (m *MockService) GetVendors(ctx context.Context) ([]*ShopVendorDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*ShopVendorDTO), args.Error(1)
}

func (m *MockService) GetCategory(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopCategoryDTO), args.Error(1)
}

func (m *MockService) GetVendor(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopVendorDTO), args.Error(1)
}

func TestController_SearchProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "successful search with no filters",
			queryParams: "",
			mockSetup: func(ms *MockService) {
				response := &ShopSearchResponse{
					Products: []*ShopProductDTO{
						{
							ID:          uuid.New(),
							Name:        "Test Product",
							Price:       10.99,
							Description: "Test Description",
							ImageURL:    "https://example.com/image.jpg",
							UnitLabel:   "kg",
							IsActive:    "active",
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
					},
					Total:   1,
					Limit:   20,
					Offset:  0,
					HasMore: false,
				}
				ms.On("SearchProducts", mock.Anything, mock.Anything).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "successful search with filters",
			queryParams: "?category_id=123e4567-e89b-12d3-a456-426614174000&min_price=5.0&max_price=20.0&search=apple&in_stock=true&limit=10&offset=5",
			mockSetup: func(ms *MockService) {
				response := &ShopSearchResponse{
					Products: []*ShopProductDTO{},
					Total:    0,
					Limit:    10,
					Offset:   5,
					HasMore:  false,
				}
				ms.On("SearchProducts", mock.Anything, mock.Anything).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "invalid category_id",
			queryParams: "?category_id=invalid-uuid",
			mockSetup: func(ms *MockService) {
				response := &ShopSearchResponse{
					Products: []*ShopProductDTO{},
					Total:    0,
					Limit:    20,
					Offset:   0,
					HasMore:  false,
				}
				ms.On("SearchProducts", mock.Anything, mock.Anything).Return(response, nil)
			},
			expectedStatus: http.StatusOK, // Should still work, just ignore invalid UUID
			expectedError:  false,
		},
		{
			name:        "service error",
			queryParams: "",
			mockSetup: func(ms *MockService) {
				ms.On("SearchProducts", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := new(MockService)
			tt.mockSetup(mockService)

			controller := NewController(mockService)
			app.Get("/products", controller.SearchProducts)

			req := httptest.NewRequest(http.MethodGet, "/products"+tt.queryParams, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			} else {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "message")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestController_GetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "successful get product",
			productID: uuid.New().String(),
			mockSetup: func(ms *MockService) {
				product := &ShopProductDTO{
					ID:          uuid.New(),
					Name:        "Test Product",
					Price:       10.99,
					Description: "Test Description",
					ImageURL:    "https://example.com/image.jpg",
					UnitLabel:   "kg",
					IsActive:    "active",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				ms.On("GetProduct", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(product, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid product ID",
			productID:      "invalid-uuid",
			mockSetup:      func(ms *MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "product not found",
			productID: uuid.New().String(),
			mockSetup: func(ms *MockService) {
				ms.On("GetProduct", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := new(MockService)
			tt.mockSetup(mockService)

			controller := NewController(mockService)
			app.Get("/products/:id", controller.GetProduct)

			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			} else {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "message")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestController_GetCategories(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful get categories",
			mockSetup: func(ms *MockService) {
				categories := []*ShopCategoryDTO{
					{
						ID:          uuid.New(),
						Name:        "Fruits",
						Description: "Fresh fruits",
					},
				}
				ms.On("GetCategories", mock.Anything).Return(categories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "service error",
			mockSetup: func(ms *MockService) {
				ms.On("GetCategories", mock.Anything).Return([]*ShopCategoryDTO(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := new(MockService)
			tt.mockSetup(mockService)

			controller := NewController(mockService)
			app.Get("/categories", controller.GetCategories)

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			} else {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "message")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestController_GetVendors(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful get vendors",
			mockSetup: func(ms *MockService) {
				vendors := []*ShopVendorDTO{
					{
						ID:        uuid.New(),
						Name:      "Test Vendor",
						Email:     "vendor@test.com",
						IsActive:  "active",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				ms.On("GetVendors", mock.Anything).Return(vendors, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "service error",
			mockSetup: func(ms *MockService) {
				ms.On("GetVendors", mock.Anything).Return([]*ShopVendorDTO(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := new(MockService)
			tt.mockSetup(mockService)

			controller := NewController(mockService)
			app.Get("/vendors", controller.GetVendors)

			req := httptest.NewRequest(http.MethodGet, "/vendors", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			} else {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "message")
			}

			mockService.AssertExpectations(t)
		})
	}
}
