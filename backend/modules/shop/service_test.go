package shop

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetActiveProducts(ctx context.Context, filters ShopSearchFilters) ([]*ShopProductDTO, int, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*ShopProductDTO), args.Int(1), args.Error(2)
}

func (m *MockRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*ShopProductDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopProductDTO), args.Error(1)
}

func (m *MockRepository) GetActiveCategories(ctx context.Context) ([]*ShopCategoryDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*ShopCategoryDTO), args.Error(1)
}

func (m *MockRepository) GetActiveVendors(ctx context.Context) ([]*ShopVendorDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*ShopVendorDTO), args.Error(1)
}

func (m *MockRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*ShopCategoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopCategoryDTO), args.Error(1)
}

func (m *MockRepository) GetVendorByID(ctx context.Context, id uuid.UUID) (*ShopVendorDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ShopVendorDTO), args.Error(1)
}

func TestService_SearchProducts(t *testing.T) {
	tests := []struct {
		name           string
		filters        ShopSearchFilters
		mockSetup      func(*MockRepository)
		expectedResult *ShopSearchResponse
		expectedError  error
	}{
		{
			name: "successful search with default pagination",
			filters: ShopSearchFilters{
				Limit:  0,
				Offset: 0,
			},
			mockSetup: func(mr *MockRepository) {
				products := []*ShopProductDTO{
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
				}
				mr.On("GetActiveProducts", mock.Anything, mock.MatchedBy(func(f ShopSearchFilters) bool {
					return f.Limit == 20 && f.Offset == 0 // Should set default limit
				})).Return(products, 1, nil)
			},
			expectedResult: &ShopSearchResponse{
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
			},
			expectedError: nil,
		},
		{
			name: "successful search with custom pagination",
			filters: ShopSearchFilters{
				Limit:  10,
				Offset: 5,
			},
			mockSetup: func(mr *MockRepository) {
				products := []*ShopProductDTO{}
				mr.On("GetActiveProducts", mock.Anything, mock.MatchedBy(func(f ShopSearchFilters) bool {
					return f.Limit == 10 && f.Offset == 5
				})).Return(products, 0, nil)
			},
			expectedResult: &ShopSearchResponse{
				Products: []*ShopProductDTO{},
				Total:    0,
				Limit:    10,
				Offset:   5,
				HasMore:  false,
			},
			expectedError: nil,
		},
		{
			name: "search with limit cap",
			filters: ShopSearchFilters{
				Limit:  150, // Should be capped at 100
				Offset: 0,
			},
			mockSetup: func(mr *MockRepository) {
				products := []*ShopProductDTO{}
				mr.On("GetActiveProducts", mock.Anything, mock.MatchedBy(func(f ShopSearchFilters) bool {
					return f.Limit == 100 // Should be capped
				})).Return(products, 0, nil)
			},
			expectedResult: &ShopSearchResponse{
				Products: []*ShopProductDTO{},
				Total:    0,
				Limit:    100,
				Offset:   0,
				HasMore:  false,
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			filters: ShopSearchFilters{
				Limit:  20,
				Offset: 0,
			},
			mockSetup: func(mr *MockRepository) {
				mr.On("GetActiveProducts", mock.Anything, mock.Anything).Return([]*ShopProductDTO(nil), 0, errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.SearchProducts(context.Background(), tt.filters)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Total, result.Total)
				assert.Equal(t, tt.expectedResult.Limit, result.Limit)
				assert.Equal(t, tt.expectedResult.Offset, result.Offset)
				assert.Equal(t, tt.expectedResult.HasMore, result.HasMore)
				assert.Len(t, result.Products, len(tt.expectedResult.Products))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      uuid.UUID
		mockSetup      func(*MockRepository)
		expectedResult *ShopProductDTO
		expectedError  error
	}{
		{
			name:      "successful get product",
			productID: uuid.New(),
			mockSetup: func(mr *MockRepository) {
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
				mr.On("GetProductByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(product, nil)
			},
			expectedResult: &ShopProductDTO{
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
			expectedError: nil,
		},
		{
			name:      "product not found",
			productID: uuid.New(),
			mockSetup: func(mr *MockRepository) {
				mr.On("GetProductByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.GetProduct(context.Background(), tt.productID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Price, result.Price)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetCategories(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockRepository)
		expectedResult []*ShopCategoryDTO
		expectedError  error
	}{
		{
			name: "successful get categories",
			mockSetup: func(mr *MockRepository) {
				categories := []*ShopCategoryDTO{
					{
						ID:          uuid.New(),
						Name:        "Fruits",
						Description: "Fresh fruits",
					},
				}
				mr.On("GetActiveCategories", mock.Anything).Return(categories, nil)
			},
			expectedResult: []*ShopCategoryDTO{
				{
					ID:          uuid.New(),
					Name:        "Fruits",
					Description: "Fresh fruits",
				},
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			mockSetup: func(mr *MockRepository) {
				mr.On("GetActiveCategories", mock.Anything).Return([]*ShopCategoryDTO(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.GetCategories(context.Background())

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetVendors(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockRepository)
		expectedResult []*ShopVendorDTO
		expectedError  error
	}{
		{
			name: "successful get vendors",
			mockSetup: func(mr *MockRepository) {
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
				mr.On("GetActiveVendors", mock.Anything).Return(vendors, nil)
			},
			expectedResult: []*ShopVendorDTO{
				{
					ID:        uuid.New(),
					Name:      "Test Vendor",
					Email:     "vendor@test.com",
					IsActive:  "active",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			mockSetup: func(mr *MockRepository) {
				mr.On("GetActiveVendors", mock.Anything).Return([]*ShopVendorDTO(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.GetVendors(context.Background())

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
