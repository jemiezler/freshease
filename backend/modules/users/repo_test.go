package users

import (
	"context"
	"testing"

	"freshease/backend/ent"
	"freshease/backend/ent/enttest"
	"freshease/backend/internal/common/errs"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntRepo_List(t *testing.T) {
	tests := []struct {
		name          string
		setupData     func(*ent.Client) error
		expectedCount int
		expectedError bool
	}{
		{
			name: "success - returns empty list when no users",
			setupData: func(client *ent.Client) error {
				return nil // No setup needed
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "success - returns users list",
			setupData: func(client *ent.Client) error {
				_, err := client.User.Create().
					SetEmail("user1@example.com").
					SetName("User One").
					SetPassword("password123").
					Save(context.Background())
				if err != nil {
					return err
				}

				_, err = client.User.Create().
					SetEmail("user2@example.com").
					SetName("User Two").
					SetPassword("password123").
					Save(context.Background())
				return err
			},
			expectedCount: 2,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
			defer client.Close()

			err := tt.setupData(client)
			require.NoError(t, err)

			repo := NewEntRepo(client)
			ctx := context.Background()

			result, err := repo.List(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}
}

func TestEntRepo_FindByID(t *testing.T) {
	tests := []struct {
		name          string
		setupData     func(*ent.Client) (uuid.UUID, error)
		expectedError bool
		expectedUser  *GetUserDTO
	}{
		{
			name: "success - returns user by ID",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				user, err := client.User.Create().
					SetEmail("user@example.com").
					SetName("Test User").
					SetPassword("password123").
					SetPhone("+1234567890").
					SetBio("This is a longer test bio that meets the minimum length requirement").
					SetAvatar("https://example.com/avatar.jpg").
					SetCover("https://example.com/cover.jpg").
					SetStatus("active").
					Save(context.Background())
				if err != nil {
					return uuid.Nil, err
				}
				return user.ID, nil
			},
			expectedError: false,
			expectedUser: &GetUserDTO{
				Email:  "user@example.com",
				Name:   "Test User",
				Phone:  stringPtr("+1234567890"),
				Bio:    stringPtr("This is a longer test bio that meets the minimum length requirement"),
				Avatar: stringPtr("https://example.com/avatar.jpg"),
				Cover:  stringPtr("https://example.com/cover.jpg"),
				Status: "active",
			},
		},
		{
			name: "error - user not found",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				return uuid.New(), nil // Non-existent ID
			},
			expectedError: true,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
			defer client.Close()

			userID, err := tt.setupData(client)
			require.NoError(t, err)

			repo := NewEntRepo(client)
			ctx := context.Background()

			result, err := repo.FindByID(ctx, userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, userID, result.ID)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.Name, result.Name)
				assert.Equal(t, tt.expectedUser.Status, result.Status)

				// Check optional fields - handle nil values safely
				if tt.expectedUser.Phone != nil {
					assert.Equal(t, tt.expectedUser.Phone, result.Phone)
				} else {
					assert.Nil(t, result.Phone)
				}
				if tt.expectedUser.Bio != nil {
					assert.Equal(t, tt.expectedUser.Bio, result.Bio)
				} else {
					assert.Nil(t, result.Bio)
				}
				if tt.expectedUser.Avatar != nil {
					assert.Equal(t, tt.expectedUser.Avatar, result.Avatar)
				} else {
					assert.Nil(t, result.Avatar)
				}
				if tt.expectedUser.Cover != nil {
					assert.Equal(t, tt.expectedUser.Cover, result.Cover)
				} else {
					assert.Nil(t, result.Cover)
				}
			}
		})
	}
}

func TestEntRepo_Create(t *testing.T) {
	tests := []struct {
		name          string
		createDTO     CreateUserDTO
		expectedError bool
		expectedUser  *GetUserDTO
	}{
		{
			name: "success - creates new user",
			createDTO: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
				Phone:    stringPtr("+1234567890"),
				Bio:      stringPtr("New user bio"),
				Avatar:   stringPtr("https://example.com/avatar.jpg"),
				Cover:    stringPtr("https://example.com/cover.jpg"),
				Status:   stringPtr("active"),
			},
			expectedError: false,
			expectedUser: &GetUserDTO{
				Email:  "newuser@example.com",
				Name:   "New User",
				Phone:  stringPtr("+1234567890"),
				Bio:    stringPtr("New user bio"),
				Avatar: stringPtr("https://example.com/avatar.jpg"),
				Cover:  stringPtr("https://example.com/cover.jpg"),
				Status: "active",
			},
		},
		{
			name: "success - creates user with minimal data",
			createDTO: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "minimal@example.com",
				Password: "password123",
				Name:     "Minimal User",
			},
			expectedError: false,
			expectedUser: &GetUserDTO{
				Email:  "minimal@example.com",
				Name:   "Minimal User",
				Status: "active", // Default status
			},
		},
		{
			name: "error - duplicate email",
			createDTO: CreateUserDTO{
				ID:       uuid.New(),
				Email:    "duplicate@example.com",
				Password: "password123",
				Name:     "Duplicate User",
			},
			expectedError: true,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
			defer client.Close()

			// Setup: create a user with duplicate email for the duplicate test
			if tt.name == "error - duplicate email" {
				_, err := client.User.Create().
					SetEmail("duplicate@example.com").
					SetName("Existing User").
					SetPassword("password123").
					Save(context.Background())
				require.NoError(t, err)
			}

			repo := NewEntRepo(client)
			ctx := context.Background()

			result, err := repo.Create(ctx, &tt.createDTO)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.ID) // Just check it's not nil
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.Name, result.Name)
				assert.Equal(t, tt.expectedUser.Status, result.Status)

				// Check optional fields - handle nil values safely
				if tt.expectedUser.Phone != nil {
					assert.Equal(t, tt.expectedUser.Phone, result.Phone)
				} else {
					assert.Nil(t, result.Phone)
				}
				if tt.expectedUser.Bio != nil {
					assert.Equal(t, tt.expectedUser.Bio, result.Bio)
				} else {
					assert.Nil(t, result.Bio)
				}
				if tt.expectedUser.Avatar != nil {
					assert.Equal(t, tt.expectedUser.Avatar, result.Avatar)
				} else {
					assert.Nil(t, result.Avatar)
				}
				if tt.expectedUser.Cover != nil {
					assert.Equal(t, tt.expectedUser.Cover, result.Cover)
				} else {
					assert.Nil(t, result.Cover)
				}
			}
		})
	}
}

