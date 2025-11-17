package reviews

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

func (m *MockRepository) List(ctx context.Context) ([]*GetReviewDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetReviewDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetReviewDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReviewDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateReviewDTO) (*GetReviewDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReviewDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateReviewDTO) (*GetReviewDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReviewDTO), args.Error(1)
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
			name: "success - returns reviews list",
			mockSetup: func(mockRepo *MockRepository) {
				userID := uuid.New()
				productID := uuid.New()
				comment := "Great product!"
				mockRepo.On("List", context.Background()).Return([]*GetReviewDTO{
					{
						ID:        uuid.New(),
						Rating:    5,
						Comment:   &comment,
						UserID:    userID,
						ProductID: productID,
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Rating:    4,
						Comment:   &comment,
						UserID:    userID,
						ProductID: productID,
						CreatedAt: time.Now(),
					},
				}, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", context.Background()).Return([]*GetReviewDTO(nil), errors.New("database error"))
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
		reviewID      uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name:     "success - returns review by ID",
			reviewID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				userID := uuid.New()
				productID := uuid.New()
				comment := "Great product!"
				mockRepo.On("FindByID", context.Background(), id).Return(&GetReviewDTO{
					ID:        id,
					Rating:    5,
					Comment:   &comment,
					UserID:    userID,
					ProductID: productID,
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedError: false,
		},
		{
			name:     "error - review not found",
			reviewID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", context.Background(), id).Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.reviewID)

			svc := NewService(mockRepo)
			result, err := svc.Get(context.Background(), tt.reviewID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.reviewID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	comment := "New review"
	now := time.Now()
	tests := []struct {
		name          string
		dto           CreateReviewDTO
		mockSetup     func(*MockRepository, CreateReviewDTO)
		expectedError bool
	}{
		{
			name: "success - creates review",
			dto: CreateReviewDTO{
				ID:        uuid.New(),
				Rating:    5,
				Comment:   &comment,
				UserID:    userID,
				ProductID: productID,
				CreatedAt: &now,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateReviewDTO) {
				mockRepo.On("Create", context.Background(), &dto).Return(&GetReviewDTO{
					ID:        dto.ID,
					Rating:    dto.Rating,
					Comment:   dto.Comment,
					UserID:    dto.UserID,
					ProductID: dto.ProductID,
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateReviewDTO{
				ID:        uuid.New(),
				Rating:    5,
				UserID:    userID,
				ProductID: productID,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateReviewDTO) {
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
				assert.Equal(t, tt.dto.Rating, result.Rating)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	reviewID := uuid.New()
	newRating := 5
	newComment := "Updated comment"

	tests := []struct {
		name          string
		dto           UpdateReviewDTO
		mockSetup     func(*MockRepository, UpdateReviewDTO)
		expectedError bool
	}{
		{
			name: "success - updates review",
			dto: UpdateReviewDTO{
				ID:      reviewID,
				Rating:  &newRating,
				Comment: &newComment,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateReviewDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*reviews.UpdateReviewDTO")).Return(&GetReviewDTO{
					ID:        reviewID,
					Rating:    newRating,
					Comment:   &newComment,
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: UpdateReviewDTO{
				ID:     reviewID,
				Rating: &newRating,
			},
			mockSetup: func(mockRepo *MockRepository, dto UpdateReviewDTO) {
				mockRepo.On("Update", context.Background(), mock.AnythingOfType("*reviews.UpdateReviewDTO")).Return(nil, errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			svc := NewService(mockRepo)
			result, err := svc.Update(context.Background(), reviewID, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, reviewID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	reviewID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes review",
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
			tt.mockSetup(mockRepo, reviewID)

			svc := NewService(mockRepo)
			err := svc.Delete(context.Background(), reviewID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

