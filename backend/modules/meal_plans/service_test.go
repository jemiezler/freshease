package meal_plans

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

func (m *MockRepository) List(ctx context.Context) ([]*GetMeal_planDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetMeal_planDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetMeal_planDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_planDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateMeal_planDTO) (*GetMeal_planDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_planDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateMeal_planDTO) (*GetMeal_planDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_planDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockRepository)
		want      []*GetMeal_planDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns meal plans list",
			mockSetup: func(mockRepo *MockRepository) {
				goal1 := "Weight loss"
				goal2 := "Muscle gain"
				expectedItems := []*GetMeal_planDTO{
					{
						ID:        uuid.New(),
						WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						Goal:      &goal1,
						UserID:    uuid.New(),
					},
					{
						ID:        uuid.New(),
						WeekStart: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
						Goal:      &goal2,
						UserID:    uuid.New(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedItems, nil)
			},
			want: []*GetMeal_planDTO{
				{
					ID:        uuid.New(),
					WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Goal:      stringPtr("Weight loss"),
					UserID:    uuid.New(),
				},
				{
					ID:        uuid.New(),
					WeekStart: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
					Goal:      stringPtr("Muscle gain"),
					UserID:    uuid.New(),
				},
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetMeal_planDTO(nil), errors.New("database error"))
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
		want      *GetMeal_planDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns meal plan by ID",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				goal := "Weight loss"
				expectedItem := &GetMeal_planDTO{
					ID:        id,
					WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Goal:      &goal,
					UserID:    uuid.New(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedItem, nil)
			},
			want: &GetMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Goal:      stringPtr("Weight loss"),
				UserID:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - meal plan not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetMeal_planDTO)(nil), errors.New("not found"))
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
				assert.Equal(t, tt.want.WeekStart, got.WeekStart)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		dto       CreateMeal_planDTO
		mockSetup func(*MockRepository, CreateMeal_planDTO)
		want      *GetMeal_planDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - creates new meal plan",
			dto: CreateMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Goal:      stringPtr("Weight loss"),
				UserID:    uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateMeal_planDTO) {
				expectedItem := &GetMeal_planDTO{
					ID:        dto.ID,
					WeekStart: dto.WeekStart,
					Goal:      dto.Goal,
					UserID:    dto.UserID,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateMeal_planDTO) bool {
					return actual.ID == dto.ID && actual.UserID == dto.UserID
				})).Return(expectedItem, nil)
			},
			want: &GetMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Goal:      stringPtr("Weight loss"),
				UserID:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UserID:    uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateMeal_planDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateMeal_planDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetMeal_planDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.WeekStart, got.WeekStart)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		dto       UpdateMeal_planDTO
		mockSetup func(*MockRepository, uuid.UUID, UpdateMeal_planDTO)
		want      *GetMeal_planDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - updates meal plan",
			id:   uuid.New(),
			dto: UpdateMeal_planDTO{
				ID:   uuid.New(),
				Goal: stringPtr("Muscle gain"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateMeal_planDTO) {
				expectedItem := &GetMeal_planDTO{
					ID:        id,
					WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Goal:      dto.Goal,
					UserID:    uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateMeal_planDTO) bool {
					return actual.ID == id
				})).Return(expectedItem, nil)
			},
			want: &GetMeal_planDTO{
				ID:        uuid.New(),
				WeekStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Goal:      stringPtr("Muscle gain"),
				UserID:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateMeal_planDTO{
				ID:   uuid.New(),
				Goal: stringPtr("Muscle gain"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateMeal_planDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateMeal_planDTO) bool {
					return actual.ID == id
				})).Return((*GetMeal_planDTO)(nil), errors.New("update failed"))
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
				if tt.dto.Goal != nil {
					assert.Equal(t, *tt.dto.Goal, *got.Goal)
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
			name: "success - deletes meal plan",
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

