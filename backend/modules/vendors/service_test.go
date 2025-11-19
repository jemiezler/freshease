package vendors

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"freshease/backend/modules/uploads"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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

func (m *MockUploadsService) GetImage(ctx context.Context, objectName string) (io.ReadCloser, *minio.ObjectInfo, error) {
	args := m.Called(ctx, objectName)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(io.ReadCloser), nil, args.Error(2)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(*minio.ObjectInfo), args.Error(2)
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

func TestService_UploadVendorLogo(t *testing.T) {
	tests := []struct {
		name           string
		file           *multipart.FileHeader
		mockSetup      func(*MockUploadsService, *multipart.FileHeader)
		expectedResult string
		expectedError  error
	}{
		{
			name: "success - uploads vendor logo",
			file: &multipart.FileHeader{
				Filename: "logo.jpg",
				Size:     1536,
			},
			mockSetup: func(mockUploads *MockUploadsService, file *multipart.FileHeader) {
				mockUploads.On("UploadImage", mock.Anything, file, "vendors/logos").Return("vendors/logos/logo.jpg", nil)
			},
			expectedResult: "vendors/logos/logo.jpg",
			expectedError:  nil,
		},
		{
			name: "error - upload service returns error",
			file: &multipart.FileHeader{
				Filename: "logo.jpg",
				Size:     1536,
			},
			mockSetup: func(mockUploads *MockUploadsService, file *multipart.FileHeader) {
				mockUploads.On("UploadImage", mock.Anything, file, "vendors/logos").Return("", errors.New("upload failed"))
			},
			expectedResult: "",
			expectedError:  errors.New("upload failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockUploads, tt.file)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.UploadVendorLogo(ctx, tt.file)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockUploads.AssertExpectations(t)
		})
	}
}

func TestService_GetVendorImageURL(t *testing.T) {
	tests := []struct {
		name           string
		objectName     string
		mockSetup      func(*MockUploadsService, string)
		expectedResult string
		expectedError  error
	}{
		{
			name:       "success - returns vendor image URL",
			objectName: "vendors/logos/logo.jpg",
			mockSetup: func(mockUploads *MockUploadsService, objectName string) {
				mockUploads.On("GetImageURL", mock.Anything, objectName).Return("https://example.com/vendors/logos/logo.jpg", nil)
			},
			expectedResult: "https://example.com/vendors/logos/logo.jpg",
			expectedError:  nil,
		},
		{
			name:       "error - upload service returns error",
			objectName: "vendors/logos/logo.jpg",
			mockSetup: func(mockUploads *MockUploadsService, objectName string) {
				mockUploads.On("GetImageURL", mock.Anything, objectName).Return("", errors.New("failed to generate URL"))
			},
			expectedResult: "",
			expectedError:  errors.New("failed to generate URL"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockUploads, tt.objectName)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.GetVendorImageURL(ctx, tt.objectName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockUploads.AssertExpectations(t)
		})
	}
}
