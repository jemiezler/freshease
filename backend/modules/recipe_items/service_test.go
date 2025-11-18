package recipe_items

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

func (m *MockRepository) List(ctx context.Context) ([]*GetRecipe_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetRecipe_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateRecipe_itemDTO) (*GetRecipe_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateRecipe_itemDTO) (*GetRecipe_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipe_itemDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockRepository)
		want      []*GetRecipe_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns recipe items list",
			mockSetup: func(mockRepo *MockRepository) {
				expectedItems := []*GetRecipe_itemDTO{
					{
						ID:        uuid.New(),
						Amount:    2.5,
						Unit:      "cups",
						RecipeID:  uuid.New(),
						ProductID: uuid.New(),
					},
					{
						ID:        uuid.New(),
						Amount:    1.0,
						Unit:      "tbsp",
						RecipeID:  uuid.New(),
						ProductID: uuid.New(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedItems, nil)
			},
			want: []*GetRecipe_itemDTO{
				{
					ID:        uuid.New(),
					Amount:    2.5,
					Unit:      "cups",
					RecipeID:  uuid.New(),
					ProductID: uuid.New(),
				},
				{
					ID:        uuid.New(),
					Amount:    1.0,
					Unit:      "tbsp",
					RecipeID:  uuid.New(),
					ProductID: uuid.New(),
				},
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetRecipe_itemDTO(nil), errors.New("database error"))
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
		want      *GetRecipe_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - returns recipe item by ID",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				expectedItem := &GetRecipe_itemDTO{
					ID:        id,
					Amount:    2.5,
					Unit:      "cups",
					RecipeID:  uuid.New(),
					ProductID: uuid.New(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedItem, nil)
			},
			want: &GetRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    2.5,
				Unit:      "cups",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - recipe item not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetRecipe_itemDTO)(nil), errors.New("not found"))
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
				assert.Equal(t, tt.want.Amount, got.Amount)
				assert.Equal(t, tt.want.Unit, got.Unit)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		dto       CreateRecipe_itemDTO
		mockSetup func(*MockRepository, CreateRecipe_itemDTO)
		want      *GetRecipe_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - creates new recipe item",
			dto: CreateRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    2.5,
				Unit:      "cups",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateRecipe_itemDTO) {
				expectedItem := &GetRecipe_itemDTO{
					ID:        dto.ID,
					Amount:    dto.Amount,
					Unit:      dto.Unit,
					RecipeID:  dto.RecipeID,
					ProductID: dto.ProductID,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateRecipe_itemDTO) bool {
					return actual.ID == dto.ID && actual.Amount == dto.Amount
				})).Return(expectedItem, nil)
			},
			want: &GetRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    2.5,
				Unit:      "cups",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    2.5,
				Unit:      "cups",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateRecipe_itemDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateRecipe_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetRecipe_itemDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.Amount, got.Amount)
				assert.Equal(t, tt.dto.Unit, got.Unit)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		dto       UpdateRecipe_itemDTO
		mockSetup func(*MockRepository, uuid.UUID, UpdateRecipe_itemDTO)
		want      *GetRecipe_itemDTO
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - updates recipe item",
			id:   uuid.New(),
			dto: UpdateRecipe_itemDTO{
				ID:     uuid.New(),
				Amount: float64Ptr(3.0),
				Unit:   stringPtr("tbsp"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateRecipe_itemDTO) {
				expectedItem := &GetRecipe_itemDTO{
					ID:        id,
					Amount:    *dto.Amount,
					Unit:      *dto.Unit,
					RecipeID:  uuid.New(),
					ProductID: uuid.New(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateRecipe_itemDTO) bool {
					return actual.ID == id
				})).Return(expectedItem, nil)
			},
			want: &GetRecipe_itemDTO{
				ID:        uuid.New(),
				Amount:    3.0,
				Unit:      "tbsp",
				RecipeID:  uuid.New(),
				ProductID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateRecipe_itemDTO{
				ID:     uuid.New(),
				Amount: float64Ptr(3.0),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateRecipe_itemDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateRecipe_itemDTO) bool {
					return actual.ID == id
				})).Return((*GetRecipe_itemDTO)(nil), errors.New("update failed"))
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
				if tt.dto.Amount != nil {
					assert.Equal(t, *tt.dto.Amount, got.Amount)
				}
				if tt.dto.Unit != nil {
					assert.Equal(t, *tt.dto.Unit, got.Unit)
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
			name: "success - deletes recipe item",
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

