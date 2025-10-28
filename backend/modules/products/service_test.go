package products

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

func (m *MockRepository) List(ctx context.Context) ([]*GetProductDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetProductDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetProductDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, p *CreateProductDTO) (*GetProductDTO, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, p *UpdateProductDTO) (*GetProductDTO, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockRepository)
		expectedResult []*GetProductDTO
		expectedError  error
	}{
		{
			name: "success - returns products list",
			mockSetup: func(mockRepo *MockRepository) {
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
				mockRepo.On("List", mock.Anything).Return(expectedProducts, nil)
			},
			expectedResult: []*GetProductDTO{
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
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetProductDTO(nil), errors.New("database error"))
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
			ctx := context.Background()

			result, err := service.List(ctx)

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

func TestService_Get(t *testing.T) {
	tests := []struct {
		name           string
		productID      uuid.UUID
		mockSetup      func(*MockRepository, uuid.UUID)
		expectedResult *GetProductDTO
		expectedError  error
	}{
		{
			name:      "success - returns product by ID",
			productID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
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
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedProduct, nil)
			},
			expectedResult: &GetProductDTO{
				ID:          uuid.New(),
				Name:        "Test Product",
				Price:       99.99,
				Description: "Test product description",
				ImageURL:    "https://example.com/image.jpg",
				UnitLabel:   "kg",
				IsActive:    "true",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectedError: nil,
		},
		{
			name:      "error - product not found",
			productID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetProductDTO)(nil), errors.New("product not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("product not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.productID)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Get(ctx, tt.productID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.productID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name           string
		createDTO      CreateProductDTO
		mockSetup      func(*MockRepository, CreateProductDTO)
		expectedResult *GetProductDTO
		expectedError  error
	}{
		{
			name: "success - creates new product",
			createDTO: CreateProductDTO{
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
			mockSetup: func(mockRepo *MockRepository, dto CreateProductDTO) {
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
				mockRepo.On("Create", mock.Anything, &dto).Return(expectedProduct, nil)
			},
			expectedResult: &GetProductDTO{
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
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateProductDTO{
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
			mockSetup: func(mockRepo *MockRepository, dto CreateProductDTO) {
				mockRepo.On("Create", mock.Anything, &dto).Return((*GetProductDTO)(nil), errors.New("name already exists"))
			},
			expectedResult: nil,
			expectedError:  errors.New("name already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.createDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.createDTO.Name, result.Name)
				assert.Equal(t, tt.createDTO.Price, result.Price)
				assert.Equal(t, tt.createDTO.Description, result.Description)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name           string
		productID      uuid.UUID
		updateDTO      UpdateProductDTO
		mockSetup      func(*MockRepository, uuid.UUID, UpdateProductDTO)
		expectedResult *GetProductDTO
		expectedError  error
	}{
		{
			name:      "success - updates product",
			productID: uuid.New(),
			updateDTO: UpdateProductDTO{
				Name:        stringPtr("Updated Product"),
				Price:       float64Ptr(299.99),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductDTO) {
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
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *UpdateProductDTO) bool {
					return p.ID == id && *p.Name == *dto.Name && *p.Price == *dto.Price
				})).Return(expectedProduct, nil)
			},
			expectedResult: &GetProductDTO{
				ID:          uuid.New(),
				Name:        "Updated Product",
				Price:       299.99,
				Description: "Updated description",
				ImageURL:    "https://example.com/image.jpg",
				UnitLabel:   "kg",
				IsActive:    "true",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectedError: nil,
		},
		{
			name:      "error - repository returns error",
			productID: uuid.New(),
			updateDTO: UpdateProductDTO{
				Name: stringPtr("Updated Product"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductDTO) {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return((*GetProductDTO)(nil), errors.New("product not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("product not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.productID, tt.updateDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Update(ctx, tt.productID, tt.updateDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.productID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		productID     uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError error
	}{
		{
			name:      "success - deletes product",
			productID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "error - repository returns error",
			productID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("product not found"))
			},
			expectedError: errors.New("product not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.productID)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.productID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
