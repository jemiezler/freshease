package cart_items

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context) ([]*GetCart_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetCart_itemDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetCart_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCart_itemDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateCart_itemDTO) (*GetCart_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCart_itemDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateCart_itemDTO) (*GetCart_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCart_itemDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockRepository)
		expectedResult []*GetCart_itemDTO
		expectedError  error
	}{
		{
			name: "success - returns cart items list",
			mockSetup: func(mockRepo *MockRepository) {
				expectedItems := []*GetCart_itemDTO{
					{
						ID:        uuid.New(),
						Qty:       2,
						UnitPrice: 1.99,
						LineTotal: 3.98,
						CartID:    uuid.New(),
						ProductID: uuid.New(),
					},
					{
						ID:        uuid.New(),
						Qty:       3,
						UnitPrice: 0.99,
						LineTotal: 2.97,
						CartID:    uuid.New(),
						ProductID: uuid.New(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedResult: []*GetCart_itemDTO{
				{
					ID:        uuid.New(),
					Qty:       2,
					UnitPrice: 1.99,
					LineTotal: 3.98,
					CartID:    uuid.New(),
					ProductID: uuid.New(),
				},
				{
					ID:        uuid.New(),
					Qty:       3,
					UnitPrice: 0.99,
					LineTotal: 2.97,
					CartID:    uuid.New(),
					ProductID: uuid.New(),
				},
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetCart_itemDTO(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.List(ctx)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name           string
		cartItemID     uuid.UUID
		mockSetup      func(*MockRepository, uuid.UUID)
		expectedResult *GetCart_itemDTO
		expectedError  error
	}{
		{
			name:       "success - returns cart item by ID",
			cartItemID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				expectedItem := &GetCart_itemDTO{
					ID:        id,
					Qty:       2,
					UnitPrice: 1.99,
					LineTotal: 3.98,
					CartID:    uuid.New(),
					ProductID: uuid.New(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedItem, nil)
			},
			expectedResult: &GetCart_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 1.99,
				LineTotal: 3.98,
				CartID:    uuid.New(),
				ProductID: uuid.New(),
			},
			expectedError: nil,
		},
		{
			name:       "error - cart item not found",
			cartItemID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetCart_itemDTO)(nil), errors.New("cart item not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("cart item not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.cartItemID)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Get(ctx, tt.cartItemID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.cartItemID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name           string
		createDTO      CreateCart_itemDTO
		mockSetup      func(*MockRepository, CreateCart_itemDTO)
		expectedResult *GetCart_itemDTO
		expectedError  error
	}{
		{
			name: "success - creates new cart item",
			createDTO: CreateCart_itemDTO{
				ID:        uuid.New(),
				Qty:       5,
				UnitPrice: 2.49,
				LineTotal: 12.45,
				CartID:    uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateCart_itemDTO) {
				expectedItem := &GetCart_itemDTO{
					ID:        dto.ID,
					Qty:       dto.Qty,
					UnitPrice: dto.UnitPrice,
					LineTotal: dto.LineTotal,
					CartID:    dto.CartID,
					ProductID: dto.ProductID,
				}
				mockRepo.On("Create", mock.Anything, &dto).Return(expectedItem, nil)
			},
			expectedResult: &GetCart_itemDTO{
				ID:        uuid.New(),
				Qty:       5,
				UnitPrice: 2.49,
				LineTotal: 12.45,
				CartID:    uuid.New(),
				ProductID: uuid.New(),
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateCart_itemDTO{
				ID:        uuid.New(),
				Qty:       5,
				UnitPrice: 2.49,
				LineTotal: 12.45,
				CartID:    uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateCart_itemDTO) {
				mockRepo.On("Create", mock.Anything, &dto).Return((*GetCart_itemDTO)(nil), errors.New("cart item already exists"))
			},
			expectedResult: nil,
			expectedError:  errors.New("cart item already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.createDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.createDTO.Qty, result.Qty)
				assert.Equal(t, tt.createDTO.UnitPrice, result.UnitPrice)
				assert.Equal(t, tt.createDTO.LineTotal, result.LineTotal)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name           string
		cartItemID     uuid.UUID
		updateDTO      UpdateCart_itemDTO
		mockSetup      func(*MockRepository, uuid.UUID, UpdateCart_itemDTO)
		expectedResult *GetCart_itemDTO
		expectedError  error
	}{
		{
			name:       "success - updates cart item",
			cartItemID: uuid.New(),
			updateDTO: UpdateCart_itemDTO{
				Qty:       intPtr(10),
				UnitPrice: float64Ptr(1.50),
				LineTotal: float64Ptr(15.00),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateCart_itemDTO) {
				expectedItem := &GetCart_itemDTO{
					ID:        id,
					Qty:       *dto.Qty,
					UnitPrice: *dto.UnitPrice,
					LineTotal: *dto.LineTotal,
					CartID:    uuid.New(),
					ProductID: uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *UpdateCart_itemDTO) bool {
					return u.ID == id && *u.Qty == *dto.Qty && *u.UnitPrice == *dto.UnitPrice
				})).Return(expectedItem, nil)
			},
			expectedResult: &GetCart_itemDTO{
				ID:        uuid.New(),
				Qty:       10,
				UnitPrice: 1.50,
				LineTotal: 15.00,
				CartID:    uuid.New(),
				ProductID: uuid.New(),
			},
			expectedError: nil,
		},
		{
			name:       "error - repository returns error",
			cartItemID: uuid.New(),
			updateDTO: UpdateCart_itemDTO{
				Qty: intPtr(10),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateCart_itemDTO) {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return((*GetCart_itemDTO)(nil), errors.New("cart item not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("cart item not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.cartItemID, tt.updateDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Update(ctx, tt.cartItemID, tt.updateDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.cartItemID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		cartItemID    uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError error
	}{
		{
			name:       "success - deletes cart item",
			cartItemID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "error - repository returns error",
			cartItemID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("cart item not found"))
			},
			expectedError: errors.New("cart item not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.cartItemID)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.cartItemID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
