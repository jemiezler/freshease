package users

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context) ([]*GetUserDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetUserDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetUserDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, u *CreateUserDTO) (*GetUserDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *UpdateUserDTO) (*GetUserDTO, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserDTO), args.Error(1)
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
		name           string
		mockSetup      func(*MockRepository)
		expectedResult []*GetUserDTO
		expectedError  error
	}{
		{
			name: "success - returns users list",
			mockSetup: func(mockRepo *MockRepository) {
				expectedUsers := []*GetUserDTO{
					{
						ID:     uuid.New(),
						Email:  "user1@example.com",
						Name:   "User One",
						Status: stringPtr("active"),
					},
					{
						ID:     uuid.New(),
						Email:  "user2@example.com",
						Name:   "User Two",
						Status: stringPtr("active"),
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedUsers, nil)
			},
			expectedResult: []*GetUserDTO{
				{
					ID:     uuid.New(),
					Email:  "user1@example.com",
					Name:   "User One",
					Status: stringPtr("active"),
				},
				{
					ID:     uuid.New(),
					Email:  "user2@example.com",
					Name:   "User Two",
					Status: stringPtr("active"),
				},
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetUserDTO(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
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

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		mockSetup      func(*MockRepository, uuid.UUID)
		expectedResult *GetUserDTO
		expectedError  error
	}{
		{
			name:   "success - returns user by ID",
			userID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				expectedUser := &GetUserDTO{
					ID:     id,
					Email:  "user@example.com",
					Name:   "Test User",
					Status: stringPtr("active"),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedUser, nil)
			},
			expectedResult: &GetUserDTO{
				ID:     uuid.New(),
				Email:  "user@example.com",
				Name:   "Test User",
				Status: stringPtr("active"),
			},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetUserDTO)(nil), errors.New("user not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.userID)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.Get(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name           string
		createDTO      CreateUserDTO
		mockSetup      func(*MockRepository, CreateUserDTO)
		expectedResult *GetUserDTO
		expectedError  error
	}{
		{
			name: "success - creates new user",
			createDTO: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
				Status:   stringPtr("active"),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateUserDTO) {
				expectedUser := &GetUserDTO{
					ID:     dto.ID,
					Email:  dto.Email,
					Name:   dto.Name,
					Status: dto.Status,
				}
				mockRepo.On("Create", mock.Anything, &dto).Return(expectedUser, nil)
			},
			expectedResult: &GetUserDTO{
				ID:     uuid.New(),
				Email:  "newuser@example.com",
				Name:   "New User",
				Status: stringPtr("active"),
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateUserDTO) {
				mockRepo.On("Create", mock.Anything, &dto).Return((*GetUserDTO)(nil), errors.New("email already exists"))
			},
			expectedResult: nil,
			expectedError:  errors.New("email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.createDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.createDTO.Email, result.Email)
				assert.Equal(t, tt.createDTO.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		updateDTO      UpdateUserDTO
		mockSetup      func(*MockRepository, uuid.UUID, UpdateUserDTO)
		expectedResult *GetUserDTO
		expectedError  error
	}{
		{
			name:   "success - updates user",
			userID: uuid.New(),
			updateDTO: UpdateUserDTO{
				Email: stringPtr("updated@example.com"),
				Name:  stringPtr("Updated User"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateUserDTO) {
				expectedUser := &GetUserDTO{
					ID:     id,
					Email:  *dto.Email,
					Name:   *dto.Name,
					Status: stringPtr("active"),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *UpdateUserDTO) bool {
					return u.ID == id && *u.Email == *dto.Email && *u.Name == *dto.Name
				})).Return(expectedUser, nil)
			},
			expectedResult: &GetUserDTO{
				ID:     uuid.New(),
				Email:  "updated@example.com",
				Name:   "Updated User",
				Status: stringPtr("active"),
			},
			expectedError: nil,
		},
		{
			name:   "error - repository returns error",
			userID: uuid.New(),
			updateDTO: UpdateUserDTO{
				Email: stringPtr("updated@example.com"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateUserDTO) {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return((*GetUserDTO)(nil), errors.New("user not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.userID, tt.updateDTO)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			result, err := service.Update(ctx, tt.userID, tt.updateDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError error
	}{
		{
			name:   "success - deletes user",
			userID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error - repository returns error",
			userID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("user not found"))
			},
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockUploads := new(MockUploadsService)
			tt.mockSetup(mockRepo, tt.userID)

			service := NewService(mockRepo, mockUploads)
			ctx := context.Background()

			err := service.Delete(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
