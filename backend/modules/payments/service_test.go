package payments

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

func (m *MockRepository) List(ctx context.Context) ([]*GetPaymentDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetPaymentDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetPaymentDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPaymentDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreatePaymentDTO) (*GetPaymentDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPaymentDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdatePaymentDTO) (*GetPaymentDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPaymentDTO), args.Error(1)
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
			name: "success - returns payments list",
			mockSetup: func(mockRepo *MockRepository) {
				orderID := uuid.New()
				providerRef := "pay_001"
				mockRepo.On("List", context.Background()).Return([]*GetPaymentDTO{
					{
						ID:          uuid.New(),
						Provider:    "stripe",
						ProviderRef: &providerRef,
						Status:      "pending",
						Amount:      110.0,
						OrderID:     orderID,
					},
					{
						ID:       uuid.New(),
						Provider: "paypal",
						Status:   "completed",
						Amount:   110.0,
						OrderID:  orderID,
					},
				}, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetPaymentDTO(nil), errors.New("database error"))
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
		paymentID     uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:      "success - returns payment by ID",
			paymentID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				orderID := uuid.New()
				providerRef := "pay_001"
				mockRepo.On("FindByID", context.Background(), id).Return(&GetPaymentDTO{
					ID:          id,
					Provider:    "stripe",
					ProviderRef: &providerRef,
					Status:      "pending",
					Amount:      110.0,
					OrderID:     orderID,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:      "error - payment not found",
			paymentID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.paymentID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.paymentID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.paymentID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	orderID := uuid.New()
	providerRef := "pay_002"
	paidAt := time.Now()
	tests := []struct {
		name          string
		dto           CreatePaymentDTO
		mockSetup     func(*MockRepository, CreatePaymentDTO)
		expectedError bool
	}{
		{
			name: "success - creates payment",
			dto: CreatePaymentDTO{
				ID:          uuid.New(),
				Provider:    "stripe",
				ProviderRef: &providerRef,
				Status:      "completed",
				Amount:      110.0,
				PaidAt:      &paidAt,
				OrderID:     orderID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreatePaymentDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetPaymentDTO{
					ID:          dto.ID,
					Provider:    dto.Provider,
					ProviderRef: dto.ProviderRef,
					Status:      dto.Status,
					Amount:      dto.Amount,
					PaidAt:      dto.PaidAt,
					OrderID:     dto.OrderID,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreatePaymentDTO{
				ID:       uuid.New(),
				Provider: "stripe",
				Status:   "pending",
				Amount:   110.0,
				OrderID:  orderID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreatePaymentDTO) {
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
				assert.Equal(t, tt.dto.Provider, result.Provider)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	paymentID := uuid.New()
	newStatus := "completed"
	newProviderRef := "pay_updated"

	tests := []struct {
		name          string
		dto           UpdatePaymentDTO
		mockSetup     func(*MockRepository, UpdatePaymentDTO)
		expectedError bool
	}{
		{
			name: "success - updates payment",
			dto: UpdatePaymentDTO{
				ID:          paymentID,
				Status:      &newStatus,
				ProviderRef: &newProviderRef,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdatePaymentDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*payments.UpdatePaymentDTO")).Return(&GetPaymentDTO{
					ID:          paymentID,
					Status:      newStatus,
					ProviderRef: &newProviderRef,
					Provider:    "stripe",
					Amount:      110.0,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdatePaymentDTO{
				ID:     paymentID,
				Status: &newStatus,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdatePaymentDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*payments.UpdatePaymentDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), paymentID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, paymentID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	paymentID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes payment",
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
			tt.mockSetup(mockRepo, paymentID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), paymentID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

