package order_items

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

func (m *MockRepository) List(ctx context.Context) ([]*GetOrder_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetOrder_itemDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetOrder_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrder_itemDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateOrder_itemDTO) (*GetOrder_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrder_itemDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateOrder_itemDTO) (*GetOrder_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetOrder_itemDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockRepository)
		want      []*GetOrder_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns order items list",
			mockSetup: func(mockRepo *MockRepository) {
				expectedItems := []*GetOrder_itemDTO{
					{
						ID:        uuid.New(),
						Qty:       2,
						UnitPrice: 10.99,
						LineTotal: 21.98,
						OrderID:   uuid.New(),
						ProductID: uuid.New(),
					},
					{
						ID:        uuid.New(),
						Qty:       3,
						UnitPrice: 15.50,
						LineTotal: 46.50,
						OrderID:   uuid.New(),
						ProductID: uuid.New(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedItems, nil)
			},
			want: []*GetOrder_itemDTO{
				{
					ID:        uuid.New(),
					Qty:       2,
					UnitPrice: 10.99,
					LineTotal: 21.98,
					OrderID:   uuid.New(),
					ProductID: uuid.New(),
				},
				{
					ID:        uuid.New(),
					Qty:       3,
					UnitPrice: 15.50,
					LineTotal: 46.50,
					OrderID:   uuid.New(),
					ProductID: uuid.New(),
				},
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetOrder_itemDTO(nil), errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := NewService(mockRepo)
			ctx := context.Background()

			got, err := svc.List(ctx)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
				assert.Len(t, got, len(tt.want))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(*MockRepository, uuid.UUID)
		want      *GetOrder_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns order item by ID",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				expectedItem := &GetOrder_itemDTO{
					ID:        id,
					Qty:       2,
					UnitPrice: 10.99,
					LineTotal: 21.98,
					OrderID:   uuid.New(),
					ProductID: uuid.New(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedItem, nil)
			},
			want: &GetOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 10.99,
				LineTotal: 21.98,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - order item not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetOrder_itemDTO)(nil), errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id)

			svc := NewService(mockRepo)
			ctx := context.Background()

			got, err := svc.Get(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Qty, got.Qty)
				assert.Equal(t, tt.want.UnitPrice, got.UnitPrice)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		dto       CreateOrder_itemDTO
		mockSetup func(*MockRepository, CreateOrder_itemDTO)
		want      *GetOrder_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - creates new order item",
			dto: CreateOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 10.99,
				LineTotal: 21.98,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrder_itemDTO) {
				expectedItem := &GetOrder_itemDTO{
					ID:        dto.ID,
					Qty:       dto.Qty,
					UnitPrice: dto.UnitPrice,
					LineTotal: dto.LineTotal,
					OrderID:   dto.OrderID,
					ProductID: dto.ProductID,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateOrder_itemDTO) bool {
					return actual.ID == dto.ID && actual.Qty == dto.Qty
				})).Return(expectedItem, nil)
			},
			want: &GetOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 10.99,
				LineTotal: 21.98,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				UnitPrice: 10.99,
				LineTotal: 21.98,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateOrder_itemDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateOrder_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetOrder_itemDTO)(nil), errors.New("creation failed"))
			},
			want:    nil,
			wantErr: true,
			errMsg:  "creation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			ctx := context.Background()

			got, err := svc.Create(ctx, tt.dto)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.dto.Qty, got.Qty)
				assert.Equal(t, tt.dto.UnitPrice, got.UnitPrice)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		dto       UpdateOrder_itemDTO
		mockSetup func(*MockRepository, uuid.UUID, UpdateOrder_itemDTO)
		want      *GetOrder_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - updates order item",
			id:   uuid.New(),
			dto: UpdateOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       intPtr(5),
				UnitPrice: float64Ptr(12.99),
				LineTotal: float64Ptr(64.95),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateOrder_itemDTO) {
				expectedItem := &GetOrder_itemDTO{
					ID:        id,
					Qty:       *dto.Qty,
					UnitPrice: *dto.UnitPrice,
					LineTotal: *dto.LineTotal,
					OrderID:   uuid.New(),
					ProductID: uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateOrder_itemDTO) bool {
					return actual.ID == id
				})).Return(expectedItem, nil)
			},
			want: &GetOrder_itemDTO{
				ID:        uuid.New(),
				Qty:       5,
				UnitPrice: 12.99,
				LineTotal: 64.95,
				OrderID:   uuid.New(),
				ProductID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateOrder_itemDTO{
				ID:  uuid.New(),
				Qty: intPtr(5),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateOrder_itemDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateOrder_itemDTO) bool {
					return actual.ID == id
				})).Return((*GetOrder_itemDTO)(nil), errors.New("update failed"))
			},
			want:    nil,
			wantErr: true,
			errMsg:  "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id, tt.dto)

			svc := NewService(mockRepo)
			ctx := context.Background()

			got, err := svc.Update(ctx, tt.id, tt.dto)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
				if tt.dto.Qty != nil {
					assert.Equal(t, *tt.dto.Qty, got.Qty)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(*MockRepository, uuid.UUID)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - deletes order item",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			wantErr: true,
			errMsg:  "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id)

			svc := NewService(mockRepo)
			ctx := context.Background()

			err := svc.Delete(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

