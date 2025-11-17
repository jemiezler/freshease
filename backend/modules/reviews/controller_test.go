package reviews

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) List(ctx context.Context) ([]*GetReviewDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetReviewDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetReviewDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReviewDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateReviewDTO) (*GetReviewDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReviewDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateReviewDTO) (*GetReviewDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReviewDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListReviews(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - returns reviews list",
			mockSetup: func(mockSvc *MockService) {
				userID := uuid.New()
				productID := uuid.New()
				comment := "Great product!"
				expectedReviews := []*GetReviewDTO{
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
				}
				mockSvc.On("List", mock.Anything).Return(expectedReviews, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMessage: "Reviews Retrieved Successfully",
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetReviewDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/reviews", controller.ListReviews)

			req := httptest.NewRequest(http.MethodGet, "/reviews", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetReview(t *testing.T) {
	reviewID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()

	tests := []struct {
		name           string
		reviewID       string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - returns review by ID",
			reviewID: reviewID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				comment := "Great product!"
				mockSvc.On("Get", mock.Anything, id).Return(&GetReviewDTO{
					ID:        id,
					Rating:    5,
					Comment:   &comment,
					UserID:    userID,
					ProductID: productID,
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMessage: "Review Retrieved Successfully",
		},
		{
			name:     "error - invalid UUID",
			reviewID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:     "error - review not found",
			reviewID: reviewID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedMessage: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.reviewID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.reviewID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/reviews/:id", controller.GetReview)

			req := httptest.NewRequest(http.MethodGet, "/reviews/"+tt.reviewID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateReview(t *testing.T) {
	reviewID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()
	comment := "New review"
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    CreateReviewDTO
		mockSetup      func(*MockService, CreateReviewDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - creates review",
			requestBody: CreateReviewDTO{
				ID:        reviewID,
				Rating:    5,
				Comment:   &comment,
				UserID:    userID,
				ProductID: productID,
				CreatedAt: &now,
			},
			mockSetup: func(mockSvc *MockService, dto CreateReviewDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(&GetReviewDTO{
					ID:        dto.ID,
					Rating:    dto.Rating,
					Comment:   dto.Comment,
					UserID:    dto.UserID,
					ProductID: dto.ProductID,
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Review Created Successfully",
		},
		{
			name: "error - service returns error",
			requestBody: CreateReviewDTO{
				ID:        reviewID,
				Rating:    5,
				UserID:    userID,
				ProductID: productID,
			},
			mockSetup: func(mockSvc *MockService, dto CreateReviewDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return(nil, errors.New("validation error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/reviews", controller.CreateReview)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateReview(t *testing.T) {
	reviewID := uuid.New()
	newRating := 5
	newComment := "Updated comment"

	tests := []struct {
		name           string
		reviewID       string
		requestBody    UpdateReviewDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateReviewDTO)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - updates review",
			reviewID: reviewID.String(),
			requestBody: UpdateReviewDTO{
				ID:      reviewID,
				Rating:  &newRating,
				Comment: &newComment,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateReviewDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return(&GetReviewDTO{
					ID:        id,
					Rating:    newRating,
					Comment:   &newComment,
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMessage: "Review Updated Successfully",
		},
		{
			name:     "error - invalid UUID",
			reviewID: "invalid-uuid",
			requestBody: UpdateReviewDTO{
				Rating: &newRating,
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateReviewDTO) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.reviewID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.reviewID)
				tt.mockSetup(mockSvc, id, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/reviews/:id", controller.UpdateReview)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/reviews/"+tt.reviewID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, responseBody, "data")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteReview(t *testing.T) {
	reviewID := uuid.New()

	tests := []struct {
		name           string
		reviewID       string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "success - deletes review",
			reviewID: reviewID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedMessage: "Review Deleted Successfully",
		},
		{
			name:     "error - invalid UUID",
			reviewID: "invalid-uuid",
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "invalid uuid",
		},
		{
			name:     "error - service returns error",
			reviewID: reviewID.String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.reviewID != "invalid-uuid" {
				id, _ := uuid.Parse(tt.reviewID)
				tt.mockSetup(mockSvc, id)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/reviews/:id", controller.DeleteReview)

			req := httptest.NewRequest(http.MethodDelete, "/reviews/"+tt.reviewID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			mockSvc.AssertExpectations(t)
		})
	}
}

