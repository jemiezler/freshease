package permissions

import (
	"context"
	"errors"
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

func (m *MockRepository) List(ctx context.Context) ([]*GetPermissionDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetPermissionDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetPermissionDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPermissionDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreatePermissionDTO) (*GetPermissionDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPermissionDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdatePermissionDTO) (*GetPermissionDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPermissionDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockRepository)
		expectedError bool
		expectedCount int
	}{
		{
			name: "success - returns list of permissions",
			mockSetup: func(mockRepo *MockRepository) {
				permissions := []*GetPermissionDTO{
					{
						ID:          uuid.New(),
						Code:        "read_users",
						Description: stringPtr("Permission to read user data"),
					},
					{
						ID:          uuid.New(),
					Code:        "write_users",
					Description: stringPtr("Permission to write user data"),
					},
				}
				mockRepo.On("List", mock.Anything).Return(permissions, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetPermissionDTO{}, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetPermissionDTO)(nil), errors.New("database error"))
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
			name: "success - returns permission",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				permission := &GetPermissionDTO{
					ID:          id,
					Code:        "read_users",
					Description: stringPtr("Permission to read user data"),
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(permission, nil)
			},
			expectedError: false,
		},
		{
			name: "error - permission not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetPermissionDTO)(nil), errors.New("not found"))
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
		dto           CreatePermissionDTO
		mockSetup     func(*MockRepository, CreatePermissionDTO)
		expectedError bool
	}{
		{
			name: "success - creates permission",
			dto: CreatePermissionDTO{
				ID:          uuid.New(),
				Code:        "read_users",
				Description: stringPtr("Permission to read user data"),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreatePermissionDTO) {
				createdPermission := &GetPermissionDTO{
					ID:          dto.ID,
					Code:        dto.Code,
					Description: dto.Description,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreatePermissionDTO) bool {
					return actual.ID == dto.ID &&
						actual.Code == dto.Code &&
						actual.Description == dto.Description
				})).Return(createdPermission, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreatePermissionDTO{
				ID:          uuid.New(),
				Code:        "read_users",
				Description: stringPtr("Permission to read user data"),
			},
			mockSetup: func(mockRepo *MockRepository, dto CreatePermissionDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreatePermissionDTO) bool {
					return actual.ID == dto.ID &&
						actual.Code == dto.Code &&
						actual.Description == dto.Description
				})).Return((*GetPermissionDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.Code, result.Code)
				assert.Equal(t, tt.dto.Description, result.Description)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		id            uuid.UUID
		dto           UpdatePermissionDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdatePermissionDTO)
		expectedError bool
	}{
		{
			name: "success - updates permission",
			id:   uuid.New(),
			dto: UpdatePermissionDTO{
				Code:        stringPtr("updated_code"),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdatePermissionDTO) {
				updatedPermission := &GetPermissionDTO{
					ID:          id,
					Code:        "updated_code",
					Description: stringPtr("Updated description"),
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdatePermissionDTO) bool {
					return actual.ID == id &&
						actual.Code != nil && *actual.Code == "updated_code" &&
						actual.Description != nil && *actual.Description == "Updated description"
				})).Return(updatedPermission, nil)
			},
			expectedError: false,
		},
		{
			name: "success - partial update",
			id:   uuid.New(),
			dto: UpdatePermissionDTO{
				Code: stringPtr("updated_code"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdatePermissionDTO) {
				updatedPermission := &GetPermissionDTO{
					ID:          id,
					Code:        "updated_code",
					Description: stringPtr("Original description"), // unchanged
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdatePermissionDTO) bool {
					return actual.ID == id &&
						actual.Code != nil && *actual.Code == "updated_code" &&
						actual.Description == nil
				})).Return(updatedPermission, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdatePermissionDTO{
				Code: stringPtr("updated_code"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdatePermissionDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdatePermissionDTO) bool {
					return actual.ID == id &&
						actual.Code != nil && *actual.Code == "updated_code"
				})).Return((*GetPermissionDTO)(nil), errors.New("update failed"))
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
			name: "success - deletes permission",
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
