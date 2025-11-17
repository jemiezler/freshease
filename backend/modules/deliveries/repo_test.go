package deliveries

import (
	"context"
	"testing"
	"time"

	"freshease/backend/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestEntRepo_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Create test deliveries
	delivery1, err := client.Delivery.Create().
		SetID(uuid.New()).
		SetProvider("fedex").
		SetStatus("pending").
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	delivery2, err := client.Delivery.Create().
		SetID(uuid.New()).
		SetProvider("ups").
		SetStatus("in_transit").
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results
	foundIDs := make(map[uuid.UUID]bool)
	for _, delivery := range result {
		foundIDs[delivery.ID] = true
		assert.NotEmpty(t, delivery.Provider)
		assert.NotEmpty(t, delivery.Status)
		assert.Equal(t, order.ID, delivery.OrderID)
	}

	assert.True(t, foundIDs[delivery1.ID])
	assert.True(t, foundIDs[delivery2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Create test delivery
	trackingNo := "TRACK-001"
	createdDelivery, err := client.Delivery.Create().
		SetID(uuid.New()).
		SetProvider("fedex").
		SetNillableTrackingNo(&trackingNo).
		SetStatus("pending").
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdDelivery.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdDelivery.ID, result.ID)
	assert.Equal(t, "fedex", result.Provider)
	assert.Equal(t, "pending", result.Status)
	assert.NotNil(t, result.TrackingNo)
	assert.Equal(t, "TRACK-001", *result.TrackingNo)
	assert.Equal(t, order.ID, result.OrderID)

	// Test FindByID - not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestEntRepo_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	trackingNo := "TRACK-002"
	eta := time.Now().Add(24 * time.Hour)
	dto := &CreateDeliveryDTO{
		ID:         uuid.New(),
		Provider:   "fedex",
		TrackingNo: &trackingNo,
		Status:     "pending",
		Eta:        &eta,
		OrderID:    order.ID,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Provider, result.Provider)
	assert.Equal(t, dto.Status, result.Status)
	assert.NotNil(t, result.TrackingNo)
	assert.Equal(t, trackingNo, *result.TrackingNo)
	assert.NotNil(t, result.Eta)
	assert.Equal(t, order.ID, result.OrderID)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Create test delivery
	createdDelivery, err := client.Delivery.Create().
		SetID(uuid.New()).
		SetProvider("fedex").
		SetStatus("pending").
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Update delivery
	newStatus := "in_transit"
	newTrackingNo := "TRACK-UPDATED"
	dto := &UpdateDeliveryDTO{
		ID:         createdDelivery.ID,
		Status:     &newStatus,
		TrackingNo: &newTrackingNo,
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdDelivery.ID, result.ID)
	assert.Equal(t, "in_transit", result.Status)
	assert.NotNil(t, result.TrackingNo)
	assert.Equal(t, newTrackingNo, *result.TrackingNo)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Create test delivery
	createdDelivery, err := client.Delivery.Create().
		SetID(uuid.New()).
		SetProvider("fedex").
		SetStatus("pending").
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Delete delivery
	err = repo.Delete(ctx, createdDelivery.ID)
	require.NoError(t, err)

	// Verify delivery is deleted
	_, err = repo.FindByID(ctx, createdDelivery.ID)
	assert.Error(t, err)
}

