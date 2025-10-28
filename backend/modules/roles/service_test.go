package roles

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

func (m *MockRepository) List(ctx context.Context) ([]*GetRoleDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GetRoleDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetRoleDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRoleDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateRoleDTO) (*GetRoleDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRoleDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateRoleDTO) (*GetRoleDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRoleDTO), args.Error(1)
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
			name: "success - returns list of roles",
			mockSetup: func(mockRepo *MockRepository) {
				roles := []*GetRoleDTO{
					{
						ID:          uuid.New(),
						Name:        "admin",
						Description: "Administrator role",
					},
					{
						ID:          uuid.New(),
						Name:        "user",
						Description: "Regular user role",
					},
				}
				mockRepo.On("List", mock.Anything).Return(roles, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "success - returns empty list",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetRoleDTO{}, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return(([]*GetRoleDTO)(nil), errors.New("database error"))
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
			name: "success - returns role",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				role := &GetRoleDTO{
					ID:          id,
					Name:        "admin",
					Description: "Administrator role",
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(role, nil)
			},
			expectedError: false,
		},
		{
			name: "error - role not found",
			id:   uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetRoleDTO)(nil), errors.New("not found"))
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
		dto           CreateRoleDTO
		mockSetup     func(*MockRepository, CreateRoleDTO)
		expectedError bool
	}{
		{
			name: "success - creates role",
			dto: CreateRoleDTO{
				ID:          uuid.New(),
				Name:        "admin",
				Description: "Administrator role",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateRoleDTO) {
				createdRole := &GetRoleDTO{
					ID:          dto.ID,
					Name:        dto.Name,
					Description: dto.Description,
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateRoleDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description
				})).Return(createdRole, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			dto: CreateRoleDTO{
				ID:          uuid.New(),
				Name:        "admin",
				Description: "Administrator role",
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateRoleDTO) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(actual *CreateRoleDTO) bool {
					return actual.ID == dto.ID &&
						actual.Name == dto.Name &&
						actual.Description == dto.Description
				})).Return((*GetRoleDTO)(nil), errors.New("creation failed"))
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
				assert.Equal(t, tt.dto.Name, result.Name)
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
		dto           UpdateRoleDTO
		mockSetup     func(*MockRepository, uuid.UUID, UpdateRoleDTO)
		expectedError bool
	}{
		{
			name: "success - updates role",
			id:   uuid.New(),
			dto: UpdateRoleDTO{
				Name:        stringPtr("updated_admin"),
				Description: stringPtr("Updated description"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateRoleDTO) {
				updatedRole := &GetRoleDTO{
					ID:          id,
					Name:        "updated_admin",
					Description: "Updated description",
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateRoleDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "updated_admin" &&
						actual.Description != nil && *actual.Description == "Updated description"
				})).Return(updatedRole, nil)
			},
			expectedError: false,
		},
		{
			name: "success - partial update",
			id:   uuid.New(),
			dto: UpdateRoleDTO{
				Name: stringPtr("updated_admin"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateRoleDTO) {
				updatedRole := &GetRoleDTO{
					ID:          id,
					Name:        "updated_admin",
					Description: "Original description", // unchanged
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateRoleDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "updated_admin" &&
						actual.Description == nil
				})).Return(updatedRole, nil)
			},
			expectedError: false,
		},
		{
			name: "error - repository returns error",
			id:   uuid.New(),
			dto: UpdateRoleDTO{
				Name: stringPtr("updated_admin"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateRoleDTO) {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(actual *UpdateRoleDTO) bool {
					return actual.ID == id &&
						actual.Name != nil && *actual.Name == "updated_admin"
				})).Return((*GetRoleDTO)(nil), errors.New("update failed"))
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
			name: "success - deletes role",
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
