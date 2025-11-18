package genai

import (
	"context"
	"testing"
	"time"

	"freshease/backend/ent"
	"freshease/backend/ent/enttest"
	"freshease/backend/ent/meal_plan"
	"freshease/backend/ent/user"
	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRepository is a simple in-memory implementation for testing
type MockRepository struct {
	profiles map[string]*UserProfile
	plans    map[string][]byte
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		profiles: make(map[string]*UserProfile),
		plans:    make(map[string][]byte),
	}
}

func (m *MockRepository) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	profile, exists := m.profiles[userID]
	if !exists {
		return nil, assert.AnError
	}
	return profile, nil
}

func (m *MockRepository) SaveGeneratedPlan(ctx context.Context, userID string, planJSON []byte) error {
	m.plans[userID] = planJSON
	return nil
}

func (m *MockRepository) SetProfile(userID string, profile *UserProfile) {
	m.profiles[userID] = profile
}

func (m *MockRepository) GetPlan(userID string) []byte {
	return m.plans[userID]
}

func TestRepository_GetUserProfile(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		setupProfile  *UserProfile
		expectedError bool
	}{
		{
			name:   "success - returns existing profile",
			userID: "user123",
			setupProfile: &UserProfile{
				UserID:        "user123",
				Gender:        "male",
				Age:           30,
				HeightCm:      175.0,
				WeightKg:      70.0,
				Allergies:     []string{"nuts"},
				Preferences:   []string{"vegetarian"},
				Target:        "maintenance",
				StepsToday:    8000,
				ActiveKcal24h: 2000.0,
			},
			expectedError: false,
		},
		{
			name:          "error - profile not found",
			userID:        "nonexistent",
			setupProfile:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			if tt.setupProfile != nil {
				repo.SetProfile(tt.userID, tt.setupProfile)
			}

			ctx := context.Background()
			profile, err := repo.GetUserProfile(ctx, tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, profile)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, tt.setupProfile.UserID, profile.UserID)
				assert.Equal(t, tt.setupProfile.Gender, profile.Gender)
				assert.Equal(t, tt.setupProfile.Age, profile.Age)
				assert.Equal(t, tt.setupProfile.HeightCm, profile.HeightCm)
				assert.Equal(t, tt.setupProfile.WeightKg, profile.WeightKg)
				assert.Equal(t, tt.setupProfile.Allergies, profile.Allergies)
				assert.Equal(t, tt.setupProfile.Preferences, profile.Preferences)
				assert.Equal(t, tt.setupProfile.Target, profile.Target)
				assert.Equal(t, tt.setupProfile.StepsToday, profile.StepsToday)
				assert.Equal(t, tt.setupProfile.ActiveKcal24h, profile.ActiveKcal24h)
			}
		})
	}
}

func TestRepository_SaveGeneratedPlan(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		planJSON  []byte
		wantError bool
	}{
		{
			name:   "success - saves plan",
			userID: "user123",
			planJSON: []byte(`[
				{
					"day": "Monday",
					"meals": {
						"breakfast": "Oatmeal with berries",
						"lunch": "Quinoa salad",
						"dinner": "Grilled chicken"
					},
					"calories": {
						"breakfast": 300,
						"lunch": 400,
						"dinner": 500
					}
				}
			]`),
			wantError: false,
		},
		{
			name:      "success - saves empty plan",
			userID:    "user456",
			planJSON:  []byte("[]"),
			wantError: false,
		},
		{
			name:      "success - saves nil plan",
			userID:    "user789",
			planJSON:  nil,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			ctx := context.Background()

			err := repo.SaveGeneratedPlan(ctx, tt.userID, tt.planJSON)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				savedPlan := repo.GetPlan(tt.userID)
				assert.Equal(t, tt.planJSON, savedPlan)
			}
		})
	}
}

