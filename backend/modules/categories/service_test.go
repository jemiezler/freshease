package categories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context) ([]*GetCategoryDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetCategoryDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetCategoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCategoryDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateCategoryDTO) (*GetCategoryDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCategoryDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateCategoryDTO) (*GetCategoryDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCategoryDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name: "success - returns categories list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetCategoryDTO{
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
				}, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetCategoryDTO(nil), errors.New("database error"))
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := NewService(mockRepo)
			result, err := svc.List(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name          string
		categoryID    uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:       "success - returns category by ID",
			categoryID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(&GetCategoryDTO{
					ID:   id,
					Name: "Test Category",
					Slug: "test-category",
				}, nil)
			},
			expectedError: false,
		},
		{
			name:       "error - category not found",
			categoryID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.categoryID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.categoryID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.categoryID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		dto           CreateCategoryDTO
		mockSetup     func(*MockRepository, CreateCategoryDTO)
		expectedError bool
	}{
		{
			name: "success - creates category",
			dto: CreateCategoryDTO{
				ID:        uuid.New(),
				Name:      "New Category",
				Slug:      "new-category",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateCategoryDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetCategoryDTO{
					ID:   dto.ID,
					Name: dto.Name,
					Slug: dto.Slug,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateCategoryDTO{
				ID:        uuid.New(),
				Name:      "New Category",
				Slug:      "new-category",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateCategoryDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Create(context.Background(), tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.dto.ID, result.ID)
				assert.Equal(t, tt.dto.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	categoryID := uuid.New()
	newName := "Updated Category"
	newSlug := "updated-category"

	tests := []struct {
		name          string
		dto           UpdateCategoryDTO
		mockSetup     func(*MockRepository, UpdateCategoryDTO)
		expectedError bool
	}{
		{
			name: "success - updates category",
			dto: UpdateCategoryDTO{
				ID:        categoryID,
				Name:      &newName,
				Slug:      &newSlug,
				UpdatedAt: time.Now(),
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateCategoryDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*categories.UpdateCategoryDTO")).Return(&GetCategoryDTO{
					ID:   categoryID,
					Name: newName,
					Slug: newSlug,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdateCategoryDTO{
				ID:        categoryID,
				Name:      &newName,
				UpdatedAt: time.Now(),
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateCategoryDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*categories.UpdateCategoryDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), categoryID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, categoryID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes category",
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", context.Background(), id).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", context.Background(), id).Return(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, categoryID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), categoryID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

