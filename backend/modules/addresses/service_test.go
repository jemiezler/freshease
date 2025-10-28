package addresses

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context) ([]*GetAddressDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*GetAddressDTO), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*GetAddressDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAddressDTO), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, dto *CreateAddressDTO) (*GetAddressDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAddressDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, dto *UpdateAddressDTO) (*GetAddressDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAddressDTO), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockRepository)
		expectedResult []*GetAddressDTO
		expectedError  error
	}{
		{
			name: "success - returns addresses list",
			mockSetup: func(mockRepo *MockRepository) {
				expectedAddresses := []*GetAddressDTO{
					{
						ID:        uuid.New(),
						Line1:     "123 Main St",
						Line2:     "Apt 4B",
						City:      "New York",
						Province:  "NY",
						Country:   "USA",
						Zip:       "10001",
						IsDefault: true,
					},
					{
						ID:        uuid.New(),
						Line1:     "456 Oak Ave",
						Line2:     "",
						City:      "Los Angeles",
						Province:  "CA",
						Country:   "USA",
						Zip:       "90210",
						IsDefault: false,
					},
				}
				mockRepo.On("List", mock.Anything).Return(expectedAddresses, nil)
			},
			expectedResult: []*GetAddressDTO{
				{
					ID:        uuid.New(),
					Line1:     "123 Main St",
					Line2:     "Apt 4B",
					City:      "New York",
					Province:  "NY",
					Country:   "USA",
					Zip:       "10001",
					IsDefault: true,
				},
				{
					ID:        uuid.New(),
					Line1:     "456 Oak Ave",
					Line2:     "",
					City:      "Los Angeles",
					Province:  "CA",
					Country:   "USA",
					Zip:       "90210",
					IsDefault: false,
				},
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			mockSetup: func(mockRepo *MockRepository) {
				mockRepo.On("List", mock.Anything).Return([]*GetAddressDTO(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
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
		addressID      uuid.UUID
		mockSetup      func(*MockRepository, uuid.UUID)
		expectedResult *GetAddressDTO
		expectedError  error
	}{
		{
			name:      "success - returns address by ID",
			addressID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				expectedAddress := &GetAddressDTO{
					ID:        id,
					Line1:     "123 Main St",
					Line2:     "Apt 4B",
					City:      "New York",
					Province:  "NY",
					Country:   "USA",
					Zip:       "10001",
					IsDefault: true,
				}
				mockRepo.On("FindByID", mock.Anything, id).Return(expectedAddress, nil)
			},
			expectedResult: &GetAddressDTO{
				ID:        uuid.New(),
				Line1:     "123 Main St",
				Line2:     "Apt 4B",
				City:      "New York",
				Province:  "NY",
				Country:   "USA",
				Zip:       "10001",
				IsDefault: true,
			},
			expectedError: nil,
		},
		{
			name:      "error - address not found",
			addressID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("FindByID", mock.Anything, id).Return((*GetAddressDTO)(nil), errors.New("address not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("address not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.addressID)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Get(ctx, tt.addressID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.addressID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name           string
		createDTO      CreateAddressDTO
		mockSetup      func(*MockRepository, CreateAddressDTO)
		expectedResult *GetAddressDTO
		expectedError  error
	}{
		{
			name: "success - creates new address",
			createDTO: CreateAddressDTO{
				ID:        uuid.New(),
				Line1:     "789 Pine St",
				Line2:     stringPtr("Unit 2"),
				City:      "Seattle",
				Province:  "WA",
				Country:   "USA",
				Zip:       "98101",
				IsDefault: false,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateAddressDTO) {
				expectedAddress := &GetAddressDTO{
					ID:        dto.ID,
					Line1:     dto.Line1,
					Line2:     *dto.Line2,
					City:      dto.City,
					Province:  dto.Province,
					Country:   dto.Country,
					Zip:       dto.Zip,
					IsDefault: dto.IsDefault,
				}
				mockRepo.On("Create", mock.Anything, &dto).Return(expectedAddress, nil)
			},
			expectedResult: &GetAddressDTO{
				ID:        uuid.New(),
				Line1:     "789 Pine St",
				Line2:     "Unit 2",
				City:      "Seattle",
				Province:  "WA",
				Country:   "USA",
				Zip:       "98101",
				IsDefault: false,
			},
			expectedError: nil,
		},
		{
			name: "error - repository returns error",
			createDTO: CreateAddressDTO{
				ID:        uuid.New(),
				Line1:     "789 Pine St",
				City:      "Seattle",
				Province:  "WA",
				Country:   "USA",
				Zip:       "98101",
				IsDefault: false,
			},
			mockSetup: func(mockRepo *MockRepository, dto CreateAddressDTO) {
				mockRepo.On("Create", mock.Anything, &dto).Return((*GetAddressDTO)(nil), errors.New("address already exists"))
			},
			expectedResult: nil,
			expectedError:  errors.New("address already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.createDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Create(ctx, tt.createDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.createDTO.Line1, result.Line1)
				assert.Equal(t, tt.createDTO.City, result.City)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name           string
		addressID      uuid.UUID
		updateDTO      UpdateAddressDTO
		mockSetup      func(*MockRepository, uuid.UUID, UpdateAddressDTO)
		expectedResult *GetAddressDTO
		expectedError  error
	}{
		{
			name:      "success - updates address",
			addressID: uuid.New(),
			updateDTO: UpdateAddressDTO{
				Line1:     stringPtr("Updated Street"),
				City:      stringPtr("Updated City"),
				IsDefault: boolPtr(true),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateAddressDTO) {
				expectedAddress := &GetAddressDTO{
					ID:        id,
					Line1:     *dto.Line1,
					Line2:     "",
					City:      *dto.City,
					Province:  "NY",
					Country:   "USA",
					Zip:       "10001",
					IsDefault: *dto.IsDefault,
				}
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *UpdateAddressDTO) bool {
					return u.ID == id && *u.Line1 == *dto.Line1 && *u.City == *dto.City
				})).Return(expectedAddress, nil)
			},
			expectedResult: &GetAddressDTO{
				ID:        uuid.New(),
				Line1:     "Updated Street",
				Line2:     "",
				City:      "Updated City",
				Province:  "NY",
				Country:   "USA",
				Zip:       "10001",
				IsDefault: true,
			},
			expectedError: nil,
		},
		{
			name:      "error - repository returns error",
			addressID: uuid.New(),
			updateDTO: UpdateAddressDTO{
				Line1: stringPtr("Updated Street"),
			},
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID, dto UpdateAddressDTO) {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return((*GetAddressDTO)(nil), errors.New("address not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("address not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.addressID, tt.updateDTO)

			service := NewService(mockRepo)
			ctx := context.Background()

			result, err := service.Update(ctx, tt.addressID, tt.updateDTO)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.addressID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		addressID     uuid.UUID
		mockSetup     func(*MockRepository, uuid.UUID)
		expectedError error
	}{
		{
			name:      "success - deletes address",
			addressID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "error - repository returns error",
			addressID: uuid.New(),
			mockSetup: func(mockRepo *MockRepository, id uuid.UUID) {
				mockRepo.On("Delete", mock.Anything, id).Return(errors.New("address not found"))
			},
			expectedError: errors.New("address not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo, tt.addressID)

			service := NewService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.addressID)

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
