package vendors

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

func (m *MockRepository) List(ctx context.Context) ([]*GetVendorDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetVendorDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetVendorDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVendorDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateVendorDTO) (*GetVendorDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVendorDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateVendorDTO) (*GetVendorDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVendorDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockRepository)
		expectedError bool
		expectedCount int
	}{
		{
			name: "success - returns list of vendors",
			mockSetup: func(mockRepo *MockRepository) {
				vendors := []*GetVendorDTO{
					{
						ID:        uuid.New(),
						Name:      stringPtr("Vendor 1"),
						Email:     stringPtr("vendor1@example.com"),
						IsActive:  "true",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Name:      stringPtr("Vendor 2"),
						Email:     stringPtr("vendor2@example.com"),
						IsActive:  "true",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				mockRepo.On("List", mock.Anything).Return(vendors, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetVendorDTO{}, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetVendorDTO)(nil), errors.New("database error"))
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.List(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - returns vendor",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				vendor := &GetVendorDTO{
					ID:        id,
					Name:      stringPtr("Test Vendor"),
					Email:     stringPtr("test@example.com"),
					IsActive:  "true",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(vendor, nil)
			},
			expectedError: false,
		},
		{
			name: "error - vendor not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetVendorDTO)(nil), errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Get(ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		dto           CreateVendorDTO
		mockSetup     func(*MockRepository, CreateVendorDTO)
		expectedError bool
	}{
		{
			name: "success - creates vendor with minimal data",
			dto: CreateVendorDTO{
				ID:       uuid.New(),
				IsActive: "true",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateVendorDTO) {
				createdVendor := &GetVendorDTO{
					ID:        dto.ID,
					IsActive:  dto.IsActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.IsActive == dto.IsActive
				})).Return(createdVendor, nil)
			},
			expectedError: false,
		},
		{
			name: "success - creates vendor with full data",
			dto: CreateVendorDTO{
				ID:          uuid.New(),
				Name:        stringPtr("Test Vendor"),
				Email:       stringPtr("test@example.com"),
				Phone:       stringPtr("123-456-7890"),
				Address:     stringPtr("123 Main St"),
				City:        stringPtr("Test City"),
				State:       stringPtr("Test State"),
				Country:     stringPtr("Test Country"),
				PostalCode:  stringPtr("12345"),
				Website:     stringPtr("https://test.com"),
				LogoURL:     stringPtr("https://test.com/logo.png"),
				Description: stringPtr("Test vendor description"),
				IsActive:    "true",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateVendorDTO) {
				createdVendor := &GetVendorDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Email:       dto.Email,
					Phone:       dto.Phone,
					Address:     dto.Address,
					City:        dto.City,
					State:       dto.State,
					Country:     dto.Country,
					PostalCode:  dto.PostalCode,
					Website:     dto.Website,
					LogoURL:     dto.LogoURL,
					Description: dto.Description,
					IsActive:    dto.IsActive,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.IsActive == dto.IsActive
				})).Return(createdVendor, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateVendorDTO{
				ID:       uuid.New(),
				IsActive: "true",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateVendorDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.IsActive == dto.IsActive
				})).Return((*GetVendorDTO)(nil), errors.New("creation failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.dto)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.dto.ID, result.ID)
				assert.Equal(t, tt.dto.IsActive, result.IsActive)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		dto           UpdateVendorDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdateVendorDTO)
		expectedError bool
	}{
		{
			name: "success - updates vendor",
			id:   uuid.New(),
			dto: UpdateVendorDTO{
				Name:     stringPtr("Updated Vendor"),
				Email:    stringPtr("updated@example.com"),
				IsActive: stringPtr("false"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateVendorDTO) {
				updatedVendor := &GetVendorDTO{
					ID:        id,
					Name:      stringPtr("Updated Vendor"),
					Email:     stringPtr("updated@example.com"),
					IsActive:  "false",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateVendorDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Vendor" &&
						actual.Email != nil && *actual.Email == "updated@example.com" &&
						actual.IsActive != nil && *actual.IsActive == "false"
				})).Return(updatedVendor, nil)
			},
			expectedError: false,
		},
		{
			name: "success - partial update",
			id:   uuid.New(),
			dto: UpdateVendorDTO{
				Name: stringPtr("Updated Vendor"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateVendorDTO) {
				updatedVendor := &GetVendorDTO{
					ID:        id,
					Name:      stringPtr("Updated Vendor"),
					IsActive:  "true", // unchanged
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateVendorDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Vendor" &&
						actual.Email == nil &&
						actual.IsActive == nil
				})).Return(updatedVendor, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateVendorDTO{
				Name: stringPtr("Updated Vendor"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateVendorDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateVendorDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Vendor"
				})).Return((*GetVendorDTO)(nil), errors.New("update failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id, tt.dto)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Update(ctx, tt.id, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError bool
	}{
		{
			name: "success - deletes vendor",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("delete failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