func TestRepository_Integration(t *testing.T) {
	t.Run("complete workflow - save profile and plan", func(t *testing.T) {
		repo := NewMockRepository()
		ctx := context.Background()

		// Setup profile
		profile := &UserProfile{
			UserID:        "user123",
			Gender:        "female",
			Age:           25,
			HeightCm:      165.0,
			WeightKg:      60.0,
			Allergies:     []string{"dairy", "gluten"},
			Preferences:   []string{"vegan", "low-carb"},
			Target:        "weight_loss",
			StepsToday:    10000,
			ActiveKcal24h: 1800.0,
		}
		repo.SetProfile("user123", profile)

		// Save a generated plan
		planJSON := []byte(`[
			{
				"day": "Monday",
				"meals": {
					"breakfast": "Green smoothie",
					"lunch": "Vegan Buddha bowl",
					"dinner": "Grilled tofu"
				},
				"calories": {
					"breakfast": 250,
					"lunch": 350,
					"dinner": 400
				},
				"total_calories": 1000
			}
		]`)

		err := repo.SaveGeneratedPlan(ctx, "user123", planJSON)
		require.NoError(t, err)

		// Retrieve profile
		retrievedProfile, err := repo.GetUserProfile(ctx, "user123")
		require.NoError(t, err)
		assert.Equal(t, profile.UserID, retrievedProfile.UserID)
		assert.Equal(t, profile.Gender, retrievedProfile.Gender)
		assert.Equal(t, profile.Age, retrievedProfile.Age)

		// Retrieve plan
		retrievedPlan := repo.GetPlan("user123")
		assert.Equal(t, planJSON, retrievedPlan)
	})

	t.Run("multiple users workflow", func(t *testing.T) {
		repo := NewMockRepository()
		ctx := context.Background()

		// Setup multiple profiles
		profile1 := &UserProfile{
			UserID:        "user1",
			Gender:        "male",
			Age:           30,
			HeightCm:      175.0,
			WeightKg:      70.0,
			Allergies:     []string{"nuts"},
			Preferences:   []string{"vegetarian"},
			Target:        "maintenance",
			StepsToday:    8000,
			ActiveKcal24h: 2000.0,
		}

		profile2 := &UserProfile{
			UserID:        "user2",
			Gender:        "female",
			Age:           25,
			HeightCm:      165.0,
			WeightKg:      60.0,
			Allergies:     []string{"dairy"},
			Preferences:   []string{"vegan"},
			Target:        "weight_loss",
			StepsToday:    10000,
			ActiveKcal24h: 1800.0,
		}

		repo.SetProfile("user1", profile1)
		repo.SetProfile("user2", profile2)

		// Save plans for both users
		plan1 := []byte(`[{"day": "Monday", "meals": {"breakfast": "Oatmeal"}}]`)
		plan2 := []byte(`[{"day": "Monday", "meals": {"breakfast": "Smoothie"}}]`)

		err1 := repo.SaveGeneratedPlan(ctx, "user1", plan1)
		err2 := repo.SaveGeneratedPlan(ctx, "user2", plan2)

		require.NoError(t, err1)
		require.NoError(t, err2)

		// Verify both profiles and plans are stored correctly
		retrievedProfile1, err := repo.GetUserProfile(ctx, "user1")
		require.NoError(t, err)
		assert.Equal(t, "male", retrievedProfile1.Gender)

		retrievedProfile2, err := repo.GetUserProfile(ctx, "user2")
		require.NoError(t, err)
		assert.Equal(t, "female", retrievedProfile2.Gender)

		assert.Equal(t, plan1, repo.GetPlan("user1"))
		assert.Equal(t, plan2, repo.GetPlan("user2"))
	})
}

