package addresses

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

type MockService struct {
	mock.Mock
}

func (m *MockService) List(ctx context.Context) ([]*GetAddressDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetAddressDTO), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (*GetAddressDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAddressDTO), args.Error(1)
}

func (m *MockService) Create(ctx context.Context, dto CreateAddressDTO) (*GetAddressDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAddressDTO), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, dto UpdateAddressDTO) (*GetAddressDTO, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAddressDTO), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestController_ListAddresses(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - returns addresses list",
			mockSetup: func(mockSvc *MockService) {
				expectedAddresses := []*GetAddressDTO{
					{
						ID:        uuid.New(),
						Line1:     "123 Main St",
						Line2:     "Apt 4B",
						City:      "New York",
						Province:  "NY",
						Country:   "USA",
						PostalCode: "10001",
						IsDefault: true,
					},
					{
						ID:        uuid.New(),
						Line1:     "456 Oak Ave",
						Line2:     "",
						City:      "Los Angeles",
						Province:  "CA",
						Country:   "USA",
						PostalCode: "90210",
						IsDefault: false,
					},
				}
				mockSvc.On("List", mock.Anything).Return(expectedAddresses, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*GetAddressDTO{
				{
					ID:        uuid.New(),
					Line1:     "123 Main St",
					Line2:     "Apt 4B",
					City:      "New York",
					Province:  "NY",
					Country:   "USA",
					PostalCode: "10001",
					IsDefault: true,
				},
				{
					ID:        uuid.New(),
					Line1:     "456 Oak Ave",
					Line2:     "",
					City:      "Los Angeles",
					Province:  "CA",
					Country:   "USA",
					PostalCode: "90210",
					IsDefault: false,
				},
			},
		},
		{
			name: "error - service returns error",
			mockSetup: func(mockSvc *MockService) {
				mockSvc.On("List", mock.Anything).Return([]*GetAddressDTO(nil), errors.New("database error"))
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
			app.Get("/addresses", controller.ListAddresses)

			req := httptest.NewRequest(http.MethodGet, "/addresses", nil)
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

func TestController_GetAddress(t *testing.T) {
	tests := []struct {
		name           string
		addressID      string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "success - returns address by ID",
			addressID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				expectedAddress := &GetAddressDTO{
					ID:        id,
					Line1:     "123 Main St",
					Line2:     "Apt 4B",
					City:      "New York",
					Province:  "NY",
					Country:   "USA",
					PostalCode: "10001",
					IsDefault: true,
				}
				mockSvc.On("Get", mock.Anything, id).Return(expectedAddress, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &GetAddressDTO{
				ID:        uuid.New(),
				Line1:     "123 Main St",
				Line2:     "Apt 4B",
				City:      "New York",
				Province:  "NY",
				Country:   "USA",
				PostalCode: "10001",
				IsDefault: true,
			},
		},
		{
			name:           "error - invalid UUID",
			addressID:      "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:      "error - address not found",
			addressID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Get", mock.Anything, id).Return((*GetAddressDTO)(nil), errors.New("address not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"message": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.addressID != "invalid-uuid" {
				addressID, err := uuid.Parse(tt.addressID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, addressID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Get("/addresses/:id", controller.GetAddress)

			req := httptest.NewRequest(http.MethodGet, "/addresses/"+tt.addressID, nil)
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

func TestController_CreateAddress(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateAddressDTO
		mockSetup      func(*MockService, CreateAddressDTO)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - creates new address",
			requestBody: CreateAddressDTO{
				ID:        uuid.New(),
				Line1:     "789 Pine St",
				Line2:     stringPtr("Unit 2"),
				City:      "Seattle",
				Province:  "WA",
				Country:   "USA",
					PostalCode: "98101",
				IsDefault: false,
			},
			mockSetup: func(mockSvc *MockService, dto CreateAddressDTO) {
				expectedAddress := &GetAddressDTO{
					ID:        dto.ID,
					Line1:     dto.Line1,
					Line2:     *dto.Line2,
					City:      dto.City,
					Province:  dto.Province,
					Country:   dto.Country,
					PostalCode: dto.PostalCode,
					IsDefault: dto.IsDefault,
				}
				mockSvc.On("Create", mock.Anything, dto).Return(expectedAddress, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &GetAddressDTO{
				ID:        uuid.New(),
				Line1:     "789 Pine St",
				Line2:     "Unit 2",
				City:      "Seattle",
				Province:  "WA",
				Country:   "USA",
					PostalCode: "98101",
				IsDefault: false,
			},
		},
		{
			name: "error - service returns error",
			requestBody: CreateAddressDTO{
				ID:        uuid.New(),
				Line1:     "789 Pine St",
				City:      "Seattle",
				Province:  "WA",
				Country:   "USA",
					PostalCode: "98101",
				IsDefault: false,
			},
			mockSetup: func(mockSvc *MockService, dto CreateAddressDTO) {
				mockSvc.On("Create", mock.Anything, dto).Return((*GetAddressDTO)(nil), errors.New("address already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "address already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			tt.mockSetup(mockSvc, tt.requestBody)

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Post("/addresses", controller.CreateAddress)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(jsonBody))
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

func TestController_UpdateAddress(t *testing.T) {
	tests := []struct {
		name           string
		addressID      string
		requestBody    UpdateAddressDTO
		mockSetup      func(*MockService, uuid.UUID, UpdateAddressDTO)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "success - updates address",
			addressID: uuid.New().String(),
			requestBody: UpdateAddressDTO{
				ID:        uuid.New(), // This will be overridden by the service
				Line1:     stringPtr("Updated Street"),
				City:      stringPtr("Updated City"),
				IsDefault: boolPtr(true),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateAddressDTO) {
				expectedAddress := &GetAddressDTO{
					ID:        id,
					Line1:     *dto.Line1,
					Line2:     "",
					City:      *dto.City,
					Province:  "NY",
					Country:   "USA",
					PostalCode: "10001",
					IsDefault: *dto.IsDefault,
				}
				mockSvc.On("Update", mock.Anything, id, dto).Return(expectedAddress, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &GetAddressDTO{
				ID:        uuid.New(),
				Line1:     "Updated Street",
				Line2:     "",
				City:      "Updated City",
				Province:  "NY",
				Country:   "USA",
				PostalCode: "10001",
				IsDefault: true,
			},
		},
		{
			name:           "error - invalid UUID",
			addressID:      "invalid-uuid",
			requestBody:    UpdateAddressDTO{ID: uuid.New()},
			mockSetup:      func(mockSvc *MockService, id uuid.UUID, dto UpdateAddressDTO) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:      "error - service returns error",
			addressID: uuid.New().String(),
			requestBody: UpdateAddressDTO{
				ID:    uuid.New(),
				Line1: stringPtr("Updated Street"),
			},
			mockSetup: func(mockSvc *MockService, id uuid.UUID, dto UpdateAddressDTO) {
				mockSvc.On("Update", mock.Anything, id, dto).Return((*GetAddressDTO)(nil), errors.New("address not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "address not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.addressID != "invalid-uuid" {
				addressID, err := uuid.Parse(tt.addressID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, addressID, tt.requestBody)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Patch("/addresses/:id", controller.UpdateAddress)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/addresses/"+tt.addressID, bytes.NewBuffer(jsonBody))
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

func TestController_DeleteAddress(t *testing.T) {
	tests := []struct {
		name           string
		addressID      string
		mockSetup      func(*MockService, uuid.UUID)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "success - deletes address",
			addressID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   map[string]string{"message": "Address Deleted Successfully"},
		},
		{
			name:           "error - invalid UUID",
			addressID:      "invalid-uuid",
			mockSetup:      func(mockSvc *MockService, id uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid uuid"},
		},
		{
			name:      "error - service returns error",
			addressID: uuid.New().String(),
			mockSetup: func(mockSvc *MockService, id uuid.UUID) {
				mockSvc.On("Delete", mock.Anything, id).Return(errors.New("address not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "address not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if tt.addressID != "invalid-uuid" {
				addressID, err := uuid.Parse(tt.addressID)
				require.NoError(t, err)
				tt.mockSetup(mockSvc, addressID)
			}

			controller := NewController(mockSvc)
			app := fiber.New()
			app.Delete("/addresses/:id", controller.DeleteAddress)

			req := httptest.NewRequest(http.MethodDelete, "/addresses/"+tt.addressID, nil)
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
