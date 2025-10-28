package genai

import "context"

// Repository is where you read/write persistent stuff (user profile, prefs, logs).
type Repository interface {
	// Example hooks you can actually use:
	GetUserProfile(ctx context.Context, userID string) (*UserProfile, error)
	SaveGeneratedPlan(ctx context.Context, userID string, planJSON []byte) error
}

type UserProfile struct {
	UserID        string
	Gender        string
	Age           int
	HeightCm      float64
	WeightKg      float64
	Allergies     []string
	Preferences   []string
	Target        string
	StepsToday    int
	ActiveKcal24h float64
}
