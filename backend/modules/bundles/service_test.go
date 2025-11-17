package bundles

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

func (m *MockRepository) List(ctx context.Context) ([]*GetBundleDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetBundleDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetBundleDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundleDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateBundleDTO) (*GetBundleDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundleDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateBundleDTO) (*GetBundleDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundleDTO), args.Error(1)
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
			name: "success - returns bundles list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetBundleDTO{
					{
						ID:       uuid.New(),
						Name:     "Bundle One",
						Price:    99.99,
						IsActive: true,
					},
					{
						ID:       uuid.New(),
						Name:     "Bundle Two",
						Price:    149.99,
						IsActive: true,
					},
				}, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetBundleDTO(nil), errors.New("database error"))
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
		bundleID      uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:     "success - returns bundle by ID",
			bundleID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(&GetBundleDTO{
					ID:       id,
					Name:     "Test Bundle",
					Price:    99.99,
					IsActive: true,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:     "error - bundle not found",
			bundleID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.bundleID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.bundleID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.bundleID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		dto           CreateBundleDTO
		mockSetup     func(*MockRepository, CreateBundleDTO)
		expectedError bool
	}{
		{
			name: "success - creates bundle",
			dto: CreateBundleDTO{
				ID:       uuid.New(),
				Name:     "New Bundle",
				Price:    199.99,
				IsActive: true,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateBundleDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetBundleDTO{
					ID:       dto.ID,
					Name:     dto.Name,
					Price:    dto.Price,
					IsActive: dto.IsActive,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateBundleDTO{
				ID:       uuid.New(),
				Name:     "New Bundle",
				Price:    199.99,
				IsActive: true,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateBundleDTO) {
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
	bundleID := uuid.New()
	newName := "Updated Bundle"
	newPrice := 299.99

	tests := []struct {
		name          string
		dto           UpdateBundleDTO
		mockSetup     func(*MockRepository, UpdateBundleDTO)
		expectedError bool
	}{
		{
			name: "success - updates bundle",
			dto: UpdateBundleDTO{
				ID:    bundleID,
				Name:  &newName,
				Price: &newPrice,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateBundleDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*bundles.UpdateBundleDTO")).Return(&GetBundleDTO{
					ID:       bundleID,
					Name:     newName,
					Price:    newPrice,
					IsActive: true,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdateBundleDTO{
				ID:    bundleID,
				Name:  &newName,
				Price: &newPrice,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateBundleDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*bundles.UpdateBundleDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), bundleID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				require.NotNil(t, result)
				assert.Equal(t, bundleID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	bundleID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes bundle",
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
			tt.mockSetup(mockRepo, bundleID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), bundleID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

