package bundle_items

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

func (m *MockService) List(ctx context.Context) ([]*GetBundle_itemDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetBundle_itemDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetBundle_itemDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundle_itemDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateBundle_itemDTO) (*GetBundle_itemDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundle_itemDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateBundle_itemDTO) (*GetBundle_itemDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBundle_itemDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListBundle_items(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name: "success - returns bundle items list",
			mockSetup: func(mockSvc *MockService) {
				expectedItems := []*GetBundle_itemDTO{
					{
						ID:        uuid.New(),
						Qty:       2,
						BundleID:  uuid.New(),
						ProductID: uuid.New(),
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedItems, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return(([]*GetBundle_itemDTO)(nil), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/bundle-items", controller.ListBundle_items)

			req := httptest.NewRequest(http.MethodGet, "/bundle-items", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_GetBundle_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - returns bundle item by ID",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				item := &GetBundle_itemDTO{
					ID:        id,
					Qty:       2,
					BundleID:  uuid.New(),
					ProductID: uuid.New(),
				}
				mockSvc.On("Get", mock.Anything, id).Return(item, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - bundle item not found",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetBundle_itemDTO)(nil), errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.itemID != "invalid-uuid" {
				itemID, err := uuid.Parse(tt.itemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, itemID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/bundle-items/:id", controller.GetBundle_item)

			req := httptest.NewRequest(http.MethodGet, "/bundle-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_CreateBundle_item(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateBundle_itemDTO
		mockSetup      func(*MockService, CreateBundle_itemDTO)
		expectedStatus int
	}{
		{
			name: "success - creates new bundle item",
			requestBody: CreateBundle_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				BundleID:  uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateBundle_itemDTO) {
				expectedItem := &GetBundle_itemDTO{
					ID:        dto.ID,
					Qty:       dto.Qty,
					BundleID:  dto.BundleID,
					ProductID: dto.ProductID,
				}
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateBundle_itemDTO) bool {
					return actual.ID == dto.ID && actual.Qty == dto.Qty
				})).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - service returns error",
			requestBody: CreateBundle_itemDTO{
				ID:        uuid.New(),
				Qty:       2,
				BundleID:  uuid.New(),
				ProductID: uuid.New(),
			},
			mockSetup: func(mockSvc *MockService, dto CreateBundle_itemDTO) {
				mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(actual CreateBundle_itemDTO) bool {
					return actual.ID == dto.ID
				})).Return((*GetBundle_itemDTO)(nil), errors.New("creation failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/bundle-items", controller.CreateBundle_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/bundle-items", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_UpdateBundle_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    UpdateBundle_itemDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateBundle_itemDTO)
		expectedStatus int
	}{
		{
			name:   "success - updates bundle item",
			itemID: uuid.New().String(),
			requestBody: UpdateBundle_itemDTO{
				ID:  uuid.New(), // This will be overwritten by controller
				Qty: intPtr(5),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateBundle_itemDTO) {
				expectedItem := &GetBundle_itemDTO{
					ID:        id,
					Qty:       5,
					BundleID:  uuid.New(),
					ProductID: uuid.New(),
				}
				// Service sets dto.ID = id, but mock receives DTO as controller passes it
				// Use mock.Anything for DTO since ID will be overwritten by service anyway
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return(expectedItem, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			requestBody:    UpdateBundle_itemDTO{},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateBundle_itemDTO) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			requestBody: UpdateBundle_itemDTO{
				ID:  uuid.New(), // This will be overwritten by controller
				Qty: intPtr(5),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateBundle_itemDTO) {
				mockSvc.On("Update", mock.Anything, id, mock.Anything).Return((*GetBundle_itemDTO)(nil), errors.New("update failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.itemID != "invalid-uuid" {
				itemID, err := uuid.Parse(tt.itemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, itemID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/bundle-items/:id", controller.UpdateBundle_item)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/bundle-items/"+tt.itemID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestController_DeleteBundle_item(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success - deletes bundle item",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "error - invalid UUID",
			itemID:         "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - service returns error",
			itemID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.itemID != "invalid-uuid" {
				itemID, err := uuid.Parse(tt.itemID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, itemID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/bundle-items/:id", controller.DeleteBundle_item)

			req := httptest.NewRequest(http.MethodDelete, "/bundle-items/"+tt.itemID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}

