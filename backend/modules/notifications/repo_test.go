package notifications

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
	notifications, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, notifications)

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create notifications
	notification1, err := client.Notification.Create().
		SetID(uuid.New()).
		SetTitle("Notification 1").
		SetBody("Body 1").
		SetChannel("email").
		SetStatus("unread").
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	notification2, err := client.Notification.Create().
		SetID(uuid.New()).
		SetTitle("Notification 2").
		SetBody("Body 2").
		SetChannel("push").
		SetStatus("read").
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	notifications, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, notifications, 2)

	// Verify notifications are returned
	notificationMap := make(map[uuid.UUID]*GetNotificationDTO)
	for _, notification := range notifications {
		notificationMap[notification.ID] = notification
	}

	assert.Contains(t, notificationMap, notification1.ID)
	assert.Contains(t, notificationMap, notification2.ID)
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

	// Create test notification
	createDTO := &CreateNotificationDTO{
		ID:      uuid.New(),
		Title:   "Test Notification",
		Body:    stringPtr("Test Body"),
		Channel: "email",
		Status:  "unread",
		UserID:  user.ID,
	}
	notification, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing notification
	foundNotification, err := repo.FindByID(ctx, notification.ID)
	require.NoError(t, err)
	assert.Equal(t, notification.ID, foundNotification.ID)
	assert.Equal(t, notification.Title, foundNotification.Title)
	assert.Equal(t, notification.Body, foundNotification.Body)
	assert.Equal(t, notification.Channel, foundNotification.Channel)
	assert.Equal(t, notification.Status, foundNotification.Status)
	assert.Equal(t, notification.UserID, foundNotification.UserID)

	// Test notification not found
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

	// Test creating new notification
	now := time.Now()
	createDTO := &CreateNotificationDTO{
		ID:        uuid.New(),
		Title:     "New Notification",
		Body:      stringPtr("New Body"),
		Channel:   "email",
		Status:    "unread",
		UserID:    user.ID,
		CreatedAt: &now,
	}
	createdNotification, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.NotNil(t, createdNotification)
	assert.Equal(t, createDTO.ID, createdNotification.ID)
	assert.Equal(t, createDTO.Title, createdNotification.Title)
	assert.Equal(t, createDTO.Body, createdNotification.Body)
	assert.Equal(t, createDTO.Channel, createdNotification.Channel)
	assert.Equal(t, createDTO.Status, createdNotification.Status)
	assert.Equal(t, createDTO.UserID, createdNotification.UserID)

	// Test Create - without Body
	createDTO2 := &CreateNotificationDTO{
		ID:        uuid.New(),
		Title:     "Notification Without Body",
		Body:      nil,
		Channel:   "push",
		Status:    "unread",
		UserID:    user.ID,
		CreatedAt: nil,
	}
	createdNotification2, err := repo.Create(ctx, createDTO2)
	require.NoError(t, err)
	assert.NotNil(t, createdNotification2)
	assert.Nil(t, createdNotification2.Body)

	// Test Create - without CreatedAt
	createDTO3 := &CreateNotificationDTO{
		ID:        uuid.New(),
		Title:     "Notification Without CreatedAt",
		Body:      stringPtr("Body text"),
		Channel:   "sms",
		Status:    "read",
		UserID:    user.ID,
		CreatedAt: nil,
	}
	createdNotification3, err := repo.Create(ctx, createDTO3)
	require.NoError(t, err)
	assert.NotNil(t, createdNotification3)
	assert.NotZero(t, createdNotification3.CreatedAt)

	// Test Create - error: user not found
	nonExistentUserID := uuid.New()
	createDTO4 := &CreateNotificationDTO{
		ID:        uuid.New(),
		Title:     "Notification",
		Body:      stringPtr("Body"),
		Channel:   "email",
		Status:    "unread",
		UserID:    nonExistentUserID,
		CreatedAt: nil,
	}
	_, err = repo.Create(ctx, createDTO4)
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

	// Create test notification
	createDTO := &CreateNotificationDTO{
		ID:      uuid.New(),
		Title:   "Original Notification",
		Body:    stringPtr("Original Body"),
		Channel: "email",
		Status:  "unread",
		UserID:  user.ID,
	}
	notification, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test updating notification - update Title, Body, Channel, Status
	newTitle := "Updated Notification"
	newBody := "Updated Body"
	newChannel := "push"
	newStatus := "read"
	updateDTO := &UpdateNotificationDTO{
		ID:      notification.ID,
		Title:   &newTitle,
		Body:    &newBody,
		Channel: &newChannel,
		Status:  &newStatus,
	}
	updatedNotification, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.NotNil(t, updatedNotification)
	assert.Equal(t, *updateDTO.Title, updatedNotification.Title)
	require.NotNil(t, updatedNotification.Body)
	assert.Equal(t, *updateDTO.Body, *updatedNotification.Body)
	assert.Equal(t, *updateDTO.Channel, updatedNotification.Channel)
	assert.Equal(t, *updateDTO.Status, updatedNotification.Status)

	// Test no fields to update
	noUpdateDTO := &UpdateNotificationDTO{ID: notification.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())
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

	// Create test notification
	createDTO := &CreateNotificationDTO{
		ID:      uuid.New(),
		Title:   "To Delete",
		Body:    stringPtr("Delete Body"),
		Channel: "email",
		Status:  "unread",
		UserID:  user.ID,
	}
	notification, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test deleting notification
	err = repo.Delete(ctx, notification.ID)
	require.NoError(t, err)

	// Verify notification is deleted
	_, err = repo.FindByID(ctx, notification.ID)
	assert.Error(t, err)
}

