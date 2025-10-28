package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"freshease/backend/internal/common/config"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func DailyMealsGenerator(gender string, age int, heightCm, weightKg float64, target string, steps_today int, active_kcal_24h int) ([]byte, error) {
	// Log function start and input parameters (best practice)
	log.Printf("INFO: DailyMealsGenerator started with target=%s. Profile: gender=%s, age=%d, weight_kg=%.1f.",
		target, gender, age, weightKg)

	cfg := config.Load()
	ctx := context.Background()
	apiKey := cfg.GENAI_APIKEY
	if apiKey == "" {
		log.Println("ERROR: GOOGLE_API_KEY is empty")
		return nil, fmt.Errorf("GOOGLE_API_KEY is empty")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("ERROR: Failed to create genai client: %v", err)
		return nil, fmt.Errorf("genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")

	// Key Logic: Determine the calorie adjustment instruction based on the target.
	var calorieInstruction string
	normalizedTarget := strings.ToLower(target)

	switch {
	case strings.Contains(normalizedTarget, "loss"):
		calorieInstruction = "Create a meal plan for MODERATE CALORIE DEFICIT (approx. 500 kcal/day). Prioritize high protein and fiber for satiety."
	case strings.Contains(normalizedTarget, "gain"):
		calorieInstruction = "Create a meal plan for a MODERATE CALORIE SURPLUS (approx. 300-500 kcal/day). Prioritize high protein and complex carbohydrates for muscle gain."
	default:
		calorieInstruction = "Create a meal plan for MAINTENANCE. Maintain a balanced calorie intake based on the profile and activity."
	}

	prompt := fmt.Sprintf(`
You are a professional nutrition planner. Your response must be a single JSON array with EXACTLY 1 item (for Monday).
Do not include markdown, code fences, comments, or any extra prose. Output JSON only.

The item must strictly follow this schema:
{
  "day": "Monday",
  "meals": { "breakfast": "...", "lunch": "...", "dinner": "..." },
  "calories": { "breakfast": 0, "lunch": 0, "dinner": 0 },
  "total_calories": 0
}

---
User Profile and Goal:
- Gender: %s, Age: %d, Height: %.1f cm, Weight: %.1f kg
- Activity: Steps Today=%d, Active Kcal (24h)=%d
- Goal: %s

***Instruction: %s Adjust calories and portions based on the profile and goal.***
`, gender, age, heightCm, weightKg, steps_today, active_kcal_24h, target, calorieInstruction)

	log.Println("INFO: Calling Gemini API with adjusted prompt.")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("ERROR: Gemini API call failed: %v", err)
		return nil, fmt.Errorf("generate content: %w", err)
	}

	var sb strings.Builder
	for i, cand := range resp.Candidates {
		if cand == nil || cand.Content == nil {
			log.Printf("WARN: Candidate %d is nil or has nil content", i)
			continue
		}
		for _, part := range cand.Content.Parts {
			if t, ok := part.(genai.Text); ok {
				sb.WriteString(string(t))
			}
		}
	}
	out := strings.TrimSpace(sb.String())
	if out == "" {
		log.Println("ERROR: No text returned from model, empty output.")
		return nil, fmt.Errorf("no text returned from model")
	}

	// JSON Cleanup (Crucial for model stability)
	out = strings.TrimPrefix(out, "```json")
	out = strings.TrimPrefix(out, "```")
	out = strings.TrimSuffix(out, "```")
	out = strings.TrimSpace(out)

	// Robustness fix: If the model returns a single object (not an array), wrap it.
	if strings.HasPrefix(out, "{") && strings.HasSuffix(out, "}") {
		log.Println("WARN: Model returned a single JSON object. Wrapping in brackets.")
		out = fmt.Sprintf("[%s]", out)
	}

	// Validate it parses as array
	var week []map[string]any
	if err := json.Unmarshal([]byte(out), &week); err != nil {
		log.Printf("FATAL: Model JSON validation failed: %v", err)
		return nil, fmt.Errorf("model did not return valid JSON: %w\nraw: %s", err, out)
	}

	log.Println("INFO: DailyMealsGenerator finished successfully.")
	return []byte(out), nil
}

func WeeklyMealsGenerator(gender string, age int, heightCm, weightKg float64) ([]byte, error) {
	cfg := config.Load()
	ctx := context.Background()
	apiKey := cfg.GENAI_APIKEY
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY is empty")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	prompt := fmt.Sprintf(`
You are a nutrition planner. Return a strict JSON array of 7 items (Mon..Sun), each:
{
  "day": "Monday",
  "meals": { "breakfast": "...", "lunch": "...", "dinner": "..." },
  "calories": { "breakfast": 0, "lunch": 0, "dinner": 0 }
}
Profile: gender=%s, age=%d, height_cm=%.1f, weight_kg=%.1f.
Do not include markdown, comments, or prose. Output JSON only.
`, gender, age, heightCm, weightKg)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("generate content: %w", err)
	}

	var sb strings.Builder
	for _, cand := range resp.Candidates {
		if cand == nil || cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			if t, ok := part.(genai.Text); ok {
				sb.WriteString(string(t))
			}
		}
	}
	out := strings.TrimSpace(sb.String())
	if out == "" {
		return nil, fmt.Errorf("no text returned from model")
	}

	// Validate it parses as array
	var week []map[string]any
	if err := json.Unmarshal([]byte(out), &week); err != nil {
		return nil, fmt.Errorf("model did not return valid JSON: %w\nraw: %s", err, out)
	}
	return []byte(out), nil
}
