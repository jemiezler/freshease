package cart_items

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (m *MockService) List(ctx context.Context) ([]*GetCart_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetCart_itemDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetCart_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCart_itemDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateCart_itemDTO) (*GetCart_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCart_itemDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateCart_itemDTO) (*GetCart_itemDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCart_itemDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListCart_items(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - returns cart items list",
			mockSetup: func(mockSvc *MockService) {
				expectedItems := []*GetCart_itemDTO{
					{
						ID:          uuid.New(),
						Name:        "Apple",
						Description: "Fresh red apples",
						Cart:        uuid.New(),
					},
					{
						ID:          uuid.New(),
						Name:        "Banana",
						Description: "Yellow bananas",
						Cart:        uuid.New(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*GetCart_itemDTO{
				{
					ID:          uuid.New(),
					Name:        "Apple",
					Description: "Fresh red apples",
					Cart:        uuid.New(),
				},
				{
					ID:          uuid.New(),
					Name:        "Banana",
					Description: "Yellow bananas",
					Cart:        uuid.New(),
				},
			},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetCart_itemDTO(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"message": "database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/cart_items", controller.ListCart_items)

			req := httptest.NewRequest(http.MethodGet, "/cart_items", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.IsType(t, map[string]interface{}{}, responseBody)
			} else {
				assert.IsType(t, map[string]interface{}{}, responseBody)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetCart_item(t *testing.T) {
	tests := []struct {
		name           string
		cartItemID     string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "success - returns cart item by ID",
			cartItemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				expectedItem := &GetCart_itemDTO{
					ID:          id,
					Name:        "Apple",
					Description: "Fresh red apples",
					Cart:        uuid.New(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &GetCart_itemDTO{
				ID:          uuid.New(),
				Name:        "Apple",
				Description: "Fresh red apples",
				Cart:        uuid.New(),
			},
		},
		{
			name:           "error - invalid UUID",
			cartItemID:     "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:       "error - cart item not found",
			cartItemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetCart_itemDTO)(nil), errors.New("cart item not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"message": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.cartItemID != "invalid-uuid" {
				cartItemID, err := uuid.Parse(tt.cartItemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, cartItemID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/cart_items/:id", controller.GetCart_item)

			req := httptest.NewRequest(http.MethodGet, "/cart_items/"+tt.cartItemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.IsType(t, map[string]interface{}{}, responseBody)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateCart_item(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateCart_itemDTO
		mockSetup      func(*MockService, CreateCart_itemDTO)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - creates new cart item",
			requestBody: CreateCart_itemDTO{
				ID:          uuid.New(),
				Name:        "Orange",
				Description: "Fresh oranges",
				Cart:        uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateCart_itemDTO) {
				expectedItem := &GetCart_itemDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Description: dto.Description,
					Cart:        dto.Cart,
				}
				mockSvc.On("Create", mock.Anything, dto).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &GetCart_itemDTO{
				ID:          uuid.New(),
				Name:        "Orange",
				Description: "Fresh oranges",
				Cart:        uuid.New(),
			},
		},
		{
			name: "error - service returns error",
			requestBody: CreateCart_itemDTO{
				ID:          uuid.New(),
				Name:        "Orange",
				Description: "Fresh oranges",
				Cart:        uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateCart_itemDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return((*GetCart_itemDTO)(nil), errors.New("cart item already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "cart item already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/cart_items", controller.CreateCart_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/cart_items", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.IsType(t, map[string]interface{}{}, responseBody)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateCart_item(t *testing.T) {
	tests := []struct {
		name           string
		cartItemID     string
		requestBody    UpdateCart_itemDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateCart_itemDTO)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "success - updates cart item",
			cartItemID: uuid.New().String(),
			requestBody: UpdateCart_itemDTO{
				ID:          uuid.New(), // This will be overridden by the service
				Name:        stringPtr("Updated Apple"),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCart_itemDTO) {
				expectedItem := &GetCart_itemDTO{
					ID:          id,
					Name:        *dto.Name,
					Description: *dto.Description,
					Cart:        uuid.New(),
				}
				mockSvc.On("Update", mock.Anything, id, dto).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &GetCart_itemDTO{
				ID:          uuid.New(),
				Name:        "Updated Apple",
				Description: "Updated description",
				Cart:        uuid.New(),
			},
		},
		{
			name:           "error - invalid UUID",
			cartItemID:     "invalid-uuid",
			requestBody:    UpdateCart_itemDTO{ID: uuid.New()},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateCart_itemDTO) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:       "error - service returns error",
			cartItemID: uuid.New().String(),
			requestBody: UpdateCart_itemDTO{
				ID:   uuid.New(),
				Name: stringPtr("Updated Apple"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateCart_itemDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return((*GetCart_itemDTO)(nil), errors.New("cart item not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "cart item not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.cartItemID != "invalid-uuid" {
				cartItemID, err := uuid.Parse(tt.cartItemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, cartItemID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/cart_items/:id", controller.UpdateCart_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/cart_items/"+tt.cartItemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.IsType(t, map[string]interface{}{}, responseBody)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteCart_item(t *testing.T) {
	tests := []struct {
		name           string
		cartItemID     string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "success - deletes cart item",
			cartItemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   map[string]string{"message": "Cart_item Deleted Successfully"},
		},
		{
			name:           "error - invalid UUID",
			cartItemID:     "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:       "error - service returns error",
			cartItemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("cart item not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "cart item not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.cartItemID != "invalid-uuid" {
				cartItemID, err := uuid.Parse(tt.cartItemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, cartItemID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/cart_items/:id", controller.DeleteCart_item)

			req := httptest.NewRequest(http.MethodDelete, "/cart_items/"+tt.cartItemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.IsType(t, map[string]interface{}{}, responseBody)

			mockSvc.AssertExpectations(t)
		})
	}
}
