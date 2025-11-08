package inventories

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

func (m *MockRepository) List(ctx context.Context) ([]*GetInventoryDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetInventoryDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetInventoryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetInventoryDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateInventoryDTO) (*GetInventoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetInventoryDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateInventoryDTO) (*GetInventoryDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetInventoryDTO), args.Error(1)
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
			name: "success - returns list of inventories",
			mockSetup: func(mockRepo *MockRepository) {
				inventories := []*GetInventoryDTO{
					{
						ID:            uuid.New(),
						Quantity:      100,
						ReorderLevel: 50,
						UpdatedAt:     time.Now(),
					},
					{
						ID:            uuid.New(),
						Quantity:      200,
						ReorderLevel: 75,
						UpdatedAt:     time.Now(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(inventories, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetInventoryDTO{}, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetInventoryDTO)(nil), errors.New("database error"))
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
			name: "success - returns inventory",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				inventory := &GetInventoryDTO{
					ID:            id,
					Quantity:      150,
					ReorderLevel: 60,
					UpdatedAt:     time.Now(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(inventory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - inventory not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetInventoryDTO)(nil), errors.New("not found"))
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
		dto           CreateInventoryDTO
		mockSetup     func(*MockRepository, CreateInventoryDTO)
		expectedError bool
	}{
		{
			name: "success - creates inventory",
			dto: CreateInventoryDTO{
				Quantity:      100,
				ReorderLevel: 50,
				UpdatedAt:     time.Now(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateInventoryDTO) {
				createdInventory := &GetInventoryDTO{
					ID:            uuid.New(),
					Quantity:      dto.Quantity,
					ReorderLevel: dto.ReorderLevel,
					UpdatedAt:     time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateInventoryDTO) bool {
					return actual.Quantity == dto.Quantity &&
						actual.ReorderLevel == dto.ReorderLevel
				})).Return(createdInventory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateInventoryDTO{
				Quantity:      100,
				ReorderLevel: 50,
				UpdatedAt:     time.Now(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateInventoryDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateInventoryDTO) bool {
					return actual.Quantity == dto.Quantity &&
						actual.ReorderLevel == dto.ReorderLevel
				})).Return((*GetInventoryDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.Quantity, result.Quantity)
				assert.Equal(t, tt.dto.ReorderLevel, result.ReorderLevel)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		dto           UpdateInventoryDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdateInventoryDTO)
		expectedError bool
	}{
		{
			name: "success - updates inventory",
			id:   uuid.New(),
			dto: UpdateInventoryDTO{
				Quantity:      intPtr(200),
				ReorderLevel: intPtr(75),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateInventoryDTO) {
				updatedInventory := &GetInventoryDTO{
					ID:            id,
					Quantity:      200,
					ReorderLevel: 75,
					UpdatedAt:     time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateInventoryDTO) bool {
					return actual.ID == id &&
						actual.Quantity != nil && *actual.Quantity == 200 &&
						actual.ReorderLevel != nil && *actual.ReorderLevel == 75
				})).Return(updatedInventory, nil)
			},
			expectedError: false,
		},
		{
			name: "success - partial update",
			id:   uuid.New(),
			dto: UpdateInventoryDTO{
				Quantity: intPtr(300),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateInventoryDTO) {
				updatedInventory := &GetInventoryDTO{
					ID:            id,
					Quantity:      300,
					ReorderLevel: 50, // unchanged
					UpdatedAt:     time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateInventoryDTO) bool {
					return actual.ID == id &&
						actual.Quantity != nil && *actual.Quantity == 300 &&
						actual.ReorderLevel == nil
				})).Return(updatedInventory, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateInventoryDTO{
				Quantity: intPtr(200),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateInventoryDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateInventoryDTO) bool {
					return actual.ID == id &&
						actual.Quantity != nil && *actual.Quantity == 200
				})).Return((*GetInventoryDTO)(nil), errors.New("update failed"))
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
			name: "success - deletes inventory",
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
