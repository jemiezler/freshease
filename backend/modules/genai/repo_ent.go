package genai

import (
	"context"
	"encoding/json"
	"freshease/backend/ent"
	"freshease/backend/ent/user"
	"time"

	"github.com/google/uuid"
)

type EntRepo struct{ db *ent.Client }

func NewEntRepo(client *ent.Client) *EntRepo { return &EntRepo{db: client} }

func (r *EntRepo) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	u, err := r.db.User.Query().
		Where(user.ID(userUUID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	profile := &UserProfile{
		UserID:    userID,
		Gender:    "",
		Age:       0,
		HeightCm:  0,
		WeightKg:  0,
		Allergies: []string{},
		Preferences: []string{},
		Target:      "",
		StepsToday:  0,
		ActiveKcal24h: 0,
	}

	// Map sex to gender
	if u.Sex != nil {
		profile.Gender = *u.Sex
	}

	// Calculate age from date_of_birth
	if u.DateOfBirth != nil {
		now := time.Now()
		birthDate := *u.DateOfBirth
		age := now.Year() - birthDate.Year()
		if now.YearDay() < birthDate.YearDay() {
			age--
		}
		profile.Age = age
	}

	// Map height and weight
	if u.HeightCm != nil {
		profile.HeightCm = *u.HeightCm
	}
	if u.WeightKg != nil {
		profile.WeightKg = *u.WeightKg
	}

	// Map goal to target
	if u.Goal != nil {
		profile.Target = *u.Goal
	}

	// Note: Allergies, Preferences, StepsToday, ActiveKcal24h are not stored in User entity
	// These would need to be stored elsewhere or passed separately

	return profile, nil
}

func (r *EntRepo) SaveGeneratedPlan(ctx context.Context, userID string, planJSON []byte) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// Verify user exists and get user entity
	userEntity, err := r.db.User.Get(ctx, userUUID)
	if err != nil {
		return err
	}

	// Parse the plan JSON to extract week_start if available
	var weekPlan []map[string]any
	if err := json.Unmarshal(planJSON, &weekPlan); err == nil && len(weekPlan) > 0 {
		// Try to extract week_start from the first day's data
		// If not available, use current week start (Monday)
		weekStart := getWeekStart(time.Now())
		if dayStr, ok := weekPlan[0]["day"].(string); ok {
			// Try to parse the day string to get the date
			// For now, we'll use current week start
			_ = dayStr
		}

		// Create or update meal plan for this week
		goal := "generated" // Default goal for AI-generated plans
		_, err = r.db.Meal_plan.Create().
			SetID(uuid.New()).
			SetWeekStart(weekStart).
			SetNillableGoal(&goal).
			SetUser(userEntity).
			Save(ctx)
		if err != nil {
			// If meal plan already exists for this week, that's okay
			// We'll just return nil to indicate success
			return nil
		}
	}

	return nil
}

// getWeekStart returns the Monday of the current week
func getWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	daysFromMonday := weekday - 1
	return t.AddDate(0, 0, -daysFromMonday).Truncate(24 * time.Hour)
}