func TestEntRepo_Update(t *testing.T) {
	tests := []struct {
		name          string
		setupData     func(*ent.Client) (uuid.UUID, error)
		updateDTO     UpdateUserDTO
		expectedError bool
		expectedUser  *GetUserDTO
	}{
		{
			name: "success - updates user",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				user, err := client.User.Create().
					SetEmail("user@example.com").
					SetName("Original User").
					SetPassword("password123").
					SetPhone("+1234567890").
					SetStatus("active").
					Save(context.Background())
				if err != nil {
					return uuid.Nil, err
				}
				return user.ID, nil
			},
			updateDTO: UpdateUserDTO{
				Email: stringPtr("updated@example.com"),
				Name:  stringPtr("Updated User"),
				Phone: stringPtr("+9876543210"),
				Bio:   stringPtr("Updated bio"),
			},
			expectedError: false,
			expectedUser: &GetUserDTO{
				Email:  "updated@example.com",
				Name:   "Updated User",
				Phone:  stringPtr("+9876543210"),
				Bio:    stringPtr("Updated bio"),
				Status: "active",
			},
		},
		{
			name: "success - partial update",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				user, err := client.User.Create().
					SetEmail("user@example.com").
					SetName("Original User").
					SetPassword("password123").
					SetStatus("active").
					Save(context.Background())
				if err != nil {
					return uuid.Nil, err
				}
				return user.ID, nil
			},
			updateDTO: UpdateUserDTO{
				Email: stringPtr("updated@example.com"),
			},
			expectedError: false,
			expectedUser: &GetUserDTO{
				Email:  "updated@example.com",
				Name:   "Original User", // Should remain unchanged
				Status: "active",
			},
		},
		{
			name: "error - no fields to update",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				user, err := client.User.Create().
					SetEmail("user@example.com").
					SetName("Original User").
					SetPassword("password123").
					SetStatus("active").
					Save(context.Background())
				if err != nil {
					return uuid.Nil, err
				}
				return user.ID, nil
			},
			updateDTO: UpdateUserDTO{
				// Empty update DTO
			},
			expectedError: true,
			expectedUser:  nil,
		},
		{
			name: "error - user not found",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				return uuid.New(), nil // Non-existent ID
			},
			updateDTO: UpdateUserDTO{
				Email: stringPtr("updated@example.com"),
			},
			expectedError: true,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
			defer client.Close()

			userID, err := tt.setupData(client)
			require.NoError(t, err)

			tt.updateDTO.ID = userID
			repo := NewEntRepo(client)
			ctx := context.Background()

			result, err := repo.Update(ctx, &tt.updateDTO)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)

				if tt.name == "error - no fields to update" {
					assert.Equal(t, errs.NoFieldsToUpdate, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, userID, result.ID)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.Name, result.Name)
				assert.Equal(t, tt.expectedUser.Status, result.Status)

				// Check optional fields
				if tt.expectedUser.Phone != nil {
					assert.Equal(t, tt.expectedUser.Phone, result.Phone)
				}
				if tt.expectedUser.Bio != nil {
					assert.Equal(t, tt.expectedUser.Bio, result.Bio)
				}
			}
		})
	}
}

func TestEntRepo_Delete(t *testing.T) {
	tests := []struct {
		name          string
		setupData     func(*ent.Client) (uuid.UUID, error)
		expectedError bool
	}{
		{
			name: "success - deletes user",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				user, err := client.User.Create().
					SetEmail("user@example.com").
					SetName("Test User").
					SetPassword("password123").
					SetStatus("active").
					Save(context.Background())
				if err != nil {
					return uuid.Nil, err
				}
				return user.ID, nil
			},
			expectedError: false,
		},
		{
			name: "error - user not found",
			setupData: func(client *ent.Client) (uuid.UUID, error) {
				return uuid.New(), nil // Non-existent ID
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
			defer client.Close()

			userID, err := tt.setupData(client)
			require.NoError(t, err)

			repo := NewEntRepo(client)
			ctx := context.Background()

			err = repo.Delete(ctx, userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify user is deleted
				_, err = client.User.Get(ctx, userID)
				assert.Error(t, err) // Should not exist anymore
			}
		})
	}
}
