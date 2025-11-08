package vendors

import (
	"context"
	"errors"
	"mime/multipart"
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

// MockUploadsService is a mock implementation of uploads.Service
type MockUploadsService struct {
	mock.Mock
}

func (m *MockUploadsService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, folder)
	return args.String(0), args.Error(1)
}

func (m *MockUploadsService) DeleteImage(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

func (m *MockUploadsService) GetImageURL(ctx context.Context, objectName string) (string, error) {
	args := m.Called(ctx, objectName)
	return args.String(0), args.Error(1)
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
					Contact: stringPtr("vendor1@example.com"),
					},
					{
						ID:        uuid.New(),
						Name:      stringPtr("Vendor 2"),
					Contact: stringPtr("vendor2@example.com"),
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo, mockUploads)
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
					Contact: stringPtr("test@example.com"),
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo, mockUploads)
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
				Contact: stringPtr("test@example.com"),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateVendorDTO) {
				createdVendor := &GetVendorDTO{
					ID:        dto.ID,
					Contact: dto.Contact,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.Contact != nil
				})).Return(createdVendor, nil)
			},
			expectedError: false,
		},
		{
			name: "success - creates vendor with full data",
			dto: CreateVendorDTO{
				ID:      uuid.New(),
				Name:    stringPtr("Test Vendor"),
				Contact: stringPtr("test@example.com"),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateVendorDTO) {
				createdVendor := &GetVendorDTO{
					ID:      dto.ID,
					Name:    dto.Name,
					Contact: dto.Contact,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.Contact != nil
				})).Return(createdVendor, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateVendorDTO{
				ID:       uuid.New(),
				Contact: stringPtr("test@example.com"),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateVendorDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateVendorDTO) bool {
					return actual.ID == dto.ID && actual.Contact != nil
				})).Return((*GetVendorDTO)(nil), errors.New("creation failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.dto)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.dto)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.dto.ID, result.ID)
				assert.Equal(t, tt.dto.Contact, result.Contact)
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
				Contact: stringPtr("updated@example.com"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateVendorDTO) {
				updatedVendor := &GetVendorDTO{
					ID:        id,
					Name:      stringPtr("Updated Vendor"),
					Contact: stringPtr("updated@example.com"),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateVendorDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Vendor" &&
						actual.Contact != nil && *actual.Contact == "updated@example.com"
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
					Contact: stringPtr("test@example.com"), // unchanged
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateVendorDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "Updated Vendor" &&
						actual.Contact == nil
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.id, tt.dto)

			service := NewService(mockRepo, mockUploads)
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
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo, mockUploads)
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
