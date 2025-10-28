package genai

import (
	"context"
	"freshease/backend/ent"
)

type EntRepo struct{ db *ent.Client }

func NewEntRepo(client *ent.Client) *EntRepo { return &EntRepo{db: client} }

func (r *EntRepo) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	// TODO: query User table with ent
	// return &UserProfile{UserID: userID, Gender: "male", Age: 22, HeightCm: 175, WeightKg: 70}, nil
	return nil, nil
}

func (r *EntRepo) SaveGeneratedPlan(ctx context.Context, userID string, planJSON []byte) error {
	// TODO: insert into a 'meal_plans' table with ent
	return nil
}
