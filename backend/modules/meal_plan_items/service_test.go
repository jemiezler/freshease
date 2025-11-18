package meal_plan_items

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

func (m *MockRepository) List(ctx context.Context) ([]*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateMeal_plan_itemDTO) (*GetMeal_plan_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMeal_plan_itemDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockRepository)
		want      []*GetMeal_plan_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns meal plan items list",
			mockSetup: func(mockRepo *MockRepository) {
				expectedItems := []*GetMeal_plan_itemDTO{
					{
						ID:         uuid.New(),
						Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						Slot:       "breakfast",
						MealPlanID: uuid.New(),
						RecipeID:   uuid.New(),
					},
					{
						ID:         uuid.New(),
						Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						Slot:       "lunch",
						MealPlanID: uuid.New(),
						RecipeID:   uuid.New(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedItems, nil)
			},
			want: []*GetMeal_plan_itemDTO{
				{
					ID:         uuid.New(),
					Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Slot:       "breakfast",
					MealPlanID: uuid.New(),
					RecipeID:   uuid.New(),
				},
				{
					ID:         uuid.New(),
					Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Slot:       "lunch",
					MealPlanID: uuid.New(),
					RecipeID:   uuid.New(),
				},
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetMeal_plan_itemDTO(nil), errors.New("database error"))
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
		want      *GetMeal_plan_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns meal plan item by ID",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				expectedItem := &GetMeal_plan_itemDTO{
					ID:         id,
					Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Slot:       "breakfast",
					MealPlanID: uuid.New(),
					RecipeID:   uuid.New(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedItem, nil)
			},
			want: &GetMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "breakfast",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - meal plan item not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetMeal_plan_itemDTO)(nil), errors.New("not found"))
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
				assert.Equal(t, tt.want.Slot, got.Slot)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		dto       CreateMeal_plan_itemDTO
		mockSetup func(*MockRepository, CreateMeal_plan_itemDTO)
		want      *GetMeal_plan_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - creates new meal plan item",
			dto: CreateMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "breakfast",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateMeal_plan_itemDTO) {
				expectedItem := &GetMeal_plan_itemDTO{
					ID:         dto.ID,
					Day:        dto.Day,
					Slot:       dto.Slot,
					MealPlanID: dto.MealPlanID,
					RecipeID:   dto.RecipeID,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateMeal_plan_itemDTO) bool {
					return actual.ID == dto.ID && actual.Slot == dto.Slot
				})).Return(expectedItem, nil)
			},
			want: &GetMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "breakfast",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "breakfast",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateMeal_plan_itemDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateMeal_plan_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetMeal_plan_itemDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.Slot, got.Slot)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		dto       UpdateMeal_plan_itemDTO
		mockSetup func(*MockRepository, uuid.UUID, UpdateMeal_plan_itemDTO)
		want      *GetMeal_plan_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - updates meal plan item",
			id:   uuid.New(),
			dto: UpdateMeal_plan_itemDTO{
				ID:   uuid.New(),
				Slot: stringPtr("lunch"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateMeal_plan_itemDTO) {
				expectedItem := &GetMeal_plan_itemDTO{
					ID:         id,
					Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Slot:       *dto.Slot,
					MealPlanID: uuid.New(),
					RecipeID:   uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateMeal_plan_itemDTO) bool {
					return actual.ID == id
				})).Return(expectedItem, nil)
			},
			want: &GetMeal_plan_itemDTO{
				ID:         uuid.New(),
				Day:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Slot:       "lunch",
				MealPlanID: uuid.New(),
				RecipeID:   uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateMeal_plan_itemDTO{
				ID:   uuid.New(),
				Slot: stringPtr("lunch"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateMeal_plan_itemDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateMeal_plan_itemDTO) bool {
					return actual.ID == id
				})).Return((*GetMeal_plan_itemDTO)(nil), errors.New("update failed"))
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
				if tt.dto.Slot != nil {
					assert.Equal(t, *tt.dto.Slot, got.Slot)
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
			name: "success - deletes meal plan item",
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

