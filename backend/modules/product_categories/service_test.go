package product_categories

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context) ([]*GetProductCategoryDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetProductCategoryDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductCategoryDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductCategoryDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetProductCategoryDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockRepository)
		expectedError bool
		expectedCount int
	}{
		{
			name: "success - returns list of product categories",
			mockSetup: func(mockRepo *MockRepository) {
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
				mockRepo.On("List", mock.Anything).Return(categories, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetProductCategoryDTO{}, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetProductCategoryDTO)(nil), errors.New("database error"))
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.List(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - returns product category",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				category := &GetProductCategoryDTO{
					ID:          id,
					Name:        "Fruits",
					Description: "Fresh fruits and vegetables",
					Slug:        "fruits",
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(category, nil)
			},
			expectedError: false,
		},
		{
			name: "error - product category not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetProductCategoryDTO)(nil), errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Get(ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		dto           CreateProductCategoryDTO
		mockSetup     func(*MockRepository, CreateProductCategoryDTO)
		expectedError bool
	}{
		{
			name: "success - creates product category",
			dto: CreateProductCategoryDTO{
				ID:          uuid.New(),
				Name:        "Fruits",
				Description: "Fresh fruits and vegetables",
				Slug:        "fruits",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateProductCategoryDTO) {
				createdCategory := &GetProductCategoryDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Description: dto.Description,
					Slug:        dto.Slug,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateProductCategoryDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description &&
						actual.Slug == dto.Slug
				})).Return(createdCategory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateProductCategoryDTO{
				ID:          uuid.New(),
				Name:        "Fruits",
				Description: "Fresh fruits and vegetables",
				Slug:        "fruits",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateProductCategoryDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateProductCategoryDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description &&
						actual.Slug == dto.Slug
				})).Return((*GetProductCategoryDTO)(nil), errors.New("creation failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.dto.ID, result.ID)
				assert.Equal(t, tt.dto.Name, result.Name)
				assert.Equal(t, tt.dto.Description, result.Description)
				assert.Equal(t, tt.dto.Slug, result.Slug)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		dto           UpdateProductCategoryDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdateProductCategoryDTO)
		expectedError bool
	}{
		{
			name: "success - updates product category",
			id:   uuid.New(),
			dto: UpdateProductCategoryDTO{
				Name:        stringPtr("Updated Fruits"),
				Description: stringPtr("Updated description"),
				Slug:        stringPtr("updated-fruits"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductCategoryDTO) {
				updatedCategory := &GetProductCategoryDTO{
					ID:          id,
					Name:        "Updated Fruits",
					Description: "Updated description",
					Slug:        "updated-fruits",
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateProductCategoryDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Fruits" &&
						actual.Description != nil && *actual.Description == "Updated description" &&
						actual.Slug != nil && *actual.Slug == "updated-fruits"
				})).Return(updatedCategory, nil)
			},
			expectedError: false,
		},
		{
			name: "success - partial update",
			id:   uuid.New(),
			dto: UpdateProductCategoryDTO{
				Name: stringPtr("Updated Fruits"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductCategoryDTO) {
				updatedCategory := &GetProductCategoryDTO{
					ID:          id,
					Name:        "Updated Fruits",
					Description: "Original description", // unchanged
					Slug:        "original-slug",        // unchanged
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateProductCategoryDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Fruits" &&
						actual.Description == nil &&
						actual.Slug == nil
				})).Return(updatedCategory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateProductCategoryDTO{
				Name: stringPtr("Updated Fruits"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductCategoryDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateProductCategoryDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Fruits"
				})).Return((*GetProductCategoryDTO)(nil), errors.New("update failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id, tt.dto)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Update(ctx, tt.id, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes product category",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
