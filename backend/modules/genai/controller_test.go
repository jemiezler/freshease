package genai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) GenerateWeeklyMeals(ctx context.Context, req *GenerateMealsReq) (any, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func (m *MockService) GenerateDailyMeals(ctx context.Context, req *GenerateMealsReq) (any, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestController_GenerateWeeklyMeals(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    GenerateMealsReq
		mockSetup      func(*MockService, GenerateMealsReq)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - generates weekly meals",
			requestBody: GenerateMealsReq{
				UserID:        "user123",
				Gender:        "male",
				Age:           30,
				HeightCm:      175.0,
				WeightKg:      70.0,
				StepsToday:    8000,
				ActiveKcal24h: 2000.0,
				Allergies:     []string{"nuts"},
				Preferences:   []string{"vegetarian"},
				Target:        "maintenance",
			},
			mockSetup: func(mockSvc *MockService, req GenerateMealsReq) {
				expectedResult := map[string]any{
					"steps_today":     req.StepsToday,
					"active_kcal_24h": req.ActiveKcal24h,
					"plan": []map[string]any{
						{
							"day": "Monday",
							"meals": map[string]string{
								"breakfast": "Oatmeal with berries",
								"lunch":     "Quinoa salad",
								"dinner":    "Grilled chicken",
							},
							"calories": map[string]int{
								"breakfast": 300,
								"lunch":     400,
								"dinner":    500,
							},
						},
					},
				}
				mockSvc.On("GenerateWeeklyMeals", mock.Anything, mock.MatchedBy(func(actual *GenerateMealsReq) bool {
					return actual.UserID == req.UserID &&
						actual.Gender == req.Gender &&
						actual.Age == req.Age &&
						actual.HeightCm == req.HeightCm &&
						actual.WeightKg == req.WeightKg &&
						actual.StepsToday == req.StepsToday &&
						actual.ActiveKcal24h == req.ActiveKcal24h &&
						actual.Target == req.Target
				})).Return(expectedResult, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: GenerateMealsReq{
				UserID:        "user123",
				Gender:        "male",
				Age:           30,
				HeightCm:      175.0,
				WeightKg:      70.0,
				StepsToday:    8000,
				ActiveKcal24h: 2000.0,
				Allergies:     []string{"nuts"},
				Preferences:   []string{"vegetarian"},
				Target:        "maintenance",
			},
			mockSetup: func(mockSvc *MockService, req GenerateMealsReq) {
				mockSvc.On("GenerateWeeklyMeals", mock.Anything, mock.MatchedBy(func(actual *GenerateMealsReq) bool {
					return actual.UserID == req.UserID &&
						actual.Gender == req.Gender &&
						actual.Age == req.Age &&
						actual.HeightCm == req.HeightCm &&
						actual.WeightKg == req.WeightKg &&
						actual.StepsToday == req.StepsToday &&
						actual.ActiveKcal24h == req.ActiveKcal24h &&
						actual.Target == req.Target
				})).Return(nil, errors.New("generation failed"))
			},
			expectedStatus: http.StatusBadGateway,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/weekly-meals", controller.GenerateWeeklyMeals)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/weekly-meals", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "steps_today")
				assert.Contains(t, responseBody, "active_kcal_24h")
				assert.Contains(t, responseBody, "plan")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GenerateDailyMeals(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    GenerateMealsReq
		mockSetup      func(*MockService, GenerateMealsReq)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success - generates daily meals",
			requestBody: GenerateMealsReq{
				UserID:        "user123",
				Gender:        "female",
				Age:           25,
				HeightCm:      165.0,
				WeightKg:      60.0,
				StepsToday:    10000,
				ActiveKcal24h: 1800.0,
				Allergies:     []string{"dairy"},
				Preferences:   []string{"vegan"},
				Target:        "weight_loss",
			},
			mockSetup: func(mockSvc *MockService, req GenerateMealsReq) {
				expectedResult := map[string]any{
					"steps_today":     req.StepsToday,
					"active_kcal_24h": req.ActiveKcal24h,
					"plan": []map[string]any{
						{
							"day": "Monday",
							"meals": map[string]string{
								"breakfast": "Green smoothie",
								"lunch":     "Vegan Buddha bowl",
								"dinner":    "Grilled tofu",
							},
							"calories": map[string]int{
								"breakfast": 250,
								"lunch":     350,
								"dinner":    400,
							},
							"total_calories": 1000,
						},
					},
				}
				mockSvc.On("GenerateDailyMeals", mock.Anything, mock.MatchedBy(func(actual *GenerateMealsReq) bool {
					return actual.UserID == req.UserID &&
						actual.Gender == req.Gender &&
						actual.Age == req.Age &&
						actual.HeightCm == req.HeightCm &&
						actual.WeightKg == req.WeightKg &&
						actual.StepsToday == req.StepsToday &&
						actual.ActiveKcal24h == req.ActiveKcal24h &&
						actual.Target == req.Target
				})).Return(expectedResult, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "error - service returns error",
			requestBody: GenerateMealsReq{
				UserID:        "user123",
				Gender:        "female",
				Age:           25,
				HeightCm:      165.0,
				WeightKg:      60.0,
				StepsToday:    10000,
				ActiveKcal24h: 1800.0,
				Allergies:     []string{"dairy"},
				Preferences:   []string{"vegan"},
				Target:        "weight_loss",
			},
			mockSetup: func(mockSvc *MockService, req GenerateMealsReq) {
				mockSvc.On("GenerateDailyMeals", mock.Anything, mock.MatchedBy(func(actual *GenerateMealsReq) bool {
					return actual.UserID == req.UserID &&
						actual.Gender == req.Gender &&
						actual.Age == req.Age &&
						actual.HeightCm == req.HeightCm &&
						actual.WeightKg == req.WeightKg &&
						actual.StepsToday == req.StepsToday &&
						actual.ActiveKcal24h == req.ActiveKcal24h &&
						actual.Target == req.Target
				})).Return(nil, errors.New("generation failed"))
			},
			expectedStatus: http.StatusBadGateway,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/daily-meals", controller.GenerateDailyMeals)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/daily-meals", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Contains(t, responseBody, "steps_today")
				assert.Contains(t, responseBody, "active_kcal_24h")
				assert.Contains(t, responseBody, "plan")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_RequestParsing(t *testing.T) {
	t.Run("handles malformed JSON", func(t *testing.T) {
		mockSvc := new(MockService)
		// No mock setup needed - should fail before service call

		controller := NewController(mockSvc)
		app := fiber.New()
		app.Post("/weekly-meals", controller.GenerateWeeklyMeals)

		// Send malformed JSON
		req := httptest.NewRequest(http.MethodPost, "/weekly-meals", bytes.NewBufferString(`{"invalid": json}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("handles empty request body", func(t *testing.T) {
		mockSvc := new(MockService)
		// No mock setup needed - should fail before service call

		controller := NewController(mockSvc)
		app := fiber.New()
		app.Post("/daily-meals", controller.GenerateDailyMeals)

		// Send empty body
		req := httptest.NewRequest(http.MethodPost, "/daily-meals", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}
