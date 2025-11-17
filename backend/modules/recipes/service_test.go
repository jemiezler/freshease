package recipes

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

func (m *MockRepository) List(ctx context.Context) ([]*GetRecipeDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetRecipeDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetRecipeDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipeDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateRecipeDTO) (*GetRecipeDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipeDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateRecipeDTO) (*GetRecipeDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecipeDTO), args.Error(1)
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
			name: "success - returns recipes list",
			mockSetup: func(mockRepo *MockRepository) {
				instructions := "Test instructions"
				mockRepo.On("List", context.Background()).Return([]*GetRecipeDTO{
					{
						ID:           uuid.New(),
						Name:         "Recipe One",
						Instructions: &instructions,
						Kcal:         500,
					},
					{
						ID:           uuid.New(),
						Name:         "Recipe Two",
						Instructions: &instructions,
						Kcal:         600,
					},
				}, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetRecipeDTO(nil), errors.New("database error"))
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
		recipeID      uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:     "success - returns recipe by ID",
			recipeID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				instructions := "Test instructions"
				mockRepo.On("FindByID", context.Background(), id).Return(&GetRecipeDTO{
					ID:           id,
					Name:         "Test Recipe",
					Instructions: &instructions,
					Kcal:         500,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:     "error - recipe not found",
			recipeID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.recipeID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.recipeID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.recipeID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		dto           CreateRecipeDTO
		mockSetup     func(*MockRepository, CreateRecipeDTO)
		expectedError bool
	}{
		{
			name: "success - creates recipe",
			dto: CreateRecipeDTO{
				ID:   uuid.New(),
				Name: "New Recipe",
				Kcal: 550,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateRecipeDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetRecipeDTO{
					ID:   dto.ID,
					Name: dto.Name,
					Kcal: dto.Kcal,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateRecipeDTO{
				ID:   uuid.New(),
				Name: "New Recipe",
				Kcal: 550,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateRecipeDTO) {
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
				assert.Equal(t, tt.dto.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	recipeID := uuid.New()
	newName := "Updated Recipe"
	newKcal := 600

	tests := []struct {
		name          string
		dto           UpdateRecipeDTO
		mockSetup     func(*MockRepository, UpdateRecipeDTO)
		expectedError bool
	}{
		{
			name: "success - updates recipe",
			dto: UpdateRecipeDTO{
				ID:   recipeID,
				Name: &newName,
				Kcal: &newKcal,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateRecipeDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*recipes.UpdateRecipeDTO")).Return(&GetRecipeDTO{
					ID:   recipeID,
					Name: newName,
					Kcal: newKcal,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdateRecipeDTO{
				ID:   recipeID,
				Name: &newName,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateRecipeDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*recipes.UpdateRecipeDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), recipeID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, recipeID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	recipeID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes recipe",
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
			tt.mockSetup(mockRepo, recipeID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), recipeID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

