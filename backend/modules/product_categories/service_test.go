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
						ID:         uuid.New(),
						ProductID:  uuid.New(),
						CategoryID: uuid.New(),
					},
					{
						ID:         uuid.New(),
						ProductID:  uuid.New(),
						CategoryID: uuid.New(),
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
					ID:         id,
					ProductID:  uuid.New(),
					CategoryID: uuid.New(),
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
				ID:         uuid.New(),
				ProductID:  uuid.New(),
				CategoryID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateProductCategoryDTO) {
				createdCategory := &GetProductCategoryDTO{
					ID:         dto.ID,
					ProductID:  dto.ProductID,
					CategoryID: dto.CategoryID,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateProductCategoryDTO) bool {
					return actual.ID == dto.ID &&
						actual.ProductID == dto.ProductID &&
						actual.CategoryID == dto.CategoryID
				})).Return(createdCategory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateProductCategoryDTO{
				ID:         uuid.New(),
				ProductID:  uuid.New(),
				CategoryID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateProductCategoryDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateProductCategoryDTO) bool {
					return actual.ID == dto.ID &&
						actual.ProductID == dto.ProductID &&
						actual.CategoryID == dto.CategoryID
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
				assert.Equal(t, tt.dto.ProductID, result.ProductID)
				assert.Equal(t, tt.dto.CategoryID, result.CategoryID)
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
				ID:         uuid.New(),
				ProductID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
				CategoryID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductCategoryDTO) {
				updatedCategory := &GetProductCategoryDTO{
					ID:         id,
					ProductID:  uuid.New(),
					CategoryID: uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateProductCategoryDTO) bool {
					return actual.ID == id &&
						actual.ProductID != nil && actual.CategoryID != nil
				})).Return(updatedCategory, nil)
			},
			expectedError: false,
		},
		{
			name: "success - partial update",
			id:   uuid.New(),
			dto: UpdateProductCategoryDTO{
				ID:        uuid.New(),
				ProductID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductCategoryDTO) {
				updatedCategory := &GetProductCategoryDTO{
					ID:         id,
					ProductID:  uuid.New(),
					CategoryID: uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateProductCategoryDTO) bool {
					return actual.ID == id &&
						actual.ProductID != nil && actual.CategoryID == nil
				})).Return(updatedCategory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateProductCategoryDTO{
				ID:        uuid.New(),
				ProductID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateProductCategoryDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateProductCategoryDTO) bool {
					return actual.ID == id &&
						actual.ProductID != nil
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
