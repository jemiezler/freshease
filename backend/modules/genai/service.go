package genai

import (
	"context"
	"encoding/json"
	"fmt"
)

type Service interface {
	GenerateWeeklyMeals(ctx context.Context, req *GenerateMealsReq) (any, error)
	GenerateDailyMeals(ctx context.Context, req *GenerateMealsReq) (any, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

type GenerateMealsReq struct {
	UserID        string   `json:"user_id,omitempty"`
	Gender        string   `json:"gender"`
	Age           int      `json:"age"`
	HeightCm      float64  `json:"height_cm"`
	WeightKg      float64  `json:"weight_kg"`
	StepsToday    int      `json:"steps_today"`
	ActiveKcal24h float64  `json:"active_kcal_24h"`
	Allergies     []string `json:"allergies,omitempty"`
	Preferences   []string `json:"preferences,omitempty"`
	Target        string   `json:"target"`
}

func (s *service) GenerateWeeklyMeals(ctx context.Context, req *GenerateMealsReq) (any, error) {
	// If any profile fields missing, try loading from repo:
	if (req.Gender == "" || req.Age == 0 || req.HeightCm == 0 || req.WeightKg == 0) && s.repo != nil && req.UserID != "" {
		if prof, err := s.repo.GetUserProfile(ctx, req.UserID); err == nil && prof != nil {
			if req.Gender == "" {
				req.Gender = prof.Gender
			}
			if req.Age == 0 {
				req.Age = prof.Age
			}
			if req.HeightCm == 0 {
				req.HeightCm = prof.HeightCm
			}
			if req.WeightKg == 0 {
				req.WeightKg = prof.WeightKg
			}
			if len(req.Allergies) == 0 {
				req.Allergies = prof.Allergies
			}
			if len(req.Preferences) == 0 {
				req.Preferences = prof.Preferences
			}
		}
	}

	raw, err := WeeklyMealsGenerator(req.Gender, req.Age, req.HeightCm, req.WeightKg)
	if err != nil {
		return nil, fmt.Errorf("generate failed: %w", err)
	}

	// validate JSON array of 7 days
	var week []map[string]any
	if err := json.Unmarshal(raw, &week); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Optionally persist
	if s.repo != nil && req.UserID != "" {
		_ = s.repo.SaveGeneratedPlan(ctx, req.UserID, raw)
	}

	return map[string]any{
		"steps_today":     req.StepsToday,
		"active_kcal_24h": req.ActiveKcal24h,
		"plan":            week,
	}, nil
}

func (s *service) GenerateDailyMeals(ctx context.Context, req *GenerateMealsReq) (any, error) {
	// If any profile fields missing, try loading from repo:
	if (req.Gender == "" || req.Age == 0 || req.HeightCm == 0 || req.WeightKg == 0 || req.ActiveKcal24h == 0 || req.Target == "" || req.StepsToday == 0) && s.repo != nil && req.UserID != "" {
		if prof, err := s.repo.GetUserProfile(ctx, req.UserID); err == nil && prof != nil {
			if req.Gender == "" {
				req.Gender = prof.Gender
			}
			if req.Age == 0 {
				req.Age = prof.Age
			}
			if req.HeightCm == 0 {
				req.HeightCm = prof.HeightCm
			}
			if req.WeightKg == 0 {
				req.WeightKg = prof.WeightKg
			}
			if req.StepsToday == 0 {
				req.StepsToday = prof.StepsToday
			}
			if req.ActiveKcal24h == 0 {
				req.ActiveKcal24h = prof.ActiveKcal24h
			}

			if len(req.Allergies) == 0 {
				req.Allergies = prof.Allergies
			}
			if len(req.Preferences) == 0 {
				req.Preferences = prof.Preferences
			}
		}
	}

	raw, err := DailyMealsGenerator(req.Gender, req.Age, req.HeightCm, req.WeightKg, req.Target, req.StepsToday, int(req.ActiveKcal24h))
	if err != nil {
		return nil, fmt.Errorf("generate failed: %w", err)
	}

	// validate JSON array of 7 days
	var week []map[string]any
	if err := json.Unmarshal(raw, &week); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Optionally persist
	if s.repo != nil && req.UserID != "" {
		_ = s.repo.SaveGeneratedPlan(ctx, req.UserID, raw)
	}

	return map[string]any{
		"steps_today":     req.StepsToday,
		"active_kcal_24h": req.ActiveKcal24h,
		"plan":            week,
	}, nil
}