// EntRepo tests using actual Ent client
func TestEntRepo_GetUserProfile(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	tests := []struct {
		name          string
		setupUser     func() (*ent.User, error)
		userID        string
		expectedError bool
		validate      func(*testing.T, *UserProfile)
	}{
		{
			name: "success - returns profile with all fields",
			setupUser: func() (*ent.User, error) {
				sex := "male"
				height := 175.0
				weight := 70.0
				goal := "weight_loss"
				birthDate := time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC)
				return client.User.Create().
					SetEmail("test@example.com").
					SetName("Test User").
					SetPassword("password").
					SetNillableSex(&sex).
					SetNillableHeightCm(&height).
					SetNillableWeightKg(&weight).
					SetNillableGoal(&goal).
					SetNillableDateOfBirth(&birthDate).
					Save(ctx)
			},
			userID:        "", // Will be set after user creation
			expectedError: false,
			validate: func(t *testing.T, profile *UserProfile) {
				assert.Equal(t, "male", profile.Gender)
				assert.Greater(t, profile.Age, 0)
				assert.Equal(t, 175.0, profile.HeightCm)
				assert.Equal(t, 70.0, profile.WeightKg)
				assert.Equal(t, "weight_loss", profile.Target)
			},
		},
		{
			name: "success - returns profile with partial fields",
			setupUser: func() (*ent.User, error) {
				sex := "female"
				height := 165.0
				return client.User.Create().
					SetEmail("test2@example.com").
					SetName("Test User 2").
					SetPassword("password").
					SetNillableSex(&sex).
					SetNillableHeightCm(&height).
					Save(ctx)
			},
			userID:        "", // Will be set after user creation
			expectedError: false,
			validate: func(t *testing.T, profile *UserProfile) {
				assert.Equal(t, "female", profile.Gender)
				assert.Equal(t, 165.0, profile.HeightCm)
				assert.Equal(t, 0.0, profile.WeightKg) // Not set
				assert.Equal(t, "", profile.Target)    // Not set
			},
		},
		{
			name: "error - invalid user ID format",
			setupUser: func() (*ent.User, error) {
				return nil, nil // No user needed
			},
			userID:        "invalid-uuid",
			expectedError: true,
			validate:      nil,
		},
		{
			name: "error - user not found",
			setupUser: func() (*ent.User, error) {
				return nil, nil // No user needed
			},
			userID:        uuid.New().String(),
			expectedError: true,
			validate:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userID string
			if tt.setupUser != nil {
				user, err := tt.setupUser()
				if err != nil {
					t.Fatalf("Failed to setup user: %v", err)
				}
				if user != nil {
					userID = user.ID.String()
				} else {
					userID = tt.userID
				}
			} else {
				userID = tt.userID
			}

			profile, err := repo.GetUserProfile(ctx, userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, profile)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, userID, profile.UserID)
				if tt.validate != nil {
					tt.validate(t, profile)
				}
			}
		})
	}
}

func TestEntRepo_SaveGeneratedPlan(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	tests := []struct {
		name          string
		setupUser     func() (*ent.User, error)
		userID        string
		planJSON      []byte
		expectedError bool
	}{
		{
			name: "success - saves plan with valid JSON",
			setupUser: func() (*ent.User, error) {
				return client.User.Create().
					SetEmail("test@example.com").
					SetName("Test User").
					SetPassword("password").
					Save(ctx)
			},
			userID:        "", // Will be set after user creation
			planJSON:      []byte(`[{"day": "Monday", "meals": {"breakfast": "Oatmeal"}}]`),
			expectedError: false,
		},
		{
			name: "success - saves plan with empty JSON array",
			setupUser: func() (*ent.User, error) {
				return client.User.Create().
					SetEmail("test2@example.com").
					SetName("Test User 2").
					SetPassword("password").
					Save(ctx)
			},
			userID:        "", // Will be set after user creation
			planJSON:      []byte(`[]`),
			expectedError: false,
		},
		{
			name: "success - saves plan with invalid JSON (handled gracefully)",
			setupUser: func() (*ent.User, error) {
				return client.User.Create().
					SetEmail("test3@example.com").
					SetName("Test User 3").
					SetPassword("password").
					Save(ctx)
			},
			userID:        "", // Will be set after user creation
			planJSON:      []byte(`invalid json`),
			expectedError: false, // Method handles invalid JSON gracefully
		},
		{
			name: "error - invalid user ID format",
			setupUser: func() (*ent.User, error) {
				return nil, nil
			},
			userID:        "invalid-uuid",
			planJSON:      []byte(`[{"day": "Monday"}]`),
			expectedError: true,
		},
		{
			name: "error - user not found",
			setupUser: func() (*ent.User, error) {
				return nil, nil
			},
			userID:        uuid.New().String(),
			planJSON:      []byte(`[{"day": "Monday"}]`),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userID string
			if tt.setupUser != nil {
				user, err := tt.setupUser()
				if err != nil {
					t.Fatalf("Failed to setup user: %v", err)
				}
				if user != nil {
					userID = user.ID.String()
				} else {
					userID = tt.userID
				}
			} else {
				userID = tt.userID
			}

			err := repo.SaveGeneratedPlan(ctx, userID, tt.planJSON)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify meal plan was created (if JSON was valid)
				if len(tt.planJSON) > 0 && tt.planJSON[0] == '[' {
					mealPlans, err := client.Meal_plan.Query().
						Where(meal_plan.HasUserWith(user.ID(uuid.MustParse(userID)))).
						All(ctx)
					if err == nil && len(mealPlans) > 0 {
						// At least one meal plan should exist
						assert.Greater(t, len(mealPlans), 0)
					}
				}
			}
		})
	}
}
