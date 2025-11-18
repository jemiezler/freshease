package notifications

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

func (m *MockRepository) List(ctx context.Context) ([]*GetNotificationDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetNotificationDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetNotificationDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNotificationDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateNotificationDTO) (*GetNotificationDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNotificationDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateNotificationDTO) (*GetNotificationDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNotificationDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockRepository)
		want      []*GetNotificationDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns notifications list",
			mockSetup: func(mockRepo *MockRepository) {
				body1 := "Your order has been shipped"
				body2 := "New product available"
				expectedItems := []*GetNotificationDTO{
					{
						ID:        uuid.New(),
						Title:     "Order Shipped",
						Body:      &body1,
						Channel:   "email",
						Status:    "sent",
						UserID:    uuid.New(),
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Title:     "New Product",
						Body:      &body2,
						Channel:   "push",
						Status:    "pending",
						UserID:    uuid.New(),
						CreatedAt: time.Now(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedItems, nil)
			},
			want: []*GetNotificationDTO{
				{
					ID:        uuid.New(),
					Title:     "Order Shipped",
					Body:      stringPtr("Your order has been shipped"),
					Channel:   "email",
					Status:    "sent",
					UserID:    uuid.New(),
					CreatedAt: time.Now(),
				},
				{
					ID:        uuid.New(),
					Title:     "New Product",
					Body:      stringPtr("New product available"),
					Channel:   "push",
					Status:    "pending",
					UserID:    uuid.New(),
					CreatedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetNotificationDTO(nil), errors.New("database error"))
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
		want      *GetNotificationDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns notification by ID",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				body := "Your order has been shipped"
				expectedItem := &GetNotificationDTO{
					ID:        id,
					Title:     "Order Shipped",
					Body:      &body,
					Channel:   "email",
					Status:    "sent",
					UserID:    uuid.New(),
					CreatedAt: time.Now(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedItem, nil)
			},
			want: &GetNotificationDTO{
				ID:        uuid.New(),
				Title:     "Order Shipped",
				Body:      stringPtr("Your order has been shipped"),
				Channel:   "email",
				Status:    "sent",
				UserID:    uuid.New(),
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "error - notification not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetNotificationDTO)(nil), errors.New("not found"))
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
				assert.Equal(t, tt.want.Title, got.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		dto       CreateNotificationDTO
		mockSetup func(*MockRepository, CreateNotificationDTO)
		want      *GetNotificationDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - creates new notification",
			dto: CreateNotificationDTO{
				ID:      uuid.New(),
				Title:   "Order Shipped",
				Body:    stringPtr("Your order has been shipped"),
				Channel: "email",
				Status:  "sent",
				UserID:  uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateNotificationDTO) {
				expectedItem := &GetNotificationDTO{
					ID:        dto.ID,
					Title:     dto.Title,
					Body:      dto.Body,
					Channel:   dto.Channel,
					Status:    dto.Status,
					UserID:    dto.UserID,
					CreatedAt: time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateNotificationDTO) bool {
					return actual.ID == dto.ID && actual.Title == dto.Title
				})).Return(expectedItem, nil)
			},
			want: &GetNotificationDTO{
				ID:        uuid.New(),
				Title:     "Order Shipped",
				Body:      stringPtr("Your order has been shipped"),
				Channel:   "email",
				Status:    "sent",
				UserID:    uuid.New(),
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateNotificationDTO{
				ID:      uuid.New(),
				Title:   "Order Shipped",
				Channel: "email",
				Status:  "sent",
				UserID:  uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateNotificationDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateNotificationDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetNotificationDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.Title, got.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		dto       UpdateNotificationDTO
		mockSetup func(*MockRepository, uuid.UUID, UpdateNotificationDTO)
		want      *GetNotificationDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - updates notification",
			id:   uuid.New(),
			dto: UpdateNotificationDTO{
				ID:     uuid.New(),
				Status: stringPtr("read"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateNotificationDTO) {
				body := "Your order has been shipped"
				expectedItem := &GetNotificationDTO{
					ID:        id,
					Title:     "Order Shipped",
					Body:      &body,
					Channel:   "email",
					Status:    *dto.Status,
					UserID:    uuid.New(),
					CreatedAt: time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateNotificationDTO) bool {
					return actual.ID == id
				})).Return(expectedItem, nil)
			},
			want: &GetNotificationDTO{
				ID:        uuid.New(),
				Title:     "Order Shipped",
				Body:      stringPtr("Your order has been shipped"),
				Channel:   "email",
				Status:    "read",
				UserID:    uuid.New(),
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateNotificationDTO{
				ID:     uuid.New(),
				Status: stringPtr("read"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateNotificationDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateNotificationDTO) bool {
					return actual.ID == id
				})).Return((*GetNotificationDTO)(nil), errors.New("update failed"))
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
				if tt.dto.Status != nil {
					assert.Equal(t, *tt.dto.Status, got.Status)
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
			name: "success - deletes notification",
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

