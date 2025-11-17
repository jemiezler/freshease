package deliveries

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

func (m *MockRepository) List(ctx context.Context) ([]*GetDeliveryDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetDeliveryDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetDeliveryDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDeliveryDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateDeliveryDTO) (*GetDeliveryDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDeliveryDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateDeliveryDTO) (*GetDeliveryDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetDeliveryDTO), args.Error(1)
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
			name: "success - returns deliveries list",
			mockSetup: func(mockRepo *MockRepository) {
				orderID := uuid.New()
				trackingNo := "TRACK-001"
				mockRepo.On("List", context.Background()).Return([]*GetDeliveryDTO{
					{
						ID:         uuid.New(),
						Provider:   "fedex",
						TrackingNo: &trackingNo,
						Status:     "pending",
						OrderID:    orderID,
					},
					{
						ID:       uuid.New(),
						Provider: "ups",
						Status:   "in_transit",
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
				mockRepo.On("List", context.Background()).Return([]*GetDeliveryDTO(nil), errors.New("database error"))
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
		deliveryID    uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:       "success - returns delivery by ID",
			deliveryID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				orderID := uuid.New()
				trackingNo := "TRACK-001"
				mockRepo.On("FindByID", context.Background(), id).Return(&GetDeliveryDTO{
					ID:         id,
					Provider:   "fedex",
					TrackingNo: &trackingNo,
					Status:     "pending",
					OrderID:    orderID,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:       "error - delivery not found",
			deliveryID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.deliveryID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.deliveryID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.deliveryID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	orderID := uuid.New()
	trackingNo := "TRACK-002"
	eta := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name          string
		dto           CreateDeliveryDTO
		mockSetup     func(*MockRepository, CreateDeliveryDTO)
		expectedError bool
	}{
		{
			name: "success - creates delivery",
			dto: CreateDeliveryDTO{
				ID:         uuid.New(),
				Provider:   "fedex",
				TrackingNo: &trackingNo,
				Status:     "pending",
				Eta:        &eta,
				OrderID:    orderID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateDeliveryDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetDeliveryDTO{
					ID:         dto.ID,
					Provider:   dto.Provider,
					TrackingNo: dto.TrackingNo,
					Status:     dto.Status,
					Eta:        dto.Eta,
					OrderID:    dto.OrderID,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateDeliveryDTO{
				ID:       uuid.New(),
				Provider: "fedex",
				Status:   "pending",
				OrderID:  orderID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateDeliveryDTO) {
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
	deliveryID := uuid.New()
	newStatus := "in_transit"
	newTrackingNo := "TRACK-UPDATED"

	tests := []struct {
		name          string
		dto           UpdateDeliveryDTO
		mockSetup     func(*MockRepository, UpdateDeliveryDTO)
		expectedError bool
	}{
		{
			name: "success - updates delivery",
			dto: UpdateDeliveryDTO{
				ID:         deliveryID,
				Status:     &newStatus,
				TrackingNo: &newTrackingNo,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateDeliveryDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*deliveries.UpdateDeliveryDTO")).Return(&GetDeliveryDTO{
					ID:         deliveryID,
					Status:     newStatus,
					TrackingNo: &newTrackingNo,
					Provider:   "fedex",
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdateDeliveryDTO{
				ID:     deliveryID,
				Status: &newStatus,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateDeliveryDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*deliveries.UpdateDeliveryDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), deliveryID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, deliveryID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	deliveryID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes delivery",
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
			tt.mockSetup(mockRepo, deliveryID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), deliveryID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

