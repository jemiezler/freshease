package meal_plans

import (
	"context"
	"testing"
	"time"

	"freshease/backend/ent/enttest"
	"freshease/backend/internal/common/errs"
	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Test empty list
	mealPlans, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, mealPlans)

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create meal plans
	weekStart1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal1 := "weight_loss"
	mealPlan1, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(weekStart1).
		SetGoal(goal1).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	weekStart2 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	goal2 := "muscle_gain"
	mealPlan2, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(weekStart2).
		SetGoal(goal2).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	mealPlans, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, mealPlans, 2)

	// Verify meal plans are returned
	mealPlanMap := make(map[uuid.UUID]*GetMeal_planDTO)
	for _, mealPlan := range mealPlans {
		mealPlanMap[mealPlan.ID] = mealPlan
	}

	assert.Contains(t, mealPlanMap, mealPlan1.ID)
	assert.Contains(t, mealPlanMap, mealPlan2.ID)
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test meal plan
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal := "weight_loss"
	createDTO := &CreateMeal_planDTO{
		ID:        uuid.New(),
		WeekStart: weekStart,
		Goal:      &goal,
		UserID:    user.ID,
	}
	mealPlan, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing meal plan
	foundMealPlan, err := repo.FindByID(ctx, mealPlan.ID)
	require.NoError(t, err)
	assert.Equal(t, mealPlan.ID, foundMealPlan.ID)
	assert.Equal(t, mealPlan.WeekStart, foundMealPlan.WeekStart)
	assert.Equal(t, mealPlan.Goal, foundMealPlan.Goal)
	assert.Equal(t, mealPlan.UserID, foundMealPlan.UserID)

	// Test meal plan not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestRepository_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Test creating new meal plan
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal := "weight_loss"
	createDTO := &CreateMeal_planDTO{
		ID:        uuid.New(),
		WeekStart: weekStart,
		Goal:      &goal,
		UserID:    user.ID,
	}
	createdMealPlan, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.NotNil(t, createdMealPlan)
	assert.Equal(t, createDTO.ID, createdMealPlan.ID)
	assert.Equal(t, createDTO.WeekStart, createdMealPlan.WeekStart)
	assert.Equal(t, *createDTO.Goal, *createdMealPlan.Goal)
	assert.Equal(t, createDTO.UserID, createdMealPlan.UserID)

	// Test creating meal plan without goal
	createDTO2 := &CreateMeal_planDTO{
		ID:        uuid.New(),
		WeekStart: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		Goal:      nil,
		UserID:    user.ID,
	}
	createdMealPlan2, err := repo.Create(ctx, createDTO2)
	require.NoError(t, err)
	assert.NotNil(t, createdMealPlan2)
	assert.Nil(t, createdMealPlan2.Goal)

	// Test Create - error: user not found
	nonExistentUserID := uuid.New()
	createDTO3 := &CreateMeal_planDTO{
		ID:        uuid.New(),
		WeekStart: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Goal:      &goal,
		UserID:    nonExistentUserID,
	}
	_, err = repo.Create(ctx, createDTO3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test meal plan
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal := "weight_loss"
	createDTO := &CreateMeal_planDTO{
		ID:        uuid.New(),
		WeekStart: weekStart,
		Goal:      &goal,
		UserID:    user.ID,
	}
	mealPlan, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test updating meal plan - update WeekStart and Goal
	newWeekStart := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	newGoal := "muscle_gain"
	updateDTO := &UpdateMeal_planDTO{
		ID:        mealPlan.ID,
		WeekStart: &newWeekStart,
		Goal:      &newGoal,
	}
	updatedMealPlan, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.NotNil(t, updatedMealPlan)
	assert.Equal(t, *updateDTO.WeekStart, updatedMealPlan.WeekStart)
	assert.Equal(t, *updateDTO.Goal, *updatedMealPlan.Goal)

	// Test updating meal plan - update only WeekStart
	newWeekStart2 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	updateDTO2 := &UpdateMeal_planDTO{
		ID:        mealPlan.ID,
		WeekStart: &newWeekStart2,
	}
	updatedMealPlan2, err := repo.Update(ctx, updateDTO2)
	require.NoError(t, err)
	assert.Equal(t, *updateDTO2.WeekStart, updatedMealPlan2.WeekStart)

	// Test updating meal plan - update only Goal
	newGoal2 := "maintenance"
	updateDTO3 := &UpdateMeal_planDTO{
		ID:   mealPlan.ID,
		Goal: &newGoal2,
	}
	updatedMealPlan3, err := repo.Update(ctx, updateDTO3)
	require.NoError(t, err)
	assert.Equal(t, *updateDTO3.Goal, *updatedMealPlan3.Goal)

	// Test no fields to update
	noUpdateDTO := &UpdateMeal_planDTO{ID: mealPlan.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())

	// Test Update - error: meal plan not found
	nonExistentID := uuid.New()
	updateDTO4 := &UpdateMeal_planDTO{
		ID:        nonExistentID,
		WeekStart: &newWeekStart,
	}
	_, err = repo.Update(ctx, updateDTO4)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test meal plan
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal := "weight_loss"
	createDTO := &CreateMeal_planDTO{
		ID:        uuid.New(),
		WeekStart: weekStart,
		Goal:      &goal,
		UserID:    user.ID,
	}
	mealPlan, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test deleting meal plan
	err = repo.Delete(ctx, mealPlan.ID)
	require.NoError(t, err)

	// Verify meal plan is deleted
	_, err = repo.FindByID(ctx, mealPlan.ID)
	assert.Error(t, err)
}

