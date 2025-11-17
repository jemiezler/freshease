package orders

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

func (m *MockRepository) Create(ctx context.Context, u *CreateOrderDTO) (*GetOrderDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrderDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateOrderDTO) (*GetOrderDTO, error) {
	args := m.Called(ctx, u)
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
				userID := uuid.New()
				mockRepo.On("List", context.Background()).Return([]*GetOrderDTO{
					{
						ID:      uuid.New(),
						OrderNo: "ORD-001",
						Status:  "pending",
						Total:   110.0,
						UserID:  userID,
					},
					{
						ID:      uuid.New(),
						OrderNo: "ORD-002",
						Status:  "completed",
						Total:   210.0,
						UserID:  userID,
					},
				}, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetOrderDTO(nil), errors.New("database error"))
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
		orderID       uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:    "success - returns order by ID",
			orderID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(&GetOrderDTO{
					ID:      id,
					OrderNo: "ORD-001",
					Status:  "pending",
					Total:   110.0,
					UserID:  uuid.New(),
				}, nil)
			},
			expectedError: false,
		},
		{
			name:    "error - order not found",
			orderID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.orderID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.orderID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.orderID, result.ID)
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
		dto           CreateOrderDTO
		mockSetup     func(*MockRepository, CreateOrderDTO)
		expectedError bool
	}{
		{
			name: "success - creates order",
			dto: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-003",
				Status:      "pending",
				Subtotal:    150.0,
				ShippingFee: 15.0,
				Discount:    5.0,
				Total:       160.0,
				PlacedAt:    &now,
				UserID:      userID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrderDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetOrderDTO{
					ID:      dto.ID,
					OrderNo: dto.OrderNo,
					Status:  dto.Status,
					Total:   dto.Total,
					UserID:  dto.UserID,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateOrderDTO{
				ID:          uuid.New(),
				OrderNo:     "ORD-004",
				Status:      "pending",
				Subtotal:    150.0,
				ShippingFee: 15.0,
				Discount:    5.0,
				Total:       160.0,
				UserID:      userID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrderDTO) {
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
				assert.Equal(t, tt.dto.OrderNo, result.OrderNo)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	orderID := uuid.New()
	newStatus := "completed"
	newTotal := 120.0

	tests := []struct {
		name          string
		dto           UpdateOrderDTO
		mockSetup     func(*MockRepository, UpdateOrderDTO)
		expectedError bool
	}{
		{
			name: "success - updates order",
			dto: UpdateOrderDTO{
				ID:     orderID,
				Status: &newStatus,
				Total:  &newTotal,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateOrderDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*orders.UpdateOrderDTO")).Return(&GetOrderDTO{
					ID:     orderID,
					Status: newStatus,
					Total:  newTotal,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdateOrderDTO{
				ID:     orderID,
				Status: &newStatus,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateOrderDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*orders.UpdateOrderDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), orderID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, orderID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	orderID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes order",
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
			tt.mockSetup(mockRepo, orderID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), orderID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

