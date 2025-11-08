package carts

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

func (m *MockRepository) List(ctx context.Context) ([]*GetCartDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetCartDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetCartDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCartDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateCartDTO) (*GetCartDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCartDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateCartDTO) (*GetCartDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCartDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockRepository)
		expectedCarts []*GetCartDTO
		expectedError bool
	}{
		{
			name: "success - returns carts list",
			mockSetup: func(mockRepo *MockRepository) {
				carts := []*GetCartDTO{
					{
						ID:        uuid.New(),
						Status:    "pending",
						Subtotal:  100.50,
						Discount:  0.0,
						Total:     100.50,
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Status:    "completed",
						Subtotal:  250.75,
						Discount:  10.0,
						Total:     240.75,
						UpdatedAt: time.Now(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(carts, nil)
			},
			expectedCarts: []*GetCartDTO{
				{
					ID:        uuid.New(),
					Status:    "pending",
					Subtotal:  100.50,
					Discount:  0.0,
					Total:     100.50,
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.New(),
					Status:    "completed",
					Subtotal:  250.75,
					Discount:  10.0,
					Total:     240.75,
					UpdatedAt: time.Now(),
				},
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetCartDTO)(nil), errors.New("database error"))
			},
			expectedCarts: nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			ctx := context.Background()

			carts, err := service.List(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, carts)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, carts)
				assert.Len(t, carts, 2)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name          string
		cartID        uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedCart  *GetCartDTO
		expectedError bool
	}{
		{
			name:   "success - returns cart by ID",
			cartID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				cart := &GetCartDTO{
					ID:        id,
					Status:    "pending",
					Total:     150.25,
					Subtotal:  100.50,
					Discount:  0.0,
					UpdatedAt: time.Now(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(cart, nil)
			},
			expectedCart: &GetCartDTO{
				ID:        uuid.New(),
				Status:    "pending",
				Total:     150.25,
				Subtotal:  100.50,
				Discount:  0.0,
				UpdatedAt: time.Now(),
			},
			expectedError: false,
		},
		{
			name:   "error - cart not found",
			cartID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetCartDTO)(nil), errors.New("cart not found"))
			},
			expectedCart:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.cartID)

			service := NewService(mockRepo)
			ctx := context.Background()

			cart, err := service.Get(ctx, tt.cartID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, cart)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cart)
				assert.Equal(t, tt.cartID, cart.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		createDTO     CreateCartDTO
		mockSetup     func(*MockRepository, CreateCartDTO)
		expectedCart  *GetCartDTO
		expectedError bool
	}{
		{
			name: "success - creates new cart",
			createDTO: CreateCartDTO{
				Status: stringPtr("pending"),
				Total:  float64Ptr(99.99),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateCartDTO) {
				expectedCart := &GetCartDTO{
					ID:        uuid.New(),
					Status:    *dto.Status,
					Total:     *dto.Total,
					Subtotal:  100.50,
					Discount:  0.0,
					UpdatedAt: time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateCartDTO) bool {
					return actual.Status != nil && *actual.Status == *dto.Status &&
						actual.Total != nil && *actual.Total == *dto.Total
				})).Return(expectedCart, nil)
			},
			expectedCart: &GetCartDTO{
				ID:        uuid.New(),
				Status:    "pending",
				Total:     99.99,
				Subtotal:  100.50,
				Discount:  0.0,
				UpdatedAt: time.Now(),
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateCartDTO{
				Status: stringPtr("pending"),
				Total:  float64Ptr(50.00),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateCartDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateCartDTO) bool {
					return actual.Status != nil && *actual.Status == *dto.Status &&
						actual.Total != nil && *actual.Total == *dto.Total
				})).Return((*GetCartDTO)(nil), errors.New("creation failed"))
			},
			expectedCart:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			cart, err := service.Create(ctx, tt.createDTO)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, cart)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cart)
				assert.Equal(t, *tt.createDTO.Status, cart.Status)
				assert.Equal(t, *tt.createDTO.Total, cart.Total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		cartID        uuid.UUID
		updateDTO     UpdateCartDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdateCartDTO)
		expectedCart  *GetCartDTO
		expectedError bool
	}{
		{
			name:   "success - updates cart",
			cartID: uuid.New(),
			updateDTO: UpdateCartDTO{
				Status: stringPtr("completed"),
				Total:  float64Ptr(200.00),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateCartDTO) {
				expectedCart := &GetCartDTO{
					ID:        id,
					Status:    *dto.Status,
					Total:     *dto.Total,
					Subtotal:  100.50,
					Discount:  0.0,
					UpdatedAt: time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateCartDTO) bool {
					return actual.ID == id &&
						actual.Status != nil && *actual.Status == *dto.Status &&
						actual.Total != nil && *actual.Total == *dto.Total
				})).Return(expectedCart, nil)
			},
			expectedCart: &GetCartDTO{
				ID:        uuid.New(),
				Status:    "completed",
				Total:     200.00,
				Subtotal:  100.50,
				Discount:  0.0,
				UpdatedAt: time.Now(),
			},
			expectedError: false,
		},
		{
			name:   "error - repository returns error",
			cartID: uuid.New(),
			updateDTO: UpdateCartDTO{
				Status: stringPtr("cancelled"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateCartDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateCartDTO) bool {
					return actual.ID == id &&
						actual.Status != nil && *actual.Status == *dto.Status
				})).Return((*GetCartDTO)(nil), errors.New("update failed"))
			},
			expectedCart:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.cartID, tt.updateDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			cart, err := service.Update(ctx, tt.cartID, tt.updateDTO)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, cart)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cart)
				assert.Equal(t, tt.cartID, cart.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		cartID        uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:   "success - deletes cart",
			cartID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "error - repository returns error",
			cartID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.cartID)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.cartID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
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
