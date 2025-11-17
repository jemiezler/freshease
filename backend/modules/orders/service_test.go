package orders

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

func (m *MockRepository) List(ctx context.Context) ([]*GetOrderDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetOrderDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetOrderDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateOrderDTO) (*GetOrderDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateOrderDTO) (*GetOrderDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
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
			name: "success - returns orders list",
			mockSetup: func(mockRepo *MockRepository) {
				orders := []*GetOrderDTO{
					{
						ID:          uuid.New(),
						OrderNo:     "ORD-001",
						Status:      "pending",
						Subtotal:    100.00,
						ShippingFee: 10.00,
						Discount:    5.00,
						Total:       105.00,
						UserID:      uuid.New(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						OrderNo:     "ORD-002",
						Status:      "completed",
						Subtotal:    200.00,
						ShippingFee: 15.00,
						Discount:    10.00,
						Total:       205.00,
						UserID:      uuid.New(),
						UpdatedAt:   time.Now(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(orders, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetOrderDTO)(nil), errors.New("database error"))
			},
			expectedCount: 0,
			expectedError: true,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetOrderDTO{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			ctx := context.Background()

			orders, err := service.List(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, orders)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, orders)
				assert.Len(t, orders, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name          string
		orderID       uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedOrder *GetOrderDTO
		expectedError bool
	}{
		{
			name:   "success - returns order by ID",
			orderID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				order := &GetOrderDTO{
					ID:          id,
					OrderNo:     "ORD-003",
					Status:      "pending",
					Subtotal:    150.00,
					ShippingFee: 12.00,
					Discount:    8.00,
					Total:       154.00,
					UserID:      uuid.New(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(order, nil)
			},
			expectedError: false,
		},
		{
			name:   "error - order not found",
			orderID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetOrderDTO)(nil), errors.New("order not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.orderID)

			service := NewService(mockRepo)
			ctx := context.Background()

			order, err := service.Get(ctx, tt.orderID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.orderID, order.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name          string
		createDTO     CreateOrderDTO
		mockSetup     func(*MockRepository, CreateOrderDTO)
		expectedError bool
	}{
		{
			name: "success - creates new order",
			createDTO: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-004",
				Status:      "pending",
				Subtotal:    200.00,
				ShippingFee: 15.00,
				Discount:    10.00,
				Total:       205.00,
				UserID:      userID,
				PlacedAt:    &now,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrderDTO) {
				expectedOrder := &GetOrderDTO{
					ID:          dto.ID,
					OrderNo:     dto.OrderNo,
					Status:      dto.Status,
					Subtotal:    dto.Subtotal,
					ShippingFee: dto.ShippingFee,
					Discount:    dto.Discount,
					Total:       dto.Total,
					UserID:      dto.UserID,
					PlacedAt:   dto.PlacedAt,
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateOrderDTO) bool {
					return actual.ID == dto.ID && actual.OrderNo == dto.OrderNo && actual.UserID == dto.UserID
				})).Return(expectedOrder, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-005",
				Status:      "pending",
				Subtotal:    100.00,
				ShippingFee: 10.00,
				Discount:    0.00,
				Total:       110.00,
				UserID:      userID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrderDTO) {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return((*GetOrderDTO)(nil), errors.New("creation failed"))
			},
			expectedError: true,
		},
		{
			name: "error - negative total",
			createDTO: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-006",
				Status:      "pending",
				Subtotal:    100.00,
				ShippingFee: 10.00,
				Discount:    120.00, // Discount exceeds subtotal + shipping
				Total:       -10.00,
				UserID:      userID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrderDTO) {
				// This should be caught by validation, but if it reaches repository, it should fail
				mockRepo.On("Create", mock.Anything, mock.Anything).Return((*GetOrderDTO)(nil), errors.New("invalid total"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			order, err := service.Create(ctx, tt.createDTO)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.createDTO.ID, order.ID)
				assert.Equal(t, tt.createDTO.OrderNo, order.OrderNo)
				assert.Equal(t, tt.createDTO.Total, order.Total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		orderID       uuid.UUID
		updateDTO     UpdateOrderDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdateOrderDTO)
		expectedError bool
	}{
		{
			name:    "success - updates order status",
			orderID: orderID,
			updateDTO: UpdateOrderDTO{
				Status: stringPtr("completed"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateOrderDTO) {
				expectedOrder := &GetOrderDTO{
					ID:          id,
					OrderNo:     "ORD-007",
					Status:      *dto.Status,
					Subtotal:    200.00,
					ShippingFee: 15.00,
					Discount:    10.00,
					Total:       205.00,
					UserID:      userID,
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateOrderDTO) bool {
					return actual.ID == id && actual.Status != nil && *actual.Status == *dto.Status
				})).Return(expectedOrder, nil)
			},
			expectedError: false,
		},
		{
			name:    "success - updates order total",
			orderID: orderID,
			updateDTO: UpdateOrderDTO{
				Total: float64Ptr(250.00),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateOrderDTO) {
				expectedOrder := &GetOrderDTO{
					ID:          id,
					OrderNo:     "ORD-008",
					Status:      "pending",
					Subtotal:    230.00,
					ShippingFee: 20.00,
					Discount:    0.00,
					Total:       *dto.Total,
					UserID:      userID,
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateOrderDTO) bool {
					return actual.ID == id && actual.Total != nil && *actual.Total == *dto.Total
				})).Return(expectedOrder, nil)
			},
			expectedError: false,
		},
		{
			name:    "error - repository returns error",
			orderID: orderID,
			updateDTO: UpdateOrderDTO{
				Status: stringPtr("cancelled"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateOrderDTO) {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return((*GetOrderDTO)(nil), errors.New("order not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.orderID, tt.updateDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			order, err := service.Update(ctx, tt.orderID, tt.updateDTO)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.orderID, order.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		orderID       uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:    "success - deletes order",
			orderID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "error - repository returns error",
			orderID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("order not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.orderID)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.orderID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
