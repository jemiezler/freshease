package products

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"
	"time"

	"freshease/backend/modules/uploads"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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

// MockUploadsService is a mock implementation of uploads.Service
type MockUploadsService struct {
	mock.Mock
}

func (m *MockUploadsService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, folder)
	return args.String(0), args.Error(1)
}

func (m *MockUploadsService) DeleteImage(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

func (m *MockUploadsService) GetImageURL(ctx context.Context, objectName string) (string, error) {
	args := m.Called(ctx, objectName)
	return args.String(0), args.Error(1)
}

func (m *MockUploadsService) GetImage(ctx context.Context, objectName string) (io.ReadCloser, *minio.ObjectInfo, error) {
	args := m.Called(ctx, objectName)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(io.ReadCloser), nil, args.Error(2)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(*minio.ObjectInfo), args.Error(2)
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
						SKU:         "PROD-001",
						Price:       99.99,
						Description: stringPtr("First product"),
						UnitLabel:   "kg",
						IsActive:    true,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						Name:        "Product Two",
						SKU:         "PROD-002",
						Price:       149.99,
						Description: stringPtr("Second product"),
						UnitLabel:   "piece",
						IsActive:    true,
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
					SKU:         "PROD-001",
					Description: stringPtr("First product"),
					UnitLabel:   "kg",
					IsActive:    true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          uuid.New(),
					Name:        "Product Two",
					Price:       149.99,
					SKU:         "PROD-002",
					Description: stringPtr("Second product"),
					UnitLabel:   "piece",
					IsActive:    true,
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo, mockUploads)
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
					SKU:         "PROD-003",
					Description: stringPtr("Test product description"),
					UnitLabel:   "kg",
					IsActive:    true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedProduct, nil)
			},
			expectedResult: &GetProductDTO{
				ID:          uuid.New(),
				Name:        "Test Product",
				Price:       99.99,
				SKU:         "PROD-003",
				Description: stringPtr("Test product description"),
				UnitLabel:   "kg",
				IsActive:    true,
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.productID)

			service := NewService(mockRepo, mockUploads)
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
				ID:           uuid.New(),
				Name:         "New Product",
				SKU:          "PROD-004",
				Price:        199.99,
				Description:  stringPtr("New product description"),
				UnitLabel:    "kg",
				IsActive:     true,
				Quantity:     100,
				ReorderLevel: 50,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateProductDTO) {
				expectedProduct := &GetProductDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					SKU:         dto.SKU,
					Price:       dto.Price,
					Description: dto.Description,
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
				SKU:         "PROD-004",
				Price:       199.99,
				Description: stringPtr("New product description"),
				UnitLabel:   "kg",
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateProductDTO{
				ID:           uuid.New(),
				Name:         "New Product",
				SKU:          "PROD-004",
				Price:        199.99,
				Description:  stringPtr("New product description"),
				UnitLabel:    "kg",
				IsActive:     true,
				Quantity:     100,
				ReorderLevel: 50,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo, mockUploads)
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
					Description: dto.Description,
					SKU:         "PROD-005",
					UnitLabel:   "kg",
					IsActive:    true,
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
				SKU:         "PROD-005",
				Description: stringPtr("Updated description"),
				UnitLabel:   "kg",
				IsActive:    true,
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.productID, tt.updateDTO)

			service := NewService(mockRepo, mockUploads)
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.productID)

			service := NewService(mockRepo, mockUploads)
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

func TestService_UploadProductImage(t *testing.T) {
	tests := []struct {
		name           string
		file           *multipart.FileHeader
		mockSetup      func(*MockUploadsService, *multipart.FileHeader)
		expectedResult string
		expectedError  error
	}{
		{
			name: "success - uploads product image",
			file: &multipart.FileHeader{
				Filename: "product.jpg",
				Size:     2048,
			},
			mockSetup: func(mockUploads *MockUploadsService, file *multipart.FileHeader) {
				mockUploads.On("UploadImage", mock.Anything, file, "products").Return("products/product.jpg", nil)
			},
			expectedResult: "products/product.jpg",
			expectedError:  nil,
		},
		{
			name: "error - upload service returns error",
			file: &multipart.FileHeader{
				Filename: "product.jpg",
				Size:     2048,
			},
			mockSetup: func(mockUploads *MockUploadsService, file *multipart.FileHeader) {
				mockUploads.On("UploadImage", mock.Anything, file, "products").Return("", errors.New("upload failed"))
			},
			expectedResult: "",
			expectedError:  errors.New("upload failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockUploads, tt.file)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.UploadProductImage(ctx, tt.file)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockUploads.AssertExpectations(t)
		})
	}
}

func TestService_GetProductImageURL(t *testing.T) {
	tests := []struct {
		name           string
		objectName     string
		mockSetup      func(*MockUploadsService, string)
		expectedResult string
		expectedError  error
	}{
		{
			name:       "success - returns product image URL",
			objectName: "products/product.jpg",
			mockSetup: func(mockUploads *MockUploadsService, objectName string) {
				mockUploads.On("GetImageURL", mock.Anything, objectName).Return("https://example.com/products/product.jpg", nil)
			},
			expectedResult: "https://example.com/products/product.jpg",
			expectedError:  nil,
		},
		{
			name:       "error - upload service returns error",
			objectName: "products/product.jpg",
			mockSetup: func(mockUploads *MockUploadsService, objectName string) {
				mockUploads.On("GetImageURL", mock.Anything, objectName).Return("", errors.New("failed to generate URL"))
			},
			expectedResult: "",
			expectedError:  errors.New("failed to generate URL"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockUploads, tt.objectName)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.GetProductImageURL(ctx, tt.objectName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockUploads.AssertExpectations(t)
		})
	}
}
